package entities

import (
    "math/rand"
    "sync"
    "time"

    "Estacionamiento/internal/domain/interfaces"
    "Estacionamiento/internal/domain/models"
)

type Direction int

const (
    None Direction = iota
    Entering
    Exiting
)

type ParkingLot struct {
    observers      []interfaces.Observer
    spaces         [20]bool
    wg             sync.WaitGroup
    stats          models.ParkingStats
    rng            *rand.Rand

    ParkingSem     models.Semaphore    
    GateSem        models.Semaphore   

    DirectionChan  chan Direction      
    SpacesChan     chan int           
    StatsChan      chan models.ParkingStats  
}

func NewParkingLot() *ParkingLot {
    p := &ParkingLot{
        rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
        ParkingSem:    models.NewSemaphore(20),
        GateSem:       models.NewSemaphore(1),
        DirectionChan: make(chan Direction, 1),
        SpacesChan:    make(chan int, 20),
        StatsChan:     make(chan models.ParkingStats, 1),
    }

    for i := 0; i < 20; i++ {
        p.SpacesChan <- i
    }

    p.DirectionChan <- None

    return p
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

func (p *ParkingLot) GetSpace() (int, bool) {
    select {
    case spot := <-p.SpacesChan:
        p.spaces[spot] = true
        return spot, true
    default:
        return -1, false
    }
}

func (p *ParkingLot) ReleaseSpace(spot int) {
    p.spaces[spot] = false
    p.SpacesChan <- spot
}

func (p *ParkingLot) GetDirection() Direction {
    return <-p.DirectionChan
}

func (p *ParkingLot) SetDirection(d Direction) {
    select {
    case <-p.DirectionChan: 
    default:
    }
    p.DirectionChan <- d
}

func (p *ParkingLot) UpdateStats(update func(*models.ParkingStats)) {
    update(&p.stats)
}

func (p *ParkingLot) GeneratePoissonDelay() time.Duration {
    return time.Duration(p.rng.ExpFloat64() * 1000) * time.Millisecond
}

func (p *ParkingLot) GetWaitGroup() *sync.WaitGroup {
    return &p.wg
}

func (p *ParkingLot) GetStats() *models.ParkingStats {
    return &p.stats
}

func (p *ParkingLot) GetSpaces() *[20]bool {
    return &p.spaces
}

func (p *ParkingLot) GetRng() *rand.Rand {
    return p.rng
}