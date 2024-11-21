package models

// Semaphore implementa un semáforo usando canales
type Semaphore chan struct{}

// NewSemaphore crea un nuevo semáforo con la capacidad especificada
func NewSemaphore(capacity int) Semaphore {
    return make(Semaphore, capacity)
}

// Acquire adquiere el semáforo
func (s Semaphore) Acquire() {
    s <- struct{}{}
}

// Release libera el semáforo
func (s Semaphore) Release() {
    <-s
}