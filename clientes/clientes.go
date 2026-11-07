// Package clientes gestiona a las personas que compran en la tienda.
package clientes

import (
	"errors"
	"fmt"
	"strings"
)

// Datos de los clientes, guardados por ID. En el AA2 pasarán a un struct Cliente.
var (
	nombres   = map[int]string{}
	correos   = map[int]string{}
	siguiente = 1
)

// RegistrarCliente crea un cliente nuevo y devuelve su ID.
func RegistrarCliente(nombre string, correo string) (int, error) {
	if nombre == "" {
		return 0, errors.New("el nombre no puede estar vacio")
	}
	if !correoValido(correo) {
		return 0, errors.New("el correo no es valido")
	}
	if correoRegistrado(correo) {
		return 0, errors.New("ese correo ya esta registrado")
	}
	id := siguiente
	nombres[id], correos[id] = nombre, correo
	siguiente++
	return id, nil
}

// BuscarCliente devuelve el nombre y el correo de un cliente.
func BuscarCliente(id int) (string, string, error) {
	if !Existe(id) {
		return "", "", fmt.Errorf("el cliente %d no existe", id)
	}
	return nombres[id], correos[id], nil
}

// ListarClientes devuelve una línea por cliente, ordenadas por ID.
func ListarClientes() []string {
	lista := []string{}
	for id := 1; id < siguiente; id++ {
		if Existe(id) {
			lista = append(lista, fmt.Sprintf("%d) %s <%s>", id, nombres[id], correos[id]))
		}
	}
	return lista
}

// Existe indica si hay un cliente con ese ID.
func Existe(id int) bool {
	_, ok := nombres[id]
	return ok
}

// correoValido hace una revisión simple: que tenga @ y un punto.
func correoValido(correo string) bool {
	return strings.Contains(correo, "@") && strings.Contains(correo, ".")
}

// correoRegistrado evita que dos clientes usen el mismo correo.
func correoRegistrado(correo string) bool {
	for _, c := range correos {
		if strings.EqualFold(c, correo) {
			return true
		}
	}
	return false
}
