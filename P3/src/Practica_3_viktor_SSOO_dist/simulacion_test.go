package main

import (
	"fmt"
	"testing"
	"time"
)

// TestCase1 ejecuta el test con distribución equilibrada (10-10-10)
func TestCase1(t *testing.T) {
	fmt.Println("")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("    TEST CASE 1: DISTRIBUCIÓN EQUILIBRADA")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("Categoría A (Mecánica):   10 coches - Prioridad Alta   - 5s/fase")
	fmt.Println("Categoría B (Eléctrica):  10 coches - Prioridad Media  - 3s/fase")
	fmt.Println("Categoría C (Carrocería): 10 coches - Prioridad Baja   - 1s/fase")
	fmt.Println("Total: 30 coches")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("")

	vehiculos := generarVehiculos(10, 10, 10)

	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     EJECUCIÓN CON RWMUTEX                                  │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	inicioRW := time.Now()
	SimularTallerRWMutex(vehiculos, 5, 3)
	duracionRW := time.Since(inicioRW)

	time.Sleep(2 * time.Second)

	vehiculos = generarVehiculos(10, 10, 10)

	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     EJECUCIÓN CON WAITGROUP                                │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	inicioWG := time.Now()
	SimularTallerWaitGroup(vehiculos, 5, 3)
	duracionWG := time.Since(inicioWG)

	// Tiempos finales
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("    TIEMPOS REGISTRADOS - TEST CASE 1")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Printf("• RWMutex:   %v\n", duracionRW)
	fmt.Printf("• WaitGroup: %v\n", duracionWG)
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("")
}

// TestCase2 ejecuta el test con prioridad alta dominante (20-5-5)
func TestCase2(t *testing.T) {
	fmt.Println("")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("    TEST CASE 2: PRIORIDAD ALTA DOMINANTE")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("Categoría A (Mecánica):   20 coches - Prioridad Alta   - 5s/fase")
	fmt.Println("Categoría B (Eléctrica):   5 coches - Prioridad Media  - 3s/fase")
	fmt.Println("Categoría C (Carrocería):  5 coches - Prioridad Baja   - 1s/fase")
	fmt.Println("Total: 30 coches")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("")

	vehiculos := generarVehiculos(20, 5, 5)

	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     EJECUCIÓN CON RWMUTEX                                  │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	inicioRW := time.Now()
	SimularTallerRWMutex(vehiculos, 5, 3)
	duracionRW := time.Since(inicioRW)

	time.Sleep(2 * time.Second)

	vehiculos = generarVehiculos(20, 5, 5)

	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     EJECUCIÓN CON WAITGROUP                                │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	inicioWG := time.Now()
	SimularTallerWaitGroup(vehiculos, 5, 3)
	duracionWG := time.Since(inicioWG)

	// Tiempos finales
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("    TIEMPOS REGISTRADOS - TEST CASE 2")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Printf("• RWMutex:   %v\n", duracionRW)
	fmt.Printf("• WaitGroup: %v\n", duracionWG)
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("")
}

// TestCase3 ejecuta el test con prioridad baja dominante (5-5-20)
func TestCase3(t *testing.T) {
	fmt.Println("")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("    TEST CASE 3: PRIORIDAD BAJA DOMINANTE")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("Categoría A (Mecánica):    5 coches - Prioridad Alta   - 5s/fase")
	fmt.Println("Categoría B (Eléctrica):   5 coches - Prioridad Media  - 3s/fase")
	fmt.Println("Categoría C (Carrocería): 20 coches - Prioridad Baja   - 1s/fase")
	fmt.Println("Total: 30 coches")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("")

	vehiculos := generarVehiculos(5, 5, 20)

	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     EJECUCIÓN CON RWMUTEX                                  │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	inicioRW := time.Now()
	SimularTallerRWMutex(vehiculos, 5, 3)
	duracionRW := time.Since(inicioRW)

	time.Sleep(2 * time.Second)

	vehiculos = generarVehiculos(5, 5, 20)

	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     EJECUCIÓN CON WAITGROUP                                │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	inicioWG := time.Now()
	SimularTallerWaitGroup(vehiculos, 5, 3)
	duracionWG := time.Since(inicioWG)

	// SECCIÓN MODIFICADA: Solo mostramos los tiempos finales
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("    TIEMPOS REGISTRADOS - TEST CASE 3")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Printf("• RWMutex:   %v\n", duracionRW)
	fmt.Printf("• WaitGroup: %v\n", duracionWG)
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("")
}

// generarVehiculos genera la lista de vehículos para los tests
func generarVehiculos(numA, numB, numC int) []*Vehiculo {
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
