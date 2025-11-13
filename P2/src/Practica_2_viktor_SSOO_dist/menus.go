package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func mostrarMenuPrincipal() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║            TALLER MECÁNICO             ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println("1. Gestión de Clientes")
	fmt.Println("2. Gestión de Vehículos")
	fmt.Println("3. Gestión de Incidencias")
	fmt.Println("4. Gestión de Mecánicos")
	fmt.Println("5. Estado del Taller")
	fmt.Println("0. Salir")
	fmt.Println("──────────────────────────────────────────")
}

func menuClientes(scanner *bufio.Scanner, cm *ClienteManager, vm *VehiculoManager) {
	for {
		fmt.Println("\n===== GESTIÓN DE CLIENTES =====")
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
			fmt.Println("\n----- LISTA DE CLIENTES -----")
			for i := 0; i < len(clientes); i++ {
				fmt.Printf("ID: %d | Nombre: %s | Tel: %s | Email: %s\n",
					clientes[i].ID, clientes[i].Nombre, clientes[i].Telefono, clientes[i].Email)

				vehiculosCliente := vm.ListarVehiculosPorCliente(clientes[i].ID)
				if len(vehiculosCliente) > 0 {
					fmt.Println("    Vehículos:")
					for _, v := range vehiculosCliente {
						fmt.Printf("      - ID: %d | Matrícula: %s (%s %s)\n", v.ID, v.Matricula, v.Marca, v.Modelo)
					}
				} else {
					fmt.Println("    (Sin vehículos registrados)")
				}
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
		fmt.Println("\n===== GESTIÓN DE VEHÍCULOS =====")
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

			vehiculo, err := vm.CrearVehiculo(matricula, marca, modelo, clienteID, cm)

			if err != nil {
				fmt.Printf("Error al crear vehículo: %v\n", err)
			} else {
				fmt.Printf("Vehículo creado con ID: %d\n", vehiculo.ID)
			}

		} else if opcion == "2" {
			vehiculos := vm.ListarVehiculos()
			incidenciaCounts := im.ContarTodasIncidencias()

			fmt.Println("\n----- LISTA DE VEHÍCULOS -----")
			for i := 0; i < len(vehiculos); i++ {
				count := incidenciaCounts[vehiculos[i].ID]
				fmt.Printf("ID: %d | Matrícula: %s | Marca: %s | Modelo: %s | Incidencias: %d | Tiempo acumulado: %.1fs\n",
					vehiculos[i].ID, vehiculos[i].Matricula, vehiculos[i].Marca, vehiculos[i].Modelo, count, vehiculos[i].TiempoAcumulado)
			}

		} else if opcion == "3" {
			fmt.Print("ID del vehículo a actualizar: ")
			scanner.Scan()
			id, _ := strconv.Atoi(scanner.Text())

			if _, existe := vm.ObtenerVehiculo(id); !existe {
				fmt.Printf("Error: Vehículo con ID %d no encontrado.\n", id)
				continue
			}

			fmt.Print("Nueva matrícula (vacío para no cambiar): ")
			scanner.Scan()
			matricula := scanner.Text()

			fmt.Print("Nueva marca (vacío para no cambiar): ")
			scanner.Scan()
			marca := scanner.Text()

			fmt.Print("Nuevo modelo (vacío para no cambiar): ")
			scanner.Scan()
			modelo := scanner.Text()

			err := vm.ActualizarVehiculo(id, matricula, marca, modelo)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Vehículo actualizado")
			}

		} else if opcion == "4" {
			fmt.Print("ID del vehículo a eliminar: ")
			scanner.Scan()
			vehiculoID, _ := strconv.Atoi(scanner.Text())

			if _, existe := vm.ObtenerVehiculo(vehiculoID); !existe {
				fmt.Printf("Error: Vehículo con ID %d no encontrado.\n", vehiculoID)
				continue
			}

			im.EliminarIncidenciasPorVehiculo(vehiculoID)
			fmt.Printf("Incidencias asociadas al vehículo %d eliminadas.\n", vehiculoID)

			err := vm.EliminarVehiculo(vehiculoID)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Vehículo con ID %d eliminado correctamente.\n", vehiculoID)
			}

		} else if opcion == "0" {
			break
		}
	}
}

func menuIncidencias(scanner *bufio.Scanner, im *IncidenciaManager, vm *VehiculoManager) {
	for {
		fmt.Println("\n===== GESTIÓN DE INCIDENCIAS =====")
		fmt.Println("1. Crear Incidencia")
		fmt.Println("2. Listar Incidencias")
		fmt.Println("3. Cambiar Estado")
		fmt.Println("4. Eliminar Incidencia")
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

			fmt.Print("ID del vehículo asociado: ")
			scanner.Scan()
			vehiculoID, _ := strconv.Atoi(scanner.Text())

			if _, existe := vm.ObtenerVehiculo(vehiculoID); !existe {
				fmt.Printf("ERROR: El vehículo con ID %d no existe. No se puede crear la incidencia.\n", vehiculoID)
				continue
			}

			incidencia := im.CrearIncidencia(tipo, prioridad, descripcion, vehiculoID)
			fmt.Printf("Incidencia creada con ID: %d para Vehículo ID: %d\n", incidencia.ID, incidencia.VehiculoID)

		} else if opcion == "2" {
			incidencias := im.ListarIncidencias()
			fmt.Println("\n----- LISTA DE INCIDENCIAS -----")
			for i := 0; i < len(incidencias); i++ {
				fmt.Printf("ID: %d | VehículoID: %d | Tipo: %s | Prioridad: %s | Estado: %s | Descripción: %s\n",
					incidencias[i].ID, incidencias[i].VehiculoID, incidencias[i].Tipo, incidencias[i].Prioridad,
					incidencias[i].Estado, incidencias[i].Descripcion)
			}

		} else if opcion == "3" {
			fmt.Print("ID de la incidencia: ")
			scanner.Scan()
			id, _ := strconv.Atoi(scanner.Text())

			if _, existe := im.ObtenerIncidencia(id); !existe {
				fmt.Printf("ERROR: Incidencia con ID %d no encontrada.\n", id)
				continue
			}

			fmt.Print("Nuevo estado (abierta/en proceso/cerrada): ")
			scanner.Scan()
			estadoStr := scanner.Text()
			estado := EstadoIncidencia(estadoStr)

			err := im.CambiarEstado(id, estado)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Estado actualizado")
			}

		} else if opcion == "4" {
			fmt.Print("ID de la incidencia a eliminar: ")
			scanner.Scan()
			id, _ := strconv.Atoi(scanner.Text())

			if _, existe := im.ObtenerIncidencia(id); !existe {
				fmt.Printf("ERROR: Incidencia con ID %d no encontrada.\n", id)
				continue
			}

			err := im.EliminarIncidencia(id)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Incidencia eliminada correctamente")
			}

		} else if opcion == "0" {
			break
		}
	}
}

func menuMecanicos(scanner *bufio.Scanner, mm *MecanicoManager, taller *Taller) {
	for {
		fmt.Println("\n===== GESTIÓN DE MECÁNICOS =====")
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

			// Arrancamos la rutina pasando la copia local
			go taller.ArrancarRutinaMecanico(&mecanico)

		} else if opcion == "2" {
			mecanicos := mm.ListarMecanicos()
			fmt.Println("\n----- LISTA DE MECÁNICOS -----")
			for i := 0; i < len(mecanicos); i++ {
				estado := "Activo"
				if !mecanicos[i].Activo {
					estado = "Inactivo"
				}
				fmt.Printf("ID: %d | Nombre: %s | Especialidad: %s | Exp: %d años | Estado: %s | Cola: %d/2\n",
					mecanicos[i].ID, mecanicos[i].Nombre, mecanicos[i].Especialidad,
					mecanicos[i].Experiencia, estado, len(mecanicos[i].ColaPersonal))
			}

		} else if opcion == "3" {
			fmt.Print("ID del mecánico: ")
			scanner.Scan()
			id, _ := strconv.Atoi(scanner.Text())

			mecanicoActual, existe := mm.ObtenerMecanico(id)
			if !existe {
				fmt.Printf("Error: Mecánico con ID %d no encontrado.\n", id)
				continue
			}

			fmt.Println("Seleccione una acción:")
			fmt.Println("1. Dar de ALTA")
			fmt.Println("2. Dar de BAJA")
			fmt.Print("Opción: ")
			scanner.Scan()
			subOpcion := scanner.Text()

			var err error
			var accion string

			if subOpcion == "1" {
				if mecanicoActual.Activo {
					fmt.Printf("AVISO: El mecánico %s ya está dado de ALTA.\n", mecanicoActual.Nombre)
					continue
				}
				err = mm.CambiarEstadoActivo(id, true)
				accion = "dado de alta"

			} else if subOpcion == "2" {
				if !mecanicoActual.Activo {
					fmt.Printf("AVISO: El mecánico %s ya está dado de BAJA.\n", mecanicoActual.Nombre)
					continue
				}
				if len(mecanicoActual.ColaPersonal) > 0 {
					fmt.Printf("ERROR: No se puede dar de baja a %s porque tiene %d coches en cola/proceso.\n",
						mecanicoActual.Nombre, len(mecanicoActual.ColaPersonal))
					continue
				}
				err = mm.CambiarEstadoActivo(id, false)
				accion = "dado de baja"

			} else {
				fmt.Println("Opción no válida. Operación cancelada.")
				continue
			}

			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Mecánico con ID %d ha sido %s correctamente.\n", id, accion)
			}

		} else if opcion == "0" {
			break
		}
	}
}

func menuTaller(scanner *bufio.Scanner, t *Taller, vm *VehiculoManager, im *IncidenciaManager) {
	for {
		fmt.Println("\n===== ESTADO DEL TALLER =====")
		fmt.Println("1. Ver Estado")
		fmt.Println("2. Enviar Vehículo(s) al Taller")
		fmt.Println("0. Volver")
		fmt.Print("Opción: ")
		scanner.Scan()
		opcion := scanner.Text()

		if opcion == "1" {
			t.ObtenerEstadoTaller()
		} else if opcion == "2" {
			fmt.Print("IDs de los vehículos a enviar (separados por coma): ")
			scanner.Scan()
			idsInput := scanner.Text()

			idsStr := strings.Split(idsInput, ",")

			for _, idStr := range idsStr {
				vehiculoID, err := strconv.Atoi(strings.TrimSpace(idStr))
				if err != nil {
					fmt.Printf("Error: '%s' no es un ID válido. Omitiendo.\n", idStr)
					continue
				}

				vehiculo, existe := vm.ObtenerVehiculo(vehiculoID)
				if !existe {
					fmt.Printf("Vehículo con ID %d no encontrado. Omitiendo.\n", vehiculoID)
					continue
				}

				incidencias := im.ObtenerIncidenciasPorVehiculo(vehiculoID)
				if len(incidencias) == 0 {
					fmt.Printf("Vehículo %s (ID %d) no tiene incidencias abiertas. Omitiendo.\n", vehiculo.Matricula, vehiculoID)
					continue
				}

				fmt.Printf("Enviando %d incidencia(s) del vehículo %s (ID %d) a la cola del taller...\n", len(incidencias), vehiculo.Matricula, vehiculo.ID)

				for _, inc := range incidencias {
					incCopia := inc
					t.AgregarTrabajo(vehiculo, incCopia)
				}
			}

		} else if opcion == "0" {
			break
		}
	}
}
