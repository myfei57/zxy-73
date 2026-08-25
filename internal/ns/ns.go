package ns

type Component string

const (
	ComponentPress   Component = "press"
	ComponentMold    Component = "mold"
	ComponentLift    Component = "lift"
	ComponentBatch   Component = "batch"
	ComponentConsole Component = "console"
)

type EventType string

const (
	EventMoldRegistered EventType = "mold.registered"
	EventPressClosed    EventType = "lift.closed"
	EventCureStarted    EventType = "timer.started"
	EventCureDone       EventType = "timer.done"
	EventMoldOpened     EventType = "lift.opened"
	EventMoldAssigned   EventType = "batch.mold.assigned"
	EventCycleReported  EventType = "batch.cycle.reported"
	EventShutdown       EventType = "plant.shutdown"
)

func (e EventType) String() string {
	return string(e)
}
