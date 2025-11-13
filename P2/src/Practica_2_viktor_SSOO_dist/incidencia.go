package main

import (
	"fmt"
	"time"
)

type IncidenciaManager struct {
	incidencias []Incidencia
	nextID      int
}

func NewIncidenciaManager() *IncidenciaManager {
	return &IncidenciaManager{
		incidencias: make([]Incidencia, 0),
		nextID:      1,
	}
}

func (im *IncidenciaManager) CrearIncidencia(tipo TipoIncidencia, prioridad Prioridad, descripcion string, vehiculoID int) Incidencia {
	incidencia := Incidencia{
		ID:           im.nextID,
		MecanicosIDs: make([]int, 0),
		VehiculoID:   vehiculoID,
		Tipo:         tipo,
		Prioridad:    prioridad,
		Descripcion:  descripcion,
		Estado:       Abierta,
	}
	im.incidencias = append(im.incidencias, incidencia)
	im.nextID++
	return incidencia
}

func (im *IncidenciaManager) ObtenerIncidencia(id int) (Incidencia, bool) {
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].ID == id {
			return im.incidencias[i], true
		}
	}
	return Incidencia{}, false
}

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

func (im *IncidenciaManager) EliminarIncidencia(id int) error {
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].ID == id {
			im.incidencias = append(im.incidencias[:i], im.incidencias[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("incidencia con ID %d no encontrada", id)
}

func (im *IncidenciaManager) ListarIncidencias() []Incidencia {
	return im.incidencias
}

func (im *IncidenciaManager) CambiarEstado(id int, nuevoEstado EstadoIncidencia) error {
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].ID == id {
			im.incidencias[i].Estado = nuevoEstado
			return nil
		}
	}
	return fmt.Errorf("incidencia con ID %d no encontrada", id)
}

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

func (im *IncidenciaManager) ContarTodasIncidencias() map[int]int {
	counts := make(map[int]int)
	for _, incidencia := range im.incidencias {
		if incidencia.VehiculoID > 0 {
			counts[incidencia.VehiculoID]++
		}
	}
	return counts
}

func (im *IncidenciaManager) ObtenerIncidenciasPorVehiculo(vehiculoID int) []Incidencia {
	lista := make([]Incidencia, 0)

	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].VehiculoID == vehiculoID && im.incidencias[i].Estado == Abierta {
			lista = append(lista, im.incidencias[i])
		}
	}
	return lista
}

func (im *IncidenciaManager) EliminarIncidenciasPorVehiculo(vehiculoID int) {
	nuevaLista := make([]Incidencia, 0)
	for i := 0; i < len(im.incidencias); i++ {
		if im.incidencias[i].VehiculoID != vehiculoID {
			nuevaLista = append(nuevaLista, im.incidencias[i])
		}
	}
	im.incidencias = nuevaLista
}
