package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Inicializar managers
	clienteManager := NewClienteManager()
	vehiculoManager := NewVehiculoManager()
	incidenciaManager := NewIncidenciaManager()
	mecanicoManager := NewMecanicoManager()

	// Crear mecánicos iniciales
	mecanicoManager.CrearMecanico("Carlos", EspecialidadMecanica, 5)
	mecanicoManager.CrearMecanico("Ana", EspecialidadElectrica, 3)
	mecanicoManager.CrearMecanico("Luis", EspecialidadCarroceria, 7)

	// Crear taller con 6 plazas (2 por mecánico)
	taller := NewTaller(6, mecanicoManager, vehiculoManager, incidenciaManager)
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
			menuIncidencias(scanner, incidenciaManager)
		} else if opcion == "4" {
			menuMecanicos(scanner, mecanicoManager)
		} else if opcion == "5" {
			menuTaller(scanner, taller, vehiculoManager, incidenciaManager)
		} else if opcion == "6" {
			ejecutarSimulacion(taller, clienteManager, vehiculoManager, incidenciaManager)
		} else if opcion == "0" {
			fmt.Println("Saliendo del sistema...")
			taller.DetenerTaller()
			time.Sleep(1 * time.Second)
			break
		} else {
			fmt.Println("Opción no válida")
		}
	}
}

func mostrarMenuPrincipal() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║            TALLER MECÁNICO             ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println("1. Gestión de Clientes")
	fmt.Println("2. Gestión de Vehículos")
	fmt.Println("3. Gestión de Incidencias")
	fmt.Println("4. Gestión de Mecánicos")
	fmt.Println("5. Estado del Taller")
	fmt.Println("6. Ejecutar Simulación")
	fmt.Println("0. Salir")
	fmt.Println("──────────────────────────────────────────")
}

func menuClientes(scanner *bufio.Scanner, cm *ClienteManager, vm *VehiculoManager) {
	for {
		fmt.Println("\n=== GESTIÓN DE CLIENTES ===")
		fmt.Println("1. Crear Cliente")
		fmt.Println("2. Listar Clientes")
		fmt.Println("3. Actualizar Cliente")
		fmt.Println("4. Eliminar Cliente")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		scanner.Scan()
		opcion := scanner.Text()

		if opcion == "1" {
			fmt.Print("Nombre: ")
			scanner.Scan()
			nombre := scanner.Text()
			fmt.Print("Teléfono: ")
			scanner.Scan()
			telefono := scanner.Text()
			fmt.Print("Email: ")
			scanner.Scan()
			email := scanner.Text()

			cliente := cm.CrearCliente(nombre, telefono, email)
			fmt.Printf("Cliente creado con ID: %d\n", cliente.ID)
		} else if opcion == "2" {
			clientes := cm.ListarClientes()
			fmt.Println("\n--- LISTA DE CLIENTES ---")
			for i := 0; i < len(clientes); i++ {
				fmt.Printf("ID: %d | Nombre: %s | Tel: %s | Email: %s\n",
					clientes[i].ID, clientes[i].Nombre, clientes[i].Telefono, clientes[i].Email)
			}
		} else if opcion == "3" {
			fmt.Print("ID del cliente: ")
			scanner.Scan()
			id, _ := strconv.Atoi(scanner.Text())
			fmt.Print("Nuevo nombre (vacío para no cambiar): ")
			scanner.Scan()
			nombre := scanner.Text()
			fmt.Print("Nuevo teléfono (vacío para no cambiar): ")
			scanner.Scan()
			telefono := scanner.Text()
			fmt.Print("Nuevo email (vacío para no cambiar): ")
			scanner.Scan()
			email := scanner.Text()

			err := cm.ActualizarCliente(id, nombre, telefono, email)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Cliente actualizado")
			}
		} else if opcion == "4" {
			fmt.Print("ID del cliente a eliminar: ")
			scanner.Scan()
			id, _ := strconv.Atoi(scanner.Text())
			err := cm.EliminarCliente(id)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Cliente eliminado")
			}
		} else if opcion == "0" {
			break
		}
	}
}

func menuVehiculos(scanner *bufio.Scanner, vm *VehiculoManager, cm *ClienteManager, im *IncidenciaManager) {
	for {
		fmt.Println("\n=== GESTIÓN DE VEHÍCULOS ===")
		fmt.Println("1. Crear Vehículo")
		fmt.Println("2. Listar Vehículos")
		fmt.Println("3. Actualizar Vehículo")
		fmt.Println("4. Eliminar Vehículo")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		scanner.Scan()
		opcion := scanner.Text()

		if opcion == "1" {
			fmt.Print("Matrícula: ")
			scanner.Scan()
			matricula := scanner.Text()
			fmt.Print("Marca: ")
			scanner.Scan()
			marca := scanner.Text()
			fmt.Print("Modelo: ")
			scanner.Scan()
			modelo := scanner.Text()
			fmt.Print("ID del cliente: ")
			scanner.Scan()
			clienteID, _ := strconv.Atoi(scanner.Text())

			vehiculo := vm.CrearVehiculo(matricula, marca, modelo, clienteID)
			fmt.Printf("Vehículo creado con ID: %d\n", vehiculo.ID)
		} else if opcion == "2" {
			vehiculos := vm.ListarVehiculos()
			fmt.Println("\n--- LISTA DE VEHÍCULOS ---")
			for i := 0; i < len(vehiculos); i++ {
				fmt.Printf("ID: %d | Matrícula: %s | Marca: %s | Modelo: %s | Tiempo acumulado: %.1fs\n",
					vehiculos[i].ID, vehiculos[i].Matricula, vehiculos[i].Marca, vehiculos[i].Modelo, vehiculos[i].TiempoAcumulado)
			}
		} else if opcion == "0" {
			break
		}
	}
}

func menuIncidencias(scanner *bufio.Scanner, im *IncidenciaManager) {
	for {
		fmt.Println("\n=== GESTIÓN DE INCIDENCIAS ===")
		fmt.Println("1. Crear Incidencia")
		fmt.Println("2. Listar Incidencias")
		fmt.Println("3. Cambiar Estado")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		scanner.Scan()
		opcion := scanner.Text()

		if opcion == "1" {
			fmt.Println("Tipo (mecanica/electrica/carroceria): ")
			scanner.Scan()
			tipoStr := scanner.Text()
			tipo := TipoIncidencia(tipoStr)

			fmt.Print("Prioridad (baja/media/alta): ")
			scanner.Scan()
			prioridadStr := scanner.Text()
			prioridad := Prioridad(prioridadStr)

			fmt.Print("Descripción: ")
			scanner.Scan()
			descripcion := scanner.Text()

			incidencia := im.CrearIncidencia(tipo, prioridad, descripcion)
			fmt.Printf("Incidencia creada con ID: %d\n", incidencia.ID)
		} else if opcion == "2" {
			incidencias := im.ListarIncidencias()
			fmt.Println("\n--- LISTA DE INCIDENCIAS ---")
			for i := 0; i < len(incidencias); i++ {
				fmt.Printf("ID: %d | Tipo: %s | Prioridad: %s | Estado: %s | Descripción: %s\n",
					incidencias[i].ID, incidencias[i].Tipo, incidencias[i].Prioridad,
					incidencias[i].Estado, incidencias[i].Descripcion)
			}
		} else if opcion == "0" {
			break
		}
	}
}

func menuMecanicos(scanner *bufio.Scanner, mm *MecanicoManager) {
	for {
		fmt.Println("\n=== GESTIÓN DE MECÁNICOS ===")
		fmt.Println("1. Crear Mecánico")
		fmt.Println("2. Listar Mecánicos")
		fmt.Println("3. Cambiar Estado (Alta/Baja)")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		scanner.Scan()
		opcion := scanner.Text()

		if opcion == "1" {
			fmt.Print("Nombre: ")
			scanner.Scan()
			nombre := scanner.Text()
			fmt.Print("Especialidad (mecanica/electrica/carroceria): ")
			scanner.Scan()
			espStr := scanner.Text()
			especialidad := Especialidad(espStr)
			fmt.Print("Años de experiencia: ")
			scanner.Scan()
			exp, _ := strconv.Atoi(scanner.Text())

			mecanico := mm.CrearMecanico(nombre, especialidad, exp)
			fmt.Printf("Mecánico creado con ID: %d\n", mecanico.ID)
		} else if opcion == "2" {
			mecanicos := mm.ListarMecanicos()
			fmt.Println("\n--- LISTA DE MECÁNICOS ---")
			for i := 0; i < len(mecanicos); i++ {
				estado := "Activo"
				if !mecanicos[i].Activo {
					estado = "Inactivo"
				}
				ocupado := ""
				if mecanicos[i].Ocupado {
					ocupado = " (OCUPADO)"
				}
				fmt.Printf("ID: %d | Nombre: %s | Especialidad: %s | Exp: %d años | Estado: %s%s\n",
					mecanicos[i].ID, mecanicos[i].Nombre, mecanicos[i].Especialidad,
					mecanicos[i].Experiencia, estado, ocupado)
			}
		} else if opcion == "3" {
			fmt.Print("ID del mecánico: ")
			scanner.Scan()
			id, _ := strconv.Atoi(scanner.Text())
			fmt.Print("¿Dar de alta? (s/n): ")
			scanner.Scan()
			respuesta := strings.ToLower(scanner.Text())
			activo := respuesta == "s"

			err := mm.CambiarEstadoActivo(id, activo)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Estado cambiado")
			}
		} else if opcion == "0" {
			break
		}
	}
}

func menuTaller(scanner *bufio.Scanner, t *Taller, vm *VehiculoManager, im *IncidenciaManager) {
	for {
		fmt.Println("\n=== ESTADO DEL TALLER ===")
		fmt.Println("1. Ver Estado")
		fmt.Println("2. Enviar Vehículo al Taller")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		scanner.Scan()
		opcion := scanner.Text()

		if opcion == "1" {
			t.ObtenerEstadoTaller()
		} else if opcion == "2" {
			fmt.Print("ID del vehículo: ")
			scanner.Scan()
			vehiculoID, _ := strconv.Atoi(scanner.Text())
			fmt.Print("ID de la incidencia: ")
			scanner.Scan()
			incidenciaID, _ := strconv.Atoi(scanner.Text())

			vehiculo, existe := vm.ObtenerVehiculo(vehiculoID)
			if !existe {
				fmt.Println("Vehículo no encontrado")
				continue
			}

			incidencia, existe := im.ObtenerIncidencia(incidenciaID)
			if !existe {
				fmt.Println("Incidencia no encontrada")
				continue
			}

			vm.AsignarIncidencia(vehiculoID, incidenciaID)
			t.AgregarTrabajo(vehiculo, incidencia)
			fmt.Println("Vehículo enviado al taller")
		} else if opcion == "0" {
			break
		}
	}
}

func ejecutarSimulacion(t *Taller, cm *ClienteManager, vm *VehiculoManager, im *IncidenciaManager) {
	fmt.Println("\nINICIANDO SIMULACIÓN...")
	fmt.Println("Creando clientes, vehículos e incidencias de prueba...")
	fmt.Println()

	// Crear clientes
	cliente1 := cm.CrearCliente("Juan Pérez", "666111222", "juan@email.com")
	cliente2 := cm.CrearCliente("María García", "666333444", "maria@email.com")
	cliente3 := cm.CrearCliente("Pedro López", "666555666", "pedro@email.com")

	// Crear vehículos
	vehiculo1 := vm.CrearVehiculo("1234ABC", "Toyota", "Corolla", cliente1.ID)
	vehiculo2 := vm.CrearVehiculo("5678DEF", "Ford", "Focus", cliente2.ID)
	vehiculo3 := vm.CrearVehiculo("9012GHI", "Seat", "Ibiza", cliente3.ID)

	// Crear incidencias
	inc1 := im.CrearIncidencia(Mecanica, Alta, "Cambio de aceite y filtros")
	inc2 := im.CrearIncidencia(Electrica, Media, "Problema con luces")
	inc3 := im.CrearIncidencia(Carroceria, Baja, "Reparación de abolladuras")

	// Asignar incidencias a vehículos
	vm.AsignarIncidencia(vehiculo1.ID, inc1.ID)
	vm.AsignarIncidencia(vehiculo2.ID, inc2.ID)
	vm.AsignarIncidencia(vehiculo3.ID, inc3.ID)

	// Enviar vehículos al taller
	t.AgregarTrabajo(vehiculo1, inc1)
	time.Sleep(500 * time.Millisecond)
	t.AgregarTrabajo(vehiculo2, inc2)
	time.Sleep(500 * time.Millisecond)
	t.AgregarTrabajo(vehiculo3, inc3)

	fmt.Println("\nEsperando a que se completen los trabajos...")
	time.Sleep(25 * time.Second)

	fmt.Println("\nRESULTADOS DE LA SIMULACIÓN:")
	t.ObtenerEstadoTaller()

	fmt.Println("\n--- VEHÍCULOS PROCESADOS ---")
	vehiculos := vm.ListarVehiculos()
	for i := 0; i < len(vehiculos); i++ {
		fmt.Printf("Vehículo %s: Tiempo total de atención = %.1f segundos\n",
			vehiculos[i].Matricula, vehiculos[i].TiempoAcumulado)
	}
}
