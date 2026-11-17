// Package carrito modela el carrito de compras de un cliente.
package carrito

import (
	"errors"
	"fmt"

	"ecommercego/clientes"
	"ecommercego/productos"
)

// ErrNoEnCarrito se devuelve al quitar un producto que no se había agregado.
var ErrNoEnCarrito = errors.New("ese producto no esta en el carrito")

// Item es una línea del carrito: qué producto y cuántas unidades.
type Item struct {
	producto *productos.Producto
	cantidad int
}

// Producto devuelve el producto de la línea.
func (i Item) Producto() *productos.Producto { return i.producto }

// Cantidad devuelve las unidades de la línea.
func (i Item) Cantidad() int { return i.cantidad }

// Subtotal es precio por cantidad, con el precio actual del producto.
func (i Item) Subtotal() float64 { return i.producto.Precio() * float64(i.cantidad) }

// Carrito pertenece a un cliente y guarda sus items en un SLICE, porque aquí sí
// importa el orden: el cliente ve los productos en el orden en que los agregó.
// Un carrito tiene pocas líneas, así que recorrerlo para buscar es barato.
type Carrito struct {
	cliente *clientes.Cliente
	items   []Item
}

// New crea un carrito vacío para el cliente.
func New(cliente *clientes.Cliente) *Carrito {
	return &Carrito{cliente: cliente}
}

// Cliente devuelve el dueño del carrito.
func (c *Carrito) Cliente() *clientes.Cliente { return c.cliente }

// Agregar añade unidades de un producto; si ya estaba, suma la cantidad.
// Valida contra el stock TOTAL que tendría la línea y no solo contra lo nuevo:
// si hay 5 y ya tengo 4 en el carrito, no puedo agregar 2 más.
func (c *Carrito) Agregar(p *productos.Producto, unidades int) error {
	if unidades <= 0 {
		return productos.ErrUnidades
	}
	pos := c.posicion(p.ID())
	total := unidades
	if pos >= 0 {
		total += c.items[pos].cantidad
	}
	if !p.HayStock(total) {
		return fmt.Errorf("%w: solo hay %d de %s", productos.ErrSinStock, p.Stock(), p.Nombre())
	}
	if pos >= 0 {
		c.items[pos].cantidad = total
		return nil
	}
	c.items = append(c.items, Item{producto: p, cantidad: unidades})
	return nil
}

// Quitar elimina un producto del carrito.
func (c *Carrito) Quitar(productoID int) error {
	pos := c.posicion(productoID)
	if pos < 0 {
		return ErrNoEnCarrito
	}
	// Une lo que hay antes y después de pos: así se "salta" el elemento.
	c.items = append(c.items[:pos], c.items[pos+1:]...)
	return nil
}

// Items devuelve una COPIA de las líneas: quien la recibe puede recorrerla,
// pero no modificar el carrito por fuera de sus métodos (encapsulación).
func (c *Carrito) Items() []Item {
	copia := make([]Item, len(c.items))
	copy(copia, c.items)
	return copia
}

// Subtotal suma todas las líneas (sin envío).
func (c *Carrito) Subtotal() float64 {
	total := 0.0
	for _, it := range c.items {
		total += it.Subtotal()
	}
	return total
}

// EstaVacio indica si no hay productos.
func (c *Carrito) EstaVacio() bool { return len(c.items) == 0 }

// Vaciar deja el carrito sin productos (después de confirmar el pedido).
func (c *Carrito) Vaciar() { c.items = nil }

// posicion es interna: devuelve el índice del producto en el slice o -1.
func (c *Carrito) posicion(productoID int) int {
	for i, it := range c.items {
		if it.producto.ID() == productoID {
			return i
		}
	}
	return -1
}
