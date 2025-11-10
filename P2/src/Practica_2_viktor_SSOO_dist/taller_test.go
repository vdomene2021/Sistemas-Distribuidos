// En /P2/src/Practica_2_viktor_SSOO_dist/taller_test.go

package main

import (
	"fmt"
	"sync"
	"testing"
)

// setupTest crea un taller limpio para una sola ejecución de test
func setupTest(t *testing.T, mechanicConfig map[Especialidad]int) (*Taller, *VehiculoManager, *IncidenciaManager) {
	// Crear managers limpios
	mm := NewMecanicoManager()
	vm := NewVehiculoManager() // <-- ¡AQUÍ ESTABA EL ERROR, AHORA CORREGIDO!
	im := NewIncidenciaManager()

	// Crear mecánicos según la configuración
	for especialidad, count := range mechanicConfig {
		for i := 0; i < count; i++ {
			mm.CrearMecanico(string(especialidad), especialidad, 5)
		}
	}

	// Crear el Taller
	taller := NewTaller(100, mm, vm, im)
	taller.IniciarTaller()

	// Detener el taller cuando el test termine
	t.Cleanup(func() {
		taller.DetenerTaller()
	})

	return taller, vm, im
}

// runTestSimulation simula la llegada de N coches UNA SOLA VEZ
func runTestSimulation(t *testing.T, numCars int, tipo TipoIncidencia, mechanicConfig map[Especialidad]int) {
	taller, vm, im := setupTest(t, mechanicConfig)
	var wg sync.WaitGroup
	taller.wg = &wg

	// Creamos vehículos e incidencias NUEVOS
	vehiculos := make([]Vehiculo, numCars)
	incidencias := make([]Incidencia, numCars)
	for i := 0; i < numCars; i++ {
		matricula := fmt.Sprintf("Matricula%d", i+1)
		marca := fmt.Sprintf("TEST-CAR%d", i+1)
		modelo := fmt.Sprintf("modelo%d", i+1)

		vehiculos[i] = vm.CrearVehiculo(matricula, marca, modelo, 1)
		incidencias[i] = im.CrearIncidencia(tipo, Alta, "Test incidence")
	}

	// Enviar todos los coches al taller
	for i := 0; i < numCars; i++ {
		taller.AgregarTrabajo(vehiculos[i], incidencias[i])
	}

	// Esperar a que el WaitGroup (wg) nos diga que
	// TODOS los trabajos se han completado (llegado a 0).
	wg.Wait()
}

// --- CASO 1: Test de Comparativa DUPLICANDO COCHES ---
//
// (Plantilla base: 1 mecánico de cada tipo)
var configBase = map[Especialidad]int{
	EspecialidadMecanica:   1,
	EspecialidadElectrica:  1,
	EspecialidadCarroceria: 1,
}

func Test_DuplicarCoches_3(t *testing.T) {
	runTestSimulation(t, 3, Mecanica, configBase)
}

func Test_DuplicarCoches_6(t *testing.T) {
	runTestSimulation(t, 6, Mecanica, configBase)
}

// --- CASO 2: Test de Comparativa DUPLICANDO PLANTILLA ---
//

// La configuración de 3 mecánicos (1-1-1) ya la tenemos en 'configBase'

// Nueva configuración con 6 mecánicos (2-2-2)
var configCaso2_6Mecanicos = map[Especialidad]int{
	EspecialidadMecanica:   2, // 2 de mecánica
	EspecialidadElectrica:  2, // 2 de eléctrica
	EspecialidadCarroceria: 2, // 2 de carrocería
}

// Usamos 10 coches de Mecánica como carga de trabajo estándar para comparar
func Test_DuplicarPlantilla_Con3Mecanicos(t *testing.T) {
	runTestSimulation(t, 7, Mecanica, configBase) // Usa la config de 3 (1-1-1)
}

func Test_DuplicarPlantilla_Con6Mecanicos(t *testing.T) {
	runTestSimulation(t, 7, Mecanica, configCaso2_6Mecanicos) // Usa la config de 6 (2-2-2)
}

// --- CASO 3: Test de Comparativa PROPORCIONES ---
//

// Configuración Favorable: 3 mecánicos, 1 eléctrico, 1 carrocería
var configProporcion_Favorable = map[Especialidad]int{
	EspecialidadMecanica:   3,
	EspecialidadElectrica:  1,
	EspecialidadCarroceria: 1,
}

// Configuración Desfavorable: 1 mecánico, 3 eléctricos, 3 carrocería
var configProporcion_Desfavorable = map[Especialidad]int{
	EspecialidadMecanica:   1,
	EspecialidadElectrica:  3,
	EspecialidadCarroceria: 3,
}

// Usamos 10 coches de Mecánica para ver la diferencia
func Test_Proporcion_Favorable(t *testing.T) {
	runTestSimulation(t, 10, Mecanica, configProporcion_Favorable)
}

func Test_Proporcion_Desfavorable(t *testing.T) {
	runTestSimulation(t, 10, Mecanica, configProporcion_Desfavorable)
}
