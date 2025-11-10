package main

import (
	"fmt"
	"time"
)

// IncidenciaManager gestiona las operaciones sobre incidencias
type IncidenciaManager struct {
	incidencias []Incidencia
	nextID      int
}

// NewIncidenciaManager crea un nuevo gestor de incidencias
func NewIncidenciaManager() *IncidenciaManager {
	return &IncidenciaManager{
		incidencias: make([]Incidencia, 0),
		nextID:      1,
	}
}

// CrearIncidencia crea una nueva incidencia
func (im *IncidenciaManager) CrearIncidencia(tipo TipoIncidencia, prioridad Prioridad, descripcion string) Incidencia {
	incidencia := Incidencia{
		ID:           im.nextID,
		MecanicosIDs: make([]int, 0),
		Tipo:         tipo,
		Prioridad:    prioridad,
		Descripcion:  descripcion,
		Estado:       Abierta,
	}
	im.incidencias = append(im.incidencias, incidencia)
	im.nextID++
	return incidencia
}

// ObtenerIncidencia obtiene una incidencia por su ID
func (im *IncidenciaManager) ObtenerIncidencia(id int) (Incidencia, bool) {
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].ID == id {
			return im.incidencias[i], true
		}
	}
	return Incidencia{}, false
}

// ActualizarIncidencia actualiza los datos de una incidencia
func (im *IncidenciaManager) ActualizarIncidencia(id int, tipo TipoIncidencia, prioridad Prioridad, descripcion string) error {
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].ID == id {
			if tipo != "" {
				im.incidencias[i].Tipo = tipo
			}
			if prioridad != "" {
				im.incidencias[i].Prioridad = prioridad
			}
			if descripcion != "" {
				im.incidencias[i].Descripcion = descripcion
			}
			return nil
		}
	}
	return fmt.Errorf("incidencia con ID %d no encontrada", id)
}

// EliminarIncidencia elimina una incidencia
func (im *IncidenciaManager) EliminarIncidencia(id int) error {
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].ID == id {
			im.incidencias = append(im.incidencias[:i], im.incidencias[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("incidencia con ID %d no encontrada", id)
}

// ListarIncidencias lista todas las incidencias
func (im *IncidenciaManager) ListarIncidencias() []Incidencia {
	return im.incidencias
}

// CambiarEstado cambia el estado de una incidencia
func (im *IncidenciaManager) CambiarEstado(id int, nuevoEstado EstadoIncidencia) error {
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].ID == id {
			im.incidencias[i].Estado = nuevoEstado
			return nil
		}
	}
	return fmt.Errorf("incidencia con ID %d no encontrada", id)
}

// AsignarMecanico asigna un mecánico a una incidencia
func (im *IncidenciaManager) AsignarMecanico(incidenciaID, mecanicoID int) error {
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].ID == incidenciaID {
			// Verificar si el mecánico ya está asignado
			for j := 0; j < len(im.incidencias[i].MecanicosIDs); j++ {
				if im.incidencias[i].MecanicosIDs[j] == mecanicoID {
					return fmt.Errorf("mecánico %d ya está asignado a esta incidencia", mecanicoID)
				}
			}
			im.incidencias[i].MecanicosIDs = append(im.incidencias[i].MecanicosIDs, mecanicoID)
			return nil
		}
	}
	return fmt.Errorf("incidencia con ID %d no encontrada", incidenciaID)
}

// DesasignarMecanico desasigna un mecánico de una incidencia
func (im *IncidenciaManager) DesasignarMecanico(incidenciaID, mecanicoID int) error {
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].ID == incidenciaID {
			nuevaLista := make([]int, 0)
			encontrado := false
			for j := 0; j < len(im.incidencias[i].MecanicosIDs); j++ {
				if im.incidencias[i].MecanicosIDs[j] != mecanicoID {
					nuevaLista = append(nuevaLista, im.incidencias[i].MecanicosIDs[j])
				} else {
					encontrado = true
				}
			}
			if !encontrado {
				return fmt.Errorf("mecánico %d no está asignado a esta incidencia", mecanicoID)
			}
			im.incidencias[i].MecanicosIDs = nuevaLista
			return nil
		}
	}
	return fmt.Errorf("incidencia con ID %d no encontrada", incidenciaID)
}

// ObtenerTiempoAtencion devuelve el tiempo de atención según el tipo de incidencia
func ObtenerTiempoAtencion(tipo TipoIncidencia) time.Duration {
	if tipo == Mecanica {
		return 5 * time.Second
	}
	if tipo == Electrica {
		return 7 * time.Second
	}
	if tipo == Carroceria {
		return 11 * time.Second
	}
	return 5 * time.Second
}
