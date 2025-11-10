package main

import (
	"fmt"
)

// ClienteManager gestiona las operaciones sobre clientes
type ClienteManager struct {
	clientes []Cliente
	nextID   int
}

// NewClienteManager crea un nuevo gestor de clientes
func NewClienteManager() *ClienteManager {
	return &ClienteManager{
		clientes: make([]Cliente, 0),
		nextID:   1,
	}
}

// CrearCliente crea un nuevo cliente
func (cm *ClienteManager) CrearCliente(nombre, telefono, email string) Cliente {
	cliente := Cliente{
		ID:       cm.nextID,
		Nombre:   nombre,
		Telefono: telefono,
		Email:    email,
	}
	cm.clientes = append(cm.clientes, cliente)
	cm.nextID++
	return cliente
}

// ObtenerCliente obtiene un cliente por su ID
func (cm *ClienteManager) ObtenerCliente(id int) (Cliente, bool) {
	for i := 0; i < len(cm.clientes); i++ {
		if cm.clientes[i].ID == id {
			return cm.clientes[i], true
		}
	}
	return Cliente{}, false
}

// ActualizarCliente actualiza los datos de un cliente
func (cm *ClienteManager) ActualizarCliente(id int, nombre, telefono, email string) error {
	for i := 0; i < len(cm.clientes); i++ {
		if cm.clientes[i].ID == id {
			if nombre != "" {
				cm.clientes[i].Nombre = nombre
			}
			if telefono != "" {
				cm.clientes[i].Telefono = telefono
			}
			if email != "" {
				cm.clientes[i].Email = email
			}
			return nil
		}
	}
	return fmt.Errorf("cliente con ID %d no encontrado", id)
}

// EliminarCliente elimina un cliente
func (cm *ClienteManager) EliminarCliente(id int) error {
	for i := 0; i < len(cm.clientes); i++ {
		if cm.clientes[i].ID == id {
			cm.clientes = append(cm.clientes[:i], cm.clientes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("cliente con ID %d no encontrado", id)
}

// ListarClientes lista todos los clientes
func (cm *ClienteManager) ListarClientes() []Cliente {
	return cm.clientes
}

// AsignarVehiculo asigna un vehículo a un cliente
func (cm *ClienteManager) AsignarVehiculo(clienteID, vehiculoID int) error {
	for i := 0; i < len(cm.clientes); i++ {
		if cm.clientes[i].ID == clienteID {
			cm.clientes[i].VehiculoID = vehiculoID
			return nil
		}
	}
	return fmt.Errorf("cliente con ID %d no encontrado", clienteID)
}
