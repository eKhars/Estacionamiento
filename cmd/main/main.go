package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"Estacionamiento/internal/domain/entities"
	"Estacionamiento/internal/domain/models"
	"Estacionamiento/internal/infrastructure/gui"
	"Estacionamiento/internal/usecases"
)

func main() {
	parkingLot := entities.NewParkingLot()
	guiInstance := gui.NewGUI()
	parkingLot.RegisterObserver(guiInstance)

	simulation := usecases.NewParkingSimulation(parkingLot)
	done := make(chan bool)
	allProcessed := make(chan bool)

	go func() {
		stats := parkingLot.GetStats()

		for vehicleID := 1; vehicleID <= 100; vehicleID++ {
			parkingLot.GetWaitGroup().Add(1)
			go simulation.SimulateVehicle(vehicleID)

			parkingLot.UpdateStats(func(s *models.ParkingStats) {
				atomic.AddInt32(&s.TotalGenerated, 1)
			})

			time.Sleep(parkingLot.GeneratePoissonDelay())
		}
		fmt.Println("Todos los carros han sido generados")

		go func() {
			for {
				if atomic.LoadInt32(&stats.TotalProcessed) == 100 {
					allProcessed <- true
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()

		<-allProcessed
		fmt.Println("Todos los carros han sido procesados")
		done <- true
	}()

	go func() {
		<-done
		fmt.Printf("Carros procesados al finalizar: %d\n",
			atomic.LoadInt32(&parkingLot.GetStats().TotalProcessed))
		time.Sleep(2 * time.Second)
		guiInstance.GetWindow().Close()
	}()

	guiInstance.GetWindow().ShowAndRun()
}
