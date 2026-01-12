package main

import (
	"fmt"
	"testing"
)

// ============================================================================
// BLOQUE 1: Recursos -> 6 Plazas, 3 Mecánicos (Fila 1 de Tabla Recursos)
// ============================================================================

// TEST 1: 10 Coches A, 10 Coches B, 10 Coches C
// Recursos: 6 Plazas, 3 Mecánicos
func Test_1_RecursosA_CochesEquilibrados(t *testing.T) {
	fmt.Println("\n=== TEST 1: [6 Plazas, 3 Mecánicos] | [10 A, 10 B, 10 C] ===")
	EjecutarSimulacion(6, 3, 10, 10, 10)
}

// TEST 2: 20 Coches A, 5 Coches B, 5 Coches C
// Recursos: 6 Plazas, 3 Mecánicos
func Test_2_RecursosA_MayoriaA(t *testing.T) {
	fmt.Println("\n=== TEST 2: [6 Plazas, 3 Mecánicos] | [20 A, 5 B, 5 C] ===")
	EjecutarSimulacion(6, 3, 20, 5, 5)
}

// TEST 3: 5 Coches A, 5 Coches B, 20 Coches C
// Recursos: 6 Plazas, 3 Mecánicos
func Test_3_RecursosA_MayoriaC(t *testing.T) {
	fmt.Println("\n=== TEST 3: [6 Plazas, 3 Mecánicos] | [5 A, 5 B, 20 C] ===")
	EjecutarSimulacion(6, 3, 5, 5, 20)
}

// ============================================================================
// BLOQUE 2: Recursos -> 4 Plazas, 4 Mecánicos (Fila 2 de Tabla Recursos)
// ============================================================================

// TEST 4: 10 Coches A, 10 Coches B, 10 Coches C
// Recursos: 4 Plazas, 4 Mecánicos
func Test_4_RecursosB_CochesEquilibrados(t *testing.T) {
	fmt.Println("\n=== TEST 4: [4 Plazas, 4 Mecánicos] | [10 A, 10 B, 10 C] ===")
	EjecutarSimulacion(4, 4, 10, 10, 10)
}

// TEST 5: 20 Coches A, 5 Coches B, 5 Coches C
// Recursos: 4 Plazas, 4 Mecánicos
func Test_5_RecursosB_MayoriaA(t *testing.T) {
	fmt.Println("\n=== TEST 5: [4 Plazas, 4 Mecánicos] | [20 A, 5 B, 5 C] ===")
	EjecutarSimulacion(4, 4, 20, 5, 5)
}

// TEST 6: 5 Coches A, 5 Coches B, 20 Coches C
// Recursos: 4 Plazas, 4 Mecánicos
func Test_6_RecursosB_MayoriaC(t *testing.T) {
	fmt.Println("\n=== TEST 6: [4 Plazas, 4 Mecánicos] | [5 A, 5 B, 20 C] ===")
	EjecutarSimulacion(4, 4, 5, 5, 20)
}
