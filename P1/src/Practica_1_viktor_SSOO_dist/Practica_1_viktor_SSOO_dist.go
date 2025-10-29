package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Cliente struct {
	ID        int
	Nombre    string
	Telefono  string
	Email     string
	Vehiculos []string // matrículas
}

type Vehiculo struct {
	Matricula        string
	Marca            string
	Modelo           string
	IDCliente        int
	EnTaller         bool
	FechaEntrada     string
	FechaSalidaEst   string
	IncidenciaActiva int // 0 = ninguna
}

type Incidencia struct {
	ID          int
	Matricula   string
	Mecanicos   []int
	Tipo        string // "mecanica", "electrica", "carroceria"
	Prioridad   string // "baja", "media", "alta"
	Descripcion string
	Estado      string // "abierta", "proceso", "cerrada"
}

type Mecanico struct {
	ID           int
	Nombre       string
	Especialidad string // "mecanica", "electrica", "carroceria"
	Anos         int
	Activo       bool
}

type Taller struct {
	Clientes    map[int]*Cliente
	Vehiculos   map[string]*Vehiculo
	Incidencias map[int]*Incidencia
	Mecanicos   map[int]*Mecanico
	nextCli     int
	nextInc     int
	nextMec     int
}

// ----------------------------- Utilidades de entrada simples -----------------------------

// Para leer desde consola
var in = bufio.NewReader(os.Stdin)

// Leer una línea de texto
func leer(texto string) string {
	fmt.Print(texto)
	s, _ := in.ReadString('\n')
	return strings.TrimSpace(s)
}

// Leer un entero
func leerInt(texto string) int {
	for {
		s := leer(texto)
		n, err := strconv.Atoi(s)
		if err == nil {
			return n
		}
		fmt.Println("Por favor, introduce un número entero válido.")
	}
}

// Leer una lista de IDs separados por comas
func leerIDsComa(texto string) []int {
	raw := leer(texto)
	if strings.TrimSpace(raw) == "" { // Si no introducimos ningun iD devolvemos nil
		return nil
	}
	ids := []int{}
	for _, p := range strings.Split(raw, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err == nil {
			ids = append(ids, n)
		}
	}
	return ids
}

func elegirOpcion(titulo string, opciones []string) string {
	for {
		fmt.Println(titulo)
		for i, op := range opciones {
			fmt.Printf("%d) %s\n", i+1, op)
		}
		opNum := leerInt("Opción: ")
		if opNum >= 1 && opNum <= len(opciones) {
			return opciones[opNum-1]
		}
		fmt.Println("Opción inválida, intenta de nuevo.")
	}
}

// ----------------------------- Creación del taller -----------------------------
func nuevoTaller() *Taller {
	return &Taller{
		Clientes:    map[int]*Cliente{},
		Vehiculos:   map[string]*Vehiculo{},
		Incidencias: map[int]*Incidencia{},
		Mecanicos:   map[int]*Mecanico{},
	}
}

// ----------------------------- Capacidad -----------------------------
func (t *Taller) capacidadMaxima() int {
	activos := 0
	for _, m := range t.Mecanicos {
		if m.Activo {
			activos++
		}
	}
	return activos * 2 // Porque nos piden 2 plazas por mecánico
}

func (t *Taller) plazasOcupadas() int {
	ocup := 0
	for _, v := range t.Vehiculos {
		if v.EnTaller {
			ocup++
		}
	}
	return ocup
}

func (t *Taller) plazasLibres() int {
	return t.capacidadMaxima() - t.plazasOcupadas()
}

// ----------------------------- CLIENTES -----------------------------
func (t *Taller) crearCliente() {
	n := leer("Nombre: ")
	tel := leer("Teléfono: ")
	em := leer("Email: ")
	t.nextCli++
	id := t.nextCli
	t.Clientes[id] = &Cliente{ID: id, Nombre: n, Telefono: tel, Email: em}
	fmt.Println("Cliente creado con ID:", id)
	fmt.Println()
}

func (t *Taller) verClientes() {
	if len(t.Clientes) == 0 {
		fmt.Println("No hay clientes registrados")
		return
	}
	for _, c := range t.Clientes {
		fmt.Printf("ID: %d | Nombre: %s | Telefono: %s | Email: %s | Vehículos: %v\n\n", c.ID, c.Nombre, c.Telefono, c.Email, c.Vehiculos)
	}
}

func (t *Taller) modificarCliente() {
	id := leerInt("ID del cliente: ")
	c, ok := t.Clientes[id]
	if !ok {
		fmt.Println("El cliente no existe")
		fmt.Println()
		return
	}
	if v := leer("Nuevo nombre (enter para no modificar): "); v != "" {
		c.Nombre = v
	}
	if v := leer("Nuevo teléfono (enter para no modificar): "); v != "" {
		c.Telefono = v
	}
	if v := leer("Nuevo email (enter para no modificar): "); v != "" {
		c.Email = v
	}
	fmt.Println("Cliente actualizado con exito")
	fmt.Println()
}

func (t *Taller) eliminarCliente() {
	id := leerInt("ID del cliente: ")
	c, ok := t.Clientes[id]
	if !ok {
		fmt.Println("El cliente no existe")
		fmt.Println()
		return
	}
	for _, m := range c.Vehiculos {
		if v, ok := t.Vehiculos[m]; ok && v.EnTaller {
			fmt.Println("No se puede eliminar ya que tiene vehículo/s en el taller")
			return
		}
	}
	for _, m := range c.Vehiculos {
		delete(t.Vehiculos, m)
	}
	delete(t.Clientes, id)
	fmt.Println("Cliente eliminado con exito")
	fmt.Println()
}

// ----------------------------- VEHÍCULOS -----------------------------
func (t *Taller) crearVehiculo() {
	mat := strings.ToUpper(leer("Matrícula: ")) // Las matriculas en mayúsculas
	if _, existe := t.Vehiculos[mat]; existe {
		fmt.Println("El vehiculo ya existe")
		fmt.Println()
		return
	}
	marca := leer("Marca: ")
	modelo := leer("Modelo: ")
	idc := leerInt("ID cliente: ")
	c, ok := t.Clientes[idc]
	if !ok {
		fmt.Println("El cliente no existe, registrelo primero")
		fmt.Println()
		return
	}
	v := &Vehiculo{Matricula: mat, Marca: marca, Modelo: modelo, IDCliente: idc}
	t.Vehiculos[mat] = v // El id del vehículo es la matrícula
	c.Vehiculos = append(c.Vehiculos, mat)
	fmt.Println("Vehículo creado con exito")
	fmt.Println()
}

func (t *Taller) verVehiculos() {
	if len(t.Vehiculos) == 0 {
		fmt.Println("No hay vehículos registrados")
		fmt.Println()
		return
	}
	for _, v := range t.Vehiculos {
		fmt.Printf("Matricula: %s | Modelo: %s %s | Cliente: %d | En taller: %t | Entrada: %s | Salida estimada: %s | Incidencia activa: %d\n\n",
			v.Matricula, v.Marca, v.Modelo, v.IDCliente, v.EnTaller, v.FechaEntrada, v.FechaSalidaEst, v.IncidenciaActiva)
	}
}

func (t *Taller) modificarVehiculo() {
	mat := strings.ToUpper(leer("Matrícula: "))
	v, ok := t.Vehiculos[mat]
	if !ok {
		fmt.Println("No hay vehiculo registrado con esa matrícula")
		fmt.Println()
		return
	}
	if s := leer("Nueva marca (enter para no modificar): "); s != "" {
		v.Marca = s
	}
	if s := leer("Nuevo modelo (enter para no modificar): "); s != "" {
		v.Modelo = s
	}
	fmt.Println("Vehículo actualizado con exito")
	fmt.Println()
}

func (t *Taller) eliminarVehiculo() {
	mat := strings.ToUpper(leer("Matrícula: "))
	v, ok := t.Vehiculos[mat]
	if !ok {
		fmt.Println("No hay vehiculo registrado con esa matrícula")
		fmt.Println()
		return
	}
	if v.EnTaller {
		fmt.Println("No se puede eliminar ya que está en el taller")
		fmt.Println()
		return
	}
	c := t.Clientes[v.IDCliente]
	nueva := []string{}
	for _, m := range c.Vehiculos {
		if m != mat {
			nueva = append(nueva, m)
		}
	}
	c.Vehiculos = nueva
	delete(t.Vehiculos, mat)
	fmt.Println("Vehículo eliminado con exito")
	fmt.Println()
}

// ----------------------------- INCIDENCIAS -----------------------------
func (t *Taller) crearIncidencia() {
	mat := strings.ToUpper(leer("Matrícula: "))
	v, ok := t.Vehiculos[mat]
	if !ok {
		fmt.Println("Registre antes el vehículo para asignarle la incidencia")
		fmt.Println()
		return
	}

	tipos := []string{"mecanica", "electrica", "carroceria"}
	prioridades := []string{"baja", "media", "alta"}

	tipo := elegirOpcion("Tipo de incidencia:", tipos)
	prio := elegirOpcion("Prioridad:", prioridades)
	desc := leer("Descripción: ")
	mecs := leerIDsComa("IDs de mecánicos (vacío --> ninguno): ")

	// validar mecánicos activos
	for _, id := range mecs {
		m, ok := t.Mecanicos[id]
		if !ok || !m.Activo {
			fmt.Println("Mecánico inválido o inactivo:", id)
			fmt.Println()
			return
		}
	}

	t.nextInc++
	id := t.nextInc
	inc := &Incidencia{
		ID:          id,
		Matricula:   mat,
		Mecanicos:   mecs,
		Tipo:        tipo,
		Prioridad:   prio,
		Descripcion: desc,
		Estado:      "abierta",
	}
	t.Incidencias[id] = inc
	v.IncidenciaActiva = id
	fmt.Println("Incidencia creada con ID:", id)
}

func (t *Taller) verIncidencias() {
	if len(t.Incidencias) == 0 {
		fmt.Println("No hay incidencias.")
		return
	}
	for _, i := range t.Incidencias {
		fmt.Printf("ID:%d | Matricula: %s | Tipo: %s | Prioridad: %s | Estado: %s | Mecanico: %v | Descripcion: %s\n",
			i.ID, i.Matricula, i.Tipo, i.Prioridad, i.Estado, i.Mecanicos, i.Descripcion)
	}
}

func (t *Taller) modificarIncidencia() {
	id := leerInt("ID incidencia: ")
	i, ok := t.Incidencias[id]
	if !ok {
		fmt.Println("No existe la incidencia")
		fmt.Println()
		return
	}

	estados := []string{"abierta", "proceso", "cerrada"}
	fmt.Printf("Nuevo estado (enter para no modificar: %s)\n", i.Estado)
	for idx, est := range estados {
		fmt.Printf("%d) %s\n", idx+1, est)
	}
	entrada := leer("Opción (1-3, enter para no modificar): ")
	if entrada != "" {
		if op, err := strconv.Atoi(entrada); err == nil && op >= 1 && op <= len(estados) {
			i.Estado = estados[op-1]
		} else {
			fmt.Println("Entrada inválida, se mantiene el estado actual")
		}
	}

	// Descripción
	if s := leer("Nueva descripción (enter para no modificar): "); s != "" {
		i.Descripcion = s
	}

	// Mecánicos
	if ids := leerIDsComa("Nuevos mecánicos (enter para no modificar): "); ids != nil {
		// validar activos
		for _, mid := range ids {
			m, ok := t.Mecanicos[mid]
			if !ok || !m.Activo {
				fmt.Println("Mecánico inválido o inactivo: ", mid)
				return
			}
		}
		i.Mecanicos = ids
	}

	fmt.Println("Incidencia actualizada con exito")
}

func (t *Taller) eliminarIncidencia() {
	id := leerInt("ID incidencia: ")
	i, ok := t.Incidencias[id]
	if !ok {
		fmt.Println("No existe la incidencia")
		return
	}
	v := t.Vehiculos[i.Matricula]
	if v != nil && v.IncidenciaActiva == id {
		v.IncidenciaActiva = 0
	}
	delete(t.Incidencias, id)
	fmt.Println("Incidencia eliminada con exito")
}

// ----------------------------- MECÁNICOS -----------------------------
func (t *Taller) crearMecanico() {
	n := leer("Nombre: ")

	especialidades := []string{"mecanica", "electrica", "carroceria"}
	esp := elegirOpcion("Especialidad:", especialidades)

	an := leerInt("Años de experiencia: ")

	t.nextMec++
	id := t.nextMec
	t.Mecanicos[id] = &Mecanico{
		ID:           id,
		Nombre:       n,
		Especialidad: esp,
		Anos:         an,
		Activo:       true,
	}
	fmt.Println("Mecánico creado con ID:", id)
}

func (t *Taller) verMecanicos() {
	if len(t.Mecanicos) == 0 {
		fmt.Println("No hay mecánicos registrados")
		return
	}
	for _, m := range t.Mecanicos {
		fmt.Printf("ID: %d | Nombre: %s | Especialidad: %s | Años experiencia: %d | Activo: %t\n", m.ID, m.Nombre, m.Especialidad, m.Anos, m.Activo)
	}
}

func (t *Taller) modificarMecanico() {
	id := leerInt("ID mecánico: ")
	m, ok := t.Mecanicos[id]
	if !ok {
		fmt.Println("No existe el mecanico con ese ID")
		return
	}

	// Nombre
	if s := leer("Nuevo nombre (enter para no modificar): "); s != "" {
		m.Nombre = s
	}

	// Especialidad
	especialidades := []string{"mecanica", "electrica", "carroceria"}
	fmt.Printf("Nueva especialidad (enter para no modificar: %s)\n", m.Especialidad)
	for idx, esp := range especialidades {
		fmt.Printf("%d) %s\n", idx+1, esp)
	}
	entrada := leer("Opción (1-3, enter para no modificar): ")
	if entrada != "" {
		if op, err := strconv.Atoi(entrada); err == nil && op >= 1 && op <= len(especialidades) {
			m.Especialidad = especialidades[op-1]
		} else {
			fmt.Println("Entrada inválida, se mantiene la especialidad actual")
		}
	}

	// Años de experiencia
	if s := leer("Años experiencia (enter para no modificar): "); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			m.Anos = n
		} else {
			fmt.Println("Valor inválido, se mantiene el actual")
		}
	}

	fmt.Println("Mecánico actualizado con exito")
}

func (t *Taller) eliminarMecanico() {
	id := leerInt("ID mecánico: ")
	if _, ok := t.Mecanicos[id]; !ok {
		fmt.Println("No existe el mecánico con ese ID")
		return
	}
	for _, inc := range t.Incidencias {
		if inc.Estado != "cerrada" {
			for _, mid := range inc.Mecanicos {
				if mid == id {
					fmt.Println("No se puede eliminar este mecanico, esta asignado a incidencia abierta.")
					return
				}
			}
		}
	}
	delete(t.Mecanicos, id)
	fmt.Println("Mecánico eliminado con exito")
}

func (t *Taller) altaBajaMecanico() {
	id := leerInt("ID mecánico: ")
	m, ok := t.Mecanicos[id]
	if !ok {
		fmt.Println("No existe el mecánico con ese ID")
		return
	}
	m.Activo = !m.Activo
	fmt.Printf("Mecánico %d ahora Activo = %t\n", id, m.Activo)
}

// ----------------------------- OPERACIONES -----------------------------
func (t *Taller) asignarVehiculoATaller() {
	if t.plazasLibres() <= 0 {
		fmt.Println("No hay plazas libres en este momento")
		return
	}
	mat := strings.ToUpper(leer("Matrícula: "))
	v, ok := t.Vehiculos[mat]
	if !ok {
		fmt.Println("Vehículo no registrado")
		return
	}
	if v.EnTaller {
		fmt.Println("El vehiculo ya está en el taller")
		return
	}
	v.FechaEntrada = leer("Fecha de entrada (YYYY-MM-DD): ")
	v.FechaSalidaEst = leer("Fecha de salida estimada (YYYY-MM-DD): ")
	v.EnTaller = true
	fmt.Printf("Ingresado %s. Libres: %d\n", mat, t.plazasLibres())
}

func (t *Taller) retirarVehiculoDeTaller() {
	mat := strings.ToUpper(leer("Matrícula: "))
	v, ok := t.Vehiculos[mat]
	if !ok {
		fmt.Println("Vehículo no registrado")
		return
	}
	if !v.EnTaller {
		fmt.Println("El vehículo no está en el taller")
		return
	}
	if v.IncidenciaActiva != 0 {
		if inc, ok := t.Incidencias[v.IncidenciaActiva]; ok && inc.Estado != "cerrada" {
			fmt.Println("No se puede retirar, incidencia no cerrada.")
			return
		}
	}
	v.EnTaller = false
	v.FechaEntrada, v.FechaSalidaEst = "", ""
	fmt.Printf("Retirado %s. Libres: %d\n", mat, t.plazasLibres())
}

func (t *Taller) verEstadoTaller() {
	fmt.Printf("Capacidad: %d | Ocupadas: %d | Libres: %d\n", t.capacidadMaxima(), t.plazasOcupadas(), t.plazasLibres())
}

func (t *Taller) cambiarEstadoIncidencia() {
	id := leerInt("ID incidencia: ")
	i, ok := t.Incidencias[id]
	if !ok {
		fmt.Println("Incidencia no registrada")
		return
	}

	estados := []string{"abierta", "proceso", "cerrada"}
	fmt.Printf("Estado actual: %s\n", i.Estado)
	nuevo := elegirOpcion("Nuevo estado:", estados)

	i.Estado = nuevo
	fmt.Println("Incidencia actualizada con exito")
}

func (t *Taller) listarIncidenciasDeVehiculo() {
	mat := strings.ToUpper(leer("Matrícula: "))
	encontradas := false
	for _, i := range t.Incidencias {
		if i.Matricula == mat {
			fmt.Printf("ID: %d | Estado: %s | Prioridad: %s | Descripcion: %s\n", i.ID, i.Estado, i.Prioridad, i.Descripcion)
			encontradas = true
		}
	}
	if !encontradas {
		fmt.Println("Sin incidencias o vehículo no registrado")
	}
}

func (t *Taller) listarVehiculosDeCliente() {
	id := leerInt("ID cliente: ")
	c, ok := t.Clientes[id]
	if !ok {
		fmt.Println("Cliente no registrado")
		return
	}
	if len(c.Vehiculos) == 0 {
		fmt.Println("Cliente sin vehículos registrados")
		return
	}
	for _, m := range c.Vehiculos {
		v := t.Vehiculos[m]
		fmt.Printf("Matricula: %s | Modelo: %s %s | En Taller: %t\n", v.Matricula, v.Marca, v.Modelo, v.EnTaller)
	}
}

func (t *Taller) listarMecanicosDisponibles() {
	asignados := map[int]bool{}
	for _, inc := range t.Incidencias {
		if inc.Estado != "cerrada" {
			for _, mid := range inc.Mecanicos {
				asignados[mid] = true
			}
		}
	}
	hay := false
	for id, m := range t.Mecanicos {
		if m.Activo && !asignados[id] {
			fmt.Printf("ID: %d | Nombre: %s | Especialidad: %s\n", m.ID, m.Nombre, m.Especialidad)
			hay = true
		}
	}
	if !hay {
		fmt.Println("No hay mecánicos disponibles en este momento")
	}
}

func (t *Taller) listarIncidenciasDeMecanico() {
	id := leerInt("ID mecánico: ")
	hay := false
	for _, inc := range t.Incidencias {
		for _, mid := range inc.Mecanicos {
			if mid == id {
				fmt.Printf("ID: %d | Matricula: %s | Estado: %s | Descripcion: %s\n", inc.ID, inc.Matricula, inc.Estado, inc.Descripcion)
				hay = true
				break
			}
		}
	}
	if !hay {
		fmt.Println("No hay incidencias asignadas a este mecanico")
	}
}

func (t *Taller) listarClientesConVehiculosEnTaller() {
	vistos := map[int]bool{}
	for _, v := range t.Vehiculos {
		if v.EnTaller {
			vistos[v.IDCliente] = true
		}
	}
	if len(vistos) == 0 {
		fmt.Println("No hay clientes con vehículos en el taller")
		return
	}
	for id := range vistos {
		c := t.Clientes[id]
		fmt.Printf("ID: %d | Nombre: %s\n", c.ID, c.Nombre)
	}
}

func (t *Taller) listarTodasIncidenciasConEstado() {
	if len(t.Incidencias) == 0 {
		fmt.Println("No hay incidencias en estos momentos")
		return
	}
	for _, inc := range t.Incidencias {
		fmt.Printf("ID: %d | Matricula: %s | Estado: %s | Prioridad: %s | Tipo: %s\n", inc.ID, inc.Matricula, inc.Estado, inc.Prioridad, inc.Tipo)
	}
}

// ----------------------------- Menú -----------------------------
func menu() {
	fmt.Println("\n\n======================= BIENVENIDO =======================")
	fmt.Println()
	fmt.Println("    1) Clientes")
	fmt.Println("    2) Vehículos")
	fmt.Println("    3) Incidencias")
	fmt.Println("    4) Mecánicos")
	fmt.Println("    5) Asignar vehículo a plaza")
	fmt.Println("    6) Retirar vehículo del taller")
	fmt.Println("    7) Ver estado del taller")
	fmt.Println("    8) Alta/Baja de mecánico")
	fmt.Println("    9) Cambiar estado de una incidencia")
	fmt.Println("    10) Listar incidencias de un vehículo")
	fmt.Println("    11) Listar vehículos de un cliente")
	fmt.Println("    12) Listar mecánicos disponibles")
	fmt.Println("    13) Listar incidencias de un mecánico")
	fmt.Println("    14) Listar clientes con vehículos en el taller")
	fmt.Println("    15) Listar todas las incidencias con su estado")
	fmt.Println("    0) Salir")
	fmt.Println()
}

// ----------------------------- Main -----------------------------
func main() {
	t := nuevoTaller()
	for {
		menu()
		op := leerInt("Opción: ")
		switch op {

		// CLIENTES
		case 1:
			fmt.Println("\n\n====== CLIENTES ======")
			fmt.Println()
			fmt.Println("1) Crear")
			fmt.Println("2) Ver")
			fmt.Println("3) Modificar")
			fmt.Println("4) Eliminar")
			fmt.Println("0) Volver")
			fmt.Println()
			s := leerInt("Subopción: ")
			switch s {
			case 1:
				t.crearCliente()
			case 2:
				t.verClientes()
			case 3:
				t.modificarCliente()
			case 4:
				t.eliminarCliente()
			}

		// VEHÍCULOS
		case 2:
			fmt.Println("\n\n====== VEHÍCULOS ======")
			fmt.Println()
			fmt.Println("1) Crear")
			fmt.Println("2) Ver")
			fmt.Println("3) Modificar")
			fmt.Println("4) Eliminar")
			fmt.Println("0) Volver")
			fmt.Println()
			s := leerInt("Subopción: ")
			switch s {
			case 1:
				t.crearVehiculo()
			case 2:
				t.verVehiculos()
			case 3:
				t.modificarVehiculo()
			case 4:
				t.eliminarVehiculo()
			}

		// INCIDENCIAS
		case 3:
			fmt.Println("\n\n====== INCIDENCIAS ======")
			fmt.Println()
			fmt.Println("1) Crear")
			fmt.Println("2) Ver")
			fmt.Println("3) Modificar")
			fmt.Println("4) Eliminar")
			fmt.Println("0) Volver")
			fmt.Println()
			s := leerInt("Subopción: ")
			switch s {
			case 1:
				t.crearIncidencia()
			case 2:
				t.verIncidencias()
			case 3:
				t.modificarIncidencia()
			case 4:
				t.eliminarIncidencia()
			}

		// MECÁNICOS
		case 4:
			fmt.Println("\n\n====== MECÁNICOS ======")
			fmt.Println()
			fmt.Println("1) Crear")
			fmt.Println("2) Ver")
			fmt.Println("3) Modificar")
			fmt.Println("4) Eliminar")
			fmt.Println("0) Volver")
			fmt.Println()
			s := leerInt("Subopción: ")
			switch s {
			case 1:
				t.crearMecanico()
			case 2:
				t.verMecanicos()
			case 3:
				t.modificarMecanico()
			case 4:
				t.eliminarMecanico()
			}
		case 5:
			t.asignarVehiculoATaller()
		case 6:
			t.retirarVehiculoDeTaller()
		case 7:
			t.verEstadoTaller()
		case 8:
			t.altaBajaMecanico()
		case 9:
			t.cambiarEstadoIncidencia()
		case 10:
			t.listarIncidenciasDeVehiculo()
		case 11:
			t.listarVehiculosDeCliente()
		case 12:
			t.listarMecanicosDisponibles()
		case 13:
			t.listarIncidenciasDeMecanico()
		case 14:
			t.listarClientesConVehiculosEnTaller()
		case 15:
			t.listarTodasIncidenciasConEstado()
		case 0:
			fmt.Println("\n¡Hasta luego!")
			fmt.Println()
			return
		default:
			fmt.Println("Opción inválida, inténtalo de nuevo")
		}
	}
}
