package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Parámetros de configuración
	var (
		numPlazas    int
		numMecanicos int
		testCase     int
		metodo       string
	)

	flag.IntVar(&numPlazas, "plazas", 5, "Número de plazas en el taller")
	flag.IntVar(&numMecanicos, "mecanicos", 3, "Número de mecánicos")
	flag.IntVar(&testCase, "test", 1, "Caso de test (1, 2 o 3)")
	flag.StringVar(&metodo, "metodo", "rwmutex", "Método: 'rwmutex' o 'waitgroup'")
	flag.Parse()

	fmt.Println("=================================================")
	fmt.Println("      SIMULACIÓN DEL TALLER DEL PUEBLO")
	fmt.Println("=================================================")
	fmt.Printf("Plazas disponibles: %d\n", numPlazas)
	fmt.Printf("Mecánicos disponibles: %d\n", numMecanicos)
	fmt.Printf("Test Case: %d\n", testCase)
	fmt.Printf("Método: %s\n", metodo)
	fmt.Println("=================================================\n")

	// Configurar número de coches según el test case
	var cochesA, cochesB, cochesC int
	switch testCase {
	case 1:
		cochesA, cochesB, cochesC = 10, 10, 10
	case 2:
		cochesA, cochesB, cochesC = 20, 5, 5
	case 3:
		cochesA, cochesB, cochesC = 5, 5, 20
	default:
		cochesA, cochesB, cochesC = 10, 10, 10
	}

	fmt.Printf("Test Case %d: Categoría A=%d, B=%d, C=%d (Total: %d coches)\n\n",
		testCase, cochesA, cochesB, cochesC, cochesA+cochesB+cochesC)

	// Generar vehículos
	vehiculos := generarVehiculos(cochesA, cochesB, cochesC)

	// Ejecutar simulación según el método
	inicio := time.Now()

	if metodo == "rwmutex" {
		fmt.Println(">>> Iniciando simulación con RWMutex <<<\n")
		SimularTallerRWMutex(vehiculos, numPlazas, numMecanicos)
	} else if metodo == "waitgroup" {
		fmt.Println(">>> Iniciando simulación con WaitGroup <<<\n")
		SimularTallerWaitGroup(vehiculos, numPlazas, numMecanicos)
	} else {
		fmt.Println("Método no válido. Usa 'rwmutex' o 'waitgroup'")
		return
	}

	duracion := time.Since(inicio)

	fmt.Println("\n=================================================")
	fmt.Println("           SIMULACIÓN COMPLETADA")
	fmt.Println("=================================================")
	fmt.Printf("Tiempo total de ejecución: %v\n", duracion)
	fmt.Printf("Coches procesados: %d\n", len(vehiculos))
	fmt.Printf("Tiempo promedio por coche: %v\n", duracion/time.Duration(len(vehiculos)))
	fmt.Println("=================================================")
}

// generarVehiculos genera la lista de vehículos para la simulación
func generarVehiculos(numA, numB, numC int) []*Vehiculo {
	var vehiculos []*Vehiculo
	id := 1

	// Generar coches categoría A (Mecánica - Prioridad Alta)
	for i := 0; i < numA; i++ {
		vehiculos = append(vehiculos, NewVehiculo(id, Mecanica))
		id++
	}

	// Generar coches categoría B (Eléctrica - Prioridad Media)
	for i := 0; i < numB; i++ {
		vehiculos = append(vehiculos, NewVehiculo(id, Electrica))
		id++
	}

	// Generar coches categoría C (Carrocería - Prioridad Baja)
	for i := 0; i < numC; i++ {
		vehiculos = append(vehiculos, NewVehiculo(id, Carroceria))
		id++
	}

	return vehiculos
}
