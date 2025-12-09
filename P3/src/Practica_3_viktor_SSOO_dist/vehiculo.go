package main

import (
	"fmt"
	"time"
)

// Prioridad del vehículo
type Prioridad int

const (
	PrioridadBaja  Prioridad = 1
	PrioridadMedia Prioridad = 2
	PrioridadAlta  Prioridad = 3
)

// TipoIncidencia representa el tipo de reparación necesaria
type TipoIncidencia string

const (
	Mecanica   TipoIncidencia = "Mecánica"
	Electrica  TipoIncidencia = "Eléctrica"
	Carroceria TipoIncidencia = "Carrocería"
)

// Vehiculo representa un coche en el taller
type Vehiculo struct {
	ID            int
	Incidencia    TipoIncidencia
	Prioridad     Prioridad
	TiempoFase    time.Duration
	TiempoLlegada time.Time
}

// NewVehiculo crea un nuevo vehículo según su categoría
func NewVehiculo(id int, incidencia TipoIncidencia) *Vehiculo {
	v := &Vehiculo{
		ID:            id,
		Incidencia:    incidencia,
		TiempoLlegada: time.Now(),
	}

	// Asignar prioridad y tiempo según categoría
	switch incidencia {
	case Mecanica: // Categoría A
		v.Prioridad = PrioridadAlta
		v.TiempoFase = 5 * time.Second
	case Electrica: // Categoría B
		v.Prioridad = PrioridadMedia
		v.TiempoFase = 3 * time.Second
	case Carroceria: // Categoría C
		v.Prioridad = PrioridadBaja
		v.TiempoFase = 1 * time.Second
	}

	return v
}

// LogEstado imprime el estado del vehículo en una fase
func (v *Vehiculo) LogEstado(fase string, estado string, tiempoEjecucion time.Duration) {
	fmt.Printf("Tiempo %v Coche %d Incidencia %s Fase %s Estado %s\n",
		tiempoEjecucion.Round(time.Millisecond),
		v.ID,
		v.Incidencia,
		fase,
		estado)
}
