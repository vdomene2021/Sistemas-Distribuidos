package main

import (
	"fmt"
)

type ClienteManager struct {
	clientes []Cliente
	nextID   int
}

func NewClienteManager() *ClienteManager {
	return &ClienteManager{
		clientes: make([]Cliente, 0),
		nextID:   1,
	}
}

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

func (cm *ClienteManager) ObtenerCliente(id int) (Cliente, bool) {
	for i := 0; i < len(cm.clientes); i++ {
		if cm.clientes[i].ID == id {
			return cm.clientes[i], true
		}
	}
	return Cliente{}, false
}

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

func (cm *ClienteManager) EliminarCliente(id int) error {
	for i := 0; i < len(cm.clientes); i++ {
		if cm.clientes[i].ID == id {
			cm.clientes = append(cm.clientes[:i], cm.clientes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("cliente con ID %d no encontrado", id)
}

func (cm *ClienteManager) ListarClientes() []Cliente {
	return cm.clientes
}
