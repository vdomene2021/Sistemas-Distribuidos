package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Inicializar managers
	clienteManager := NewClienteManager()
	vehiculoManager := NewVehiculoManager()
	incidenciaManager := NewIncidenciaManager()
	mecanicoManager := NewMecanicoManager()

	// Crear taller
	taller := NewTaller(mecanicoManager, vehiculoManager, incidenciaManager)
	taller.IniciarTaller()

	// Menú principal
	scanner := bufio.NewScanner(os.Stdin)

	for {
		mostrarMenuPrincipal()
		fmt.Print("Seleccione una opción: ")
		scanner.Scan()
		opcion := scanner.Text()

		if opcion == "1" {
			menuClientes(scanner, clienteManager, vehiculoManager)
		} else if opcion == "2" {
			menuVehiculos(scanner, vehiculoManager, clienteManager, incidenciaManager)
		} else if opcion == "3" {
			menuIncidencias(scanner, incidenciaManager, vehiculoManager)
		} else if opcion == "4" {
			menuMecanicos(scanner, mecanicoManager, taller)
		} else if opcion == "5" {
			menuTaller(scanner, taller, vehiculoManager, incidenciaManager)
		} else if opcion == "0" {
			fmt.Println("Saliendo del sistema...")
			taller.DetenerTaller()
			break
		} else {
			fmt.Println("Opción no válida")
		}
	}
}
