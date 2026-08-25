package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"rubbercure/internal/audit"
	"rubbercure/internal/batch"
	"rubbercure/internal/console"
	"rubbercure/internal/cool"
	"rubbercure/internal/lift"
	"rubbercure/internal/mold"
	"rubbercure/internal/press"
	"rubbercure/internal/steam"
	"rubbercure/internal/store"
	"rubbercure/internal/temp"
	"rubbercure/internal/timer"
)

func main() {
	dataDir := os.Getenv("RUBBERCURE_DATA_DIR")
	if dataDir == "" {
		dataDir = filepath.Join(os.TempDir(), "rubbercure")
	}
	st, err := store.New(dataDir, "plant")
	if err != nil {
		log.Fatal(err)
	}
	recorder := audit.NewRecorder(st)
	fleet := press.NewFleet()
	registry := mold.NewRegistry()
	header := steam.NewHeader()
	steamSystem := steam.NewSystem(header)
	coolSystem := cool.NewCoolSystem()
	tempController := temp.NewController("press-a", header, st)
	cureTimer := timer.NewCureTimer()
	liftControl := lift.NewLift(2)
	scheduler := batch.NewScheduler(fleet, registry, cureTimer, recorder)

	fleet.Add(press.New("press-a", "一号硫化机"))
	fleet.Add(press.New("press-b", "二号硫化机"))
	registry.Register(mold.Mold{ID: "mold-a", Spec: "胎面模", SensorIDs: []string{"s1", "s2"}})
	registry.Register(mold.Mold{ID: "mold-b", Spec: "密封圈模", SensorIDs: []string{"s3", "s4"}})

	server := console.NewServer(
		fleet,
		registry,
		steamSystem,
		coolSystem,
		tempController,
		cureTimer,
		liftControl,
		scheduler,
		recorder,
		st,
	)
	addr := os.Getenv("RUBBERCURE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	log.Printf("rubbercure console listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
