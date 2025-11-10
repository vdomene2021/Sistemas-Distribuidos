package main

import (
	"fmt"
)

type MecanicoManager struct {
	mecanicos []Mecanico
	nextID    int
}

func NewMecanicoManager() *MecanicoManager {
	return &MecanicoManager{
		mecanicos: make([]Mecanico, 0),
		nextID:    1,
	}
}

func (mm *MecanicoManager) CrearMecanico(nombre string, especialidad Especialidad, experiencia int) Mecanico {
	mecanico := Mecanico{
		ID:           mm.nextID,
		Nombre:       nombre,
		Especialidad: especialidad,
		Experiencia:  experiencia,
		Activo:       true,
		Ocupado:      false,
	}
	mm.mecanicos = append(mm.mecanicos, mecanico)
	mm.nextID++
	return mecanico
}

// ObtenerMecanico obtiene un mecánico por su ID
func (mm *MecanicoManager) ObtenerMecanico(id int) (Mecanico, bool) {
	for i := 0; i < len(mm.mecanicos); i++ {
		if mm.mecanicos[i].ID == id {
			return mm.mecanicos[i], true
		}
	}
	return Mecanico{}, false
}

func (mm *MecanicoManager) ActualizarMecanico(id int, nombre string, especialidad Especialidad, experiencia int) error {
	for i := 0; i < len(mm.mecanicos); i++ {
		if mm.mecanicos[i].ID == id {
			if nombre != "" {
				mm.mecanicos[i].Nombre = nombre
			}
			if especialidad != "" {
				mm.mecanicos[i].Especialidad = especialidad
			}
			if experiencia > 0 {
				mm.mecanicos[i].Experiencia = experiencia
			}
			return nil
		}
	}
	return fmt.Errorf("mecánico con ID %d no encontrado", id)
}

func (mm *MecanicoManager) EliminarMecanico(id int) error {
	for i := 0; i < len(mm.mecanicos); i++ {
		if mm.mecanicos[i].ID == id {
			mm.mecanicos = append(mm.mecanicos[:i], mm.mecanicos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("mecánico con ID %d no encontrado", id)
}

func (mm *MecanicoManager) ListarMecanicos() []Mecanico {
	return mm.mecanicos
}

// CambiarEstadoActivo cambia el estado activo de un mecánico
func (mm *MecanicoManager) CambiarEstadoActivo(id int, activo bool) error {
	for i := 0; i < len(mm.mecanicos); i++ {
		if mm.mecanicos[i].ID == id {
			mm.mecanicos[i].Activo = activo
			return nil
		}
	}
	return fmt.Errorf("mecánico con ID %d no encontrado", id)
}

func (mm *MecanicoManager) CambiarEstadoOcupado(id int, ocupado bool) error {
	for i := 0; i < len(mm.mecanicos); i++ {
		if mm.mecanicos[i].ID == id {
			mm.mecanicos[i].Ocupado = ocupado
			return nil
		}
	}
	return fmt.Errorf("mecánico con ID %d no encontrado", id)
}

func (mm *MecanicoManager) ListarMecanicosDisponibles() []Mecanico {
	lista := make([]Mecanico, 0)
	for i := 0; i < len(mm.mecanicos); i++ {
		if mm.mecanicos[i].Activo && !mm.mecanicos[i].Ocupado {
			lista = append(lista, mm.mecanicos[i])
		}
	}
	return lista
}

func (mm *MecanicoManager) ListarMecanicosPorEspecialidad(especialidad Especialidad) []Mecanico {
	lista := make([]Mecanico, 0)
	for i := 0; i < len(mm.mecanicos); i++ {
		if mm.mecanicos[i].Especialidad == especialidad && mm.mecanicos[i].Activo {
			lista = append(lista, mm.mecanicos[i])
		}
	}
	return lista
}

func (mm *MecanicoManager) ContarMecanicosPorEspecialidad() map[Especialidad]int {
	conteo := make(map[Especialidad]int)
	conteo[EspecialidadMecanica] = 0
	conteo[EspecialidadElectrica] = 0
	conteo[EspecialidadCarroceria] = 0

	for i := 0; i < len(mm.mecanicos); i++ {
		if mm.mecanicos[i].Activo {
			conteo[mm.mecanicos[i].Especialidad]++
		}
	}

	return conteo
}
