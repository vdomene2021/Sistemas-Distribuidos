package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// TIPOS BÁSICOS Y CONSTANTES
// ============================================================================

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

// EstadoIncidencia representa el estado de una incidencia
type EstadoIncidencia string

const (
	Abierta   EstadoIncidencia = "abierta"
	EnProceso EstadoIncidencia = "en proceso"
	Cerrada   EstadoIncidencia = "cerrada"
)

// Especialidad del mecánico
type Especialidad string

const (
	EspMecanica   Especialidad = "mecanica"
	EspElectrica  Especialidad = "electrica"
	EspCarroceria Especialidad = "carroceria"
)

// ============================================================================
// ESTRUCTURAS PARA SIMULACIÓN (TESTS ACADÉMICOS)
// ============================================================================

// Vehiculo representa un coche en el taller (para simulación)
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

// TallerSimulacion representa el taller con sus recursos para la simulación
type TallerSimulacion struct {
	NumPlazas    int
	NumMecanicos int
	PlazasSem    chan struct{}
	MecanicosSem chan struct{}
}

// NewTallerSimulacion crea un nuevo taller para simulación
func NewTallerSimulacion(numPlazas, numMecanicos int) *TallerSimulacion {
	return &TallerSimulacion{
		NumPlazas:    numPlazas,
		NumMecanicos: numMecanicos,
		PlazasSem:    make(chan struct{}, numPlazas),
		MecanicosSem: make(chan struct{}, numMecanicos),
	}
}

// ============================================================================
// ESTRUCTURAS PARA SISTEMA INTERACTIVO
// ============================================================================

// Cliente representa a un cliente del taller
type Cliente struct {
	ID       int
	Nombre   string
	Telefono string
	Email    string
}

// VehiculoCompleto representa un vehículo con toda su información
type VehiculoCompleto struct {
	ID              int
	Matricula       string
	Marca           string
	Modelo          string
	ClienteID       int
	TiempoAcumulado float64
}

// IncidenciaCompleta representa una reparación pendiente con toda su información
type IncidenciaCompleta struct {
	ID          int
	VehiculoID  int
	Tipo        TipoIncidencia
	Prioridad   Prioridad
	Estado      EstadoIncidencia
	Descripcion string
}

// Mecanico representa a un mecánico del taller
type Mecanico struct {
	ID           int
	Nombre       string
	Especialidad Especialidad
	Experiencia  int
	Activo       bool
	ColaPersonal []*TrabajoMecanico
	Canal        chan *TrabajoMecanico
	mutex        sync.Mutex
}

// TrabajoMecanico representa un trabajo asignado a un mecánico
type TrabajoMecanico struct {
	Vehiculo   VehiculoCompleto
	Incidencia IncidenciaCompleta
}

// Taller representa el sistema de gestión del taller (interactivo)
type Taller struct {
	ColaTrabajo       []*TrabajoMecanico
	MecanicoManager   *MecanicoManager
	VehiculoManager   *VehiculoManager
	IncidenciaManager *IncidenciaManager
	mutex             sync.Mutex
}

// ============================================================================
// MANAGERS (GESTORES DE DATOS)
// ============================================================================

// ClienteManager gestiona los clientes
type ClienteManager struct {
	clientes map[int]*Cliente
	nextID   int
	mutex    sync.RWMutex
}

// VehiculoManager gestiona los vehículos
type VehiculoManager struct {
	vehiculos map[int]*VehiculoCompleto
	nextID    int
	mutex     sync.RWMutex
}

// IncidenciaManager gestiona las incidencias
type IncidenciaManager struct {
	incidencias map[int]*IncidenciaCompleta
	nextID      int
	mutex       sync.RWMutex
}

// MecanicoManager gestiona los mecánicos
type MecanicoManager struct {
	mecanicos map[int]*Mecanico
	nextID    int
	mutex     sync.RWMutex
}

// ============================================================================
// CONSTRUCTORES
// ============================================================================

// NewClienteManager crea un nuevo gestor de clientes
func NewClienteManager() *ClienteManager {
	return &ClienteManager{
		clientes: make(map[int]*Cliente),
		nextID:   1,
	}
}

// NewVehiculoManager crea un nuevo gestor de vehículos
func NewVehiculoManager() *VehiculoManager {
	return &VehiculoManager{
		vehiculos: make(map[int]*VehiculoCompleto),
		nextID:    1,
	}
}

// NewIncidenciaManager crea un nuevo gestor de incidencias
func NewIncidenciaManager() *IncidenciaManager {
	return &IncidenciaManager{
		incidencias: make(map[int]*IncidenciaCompleta),
		nextID:      1,
	}
}

// NewMecanicoManager crea un nuevo gestor de mecánicos
func NewMecanicoManager() *MecanicoManager {
	return &MecanicoManager{
		mecanicos: make(map[int]*Mecanico),
		nextID:    1,
	}
}

// NewTaller crea un nuevo taller para el sistema interactivo
func NewTaller(mm *MecanicoManager, vm *VehiculoManager, im *IncidenciaManager) *Taller {
	return &Taller{
		ColaTrabajo:       make([]*TrabajoMecanico, 0),
		MecanicoManager:   mm,
		VehiculoManager:   vm,
		IncidenciaManager: im,
	}
}
