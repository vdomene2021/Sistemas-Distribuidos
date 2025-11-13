package main

import (
	"fmt"
	"sync"
	"time"
)

type Taller struct {
	plazasOcupadas    int
	colaTrabajos      chan TrabajoPendiente
	mecanicoManager   *MecanicoManager
	vehiculoManager   *VehiculoManager
	incidenciaManager *IncidenciaManager
	terminar          chan bool
	wg                *sync.WaitGroup
}

func NewTaller(mm *MecanicoManager, vm *VehiculoManager, im *IncidenciaManager) *Taller {
	return &Taller{
		plazasOcupadas:    0,
		colaTrabajos:      make(chan TrabajoPendiente, 100),
		mecanicoManager:   mm,
		vehiculoManager:   vm,
		incidenciaManager: im,
		terminar:          make(chan bool),
	}
}

func (t *Taller) IniciarTaller() {
	fmt.Println("=== TALLER INICIADO ===")
	fmt.Println("Esperando trabajos...")

	mecanicos := t.mecanicoManager.ListarMecanicos()
	for i := 0; i < len(mecanicos); i++ {
		go t.ArrancarRutinaMecanico(&mecanicos[i])
	}

	go t.procesarTrabajos()
}

func (t *Taller) ArrancarRutinaMecanico(m *Mecanico) {
	for trabajo := range m.ColaPersonal {
		t.atenderVehiculo(trabajo, *m)
	}
}

func (t *Taller) procesarTrabajos() {
	for {
		select {
		case trabajo := <-t.colaTrabajos:
			especialidadRequerida := t.obtenerEspecialidadPorTipo(trabajo.Incidencia.Tipo)

			mecanicoAsignado := t.buscarMecanicoConHueco(especialidadRequerida)

			if mecanicoAsignado == nil {
				time.Sleep(50 * time.Millisecond)
				t.colaTrabajos <- trabajo
				continue
			}

			t.mecanicoManager.IncrementarPlaza(mecanicoAsignado.ID)

			t.incidenciaManager.AsignarMecanico(trabajo.Incidencia.ID, mecanicoAsignado.ID)

			mecanicoAsignado.ColaPersonal <- trabajo

			fmt.Printf("-> Coche asignado a %s (Ocupación: %d/2)\n",
				mecanicoAsignado.Nombre, mecanicoAsignado.PlazasOcupadas)

		case <-t.terminar:
			fmt.Println("Taller cerrado")
			return
		}
	}
}

func (t *Taller) atenderVehiculo(trabajo TrabajoPendiente, mecanico Mecanico) {
	vehiculo := trabajo.Vehiculo
	incidencia := trabajo.Incidencia

	fmt.Printf("Mecánico %s (#%d) COMENZANDO a reparar vehículo %s\n",
		mecanico.Nombre, mecanico.ID, vehiculo.Marca)

	t.incidenciaManager.CambiarEstado(incidencia.ID, EnProceso)

	tiempoAtencion := ObtenerTiempoAtencion(incidencia.Tipo)
	time.Sleep(tiempoAtencion)

	tiempoSegundos := tiempoAtencion.Seconds()
	t.vehiculoManager.ActualizarTiempoAcumulado(vehiculo.ID, tiempoSegundos)

	tiempoAcumulado, existe := t.vehiculoManager.ObtenerTiempoAcumulado(vehiculo.ID)

	if existe && tiempoAcumulado > 15 {
		fmt.Printf("PRIORIDAD: Vehículo %s ha acumulado %.1f segundos. Asignando mecánico adicional...\n",
			vehiculo.Matricula, tiempoAcumulado)

		mecanicoAdicional := t.buscarCualquierMecanicoConHueco()

		if mecanicoAdicional == nil {
			especialidadRequerida := t.obtenerEspecialidadPorTipo(incidencia.Tipo)
			fmt.Printf("No hay mecánicos adicionales disponibles, contratando uno de %s...\n", especialidadRequerida)

			nuevoMecanico := t.mecanicoManager.CrearMecanico(
				fmt.Sprintf("Mecanico-Extra-%d", t.mecanicoManager.nextID),
				especialidadRequerida,
				1,
			)

			go t.ArrancarRutinaMecanico(&nuevoMecanico)

			mecanicoAdicional = &nuevoMecanico
		}

		t.mecanicoManager.IncrementarPlaza(mecanicoAdicional.ID)
		t.incidenciaManager.AsignarMecanico(incidencia.ID, mecanicoAdicional.ID)

		fmt.Printf("Mecánico adicional %s (#%d) asignado al vehículo %s\n",
			mecanicoAdicional.Nombre, mecanicoAdicional.ID, vehiculo.Matricula)

		time.Sleep(tiempoAtencion)
		t.vehiculoManager.ActualizarTiempoAcumulado(vehiculo.ID, tiempoSegundos)

		t.mecanicoManager.DecrementarPlaza(mecanicoAdicional.ID)
		fmt.Printf("Mecánico %s (#%d) ha terminado su parte del trabajo\n",
			mecanicoAdicional.Nombre, mecanicoAdicional.ID)
	}

	t.incidenciaManager.CambiarEstado(incidencia.ID, Cerrada)
	t.incidenciaManager.EliminarIncidencia(incidencia.ID)

	t.mecanicoManager.DecrementarPlaza(mecanico.ID)

	if t.wg != nil {
		t.wg.Done()
	}

	tiempoTotal, _ := t.vehiculoManager.ObtenerTiempoAcumulado(vehiculo.ID)
	fmt.Printf("Vehículo %s REPARADO por %s. Tiempo total: %.1fs (Plaza liberada)\n",
		vehiculo.Marca, mecanico.Nombre, tiempoTotal)
}

func (t *Taller) AgregarTrabajo(vehiculo Vehiculo, incidencia Incidencia) {
	trabajo := TrabajoPendiente{
		Vehiculo:     &vehiculo,
		Incidencia:   &incidencia,
		TiempoInicio: time.Now(),
	}

	if t.wg != nil {
		t.wg.Add(1)
	}

	fmt.Printf("Vehículo %s (Matrícula: %s) añadido a la cola GENERAL\n",
		vehiculo.Marca, vehiculo.Matricula)

	t.colaTrabajos <- trabajo
}

func (t *Taller) buscarMecanicoConHueco(especialidad Especialidad) *Mecanico {
	mecanicos := t.mecanicoManager.ListarMecanicos()

	for i := 0; i < len(mecanicos); i++ {
		m := &mecanicos[i]
		if m.Activo && m.Especialidad == especialidad && m.PlazasOcupadas < 2 {
			return m
		}
	}
	return nil
}

func (t *Taller) buscarCualquierMecanicoConHueco() *Mecanico {
	mecanicos := t.mecanicoManager.ListarMecanicos()
	for i := 0; i < len(mecanicos); i++ {
		m := &mecanicos[i]
		if m.Activo && m.PlazasOcupadas < 2 {
			return m
		}
	}
	return nil
}

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

func (t *Taller) DetenerTaller() {
	t.terminar <- true
}

func (t *Taller) ObtenerEstadoTaller() {
	fmt.Println("\n===== ESTADO DEL TALLER =====")

	mecanicosActivos := t.mecanicoManager.ContarMecanicosActivos()
	plazasTotales := mecanicosActivos * 2

	fmt.Printf("Plazas totales: %d (%d mecánicos x 2)\n", plazasTotales, mecanicosActivos)
	fmt.Printf("Trabajos en cola GENERAL: %d\n", len(t.colaTrabajos))

	mecanicos := t.mecanicoManager.ListarMecanicos()
	fmt.Println("----- Carga de Mecánicos -----")
	for _, m := range mecanicos {
		if m.Activo {
			fmt.Printf("  - %s (%s): %d/2 Plazas Ocupadas\n\n", m.Nombre, m.Especialidad, m.PlazasOcupadas)
		}
	}
	fmt.Println("=============================")
}
