package main

import (
	"time"
)

// Cliente representa un cliente del taller
type Cliente struct {
	ID         int
	Nombre     string
	Telefono   string
	Email      string
	VehiculoID int
}

// Vehiculo representa un vehículo en el taller
type Vehiculo struct {
	ID              int
	Matricula       string
	Marca           string
	Modelo          string
	FechaEntrada    time.Time
	FechaSalida     time.Time
	IncidenciaID    int
	ClienteID       int
	TiempoAcumulado float64 // Tiempo total de atención en segundos
}

// TipoIncidencia representa los tipos de incidencia posibles
type TipoIncidencia string

const (
	Mecanica   TipoIncidencia = "mecanica"
	Electrica  TipoIncidencia = "electrica"
	Carroceria TipoIncidencia = "carroceria"
)

// Prioridad representa el nivel de prioridad
type Prioridad string

const (
	Baja  Prioridad = "baja"
	Media Prioridad = "media"
	Alta  Prioridad = "alta"
)

// EstadoIncidencia representa el estado de una incidencia
type EstadoIncidencia string

const (
	Abierta   EstadoIncidencia = "abierta"
	EnProceso EstadoIncidencia = "en proceso"
	Cerrada   EstadoIncidencia = "cerrada"
)

// Incidencia representa una reparación a realizar
type Incidencia struct {
	ID           int
	MecanicosIDs []int
	Tipo         TipoIncidencia
	Prioridad    Prioridad
	Descripcion  string
	Estado       EstadoIncidencia
}

// Especialidad representa la especialidad de un mecánico
type Especialidad string

const (
	EspecialidadMecanica   Especialidad = "mecanica"
	EspecialidadElectrica  Especialidad = "electrica"
	EspecialidadCarroceria Especialidad = "carroceria"
)

// Mecanico representa un mecánico del taller
type Mecanico struct {
	ID           int
	Nombre       string
	Especialidad Especialidad
	Experiencia  int
	Activo       bool
	Ocupado      bool // Indica si está atendiendo un vehículo
}

// TrabajoPendiente representa un trabajo en la cola de espera
type TrabajoPendiente struct {
	Vehiculo     *Vehiculo
	Incidencia   *Incidencia
	TiempoInicio time.Time
}
