package main

import (
	"container/heap"
	"math/rand"
	"sync"
	"time"
)

// ColaPrioridad implementa heap.Interface para ordenar vehículos por prioridad
type ColaPrioridad []*Vehiculo

func (cp ColaPrioridad) Len() int { return len(cp) }

func (cp ColaPrioridad) Less(i, j int) bool {
	// Mayor prioridad primero (3 > 2 > 1)
	if cp[i].Prioridad != cp[j].Prioridad {
		return cp[i].Prioridad > cp[j].Prioridad
	}
	// Si tienen la misma prioridad, el que llegó primero
	return cp[i].TiempoLlegada.Before(cp[j].TiempoLlegada)
}

func (cp ColaPrioridad) Swap(i, j int) {
	cp[i], cp[j] = cp[j], cp[i]
}

func (cp *ColaPrioridad) Push(x interface{}) {
	*cp = append(*cp, x.(*Vehiculo))
}

func (cp *ColaPrioridad) Pop() interface{} {
	old := *cp
	n := len(old)
	item := old[n-1]
	*cp = old[0 : n-1]
	return item
}

// ============================================
// IMPLEMENTACIÓN CON RWMUTEX
// ============================================

// SimularTallerRWMutex simula el taller usando RWMutex
func SimularTallerRWMutex(vehiculos []*Vehiculo, numPlazas, numMecanicos int) {
	taller := NewTallerSimulacion(numPlazas, numMecanicos)
	tiempoInicio := time.Now()

	var rwMutex sync.RWMutex
	var wg sync.WaitGroup

	// Colas de prioridad para cada fase
	colaEntrada := make(ColaPrioridad, 0)
	colaReparacion := make(ColaPrioridad, 0)
	colaLimpieza := make(ColaPrioridad, 0)
	colaRevision := make(ColaPrioridad, 0)

	heap.Init(&colaEntrada)
	heap.Init(&colaReparacion)
	heap.Init(&colaLimpieza)
	heap.Init(&colaRevision)

	done := make(chan bool)
	vehiculosCompletados := 0
	totalVehiculos := len(vehiculos)

	// Goroutine para procesar Fase 1: Entrada
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				rwMutex.Lock()
				if colaEntrada.Len() == 0 {
					rwMutex.Unlock()
					time.Sleep(50 * time.Millisecond)
					continue
				}
				v := heap.Pop(&colaEntrada).(*Vehiculo)
				rwMutex.Unlock()

				taller.PlazasSem <- struct{}{}

				v.LogEstado("Entrada", "Esperando", time.Since(tiempoInicio))
				v.LogEstado("Entrada", "En Proceso", time.Since(tiempoInicio))
				time.Sleep(v.TiempoFase)
				v.LogEstado("Entrada", "Completado", time.Since(tiempoInicio))

				rwMutex.Lock()
				heap.Push(&colaReparacion, v)
				rwMutex.Unlock()
			}
		}
	}()

	// Goroutine para procesar Fase 2: Reparación
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				rwMutex.Lock()
				if colaReparacion.Len() == 0 {
					rwMutex.Unlock()
					time.Sleep(50 * time.Millisecond)
					continue
				}
				v := heap.Pop(&colaReparacion).(*Vehiculo)
				rwMutex.Unlock()

				taller.MecanicosSem <- struct{}{}

				v.LogEstado("Reparación", "Esperando", time.Since(tiempoInicio))
				v.LogEstado("Reparación", "En Proceso", time.Since(tiempoInicio))
				time.Sleep(v.TiempoFase)
				v.LogEstado("Reparación", "Completado", time.Since(tiempoInicio))

				<-taller.MecanicosSem

				rwMutex.Lock()
				heap.Push(&colaLimpieza, v)
				rwMutex.Unlock()
			}
		}
	}()

	// Goroutine para procesar Fase 3: Limpieza
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				rwMutex.Lock()
				if colaLimpieza.Len() == 0 {
					rwMutex.Unlock()
					time.Sleep(50 * time.Millisecond)
					continue
				}
				v := heap.Pop(&colaLimpieza).(*Vehiculo)
				rwMutex.Unlock()

				v.LogEstado("Limpieza", "Esperando", time.Since(tiempoInicio))
				v.LogEstado("Limpieza", "En Proceso", time.Since(tiempoInicio))
				time.Sleep(v.TiempoFase)
				v.LogEstado("Limpieza", "Completado", time.Since(tiempoInicio))

				rwMutex.Lock()
				heap.Push(&colaRevision, v)
				rwMutex.Unlock()
			}
		}
	}()

	// Goroutine para procesar Fase 4: Revisión Final
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				rwMutex.Lock()
				if colaRevision.Len() == 0 {
					rwMutex.Unlock()
					time.Sleep(50 * time.Millisecond)
					continue
				}
				v := heap.Pop(&colaRevision).(*Vehiculo)
				rwMutex.Unlock()

				v.LogEstado("Revisión Final", "Esperando", time.Since(tiempoInicio))
				v.LogEstado("Revisión Final", "En Proceso", time.Since(tiempoInicio))
				time.Sleep(v.TiempoFase)
				v.LogEstado("Revisión Final", "Completado", time.Since(tiempoInicio))

				<-taller.PlazasSem

				rwMutex.Lock()
				vehiculosCompletados++
				rwMutex.Unlock()
			}
		}
	}()

	// Agregar vehículos de forma aleatoria
	vehiculosMezclados := make([]*Vehiculo, len(vehiculos))
	copy(vehiculosMezclados, vehiculos)
	rand.Shuffle(len(vehiculosMezclados), func(i, j int) {
		vehiculosMezclados[i], vehiculosMezclados[j] = vehiculosMezclados[j], vehiculosMezclados[i]
	})

	for _, v := range vehiculosMezclados {
		rwMutex.Lock()
		heap.Push(&colaEntrada, v)
		rwMutex.Unlock()
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	// Esperar a que todos terminen
	for vehiculosCompletados < totalVehiculos {
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(500 * time.Millisecond)
	close(done)
	wg.Wait()
}

// ============================================
// IMPLEMENTACIÓN CON WAITGROUP
// ============================================

// SimularTallerWaitGroup simula el taller usando WaitGroup
func SimularTallerWaitGroup(vehiculos []*Vehiculo, numPlazas, numMecanicos int) {
	taller := NewTallerSimulacion(numPlazas, numMecanicos)
	tiempoInicio := time.Now()

	var mutex sync.Mutex
	var wg sync.WaitGroup

	colaEntrada := make(ColaPrioridad, 0)
	colaReparacion := make(ColaPrioridad, 0)
	colaLimpieza := make(ColaPrioridad, 0)
	colaRevision := make(ColaPrioridad, 0)

	heap.Init(&colaEntrada)
	heap.Init(&colaReparacion)
	heap.Init(&colaLimpieza)
	heap.Init(&colaRevision)

	done := make(chan bool)
	vehiculosCompletados := 0
	totalVehiculos := len(vehiculos)

	// Procesador de Fase 1: Entrada
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mutex.Lock()
				if colaEntrada.Len() == 0 {
					mutex.Unlock()
					time.Sleep(50 * time.Millisecond)
					continue
				}
				v := heap.Pop(&colaEntrada).(*Vehiculo)
				mutex.Unlock()

				wg.Add(1)
				go func(vehiculo *Vehiculo) {
					defer wg.Done()

					taller.PlazasSem <- struct{}{}

					vehiculo.LogEstado("Entrada", "Esperando", time.Since(tiempoInicio))
					vehiculo.LogEstado("Entrada", "En Proceso", time.Since(tiempoInicio))
					time.Sleep(vehiculo.TiempoFase)
					vehiculo.LogEstado("Entrada", "Completado", time.Since(tiempoInicio))

					mutex.Lock()
					heap.Push(&colaReparacion, vehiculo)
					mutex.Unlock()
				}(v)
			}
		}
	}()

	// Procesador de Fase 2: Reparación
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mutex.Lock()
				if colaReparacion.Len() == 0 {
					mutex.Unlock()
					time.Sleep(50 * time.Millisecond)
					continue
				}
				v := heap.Pop(&colaReparacion).(*Vehiculo)
				mutex.Unlock()

				wg.Add(1)
				go func(vehiculo *Vehiculo) {
					defer wg.Done()

					taller.MecanicosSem <- struct{}{}

					vehiculo.LogEstado("Reparación", "Esperando", time.Since(tiempoInicio))
					vehiculo.LogEstado("Reparación", "En Proceso", time.Since(tiempoInicio))
					time.Sleep(vehiculo.TiempoFase)
					vehiculo.LogEstado("Reparación", "Completado", time.Since(tiempoInicio))

					<-taller.MecanicosSem

					mutex.Lock()
					heap.Push(&colaLimpieza, vehiculo)
					mutex.Unlock()
				}(v)
			}
		}
	}()

	// Procesador de Fase 3: Limpieza
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mutex.Lock()
				if colaLimpieza.Len() == 0 {
					mutex.Unlock()
					time.Sleep(50 * time.Millisecond)
					continue
				}
				v := heap.Pop(&colaLimpieza).(*Vehiculo)
				mutex.Unlock()

				wg.Add(1)
				go func(vehiculo *Vehiculo) {
					defer wg.Done()

					vehiculo.LogEstado("Limpieza", "Esperando", time.Since(tiempoInicio))
					vehiculo.LogEstado("Limpieza", "En Proceso", time.Since(tiempoInicio))
					time.Sleep(vehiculo.TiempoFase)
					vehiculo.LogEstado("Limpieza", "Completado", time.Since(tiempoInicio))

					mutex.Lock()
					heap.Push(&colaRevision, vehiculo)
					mutex.Unlock()
				}(v)
			}
		}
	}()

	// Procesador de Fase 4: Revisión Final
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mutex.Lock()
				if colaRevision.Len() == 0 {
					mutex.Unlock()
					time.Sleep(50 * time.Millisecond)
					continue
				}
				v := heap.Pop(&colaRevision).(*Vehiculo)
				mutex.Unlock()

				wg.Add(1)
				go func(vehiculo *Vehiculo) {
					defer wg.Done()

					vehiculo.LogEstado("Revisión Final", "Esperando", time.Since(tiempoInicio))
					vehiculo.LogEstado("Revisión Final", "En Proceso", time.Since(tiempoInicio))
					time.Sleep(vehiculo.TiempoFase)
					vehiculo.LogEstado("Revisión Final", "Completado", time.Since(tiempoInicio))

					<-taller.PlazasSem

					mutex.Lock()
					vehiculosCompletados++
					mutex.Unlock()
				}(v)
			}
		}
	}()

	// Agregar vehículos de forma aleatoria
	vehiculosMezclados := make([]*Vehiculo, len(vehiculos))
	copy(vehiculosMezclados, vehiculos)
	rand.Shuffle(len(vehiculosMezclados), func(i, j int) {
		vehiculosMezclados[i], vehiculosMezclados[j] = vehiculosMezclados[j], vehiculosMezclados[i]
	})

	for _, v := range vehiculosMezclados {
		mutex.Lock()
		heap.Push(&colaEntrada, v)
		mutex.Unlock()
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	// Esperar a que todos terminen
	for vehiculosCompletados < totalVehiculos {
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(500 * time.Millisecond)
	close(done)
	wg.Wait()
}
