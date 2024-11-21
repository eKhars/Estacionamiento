package usecases

import (
    "fmt"
    "sync/atomic"
    "time"

    "Estacionamiento/internal/domain/entities"
)

type ParkingSimulation struct {
    parkingLot *entities.ParkingLot
}

func NewParkingSimulation(parkingLot *entities.ParkingLot) *ParkingSimulation {
    return &ParkingSimulation{
        parkingLot: parkingLot,
    }
}

func (ps *ParkingSimulation) SimulateVehicle(id int) {
    defer ps.parkingLot.GetWaitGroup().Done()

    stats := ps.parkingLot.GetStats()
    atomic.AddInt32(&stats.WaitingToEnter, 1)
    ps.parkingLot.NotifyStatsObservers(*stats)

    var spot int
    
    for {
        ps.parkingLot.GetGateDirectionMutex().Lock()
        if ps.parkingLot.IsExiting() {
            ps.parkingLot.GetGateDirectionMutex().Unlock()
            time.Sleep(100 * time.Millisecond)
            continue
        }
        ps.parkingLot.SetIsEntering(true)
        ps.parkingLot.GetGateDirectionMutex().Unlock()

        ps.parkingLot.GetGateMutex().Lock()
        atomic.AddInt32(&stats.WaitingToEnter, -1)
        ps.parkingLot.NotifyGateObservers(true)

        ps.parkingLot.GetMutex().Lock()
        var found bool
        spot, found = ps.parkingLot.FindEmptySpot()
        
        if !found {
            ps.parkingLot.GetMutex().Unlock()
            ps.parkingLot.GetGateMutex().Unlock()
            ps.parkingLot.NotifyGateObservers(false)
            
            ps.parkingLot.GetGateDirectionMutex().Lock()
            ps.parkingLot.SetIsEntering(false)
            ps.parkingLot.GetGateDirectionMutex().Unlock()
            
            atomic.AddInt32(&stats.WaitingToEnter, 1)
            ps.parkingLot.NotifyStatsObservers(*stats)
            
            fmt.Printf("Vehículo %d: Esperando por espacio disponible\n", id)
            time.Sleep(100 * time.Millisecond)
            continue
        }

        ps.parkingLot.GetSpaces()[spot] = true
        atomic.AddInt32(&stats.CurrentlyParked, 1)
        ps.parkingLot.GetMutex().Unlock()

        time.Sleep(500 * time.Millisecond)
        ps.parkingLot.NotifyGateObservers(false)
        ps.parkingLot.GetGateMutex().Unlock()

        ps.parkingLot.GetGateDirectionMutex().Lock()
        ps.parkingLot.SetIsEntering(false)
        ps.parkingLot.GetGateDirectionMutex().Unlock()

        ps.parkingLot.NotifyObservers(spot, true)
        ps.parkingLot.NotifyStatsObservers(*stats)
        fmt.Printf("Vehículo %d: Estacionado en espacio %d\n", id, spot)
        break
    }

    stayTime := time.Duration(ps.parkingLot.GetRng().Intn(2000)+3000) * time.Millisecond
    time.Sleep(stayTime)

    atomic.AddInt32(&stats.WaitingToExit, 1)
    ps.parkingLot.NotifyStatsObservers(*stats)

    for {
        ps.parkingLot.GetGateDirectionMutex().Lock()
        if ps.parkingLot.IsEntering() {
            ps.parkingLot.GetGateDirectionMutex().Unlock()
            time.Sleep(100 * time.Millisecond)
            continue
        }
        ps.parkingLot.SetIsExiting(true)
        ps.parkingLot.GetGateDirectionMutex().Unlock()
        break
    }

    ps.parkingLot.GetGateMutex().Lock()
    atomic.AddInt32(&stats.WaitingToExit, -1)
    ps.parkingLot.NotifyGateObservers(true)

    ps.parkingLot.GetMutex().Lock()
    ps.parkingLot.GetSpaces()[spot] = false
    atomic.AddInt32(&stats.CurrentlyParked, -1)
    atomic.AddInt32(&stats.TotalProcessed, 1)
    ps.parkingLot.GetMutex().Unlock()
    
    ps.parkingLot.NotifyObservers(spot, false)
    ps.parkingLot.NotifyStatsObservers(*stats)
    fmt.Printf("Vehículo %d: Saliendo del espacio %d\n", id, spot)
    
    time.Sleep(500 * time.Millisecond)
    ps.parkingLot.NotifyGateObservers(false)
    ps.parkingLot.GetGateMutex().Unlock()

    ps.parkingLot.GetGateDirectionMutex().Lock()
    ps.parkingLot.SetIsExiting(false)
    ps.parkingLot.GetGateDirectionMutex().Unlock()
}