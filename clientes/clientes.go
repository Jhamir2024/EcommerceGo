// Package clientes modela a los compradores (struct Cliente) y su registro
// (struct Registro).
package clientes

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Errores del paquete, comparables con errors.Is.
var (
	ErrNoExiste       = errors.New("el cliente no existe")
	ErrNombreVacio    = errors.New("el nombre no puede estar vacio")
	ErrCorreo         = errors.New("el correo no es valido")
	ErrCorreoRepetido = errors.New("ya existe un cliente con ese correo")
)

// Cliente es una persona que compra en la tienda. Campos no exportados:
// el correo solo cambia por SetCorreo, que lo valida.
type Cliente struct {
	id     int
	nombre string
	correo string
}

// NewCliente crea un cliente válido o devuelve el error que lo impide.
func NewCliente(id int, nombre, correo string) (*Cliente, error) {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return nil, ErrNombreVacio
	}
	correo = strings.ToLower(strings.TrimSpace(correo))
	if !correoValido(correo) {
		return nil, ErrCorreo
	}
	return &Cliente{id: id, nombre: nombre, correo: correo}, nil
}

// ID devuelve el identificador del cliente.
func (c Cliente) ID() int { return c.id }

// Nombre devuelve el nombre del cliente.
func (c Cliente) Nombre() string { return c.nombre }

// Correo devuelve el correo del cliente.
func (c Cliente) Correo() string { return c.correo }

// String implementa fmt.Stringer.
func (c Cliente) String() string {
	return fmt.Sprintf("#%d %s <%s>", c.id, c.nombre, c.correo)
}

// SetCorreo cambia el correo si tiene un formato válido.
func (c *Cliente) SetCorreo(nuevo string) error {
	nuevo = strings.ToLower(strings.TrimSpace(nuevo))
	if !correoValido(nuevo) {
		return ErrCorreo
	}
	c.correo = nuevo
	return nil
}

// correoValido es interna (minúscula): revisa que haya algo antes de la @,
// un punto después y ningún espacio. Es una validación simple, suficiente
// para evitar errores de tipeo; no cubre todos los casos del estándar.
func correoValido(correo string) bool {
	arroba := strings.Index(correo, "@")
	return arroba > 0 && strings.Contains(correo[arroba:], ".") && !strings.Contains(correo, " ")
}

// Registro guarda los clientes en un map por ID (búsqueda directa al armar
// carritos y pedidos).
type Registro struct {
	clientes    map[int]*Cliente
	siguienteID int
}

// NewRegistro crea un registro vacío.
func NewRegistro() *Registro {
	return &Registro{clientes: make(map[int]*Cliente), siguienteID: 1}
}

// Registrar crea un cliente nuevo. Rechaza correos repetidos.
func (r *Registro) Registrar(nombre, correo string) (*Cliente, error) {
	c, err := NewCliente(r.siguienteID, nombre, correo)
	if err != nil {
		return nil, err
	}
	for _, existente := range r.clientes {
		if existente.correo == c.correo {
			return nil, ErrCorreoRepetido
		}
	}
	r.clientes[c.id] = c
	r.siguienteID++
	return c, nil
}

// Buscar devuelve el cliente con ese ID o ErrNoExiste.
func (r *Registro) Buscar(id int) (*Cliente, error) {
	c, ok := r.clientes[id]
	if !ok {
		return nil, fmt.Errorf("%w: ID %d", ErrNoExiste, id)
	}
	return c, nil
}

// Listar devuelve los clientes ordenados por ID.
func (r *Registro) Listar() []*Cliente {
	lista := make([]*Cliente, 0, len(r.clientes))
	for _, c := range r.clientes {
		lista = append(lista, c)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].id < lista[j].id })
	return lista
}
