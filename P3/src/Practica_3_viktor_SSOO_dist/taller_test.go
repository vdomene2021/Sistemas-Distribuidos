package main

import (
	"testing"
	"time"
)

// Función auxiliar para generar vehículos de test
func generarVehiculosTest(numA, numB, numC int) []*Vehiculo {
	var vehiculos []*Vehiculo
	id := 1

	for i := 0; i < numA; i++ {
		vehiculos = append(vehiculos, NewVehiculo(id, Mecanica))
		id++
	}

	for i := 0; i < numB; i++ {
		vehiculos = append(vehiculos, NewVehiculo(id, Electrica))
		id++
	}

	for i := 0; i < numC; i++ {
		vehiculos = append(vehiculos, NewVehiculo(id, Carroceria))
		id++
	}

	return vehiculos
}

// ============================================
// TEST CASE 1: A=10, B=10, C=10
// ============================================

func TestCase1_RWMutex(t *testing.T) {
	t.Log("=== TEST CASE 1 - RWMutex ===")
	t.Log("Categoría A: 10, B: 10, C: 10")

	vehiculos := generarVehiculosTest(10, 10, 10)
	inicio := time.Now()

	SimularTallerRWMutex(vehiculos, 5, 3)

	duracion := time.Since(inicio)
	t.Logf("Duración: %v", duracion)
	t.Logf("Coches procesados: %d", len(vehiculos))
	t.Logf("Tiempo promedio por coche: %v", duracion/time.Duration(len(vehiculos)))
}

func TestCase1_WaitGroup(t *testing.T) {
	t.Log("=== TEST CASE 1 - WaitGroup ===")
	t.Log("Categoría A: 10, B: 10, C: 10")

	vehiculos := generarVehiculosTest(10, 10, 10)
	inicio := time.Now()

	SimularTallerWaitGroup(vehiculos, 5, 3)

	duracion := time.Since(inicio)
	t.Logf("Duración: %v", duracion)
	t.Logf("Coches procesados: %d", len(vehiculos))
	t.Logf("Tiempo promedio por coche: %v", duracion/time.Duration(len(vehiculos)))
}

// ============================================
// TEST CASE 2: A=20, B=5, C=5
// ============================================

func TestCase2_RWMutex(t *testing.T) {
	t.Log("=== TEST CASE 2 - RWMutex ===")
	t.Log("Categoría A: 20, B: 5, C: 5")

	vehiculos := generarVehiculosTest(20, 5, 5)
	inicio := time.Now()

	SimularTallerRWMutex(vehiculos, 5, 3)

	duracion := time.Since(inicio)
	t.Logf("Duración: %v", duracion)
	t.Logf("Coches procesados: %d", len(vehiculos))
	t.Logf("Tiempo promedio por coche: %v", duracion/time.Duration(len(vehiculos)))
}

func TestCase2_WaitGroup(t *testing.T) {
	t.Log("=== TEST CASE 2 - WaitGroup ===")
	t.Log("Categoría A: 20, B: 5, C: 5")

	vehiculos := generarVehiculosTest(20, 5, 5)
	inicio := time.Now()

	SimularTallerWaitGroup(vehiculos, 5, 3)

	duracion := time.Since(inicio)
	t.Logf("Duración: %v", duracion)
	t.Logf("Coches procesados: %d", len(vehiculos))
	t.Logf("Tiempo promedio por coche: %v", duracion/time.Duration(len(vehiculos)))
}

// ============================================
// TEST CASE 3: A=5, B=5, C=20
// ============================================

func TestCase3_RWMutex(t *testing.T) {
	t.Log("=== TEST CASE 3 - RWMutex ===")
	t.Log("Categoría A: 5, B: 5, C: 20")

	vehiculos := generarVehiculosTest(5, 5, 20)
	inicio := time.Now()

	SimularTallerRWMutex(vehiculos, 5, 3)

	duracion := time.Since(inicio)
	t.Logf("Duración: %v", duracion)
	t.Logf("Coches procesados: %d", len(vehiculos))
	t.Logf("Tiempo promedio por coche: %v", duracion/time.Duration(len(vehiculos)))
}

func TestCase3_WaitGroup(t *testing.T) {
	t.Log("=== TEST CASE 3 - WaitGroup ===")
	t.Log("Categoría A: 5, B: 5, C: 20")

	vehiculos := generarVehiculosTest(5, 5, 20)
	inicio := time.Now()

	SimularTallerWaitGroup(vehiculos, 5, 3)

	duracion := time.Since(inicio)
	t.Logf("Duración: %v", duracion)
	t.Logf("Coches procesados: %d", len(vehiculos))
	t.Logf("Tiempo promedio por coche: %v", duracion/time.Duration(len(vehiculos)))
}

// ============================================
// BENCHMARKS
// ============================================

func BenchmarkRWMutex_Case1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vehiculos := generarVehiculosTest(10, 10, 10)
		SimularTallerRWMutex(vehiculos, 5, 3)
	}
}

func BenchmarkWaitGroup_Case1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vehiculos := generarVehiculosTest(10, 10, 10)
		SimularTallerWaitGroup(vehiculos, 5, 3)
	}
}

func BenchmarkRWMutex_Case2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vehiculos := generarVehiculosTest(20, 5, 5)
		SimularTallerRWMutex(vehiculos, 5, 3)
	}
}

func BenchmarkWaitGroup_Case2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vehiculos := generarVehiculosTest(20, 5, 5)
		SimularTallerWaitGroup(vehiculos, 5, 3)
	}
}

func BenchmarkRWMutex_Case3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vehiculos := generarVehiculosTest(5, 5, 20)
		SimularTallerRWMutex(vehiculos, 5, 3)
	}
}

func BenchmarkWaitGroup_Case3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vehiculos := generarVehiculosTest(5, 5, 20)
		SimularTallerWaitGroup(vehiculos, 5, 3)
	}
}

// ============================================
// TESTS UNITARIOS
// ============================================

func TestNewVehiculo(t *testing.T) {
	// Test categoría A
	vA := NewVehiculo(1, Mecanica)
	if vA.Prioridad != PrioridadAlta {
		t.Errorf("Vehículo mecánico debería tener prioridad alta, obtuvo: %d", vA.Prioridad)
	}
	if vA.TiempoFase != 5*time.Second {
		t.Errorf("Vehículo mecánico debería tener 5s por fase, obtuvo: %v", vA.TiempoFase)
	}

	// Test categoría B
	vB := NewVehiculo(2, Electrica)
	if vB.Prioridad != PrioridadMedia {
		t.Errorf("Vehículo eléctrico debería tener prioridad media, obtuvo: %d", vB.Prioridad)
	}
	if vB.TiempoFase != 3*time.Second {
		t.Errorf("Vehículo eléctrico debería tener 3s por fase, obtuvo: %v", vB.TiempoFase)
	}

	// Test categoría C
	vC := NewVehiculo(3, Carroceria)
	if vC.Prioridad != PrioridadBaja {
		t.Errorf("Vehículo de carrocería debería tener prioridad baja, obtuvo: %d", vC.Prioridad)
	}
	if vC.TiempoFase != 1*time.Second {
		t.Errorf("Vehículo de carrocería debería tener 1s por fase, obtuvo: %v", vC.TiempoFase)
	}
}

func TestColaPrioridad(t *testing.T) {
	cola := make(ColaPrioridad, 0)

	vBaja := NewVehiculo(1, Carroceria)
	time.Sleep(10 * time.Millisecond)
	vMedia := NewVehiculo(2, Electrica)
	time.Sleep(10 * time.Millisecond)
	vAlta := NewVehiculo(3, Mecanica)

	// Agregar en orden: Baja, Media, Alta
	cola.Push(vBaja)
	cola.Push(vMedia)
	cola.Push(vAlta)

	// Verificar que Less funciona correctamente
	if !cola.Less(2, 0) { // Alta > Baja
		t.Error("Prioridad alta debería ser mayor que baja")
	}
	if !cola.Less(1, 0) { // Media > Baja
		t.Error("Prioridad media debería ser mayor que baja")
	}
	if !cola.Less(2, 1) { // Alta > Media
		t.Error("Prioridad alta debería ser mayor que media")
	}
}

func TestTaller(t *testing.T) {
	taller := NewTaller(5, 3)

	if taller.NumPlazas != 5 {
		t.Errorf("Taller debería tener 5 plazas, tiene: %d", taller.NumPlazas)
	}

	if taller.NumMecanicos != 3 {
		t.Errorf("Taller debería tener 3 mecánicos, tiene: %d", taller.NumMecanicos)
	}

	if cap(taller.PlazasSem) != 5 {
		t.Errorf("Semáforo de plazas debería tener capacidad 5, tiene: %d", cap(taller.PlazasSem))
	}

	if cap(taller.MecanicosSem) != 3 {
		t.Errorf("Semáforo de mecánicos debería tener capacidad 3, tiene: %d", cap(taller.MecanicosSem))
	}
}
