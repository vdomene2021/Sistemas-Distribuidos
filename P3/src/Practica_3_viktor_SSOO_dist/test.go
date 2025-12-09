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

	fmt.Println("")
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     RESULTADOS RWMUTEX - TEST CASE 1                       │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Printf("Tiempo total: %v\n", duracionRW)
	fmt.Printf("Coches procesados: %d\n", len(vehiculos))
	fmt.Println("")

	time.Sleep(2 * time.Second)

	vehiculos = generarVehiculos(10, 10, 10)

	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     EJECUCIÓN CON WAITGROUP                                │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	inicioWG := time.Now()
	SimularTallerWaitGroup(vehiculos, 5, 3)
	duracionWG := time.Since(inicioWG)

	fmt.Println("")
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     RESULTADOS WAITGROUP - TEST CASE 1                     │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Printf("Tiempo total: %v\n", duracionWG)
	fmt.Printf("Coches procesados: %d\n", len(vehiculos))
	fmt.Println("")

	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("    COMPARACIÓN - TEST CASE 1")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Printf("RWMutex:   %v\n", duracionRW)
	fmt.Printf("WaitGroup: %v\n", duracionWG)
	fmt.Println("")

	if duracionRW < duracionWG {
		diferencia := duracionWG - duracionRW
		porcentaje := float64(diferencia) / float64(duracionWG) * 100
		fmt.Printf("✓ RWMutex fue más rápido por %v (%.2f%%)\n", diferencia, porcentaje)
	} else {
		diferencia := duracionRW - duracionWG
		porcentaje := float64(diferencia) / float64(duracionRW) * 100
		fmt.Printf("✓ WaitGroup fue más rápido por %v (%.2f%%)\n", diferencia, porcentaje)
	}
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

	fmt.Println("")
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     RESULTADOS RWMUTEX - TEST CASE 2                       │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Printf("Tiempo total: %v\n", duracionRW)
	fmt.Printf("Coches procesados: %d\n", len(vehiculos))
	fmt.Println("")

	time.Sleep(2 * time.Second)

	vehiculos = generarVehiculos(20, 5, 5)

	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     EJECUCIÓN CON WAITGROUP                                │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	inicioWG := time.Now()
	SimularTallerWaitGroup(vehiculos, 5, 3)
	duracionWG := time.Since(inicioWG)

	fmt.Println("")
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     RESULTADOS WAITGROUP - TEST CASE 2                     │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Printf("Tiempo total: %v\n", duracionWG)
	fmt.Printf("Coches procesados: %d\n", len(vehiculos))
	fmt.Println("")

	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("    COMPARACIÓN - TEST CASE 2")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Printf("RWMutex:   %v\n", duracionRW)
	fmt.Printf("WaitGroup: %v\n", duracionWG)
	fmt.Println("")

	if duracionRW < duracionWG {
		diferencia := duracionWG - duracionRW
		porcentaje := float64(diferencia) / float64(duracionWG) * 100
		fmt.Printf("✓ RWMutex fue más rápido por %v (%.2f%%)\n", diferencia, porcentaje)
	} else {
		diferencia := duracionRW - duracionWG
		porcentaje := float64(diferencia) / float64(duracionRW) * 100
		fmt.Printf("✓ WaitGroup fue más rápido por %v (%.2f%%)\n", diferencia, porcentaje)
	}
	fmt.Println("")
	fmt.Println("Observación: 66.7% de coches con prioridad alta (5s/fase)")
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

	fmt.Println("")
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     RESULTADOS RWMUTEX - TEST CASE 3                       │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Printf("Tiempo total: %v\n", duracionRW)
	fmt.Printf("Coches procesados: %d\n", len(vehiculos))
	fmt.Println("")

	time.Sleep(2 * time.Second)

	vehiculos = generarVehiculos(5, 5, 20)

	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     EJECUCIÓN CON WAITGROUP                                │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	inicioWG := time.Now()
	SimularTallerWaitGroup(vehiculos, 5, 3)
	duracionWG := time.Since(inicioWG)

	fmt.Println("")
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│     RESULTADOS WAITGROUP - TEST CASE 3                     │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")
	fmt.Printf("Tiempo total: %v\n", duracionWG)
	fmt.Printf("Coches procesados: %d\n", len(vehiculos))
	fmt.Println("")

	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("    COMPARACIÓN - TEST CASE 3")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Printf("RWMutex:   %v\n", duracionRW)
	fmt.Printf("WaitGroup: %v\n", duracionWG)
	fmt.Println("")

	if duracionRW < duracionWG {
		diferencia := duracionWG - duracionRW
		porcentaje := float64(diferencia) / float64(duracionWG) * 100
		fmt.Printf("✓ RWMutex fue más rápido por %v (%.2f%%)\n", diferencia, porcentaje)
	} else {
		diferencia := duracionRW - duracionWG
		porcentaje := float64(diferencia) / float64(duracionRW) * 100
		fmt.Printf("✓ WaitGroup fue más rápido por %v (%.2f%%)\n", diferencia, porcentaje)
	}
	fmt.Println("")
	fmt.Println("Observación: 66.7% de coches con prioridad baja (1s/fase)")
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
