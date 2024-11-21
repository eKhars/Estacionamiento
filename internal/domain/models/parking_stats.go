package models

type ParkingStats struct {
    TotalGenerated  int32
    TotalProcessed  int32
    WaitingToEnter  int32
    WaitingToExit   int32
    CurrentlyParked int32
}