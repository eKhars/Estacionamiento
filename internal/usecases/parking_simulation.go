package usecases

import (
    "fmt"
    "sync/atomic"
    "time"

    "Estacionamiento/internal/domain/entities"
    "Estacionamiento/internal/domain/models"
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

    ps.parkingLot.UpdateStats(func(s *models.ParkingStats) {
        atomic.AddInt32(&s.WaitingToEnter, 1)
        ps.parkingLot.NotifyStatsObservers(*s)
    })

    ps.parkingLot.ParkingSem.Acquire()
    
    var spot int
    var found bool

    for {
        direction := ps.parkingLot.GetDirection()
        
        if direction == entities.Exiting {
            ps.parkingLot.SetDirection(direction)
            time.Sleep(100 * time.Millisecond)
            continue
        }

        ps.parkingLot.SetDirection(entities.Entering)
        
        ps.parkingLot.GateSem.Acquire()
        ps.parkingLot.UpdateStats(func(s *models.ParkingStats) {
            atomic.AddInt32(&s.WaitingToEnter, -1)
            ps.parkingLot.NotifyStatsObservers(*s)
        })

        ps.parkingLot.NotifyGateObservers(true)
        spot, found = ps.parkingLot.GetSpace()
        
        if !found {
            ps.parkingLot.GateSem.Release()
            ps.parkingLot.NotifyGateObservers(false)
            ps.parkingLot.SetDirection(entities.None)
            
            ps.parkingLot.UpdateStats(func(s *models.ParkingStats) {
                atomic.AddInt32(&s.WaitingToEnter, 1)
                ps.parkingLot.NotifyStatsObservers(*s)
            })
            
            fmt.Printf("Vehículo %d: Esperando por espacio disponible\n", id)
            time.Sleep(100 * time.Millisecond)
            continue
        }

        ps.parkingLot.UpdateStats(func(s *models.ParkingStats) {
            atomic.AddInt32(&s.CurrentlyParked, 1)
            ps.parkingLot.NotifyStatsObservers(*s)
        })
        time.Sleep(500 * time.Millisecond)
        
        ps.parkingLot.NotifyGateObservers(false)
        ps.parkingLot.GateSem.Release()
        ps.parkingLot.SetDirection(entities.None)

        ps.parkingLot.NotifyObservers(spot, true)
        fmt.Printf("Vehículo %d: Estacionado en espacio %d\n", id, spot)
        break
    }

    stayTime := time.Duration(ps.parkingLot.GetRng().Intn(2000)+3000) * time.Millisecond
    time.Sleep(stayTime)

    ps.parkingLot.UpdateStats(func(s *models.ParkingStats) {
        atomic.AddInt32(&s.WaitingToExit, 1)
        ps.parkingLot.NotifyStatsObservers(*s)
    })
    for {
        direction := ps.parkingLot.GetDirection()
        if direction == entities.Entering {
            ps.parkingLot.SetDirection(direction)
            time.Sleep(100 * time.Millisecond)
            continue
        }

        ps.parkingLot.SetDirection(entities.Exiting)
        break
    }

    ps.parkingLot.GateSem.Acquire()
    
    ps.parkingLot.UpdateStats(func(s *models.ParkingStats) {
        atomic.AddInt32(&s.WaitingToExit, -1)
        ps.parkingLot.NotifyStatsObservers(*s)
    })
    
    ps.parkingLot.NotifyGateObservers(true)

    ps.parkingLot.ReleaseSpace(spot)
    ps.parkingLot.UpdateStats(func(s *models.ParkingStats) {
        atomic.AddInt32(&s.CurrentlyParked, -1)
        atomic.AddInt32(&s.TotalProcessed, 1)
        ps.parkingLot.NotifyStatsObservers(*s)
    })

    ps.parkingLot.NotifyObservers(spot, false)
    fmt.Printf("Vehículo %d: Saliendo del espacio %d\n", id, spot)

    time.Sleep(500 * time.Millisecond)
    
    ps.parkingLot.NotifyGateObservers(false)
    ps.parkingLot.GateSem.Release()
    ps.parkingLot.SetDirection(entities.None)
    
    ps.parkingLot.ParkingSem.Release()
}