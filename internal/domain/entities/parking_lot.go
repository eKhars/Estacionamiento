package entities

import (
    "math/rand"
    "sync"
    "time"

    "Estacionamiento/internal/domain/interfaces"
    "Estacionamiento/internal/domain/models"
)

type ParkingLot struct {
    observers      []interfaces.Observer
    spaces         [20]bool
    mutex          sync.Mutex
    gate           sync.Mutex
    wg             sync.WaitGroup
    stats          models.ParkingStats
    rng            *rand.Rand
    gateDirection  sync.Mutex
    isExiting      bool
    isEntering     bool
}

func NewParkingLot() *ParkingLot {
    return &ParkingLot{
        rng: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (p *ParkingLot) RegisterObserver(o interfaces.Observer) {
    p.observers = append(p.observers, o)
}

func (p *ParkingLot) RemoveObserver(o interfaces.Observer) {
    for i, observer := range p.observers {
        if observer == o {
            p.observers = append(p.observers[:i], p.observers[i+1:]...)
            break
        }
    }
}

func (p *ParkingLot) NotifyObservers(parkingSpot int, isOccupied bool) {
    for _, observer := range p.observers {
        observer.Update(parkingSpot, isOccupied)
    }
}

func (p *ParkingLot) NotifyGateObservers(isOccupied bool) {
    for _, observer := range p.observers {
        observer.UpdateGate(isOccupied)
    }
}

func (p *ParkingLot) NotifyStatsObservers(stats models.ParkingStats) {
    for _, observer := range p.observers {
        observer.UpdateStats(stats)
    }
}

func (p *ParkingLot) FindEmptySpot() (int, bool) {
    for i, occupied := range p.spaces {
        if !occupied {
            return i, true
        }
    }
    return -1, false
}

func (p *ParkingLot) GeneratePoissonDelay() time.Duration {
    return time.Duration(p.rng.ExpFloat64() * 1000) * time.Millisecond
}

func (p *ParkingLot) GetWaitGroup() *sync.WaitGroup {
    return &p.wg
}

func (p *ParkingLot) GetMutex() *sync.Mutex {
    return &p.mutex
}

func (p *ParkingLot) GetGateMutex() *sync.Mutex {
    return &p.gate
}

func (p *ParkingLot) GetGateDirectionMutex() *sync.Mutex {
    return &p.gateDirection
}

func (p *ParkingLot) GetStats() *models.ParkingStats {
    return &p.stats
}

func (p *ParkingLot) GetSpaces() *[20]bool {
    return &p.spaces
}

func (p *ParkingLot) SetIsExiting(value bool) {
    p.isExiting = value
}

func (p *ParkingLot) SetIsEntering(value bool) {
    p.isEntering = value
}

func (p *ParkingLot) IsExiting() bool {
    return p.isExiting
}

func (p *ParkingLot) IsEntering() bool {
    return p.isEntering
}

func (p *ParkingLot) GetRng() *rand.Rand {
    return p.rng
}