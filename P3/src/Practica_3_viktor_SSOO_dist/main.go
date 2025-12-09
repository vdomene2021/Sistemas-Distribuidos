package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func (cm *ClienteManager) CrearCliente(nombre, telefono, email string) *Cliente {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cliente := &Cliente{
		ID:       cm.nextID,
		Nombre:   nombre,
		Telefono: telefono,
		Email:    email,
	}
	cm.clientes[cm.nextID] = cliente
	cm.nextID++
	return cliente
}

func (cm *ClienteManager) ListarClientes() []*Cliente {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	clientes := make([]*Cliente, 0, len(cm.clientes))
	for _, c := range cm.clientes {
		clientes = append(clientes, c)
	}
	return clientes
}

func (cm *ClienteManager) ActualizarCliente(id int, nombre, telefono, email string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cliente, existe := cm.clientes[id]
	if !existe {
		return fmt.Errorf("cliente no encontrado")
	}

	if nombre != "" {
		cliente.Nombre = nombre
	}
	if telefono != "" {
		cliente.Telefono = telefono
	}
	if email != "" {
		cliente.Email = email
	}
	return nil
}

func (cm *ClienteManager) EliminarCliente(id int) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, existe := cm.clientes[id]; !existe {
		return fmt.Errorf("cliente no encontrado")
	}
	delete(cm.clientes, id)
	return nil
}

func (cm *ClienteManager) ExisteCliente(id int) bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	_, existe := cm.clientes[id]
	return existe
}

func (vm *VehiculoManager) CrearVehiculo(matricula, marca, modelo string, clienteID int, cm *ClienteManager) (*VehiculoCompleto, error) {
	if !cm.ExisteCliente(clienteID) {
		return nil, fmt.Errorf("el cliente con ID %d no existe", clienteID)
	}

	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	vehiculo := &VehiculoCompleto{
		ID:        vm.nextID,
		Matricula: matricula,
		Marca:     marca,
		Modelo:    modelo,
		ClienteID: clienteID,
	}
	vm.vehiculos[vm.nextID] = vehiculo
	vm.nextID++
	return vehiculo, nil
}

func (vm *VehiculoManager) ListarVehiculos() []*VehiculoCompleto {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	vehiculos := make([]*VehiculoCompleto, 0, len(vm.vehiculos))
	for _, v := range vm.vehiculos {
		vehiculos = append(vehiculos, v)
	}
	return vehiculos
}

func (vm *VehiculoManager) ListarVehiculosPorCliente(clienteID int) []*VehiculoCompleto {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	vehiculos := make([]*VehiculoCompleto, 0)
	for _, v := range vm.vehiculos {
		if v.ClienteID == clienteID {
			vehiculos = append(vehiculos, v)
		}
	}
	return vehiculos
}

func (vm *VehiculoManager) ObtenerVehiculo(id int) (*VehiculoCompleto, bool) {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()
	v, existe := vm.vehiculos[id]
	return v, existe
}

func (vm *VehiculoManager) ActualizarVehiculo(id int, matricula, marca, modelo string) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	vehiculo, existe := vm.vehiculos[id]
	if !existe {
		return fmt.Errorf("vehículo no encontrado")
	}

	if matricula != "" {
		vehiculo.Matricula = matricula
	}
	if marca != "" {
		vehiculo.Marca = marca
	}
	if modelo != "" {
		vehiculo.Modelo = modelo
	}
	return nil
}

func (vm *VehiculoManager) EliminarVehiculo(id int) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	if _, existe := vm.vehiculos[id]; !existe {
		return fmt.Errorf("vehículo no encontrado")
	}
	delete(vm.vehiculos, id)
	return nil
}

func (vm *VehiculoManager) AgregarTiempo(id int, tiempo float64) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	if v, existe := vm.vehiculos[id]; existe {
		v.TiempoAcumulado += tiempo
	}
}

func (im *IncidenciaManager) CrearIncidencia(tipo TipoIncidencia, prioridad Prioridad, descripcion string, vehiculoID int) *IncidenciaCompleta {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	incidencia := &IncidenciaCompleta{
		ID:          im.nextID,
		VehiculoID:  vehiculoID,
		Tipo:        tipo,
		Prioridad:   prioridad,
		Estado:      Abierta,
		Descripcion: descripcion,
	}
	im.incidencias[im.nextID] = incidencia
	im.nextID++
	return incidencia
}

func (im *IncidenciaManager) ListarIncidencias() []*IncidenciaCompleta {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	incidencias := make([]*IncidenciaCompleta, 0, len(im.incidencias))
	for _, inc := range im.incidencias {
		incidencias = append(incidencias, inc)
	}
	return incidencias
}

func (im *IncidenciaManager) ObtenerIncidencia(id int) (*IncidenciaCompleta, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()
	inc, existe := im.incidencias[id]
	return inc, existe
}

func (im *IncidenciaManager) ObtenerIncidenciasPorVehiculo(vehiculoID int) []*IncidenciaCompleta {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	incidencias := make([]*IncidenciaCompleta, 0)
	for _, inc := range im.incidencias {
		if inc.VehiculoID == vehiculoID && inc.Estado == Abierta {
			incidencias = append(incidencias, inc)
		}
	}
	return incidencias
}

func (im *IncidenciaManager) CambiarEstado(id int, estado EstadoIncidencia) error {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	incidencia, existe := im.incidencias[id]
	if !existe {
		return fmt.Errorf("incidencia no encontrada")
	}
	incidencia.Estado = estado
	return nil
}

func (im *IncidenciaManager) EliminarIncidencia(id int) error {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	if _, existe := im.incidencias[id]; !existe {
		return fmt.Errorf("incidencia no encontrada")
	}
	delete(im.incidencias, id)
	return nil
}

func (im *IncidenciaManager) EliminarIncidenciasPorVehiculo(vehiculoID int) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	for id, inc := range im.incidencias {
		if inc.VehiculoID == vehiculoID {
			delete(im.incidencias, id)
		}
	}
}

func (im *IncidenciaManager) ContarTodasIncidencias() map[int]int {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	counts := make(map[int]int)
	for _, inc := range im.incidencias {
		counts[inc.VehiculoID]++
	}
	return counts
}

func (mm *MecanicoManager) CrearMecanico(nombre string, especialidad Especialidad, experiencia int) Mecanico {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	mecanico := Mecanico{
		ID:           mm.nextID,
		Nombre:       nombre,
		Especialidad: especialidad,
		Experiencia:  experiencia,
		Activo:       true,
		ColaPersonal: make([]*TrabajoMecanico, 0, 2),
		Canal:        make(chan *TrabajoMecanico, 2),
	}
	mm.mecanicos[mm.nextID] = &mecanico
	mm.nextID++
	return mecanico
}

func (mm *MecanicoManager) ListarMecanicos() []*Mecanico {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	mecanicos := make([]*Mecanico, 0, len(mm.mecanicos))
	for _, m := range mm.mecanicos {
		mecanicos = append(mecanicos, m)
	}
	return mecanicos
}

func (mm *MecanicoManager) ObtenerMecanico(id int) (*Mecanico, bool) {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()
	m, existe := mm.mecanicos[id]
	return m, existe
}

func (mm *MecanicoManager) CambiarEstadoActivo(id int, activo bool) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	mecanico, existe := mm.mecanicos[id]
	if !existe {
		return fmt.Errorf("mecánico no encontrado")
	}
	mecanico.Activo = activo
	return nil
}

func (t *Taller) AgregarTrabajo(vehiculo *VehiculoCompleto, incidencia *IncidenciaCompleta) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	trabajo := &TrabajoMecanico{
		Vehiculo:   *vehiculo,
		Incidencia: *incidencia,
	}
	t.ColaTrabajo = append(t.ColaTrabajo, trabajo)
	fmt.Printf("Trabajo agregado: Vehículo %s - Incidencia %s\n", vehiculo.Matricula, incidencia.Tipo)

	// Intentar asignar inmediatamente
	go t.AsignarTrabajosAutomaticamente()
}

func (t *Taller) AsignarTrabajosAutomaticamente() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if len(t.ColaTrabajo) == 0 {
		return
	}

	mecanicos := t.MecanicoManager.ListarMecanicos()

	for i := 0; i < len(t.ColaTrabajo); i++ {
		trabajo := t.ColaTrabajo[i]
		asignado := false

		// Buscar mecánico disponible de la especialidad adecuada
		for _, mec := range mecanicos {
			if !mec.Activo {
				continue
			}

			// Verificar especialidad
			especialidadCorrecta := false
			switch trabajo.Incidencia.Tipo {
			case Mecanica:
				especialidadCorrecta = (mec.Especialidad == EspMecanica)
			case Electrica:
				especialidadCorrecta = (mec.Especialidad == EspElectrica)
			case Carroceria:
				especialidadCorrecta = (mec.Especialidad == EspCarroceria)
			}

			if !especialidadCorrecta {
				continue
			}

			// Verificar si tiene espacio en su cola
			mec.mutex.Lock()
			if len(mec.ColaPersonal) < 2 {
				mec.ColaPersonal = append(mec.ColaPersonal, trabajo)
				fmt.Printf("Asignado a %s (especialidad: %s)\n", mec.Nombre, mec.Especialidad)

				// Enviar por canal
				select {
				case mec.Canal <- trabajo:
				default:
				}

				asignado = true
				mec.mutex.Unlock()

				// Eliminar de la cola principal
				t.ColaTrabajo = append(t.ColaTrabajo[:i], t.ColaTrabajo[i+1:]...)
				i--
				break
			}
			mec.mutex.Unlock()
		}

		if !asignado {
			break // No hay mecánicos disponibles, dejamos el resto en cola
		}
	}
}

func (t *Taller) ArrancarRutinaMecanico(m *Mecanico) {
	fmt.Printf("Mecánico %s iniciado y listo para trabajar\n", m.Nombre)

	for {
		if !m.Activo {
			time.Sleep(1 * time.Second)
			continue
		}

		select {
		case trabajo := <-m.Canal:
			fmt.Printf("%s comienza a trabajar en %s (%s)\n",
				m.Nombre, trabajo.Vehiculo.Matricula, trabajo.Incidencia.Tipo)

			// Simular trabajo (3-8 segundos según experiencia)
			tiempoTrabajo := time.Duration(10-m.Experiencia) * time.Second
			if tiempoTrabajo < 3*time.Second {
				tiempoTrabajo = 3 * time.Second
			}
			time.Sleep(tiempoTrabajo)

			// Actualizar tiempo acumulado del vehículo
			t.VehiculoManager.AgregarTiempo(trabajo.Vehiculo.ID, tiempoTrabajo.Seconds())

			// Cambiar estado de incidencia a cerrada
			t.IncidenciaManager.CambiarEstado(trabajo.Incidencia.ID, Cerrada)

			// Quitar de cola personal
			m.mutex.Lock()
			for i, tr := range m.ColaPersonal {
				if tr == trabajo {
					m.ColaPersonal = append(m.ColaPersonal[:i], m.ColaPersonal[i+1:]...)
					break
				}
			}
			m.mutex.Unlock()

			fmt.Printf("%s terminó trabajo en %s\n", m.Nombre, trabajo.Vehiculo.Matricula)

			// Intentar asignar más trabajos
			go t.AsignarTrabajosAutomaticamente()

		case <-time.After(2 * time.Second):
			// Timeout, volver a verificar
			continue
		}
	}
}

func (t *Taller) ObtenerEstadoTaller() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	fmt.Println("\n╔════════════════════════════════════════════════════╗")
	fmt.Println("║           ESTADO ACTUAL DEL TALLER                 ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")

	fmt.Printf("\nTrabajos en cola principal: %d\n", len(t.ColaTrabajo))
	if len(t.ColaTrabajo) > 0 {
		for i, trabajo := range t.ColaTrabajo {
			fmt.Printf("  %d. Vehículo: %s | Incidencia: %s | Prioridad: %d\n",
				i+1, trabajo.Vehiculo.Matricula, trabajo.Incidencia.Tipo, trabajo.Incidencia.Prioridad)
		}
	}

	fmt.Println("\nEstado de mecánicos:")
	mecanicos := t.MecanicoManager.ListarMecanicos()
	for _, mec := range mecanicos {
		estado := "Activo"
		if !mec.Activo {
			estado = "Inactivo"
		}
		fmt.Printf("  • %s (%s) - %s - Cola: %d/2\n",
			mec.Nombre, mec.Especialidad, estado, len(mec.ColaPersonal))

		if len(mec.ColaPersonal) > 0 {
			for j, tr := range mec.ColaPersonal {
				fmt.Printf("      %d. %s (%s)\n", j+1, tr.Vehiculo.Matricula, tr.Incidencia.Tipo)
			}
		}
	}

	fmt.Println("\n════════════════════════════════════════════════════════")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Inicializar managers
	clienteManager := NewClienteManager()
	vehiculoManager := NewVehiculoManager()
	incidenciaManager := NewIncidenciaManager()
	mecanicoManager := NewMecanicoManager()
	taller := NewTaller(mecanicoManager, vehiculoManager, incidenciaManager)

	fmt.Println("╔═══════════════════════════════════════════════════╗")
	fmt.Println("║     SISTEMA DE GESTIÓN DE TALLER MECÁNICO         ║")
	fmt.Println("╚═══════════════════════════════════════════════════╝")

	for {
		mostrarMenuPrincipal()
		fmt.Print("Seleccione una opción: ")
		scanner.Scan()
		opcion := scanner.Text()

		switch opcion {
		case "1":
			menuClientes(scanner, clienteManager, vehiculoManager)
		case "2":
			menuVehiculos(scanner, vehiculoManager, clienteManager, incidenciaManager)
		case "3":
			menuIncidencias(scanner, incidenciaManager, vehiculoManager)
		case "4":
			menuMecanicos(scanner, mecanicoManager, taller)
		case "5":
			menuTaller(scanner, taller, vehiculoManager, incidenciaManager)
		case "0":
			fmt.Println("\nGracias por usar el sistema. ¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
		}
	}
}
