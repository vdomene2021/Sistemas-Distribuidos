package main

import (
	"fmt"
	"time"
)

// Taller representa el taller con su sistema de concurrencia
type Taller struct {
	plazasTotales     int
	plazasOcupadas    int
	colaTrabajos      chan TrabajoPendiente
	mecanicoManager   *MecanicoManager
	vehiculoManager   *VehiculoManager
	incidenciaManager *IncidenciaManager
	terminar          chan bool
}

// NewTaller crea un nuevo taller
func NewTaller(plazas int, mm *MecanicoManager, vm *VehiculoManager, im *IncidenciaManager) *Taller {
	return &Taller{
		plazasTotales:     plazas,
		plazasOcupadas:    0,
		colaTrabajos:      make(chan TrabajoPendiente, 100), // Cola de espera ilimitada (buffered)
		mecanicoManager:   mm,
		vehiculoManager:   vm,
		incidenciaManager: im,
		terminar:          make(chan bool),
	}
}

// IniciarTaller inicia el sistema de procesamiento del taller
func (t *Taller) IniciarTaller() {
	fmt.Println("=== TALLER INICIADO ===")
	fmt.Println("Esperando trabajos...")

	go t.procesarTrabajos()
}

// procesarTrabajos procesa los trabajos de la cola
func (t *Taller) procesarTrabajos() {
	for {
		select {
		case trabajo := <-t.colaTrabajos:
			// Buscar un mecánico disponible con la especialidad adecuada
			especialidadRequerida := t.obtenerEspecialidadPorTipo(trabajo.Incidencia.Tipo)
			mecanicoAsignado := t.buscarMecanicoDisponible(especialidadRequerida)

			if mecanicoAsignado.ID == 0 {
				// No hay mecánicos disponibles, contratar uno nuevo
				fmt.Printf("No hay mecánicos disponibles de %s, contratando uno nuevo...\n", especialidadRequerida)
				nuevoMecanico := t.mecanicoManager.CrearMecanico(
					fmt.Sprintf("Mecanico-%d", t.mecanicoManager.nextID),
					especialidadRequerida,
					1,
				)
				mecanicoAsignado = nuevoMecanico
			}

			// Marcar el mecánico como ocupado
			t.mecanicoManager.CambiarEstadoOcupado(mecanicoAsignado.ID, true)

			// Asignar el mecánico a la incidencia
			t.incidenciaManager.AsignarMecanico(trabajo.Incidencia.ID, mecanicoAsignado.ID)

			// Lanzar goroutine para atender el vehículo
			go t.atenderVehiculo(trabajo, mecanicoAsignado)

		case <-t.terminar:
			fmt.Println("Taller cerrado")
			return
		}
	}
}

// atenderVehiculo atiende un vehículo (goroutine)
func (t *Taller) atenderVehiculo(trabajo TrabajoPendiente, mecanico Mecanico) {
	vehiculo := trabajo.Vehiculo
	incidencia := trabajo.Incidencia

	fmt.Printf("Mecánico %s (#%d) atendiendo vehículo %s (Matrícula: %s) - Incidencia: %s\n",
		mecanico.Nombre, mecanico.ID, vehiculo.Marca, vehiculo.Matricula, incidencia.Tipo)

	// Cambiar estado de la incidencia a "en proceso"
	t.incidenciaManager.CambiarEstado(incidencia.ID, EnProceso)

	// Obtener tiempo de atención
	tiempoAtencion := ObtenerTiempoAtencion(incidencia.Tipo)

	// Simular trabajo
	time.Sleep(tiempoAtencion)

	// Actualizar tiempo acumulado del vehículo
	tiempoSegundos := tiempoAtencion.Seconds()
	t.vehiculoManager.ActualizarTiempoAcumulado(vehiculo.ID, tiempoSegundos)

	// Verificar si el vehículo necesita otro mecánico (más de 15 segundos)
	tiempoAcumulado, existe := t.vehiculoManager.ObtenerTiempoAcumulado(vehiculo.ID)

	if existe && tiempoAcumulado > 15 {
		fmt.Printf("⚡ PRIORIDAD: Vehículo %s ha acumulado %.1f segundos. Asignando mecánico adicional...\n",
			vehiculo.Matricula, tiempoAcumulado)

		// Buscar otro mecánico disponible (de cualquier especialidad)
		mecanicoAdicional := t.buscarCualquierMecanicoDisponible()

		if mecanicoAdicional.ID == 0 {
			// No hay mecánicos disponibles, contratar uno nuevo
			especialidadRequerida := t.obtenerEspecialidadPorTipo(incidencia.Tipo)
			fmt.Printf("No hay mecánicos adicionales disponibles, contratando uno de %s...\n", especialidadRequerida)
			nuevoMecanico := t.mecanicoManager.CrearMecanico(
				fmt.Sprintf("Mecanico-Extra-%d", t.mecanicoManager.nextID),
				especialidadRequerida,
				1,
			)
			mecanicoAdicional = nuevoMecanico
		}

		// Asignar el mecánico adicional
		t.mecanicoManager.CambiarEstadoOcupado(mecanicoAdicional.ID, true)
		t.incidenciaManager.AsignarMecanico(incidencia.ID, mecanicoAdicional.ID)

		fmt.Printf("Mecánico adicional %s (#%d) asignado al vehículo %s\n",
			mecanicoAdicional.Nombre, mecanicoAdicional.ID, vehiculo.Matricula)

		// Trabajar más tiempo con el mecánico adicional
		time.Sleep(tiempoAtencion)
		t.vehiculoManager.ActualizarTiempoAcumulado(vehiculo.ID, tiempoSegundos)

		// Liberar el mecánico adicional
		t.mecanicoManager.CambiarEstadoOcupado(mecanicoAdicional.ID, false)
		fmt.Printf("Mecánico %s (#%d) ha terminado su parte del trabajo\n",
			mecanicoAdicional.Nombre, mecanicoAdicional.ID)
	}

	// Cambiar estado de la incidencia a "cerrada"
	t.incidenciaManager.CambiarEstado(incidencia.ID, Cerrada)

	// Liberar el mecánico
	t.mecanicoManager.CambiarEstadoOcupado(mecanico.ID, false)

	tiempoTotal, _ := t.vehiculoManager.ObtenerTiempoAcumulado(vehiculo.ID)
	fmt.Printf("Vehículo %s (Matrícula: %s) reparado. Tiempo total: %.1f segundos\n",
		vehiculo.Marca, vehiculo.Matricula, tiempoTotal)
}

// AgregarTrabajo agrega un trabajo a la cola
func (t *Taller) AgregarTrabajo(vehiculo Vehiculo, incidencia Incidencia) {
	trabajo := TrabajoPendiente{
		Vehiculo:     &vehiculo,
		Incidencia:   &incidencia,
		TiempoInicio: time.Now(),
	}

	fmt.Printf("Vehículo %s (Matrícula: %s) añadido a la cola - Incidencia: %s\n",
		vehiculo.Marca, vehiculo.Matricula, incidencia.Tipo)

	t.colaTrabajos <- trabajo
}

// buscarMecanicoDisponible busca un mecánico disponible con la especialidad requerida
func (t *Taller) buscarMecanicoDisponible(especialidad Especialidad) Mecanico {
	mecanicos := t.mecanicoManager.ListarMecanicos()

	for i := 0; i < len(mecanicos); i++ {
		if mecanicos[i].Activo && !mecanicos[i].Ocupado && mecanicos[i].Especialidad == especialidad {
			return mecanicos[i]
		}
	}

	return Mecanico{} // Retorna mecánico vacío si no encuentra
}

// buscarCualquierMecanicoDisponible busca cualquier mecánico disponible
func (t *Taller) buscarCualquierMecanicoDisponible() Mecanico {
	mecanicos := t.mecanicoManager.ListarMecanicos()

	for i := 0; i < len(mecanicos); i++ {
		if mecanicos[i].Activo && !mecanicos[i].Ocupado {
			return mecanicos[i]
		}
	}

	return Mecanico{} // Retorna mecánico vacío si no encuentra
}

// obtenerEspecialidadPorTipo convierte el tipo de incidencia en especialidad
func (t *Taller) obtenerEspecialidadPorTipo(tipo TipoIncidencia) Especialidad {
	if tipo == Mecanica {
		return EspecialidadMecanica
	}
	if tipo == Electrica {
		return EspecialidadElectrica
	}
	if tipo == Carroceria {
		return EspecialidadCarroceria
	}
	return EspecialidadMecanica
}

// DetenerTaller detiene el taller
func (t *Taller) DetenerTaller() {
	t.terminar <- true
}

// ObtenerEstadoTaller muestra el estado actual del taller
func (t *Taller) ObtenerEstadoTaller() {
	fmt.Println("\n=== ESTADO DEL TALLER ===")
	fmt.Printf("Plazas totales: %d\n", t.plazasTotales)
	fmt.Printf("Trabajos en cola: %d\n", len(t.colaTrabajos))

	mecanicosDisponibles := t.mecanicoManager.ListarMecanicosDisponibles()
	fmt.Printf("Mecánicos disponibles: %d\n", len(mecanicosDisponibles))

	conteo := t.mecanicoManager.ContarMecanicosPorEspecialidad()
	fmt.Printf("  - Mecánica: %d\n", conteo[EspecialidadMecanica])
	fmt.Printf("  - Eléctrica: %d\n", conteo[EspecialidadElectrica])
	fmt.Printf("  - Carrocería: %d\n", conteo[EspecialidadCarroceria])
	fmt.Println("========================")
}
