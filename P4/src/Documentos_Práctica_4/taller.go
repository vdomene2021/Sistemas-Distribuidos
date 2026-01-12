/*
ESTE ES EL ÚNICO ARCHIVO QUE SE PUEDE MODIFICAR
*/

package main

import (
	"bytes"
	"fmt"
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

// Variables globales para configuración (se sobreescribirán en los tests)
var (
	NumPlazas    = 6
	NumMecanicos = 3
	CochesA      = 10
	CochesB      = 10
	CochesC      = 10
)

// ============================================================================
// ESTRUCTURAS Y TIPOS (Sin cambios)
// ============================================================================

type Categoria int

const (
	CategoriaA Categoria = iota
	CategoriaB
	CategoriaC
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
		return 3
	case CategoriaB:
		return 2
	case CategoriaC:
		return 1
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

type EstadoTaller int

const (
	TallerInactivo EstadoTaller = iota
	SoloCategoriaA
	SoloCategoriaB
	SoloCategoriaC
	PrioridadCategoriaA
	PrioridadCategoriaB
	PrioridadCategoriaC
	NoDefinido1
	NoDefinido2
	TallerCerrado
)

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

type Coche struct {
	ID        int
	Categoria Categoria
}

type SolicitudFase struct {
	Coche     *Coche
	Respuesta chan bool
}

// ============================================================================
// LÓGICA DEL TALLER
// ============================================================================

type Taller struct {
	solPlazas       chan SolicitudFase
	solMecanicos    chan SolicitudFase
	solLimpieza     chan SolicitudFase
	solEntrega      chan SolicitudFase
	plazasLibres    chan struct{}
	mecanicosLibres chan struct{}
	estadoActual    EstadoTaller
	estadoMutex     sync.RWMutex
	tiempoInicio    time.Time
	wg              sync.WaitGroup
}

// NuevoTaller ahora acepta parámetros para configurar recursos dinámicamente
func NuevoTaller(plazas, mecanicos int) *Taller {
	t := &Taller{
		solPlazas:       make(chan SolicitudFase, 100),
		solMecanicos:    make(chan SolicitudFase, 100),
		solLimpieza:     make(chan SolicitudFase, 100),
		solEntrega:      make(chan SolicitudFase, 100),
		plazasLibres:    make(chan struct{}, plazas),
		mecanicosLibres: make(chan struct{}, mecanicos),
		estadoActual:    TallerInactivo,
		tiempoInicio:    time.Now(),
	}

	for i := 0; i < plazas; i++ {
		t.plazasLibres <- struct{}{}
	}
	for i := 0; i < mecanicos; i++ {
		t.mecanicosLibres <- struct{}{}
	}

	return t
}

func (t *Taller) CambiarEstado(nuevoEstado EstadoTaller) {
	t.estadoMutex.Lock()
	defer t.estadoMutex.Unlock()
	if nuevoEstado == NoDefinido1 || nuevoEstado == NoDefinido2 {
		return
	}
	t.estadoActual = nuevoEstado
	fmt.Printf("[TALLER] *** CAMBIO DE ESTADO A: %d ***\n", nuevoEstado)
}

func (t *Taller) ObtenerEstado() EstadoTaller {
	t.estadoMutex.RLock()
	defer t.estadoMutex.RUnlock()
	return t.estadoActual
}

// PuedeAcceder verifica si una categoría tiene permiso de entrada
func (t *Taller) PuedeAcceder(categoria Categoria) bool {
	estado := t.ObtenerEstado()
	switch estado {
	case TallerCerrado: // Estado 9: Cierre total
		return false
	// CORRECCIÓN: Quitamos TallerInactivo (0) de aquí para evitar el bloqueo (Deadlock).
	// Si mutua envía un 0 antes de que acabemos, dejamos que los coches sigan pasando.
	case SoloCategoriaA:
		return categoria == CategoriaA
	case SoloCategoriaB:
		return categoria == CategoriaB
	case SoloCategoriaC:
		return categoria == CategoriaC
	default:
		return true // En estados de prioridad o Inactivo (0), dejamos pasar para drenar la cola
	}
}

func (t *Taller) ObtenerPrioridadEfectiva(categoria Categoria) int {
	estado := t.ObtenerEstado()
	prio := categoria.Prioridad()
	switch estado {
	case PrioridadCategoriaA:
		if categoria == CategoriaA {
			prio += 10
		}
	case PrioridadCategoriaB:
		if categoria == CategoriaB {
			prio += 10
		}
	case PrioridadCategoriaC:
		if categoria == CategoriaC {
			prio += 10
		}
	}
	return prio
}

func (t *Taller) LogEvento(coche *Coche, fase Fase, estado string) {
	tiempo := time.Since(t.tiempoInicio).Seconds()
	fmt.Printf("Tiempo %06.3f Coche %d Incidencia %s Fase %s Estado %s\n",
		tiempo, coche.ID, coche.Categoria, fase, estado)
}

// ============================================================================
// GESTORES Y PROCESOS
// ============================================================================

func (t *Taller) GestorGenerico(canalEntrada chan SolicitudFase, recursoLimitado chan struct{}, nombre string) {
	var cola []*SolicitudFase
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case sol := <-canalEntrada:
			cola = append(cola, &sol)
		case <-ticker.C:
			if len(cola) == 0 {
				continue
			}
			// Ordenar
			for i := 0; i < len(cola)-1; i++ {
				for j := i + 1; j < len(cola); j++ {
					p1 := t.ObtenerPrioridadEfectiva(cola[i].Coche.Categoria)
					p2 := t.ObtenerPrioridadEfectiva(cola[j].Coche.Categoria)
					if p2 > p1 {
						cola[i], cola[j] = cola[j], cola[i]
					}
				}
			}
			// Procesar
			i := 0
			for i < len(cola) {
				sol := cola[i]
				if nombre == "Plazas" && !t.PuedeAcceder(sol.Coche.Categoria) {
					i++
					continue
				}
				asignado := false
				if recursoLimitado != nil {
					select {
					case <-recursoLimitado:
						asignado = true
					default:
						asignado = false
					}
				} else {
					asignado = true
				}

				if asignado {
					sol.Respuesta <- true
					cola = append(cola[:i], cola[i+1:]...)
					if nombre == "Limpieza" || nombre == "Entrega" {
						break
					}
				} else {
					break
				}
			}
		}
	}
}

func (t *Taller) ProcesarCoche(coche *Coche) {
	defer t.wg.Done()
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)

	// FASE 1
	t.LogEvento(coche, FaseEsperaPlaza, "Esperando")
	resp := make(chan bool)
	t.solPlazas <- SolicitudFase{coche, resp}
	<-resp
	t.LogEvento(coche, FaseEsperaPlaza, "Entrando")
	time.Sleep(coche.Categoria.TiempoFase())
	t.LogEvento(coche, FaseEsperaPlaza, "Saliendo")

	// FASE 2
	t.LogEvento(coche, FaseReparacion, "Esperando")
	t.solMecanicos <- SolicitudFase{coche, resp}
	<-resp
	t.plazasLibres <- struct{}{} // Libera Plaza
	t.LogEvento(coche, FaseReparacion, "Entrando")
	time.Sleep(coche.Categoria.TiempoFase())
	t.LogEvento(coche, FaseReparacion, "Saliendo")
	t.mecanicosLibres <- struct{}{} // Libera Mecánico

	// FASE 3
	t.LogEvento(coche, FaseLimpieza, "Esperando")
	t.solLimpieza <- SolicitudFase{coche, resp}
	<-resp
	t.LogEvento(coche, FaseLimpieza, "Entrando")
	time.Sleep(coche.Categoria.TiempoFase())
	t.LogEvento(coche, FaseLimpieza, "Saliendo")

	// FASE 4
	t.LogEvento(coche, FaseEntrega, "Esperando")
	t.solEntrega <- SolicitudFase{coche, resp}
	<-resp
	t.LogEvento(coche, FaseEntrega, "Entrando")
	time.Sleep(coche.Categoria.TiempoFase())
	t.LogEvento(coche, FaseEntrega, "Saliendo")

	t.LogEvento(coche, FaseCompletado, "Finalizado")
}

// EjecutarSimulacion encapsula la lógica principal para ser llamada desde tests
func EjecutarSimulacion(nPlazas, nMecs, nA, nB, nC int) {
	rand.Seed(time.Now().UnixNano())

	// Inicializar taller con parámetros dinámicos
	taller := NuevoTaller(nPlazas, nMecs)

	go taller.GestorGenerico(taller.solPlazas, taller.plazasLibres, "Plazas")
	go taller.GestorGenerico(taller.solMecanicos, taller.mecanicosLibres, "Mecanicos")
	go taller.GestorGenerico(taller.solLimpieza, nil, "Limpieza")
	go taller.GestorGenerico(taller.solEntrega, nil, "Entrega")

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		logger.Fatal(err)
	}

	// Importante: Cerrar conexión al terminar la función para matar al listener
	defer conn.Close()

	inicioOperaciones := make(chan bool)
	onceStart := sync.Once{}

	go func() {
		buf := make([]byte, 512)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				// CORRECCIÓN: Si hay error (como EOF o conexión cerrada), SALIR del bucle
				// para que el programa pueda terminar limpiamente.
				break
			}
			if n > 0 {
				msg = string(buf[:n])
				// LINEA DE LA PLANTILLA
				fmt.Println("len: " + strconv.Itoa(n) + " msg: " + msg)

				msg = strings.TrimSpace(msg)
				val, err := strconv.Atoi(msg)
				if err == nil {
					taller.CambiarEstado(EstadoTaller(val))
					if val >= 1 && val <= 6 {
						onceStart.Do(func() {
							close(inicioOperaciones)
						})
					}
				}
			}
		}
	}()

	fmt.Println("[TALLER] Esperando instrucciones de la Mutua...")
	<-inicioOperaciones
	fmt.Printf("[TALLER] Iniciando procesamiento con: Plazas=%d, Mecs=%d, A=%d, B=%d, C=%d\n",
		nPlazas, nMecs, nA, nB, nC)

	taller.tiempoInicio = time.Now()

	// CORRECCIÓN: ID secuencial global
	idGlobal := 1

	lanzar := func(n int, cat Categoria) {
		for i := 0; i < n; i++ {
			coche := &Coche{ID: idGlobal, Categoria: cat}
			idGlobal++
			taller.wg.Add(1)
			go taller.ProcesarCoche(coche)
		}
	}

	lanzar(nA, CategoriaA)
	lanzar(nB, CategoriaB)
	lanzar(nC, CategoriaC)

	taller.wg.Wait()
	fmt.Println("[TALLER] Fin de la jornada.")

	// Al llegar aquí, se ejecuta el 'defer conn.Close()',
	// lo que causará un error en el 'conn.Read' del bucle superior,
	// haciendo que el 'break' salte y todo termine bien.
}

func main() {
	// Ejecución por defecto (si se lanza con go run taller.go)
	EjecutarSimulacion(NumPlazas, NumMecanicos, CochesA, CochesB, CochesC)
}
