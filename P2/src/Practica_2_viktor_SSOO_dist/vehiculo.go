package main

import (
	"fmt"
	"time"
)

// VehiculoManager gestiona las operaciones sobre vehículos
type VehiculoManager struct {
	vehiculos []Vehiculo
	nextID    int
}

// NewVehiculoManager crea un nuevo gestor de vehículos
func NewVehiculoManager() *VehiculoManager {
	return &VehiculoManager{
		vehiculos: make([]Vehiculo, 0),
		nextID:    1,
	}
}

// CrearVehiculo crea un nuevo vehículo
func (vm *VehiculoManager) CrearVehiculo(matricula, marca, modelo string, clienteID int) Vehiculo {
	vehiculo := Vehiculo{
		ID:              vm.nextID,
		Matricula:       matricula,
		Marca:           marca,
		Modelo:          modelo,
		FechaEntrada:    time.Now(),
		ClienteID:       clienteID,
		TiempoAcumulado: 0,
	}
	vm.vehiculos = append(vm.vehiculos, vehiculo)
	vm.nextID++
	return vehiculo
}

// ObtenerVehiculo obtiene un vehículo por su ID
func (vm *VehiculoManager) ObtenerVehiculo(id int) (Vehiculo, bool) {
	for i := 0; i < len(vm.vehiculos); i++ {
		if vm.vehiculos[i].ID == id {
			return vm.vehiculos[i], true
		}
	}
	return Vehiculo{}, false
}

// ActualizarVehiculo actualiza los datos de un vehículo
func (vm *VehiculoManager) ActualizarVehiculo(id int, matricula, marca, modelo string) error {
	for i := 0; i < len(vm.vehiculos); i++ {
		if vm.vehiculos[i].ID == id {
			if matricula != "" {
				vm.vehiculos[i].Matricula = matricula
			}
			if marca != "" {
				vm.vehiculos[i].Marca = marca
			}
			if modelo != "" {
				vm.vehiculos[i].Modelo = modelo
			}
			return nil
		}
	}
	return fmt.Errorf("vehículo con ID %d no encontrado", id)
}

// EliminarVehiculo elimina un vehículo
func (vm *VehiculoManager) EliminarVehiculo(id int) error {
	for i := 0; i < len(vm.vehiculos); i++ {
		if vm.vehiculos[i].ID == id {
			vm.vehiculos = append(vm.vehiculos[:i], vm.vehiculos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("vehículo con ID %d no encontrado", id)
}

// ListarVehiculos lista todos los vehículos
func (vm *VehiculoManager) ListarVehiculos() []Vehiculo {
	return vm.vehiculos
}

// AsignarIncidencia asigna una incidencia a un vehículo
func (vm *VehiculoManager) AsignarIncidencia(vehiculoID, incidenciaID int) error {
	for i := 0; i < len(vm.vehiculos); i++ {
		if vm.vehiculos[i].ID == vehiculoID {
			vm.vehiculos[i].IncidenciaID = incidenciaID
			return nil
		}
	}
	return fmt.Errorf("vehículo con ID %d no encontrado", vehiculoID)
}

// ActualizarTiempoAcumulado actualiza el tiempo acumulado de atención de un vehículo
func (vm *VehiculoManager) ActualizarTiempoAcumulado(vehiculoID int, tiempo float64) error {
	for i := 0; i < len(vm.vehiculos); i++ {
		if vm.vehiculos[i].ID == vehiculoID {
			vm.vehiculos[i].TiempoAcumulado = vm.vehiculos[i].TiempoAcumulado + tiempo
			return nil
		}
	}
	return fmt.Errorf("vehículo con ID %d no encontrado", vehiculoID)
}

// ObtenerTiempoAcumulado obtiene el tiempo acumulado de un vehículo
func (vm *VehiculoManager) ObtenerTiempoAcumulado(vehiculoID int) (float64, bool) {
	for i := 0; i < len(vm.vehiculos); i++ {
		if vm.vehiculos[i].ID == vehiculoID {
			return vm.vehiculos[i].TiempoAcumulado, true
		}
	}
	return 0, false
}

// ListarVehiculosPorCliente lista todos los vehículos de un cliente
func (vm *VehiculoManager) ListarVehiculosPorCliente(clienteID int) []Vehiculo {
	lista := make([]Vehiculo, 0)
	for i := 0; i < len(vm.vehiculos); i++ {
		if vm.vehiculos[i].ClienteID == clienteID {
			lista = append(lista, vm.vehiculos[i])
		}
	}
	return lista
}
