package interfaces

import (
    "Estacionamiento/internal/domain/models"
)

type Subject interface {
    RegisterObserver(o Observer)
    RemoveObserver(o Observer)
    NotifyObservers(parkingSpot int, isOccupied bool)
    NotifyGateObservers(isOccupied bool)
    NotifyStatsObservers(stats models.ParkingStats)
}