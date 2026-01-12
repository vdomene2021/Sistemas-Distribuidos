/*
ESTE ES EL ÚNICO ARCHIVO QUE SE PUEDE MODIFICAR

RECOMENDACIÓN: Solo modicar a partir de la parte
				donde se encuentran la explicación
				de las otras variables.

*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
	msg    string
)

// ============================================================================
// CONFIGURACIÓN DEL TALLER - Variables configurables para los tests
// ============================================================================

var (
	NumPlazas    = 6
	NumMecanicos = 3
	CochesA      = 10
	CochesB      = 10
	CochesC      = 10
)

// ============================================================================
// TIPOS Y ESTRUCTURAS
// ============================================================================

// Categoria representa el tipo de incidencia del coche
type Categoria int

const (
	CategoriaA Categoria = iota // Mecánica - Alta prioridad - 5 seg
	CategoriaB                  // Eléctrica - Media prioridad - 3 seg
	CategoriaC                  // Carrocería - Baja prioridad - 1 seg
)

func (c Categoria) String() string {
	switch c {
	case CategoriaA:
		return "Mecánica"
	case CategoriaB:
		return "Eléctrica"
	case CategoriaC:
		return "Carrocería"
	}
	return "Desconocida"
}

func (c Categoria) Prioridad() int {
	switch c {
	case CategoriaA:
		return 3 // Alta
	case CategoriaB:
		return 2 // Media
	case CategoriaC:
		return 1 // Baja
	}
	return 0
}

func (c Categoria) TiempoFase() time.Duration {
	switch c {
	case CategoriaA:
		return 5 * time.Second
	case CategoriaB:
		return 3 * time.Second
	case CategoriaC:
		return 1 * time.Second
	}
	return 1 * time.Second
}

// EstadoTaller representa el estado operativo según mensajes del servidor
type EstadoTaller int

const (
	TallerInactivo      EstadoTaller = iota // 0
	SoloCategoriaA                          // 1
	SoloCategoriaB                          // 2
	SoloCategoriaC                          // 3
	PrioridadCategoriaA                     // 4
	PrioridadCategoriaB                     // 5
	PrioridadCategoriaC                     // 6
	NoDefinido1                             // 7
	NoDefinido2                             // 8
	TallerCerrado                           // 9
)

// Fase representa cada etapa del proceso de reparación
type Fase int

const (
	FaseEsperaPlaza Fase = iota
	FaseReparacion
	FaseLimpieza
	FaseEntrega
	FaseCompletado
)

func (f Fase) String() string {
	switch f {
	case FaseEsperaPlaza:
		return "EsperaPlaza"
	case FaseReparacion:
		return "Reparacion"
	case FaseLimpieza:
		return "Limpieza"
	case FaseEntrega:
		return "Entrega"
	case FaseCompletado:
		return "Completado"
	}
	return "Desconocida"
}

// Coche representa un vehículo en el taller
type Coche struct {
	ID        int
	Categoria Categoria
}

// SolicitudFase representa una petición de un coche para entrar en una fase
type SolicitudFase struct {
	Coche     *Coche
	Respuesta chan bool
}

// ============================================================================
// GESTOR DEL TALLER
// ============================================================================

type Taller struct {
	// Canales para cada fase
	solPlazas    chan SolicitudFase
	solMecanicos chan SolicitudFase
	solLimpieza  chan SolicitudFase
	solEntrega   chan SolicitudFase

	// Control de recursos
	plazasLibres    chan struct{}
	mecanicosLibres chan struct{}

	// Estado del taller
	estadoActual EstadoTaller
	estadoMutex  sync.RWMutex

	// Control de tiempo
	tiempoInicio time.Time

	// WaitGroup para esperar a que terminen todos los coches
	wg sync.WaitGroup
}

func NuevoTaller() *Taller {
	t := &Taller{
		solPlazas:       make(chan SolicitudFase, 100),
		solMecanicos:    make(chan SolicitudFase, 100),
		solLimpieza:     make(chan SolicitudFase, 100),
		solEntrega:      make(chan SolicitudFase, 100),
		plazasLibres:    make(chan struct{}, NumPlazas),
		mecanicosLibres: make(chan struct{}, NumMecanicos),
		estadoActual:    TallerInactivo, // Esperar a que mutua active el taller
		tiempoInicio:    time.Now(),
	}

	// Inicializar recursos disponibles
	for i := 0; i < NumPlazas; i++ {
		t.plazasLibres <- struct{}{}
	}
	for i := 0; i < NumMecanicos; i++ {
		t.mecanicosLibres <- struct{}{}
	}

	return t
}

func (t *Taller) CambiarEstado(nuevoEstado EstadoTaller) {
	t.estadoMutex.Lock()
	defer t.estadoMutex.Unlock()

	// Si es NoDefinido (7 u 8), mantener estado anterior
	if nuevoEstado == NoDefinido1 || nuevoEstado == NoDefinido2 {
		fmt.Printf("[TALLER] Estado %d recibido - Manteniendo estado anterior: %d\n", nuevoEstado, t.estadoActual)
		return
	}

	t.estadoActual = nuevoEstado

	// Mensaje descriptivo del estado
	var descripcion string
	switch nuevoEstado {
	case TallerInactivo:
		descripcion = "INACTIVO - No hay atención"
	case SoloCategoriaA:
		descripcion = "SOLO CATEGORÍA A (Mecánica)"
	case SoloCategoriaB:
		descripcion = "SOLO CATEGORÍA B (Eléctrica)"
	case SoloCategoriaC:
		descripcion = "SOLO CATEGORÍA C (Carrocería)"
	case PrioridadCategoriaA:
		descripcion = "PRIORIDAD CATEGORÍA A (Mecánica prioritaria)"
	case PrioridadCategoriaB:
		descripcion = "PRIORIDAD CATEGORÍA B (Eléctrica prioritaria)"
	case PrioridadCategoriaC:
		descripcion = "PRIORIDAD CATEGORÍA C (Carrocería prioritaria)"
	case TallerCerrado:
		descripcion = "CERRADO - No hay atención"
	default:
		descripcion = "DESCONOCIDO"
	}

	fmt.Printf("[TALLER] *** NUEVO ESTADO: %d - %s ***\n", nuevoEstado, descripcion)
}

func (t *Taller) ObtenerEstado() EstadoTaller {
	t.estadoMutex.RLock()
	defer t.estadoMutex.RUnlock()
	return t.estadoActual
}

func (t *Taller) PuedeAcceder(categoria Categoria) bool {
	estado := t.ObtenerEstado()

	switch estado {
	case TallerInactivo, TallerCerrado:
		return false
	case SoloCategoriaA:
		return categoria == CategoriaA
	case SoloCategoriaB:
		return categoria == CategoriaB
	case SoloCategoriaC:
		return categoria == CategoriaC
	case PrioridadCategoriaA, PrioridadCategoriaB, PrioridadCategoriaC:
		return true // Todos pueden acceder, pero con prioridad
	default:
		return true
	}
}

func (t *Taller) ObtenerPrioridadEfectiva(categoria Categoria) int {
	estado := t.ObtenerEstado()
	prioridadBase := categoria.Prioridad()

	switch estado {
	case PrioridadCategoriaA:
		if categoria == CategoriaA {
			return prioridadBase + 10
		}
	case PrioridadCategoriaB:
		if categoria == CategoriaB {
			return prioridadBase + 10
		}
	case PrioridadCategoriaC:
		if categoria == CategoriaC {
			return prioridadBase + 10
		}
	}

	return prioridadBase
}

func (t *Taller) LogEvento(coche *Coche, fase Fase, estado string) {
	tiempoTranscurrido := time.Since(t.tiempoInicio).Seconds()
	fmt.Printf("Tiempo %06.3f Coche %d Incidencia %s Fase %s Estado %s\n",
		tiempoTranscurrido, coche.ID, coche.Categoria, fase, estado)
}

// ============================================================================
// GESTORES DE FASES (Goroutines que gestionan cada fase)
// ============================================================================

func (t *Taller) GestorPlazas() {
	var cola []*SolicitudFase
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case solicitud := <-t.solPlazas:
			cola = append(cola, &solicitud)
		case <-ticker.C:
			t.procesarColaPlazas(&cola)
		}
	}
}

func (t *Taller) procesarColaPlazas(cola *[]*SolicitudFase) {
	if len(*cola) == 0 {
		return
	}

	// Ordenar por prioridad
	for i := 0; i < len(*cola)-1; i++ {
		for j := i + 1; j < len(*cola); j++ {
			priI := t.ObtenerPrioridadEfectiva((*cola)[i].Coche.Categoria)
			priJ := t.ObtenerPrioridadEfectiva((*cola)[j].Coche.Categoria)
			if priJ > priI {
				(*cola)[i], (*cola)[j] = (*cola)[j], (*cola)[i]
			}
		}
	}

	// Procesar solicitudes en orden de prioridad
	i := 0
	for i < len(*cola) {
		solicitud := (*cola)[i]

		// Verificar si puede acceder según el estado del taller
		if !t.PuedeAcceder(solicitud.Coche.Categoria) {
			i++
			continue
		}

		// Intentar asignar plaza (no bloqueante)
		select {
		case <-t.plazasLibres:
			solicitud.Respuesta <- true
			*cola = append((*cola)[:i], (*cola)[i+1:]...)
			// No incrementar i porque eliminamos el elemento
		default:
			i++
		}
	}
}

func (t *Taller) GestorMecanicos() {
	var cola []*SolicitudFase
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case solicitud := <-t.solMecanicos:
			cola = append(cola, &solicitud)
		case <-ticker.C:
			t.procesarColaMecanicos(&cola)
		}
	}
}

func (t *Taller) procesarColaMecanicos(cola *[]*SolicitudFase) {
	if len(*cola) == 0 {
		return
	}

	// Ordenar por prioridad
	for i := 0; i < len(*cola)-1; i++ {
		for j := i + 1; j < len(*cola); j++ {
			priI := t.ObtenerPrioridadEfectiva((*cola)[i].Coche.Categoria)
			priJ := t.ObtenerPrioridadEfectiva((*cola)[j].Coche.Categoria)
			if priJ > priI {
				(*cola)[i], (*cola)[j] = (*cola)[j], (*cola)[i]
			}
		}
	}

	i := 0
	for i < len(*cola) {
		solicitud := (*cola)[i]

		select {
		case <-t.mecanicosLibres:
			solicitud.Respuesta <- true
			*cola = append((*cola)[:i], (*cola)[i+1:]...)
		default:
			i++
		}
	}
}

func (t *Taller) GestorLimpieza() {
	var cola []*SolicitudFase
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case solicitud := <-t.solLimpieza:
			cola = append(cola, &solicitud)
		case <-ticker.C:
			t.procesarColaLimpieza(&cola)
		}
	}
}

func (t *Taller) procesarColaLimpieza(cola *[]*SolicitudFase) {
	if len(*cola) == 0 {
		return
	}

	// Ordenar por prioridad
	for i := 0; i < len(*cola)-1; i++ {
		for j := i + 1; j < len(*cola); j++ {
			priI := t.ObtenerPrioridadEfectiva((*cola)[i].Coche.Categoria)
			priJ := t.ObtenerPrioridadEfectiva((*cola)[j].Coche.Categoria)
			if priJ > priI {
				(*cola)[i], (*cola)[j] = (*cola)[j], (*cola)[i]
			}
		}
	}

	// Procesar en orden (limpieza es secuencial)
	if len(*cola) > 0 {
		solicitud := (*cola)[0]
		solicitud.Respuesta <- true
		*cola = (*cola)[1:]
	}
}

func (t *Taller) GestorEntrega() {
	var cola []*SolicitudFase
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case solicitud := <-t.solEntrega:
			cola = append(cola, &solicitud)
		case <-ticker.C:
			t.procesarColaEntrega(&cola)
		}
	}
}

func (t *Taller) procesarColaEntrega(cola *[]*SolicitudFase) {
	if len(*cola) == 0 {
		return
	}

	// Ordenar por prioridad
	for i := 0; i < len(*cola)-1; i++ {
		for j := i + 1; j < len(*cola); j++ {
			priI := t.ObtenerPrioridadEfectiva((*cola)[i].Coche.Categoria)
			priJ := t.ObtenerPrioridadEfectiva((*cola)[j].Coche.Categoria)
			if priJ > priI {
				(*cola)[i], (*cola)[j] = (*cola)[j], (*cola)[i]
			}
		}
	}

	// Procesar en orden
	if len(*cola) > 0 {
		solicitud := (*cola)[0]
		solicitud.Respuesta <- true
		*cola = (*cola)[1:]
	}
}

// ============================================================================
// PROCESO DEL COCHE
// ============================================================================

func (t *Taller) ProcesarCoche(coche *Coche) {
	defer t.wg.Done()

	// Esperar un tiempo aleatorio antes de llegar
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

	// FASE 1: Espera de Plaza
	t.LogEvento(coche, FaseEsperaPlaza, "Esperando")

	respuestaPlaza := make(chan bool)
	t.solPlazas <- SolicitudFase{Coche: coche, Respuesta: respuestaPlaza}
	<-respuestaPlaza // Esperar hasta que el gestor nos asigne una plaza

	t.LogEvento(coche, FaseEsperaPlaza, "Entrando")
	// Preparar documentación del coche (tiempo asignado por prioridad)
	time.Sleep(coche.Categoria.TiempoFase())
	t.LogEvento(coche, FaseEsperaPlaza, "Saliendo")

	// Liberar plaza - ya no necesitamos la plaza de espera
	t.plazasLibres <- struct{}{}

	// FASE 2: Reparación
	t.LogEvento(coche, FaseReparacion, "Esperando")

	respuesta := make(chan bool)
	t.solMecanicos <- SolicitudFase{Coche: coche, Respuesta: respuesta}
	<-respuesta

	t.LogEvento(coche, FaseReparacion, "Entrando")
	time.Sleep(coche.Categoria.TiempoFase())
	t.LogEvento(coche, FaseReparacion, "Saliendo")

	// Liberar mecánico
	t.mecanicosLibres <- struct{}{}

	// FASE 3: Limpieza
	t.LogEvento(coche, FaseLimpieza, "Esperando")

	respuesta2 := make(chan bool)
	t.solLimpieza <- SolicitudFase{Coche: coche, Respuesta: respuesta2}
	<-respuesta2

	t.LogEvento(coche, FaseLimpieza, "Entrando")
	time.Sleep(coche.Categoria.TiempoFase())
	t.LogEvento(coche, FaseLimpieza, "Saliendo")

	// FASE 4: Entrega
	t.LogEvento(coche, FaseEntrega, "Esperando")

	respuesta3 := make(chan bool)
	t.solEntrega <- SolicitudFase{Coche: coche, Respuesta: respuesta3}
	<-respuesta3

	t.LogEvento(coche, FaseEntrega, "Entrando")
	time.Sleep(coche.Categoria.TiempoFase())
	t.LogEvento(coche, FaseEntrega, "Saliendo")

	t.LogEvento(coche, FaseCompletado, "Finalizado")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Crear taller
	taller := NuevoTaller()

	// Iniciar gestores de fases
	go taller.GestorPlazas()
	go taller.GestorMecanicos()
	go taller.GestorLimpieza()
	go taller.GestorEntrega()

	fmt.Println("[TALLER] Conectando al servidor...")
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()

	// Canal para señalizar cuando se recibe el primer estado operativo
	tallerListo := make(chan bool, 1)
	primeraVez := true

	// Goroutine para procesar mensajes del servidor
	go func() {
		buf := make([]byte, 512)
		for {
			n, err := conn.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			if n > 0 {
				msg = string(buf[:n])
				msg = strings.TrimSpace(msg)
				if estado, err := strconv.Atoi(msg); err == nil {
					if estado >= 0 && estado <= 9 {
						taller.CambiarEstado(EstadoTaller(estado))

						// Señalizar que el taller está listo cuando reciba un estado operativo
						if primeraVez && estado >= 1 && estado <= 6 {
							primeraVez = false
							tallerListo <- true
						}
					}
				}
			}
		}
	}()

	// Esperar a que llegue el primer estado operativo del servidor
	<-tallerListo

	// Reiniciar el timer ahora que el sistema está activo
	taller.tiempoInicio = time.Now()

	// Crear y lanzar coches
	cocheID := 1

	// Coches Categoría A
	for i := 0; i < CochesA; i++ {
		coche := &Coche{
			ID:        cocheID,
			Categoria: CategoriaA,
		}
		cocheID++
		taller.wg.Add(1)
		go taller.ProcesarCoche(coche)
	}

	// Coches Categoría B
	for i := 0; i < CochesB; i++ {
		coche := &Coche{
			ID:        cocheID,
			Categoria: CategoriaB,
		}
		cocheID++
		taller.wg.Add(1)
		go taller.ProcesarCoche(coche)
	}

	// Coches Categoría C
	for i := 0; i < CochesC; i++ {
		coche := &Coche{
			ID:        cocheID,
			Categoria: CategoriaC,
		}
		cocheID++
		taller.wg.Add(1)
		go taller.ProcesarCoche(coche)
	}
}
