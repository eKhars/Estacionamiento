package interfaces

import (
    "Estacionamiento/internal/domain/models"
)

type Observer interface {
    Update(parkingSpot int, isOccupied bool)
    UpdateGate(isOccupied bool)
    UpdateStats(stats models.ParkingStats)
}