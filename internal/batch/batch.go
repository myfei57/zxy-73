package batch

import (
	"errors"
	"fmt"
	"sync"

	"rubbercure/internal/audit"
	"rubbercure/internal/mold"
	"rubbercure/internal/ns"
	"rubbercure/internal/press"
	"rubbercure/internal/timer"
)

type Assignment struct {
	PressID string
	MoldID  string
	Seq     int
}

type Scheduler struct {
	mu          sync.Mutex
	fleet       *press.Fleet
	registry    *mold.Registry
	timer       *timer.CureTimer
	audit       *audit.Recorder
	nextMold    int
	assignments []Assignment
	reportQueue []string
}

func NewScheduler(fleet *press.Fleet, registry *mold.Registry, cureTimer *timer.CureTimer, recorder *audit.Recorder) *Scheduler {
	return &Scheduler{
		fleet:    fleet,
		registry: registry,
		timer:    cureTimer,
		audit:    recorder,
	}
}

func (s *Scheduler) AssignMold(pressID string, moldID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.fleet.AssignMold(pressID, moldID); err != nil {
		return err
	}
	s.assignments = append(s.assignments, Assignment{PressID: pressID, MoldID: moldID, Seq: s.nextMold + 1})
	s.audit.Record(ns.ComponentBatch, ns.EventMoldAssigned, fmt.Sprintf("%s:%s", pressID, moldID))
	return nil
}

func (s *Scheduler) ReportDone(pressID string) error {
	p, err := s.fleet.Get(pressID)
	if err != nil {
		return err
	}
	if !p.Busy() {
		return errors.New("press is not running a cycle")
	}
	s.timer.Complete()
	s.timer.ResetLatch()
	if err := s.fleet.ReleaseMold(pressID); err != nil {
		return err
	}
	moldID := fmt.Sprintf("M%03d", s.nextMold)
	s.nextMold++
	if err := s.fleet.AssignMold(pressID, moldID); err != nil {
		return err
	}
	s.assignments = append(s.assignments, Assignment{PressID: pressID, MoldID: moldID, Seq: s.nextMold})
	s.audit.Record(ns.ComponentBatch, ns.EventCycleReported, pressID)
	return nil
}

func (s *Scheduler) Assignments() []Assignment {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Assignment, len(s.assignments))
	copy(out, s.assignments)
	return out
}

func (s *Scheduler) Registry() *mold.Registry {
	return s.registry
}

func (s *Scheduler) Fleet() *press.Fleet {
	return s.fleet
}

func (s *Scheduler) Timer() *timer.CureTimer {
	return s.timer
}
