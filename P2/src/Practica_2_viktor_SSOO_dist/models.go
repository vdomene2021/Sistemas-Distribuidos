package main

import (
	"time"
)

type Cliente struct {
	ID         int
	Nombre     string
	Telefono   string
	Email      string
	VehiculoID int
}

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

type TipoIncidencia string

const (
	Mecanica   TipoIncidencia = "mecanica"
	Electrica  TipoIncidencia = "electrica"
	Carroceria TipoIncidencia = "carroceria"
)

type Prioridad string

const (
	Baja  Prioridad = "baja"
	Media Prioridad = "media"
	Alta  Prioridad = "alta"
)

type EstadoIncidencia string

const (
	Abierta   EstadoIncidencia = "abierta"
	EnProceso EstadoIncidencia = "en proceso"
	Cerrada   EstadoIncidencia = "cerrada"
)

type Incidencia struct {
	ID           int
	MecanicosIDs []int
	Tipo         TipoIncidencia
	Prioridad    Prioridad
	Descripcion  string
	Estado       EstadoIncidencia
}

type Especialidad string

const (
	EspecialidadMecanica   Especialidad = "mecanica"
	EspecialidadElectrica  Especialidad = "electrica"
	EspecialidadCarroceria Especialidad = "carroceria"
)

type Mecanico struct {
	ID           int
	Nombre       string
	Especialidad Especialidad
	Experiencia  int
	Activo       bool
	Ocupado      bool // Indica si está atendiendo un vehículo
}

type TrabajoPendiente struct {
	Vehiculo     *Vehiculo
	Incidencia   *Incidencia
	TiempoInicio time.Time
}
