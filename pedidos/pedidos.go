// Package pedidos convierte un carrito en una compra (struct Pedido) y lleva el
// historial de ventas (struct Gestor).
package pedidos

import (
	"errors"
	"fmt"

	"ecommercego/carrito"
	"ecommercego/clientes"
	"ecommercego/productos"
)

// Errores del paquete.
var (
	ErrNoExiste     = errors.New("el pedido no existe")
	ErrCarritoVacio = errors.New("el carrito esta vacio")
	ErrSinEnvio     = errors.New("debe elegir un metodo de envio")
	ErrYaCancelado  = errors.New("el pedido ya estaba cancelado")
)

// Estado es un tipo propio (y no un string suelto) para que solo se usen los
// valores definidos abajo.
type Estado string

// Estados posibles de un pedido.
const (
	Confirmado Estado = "confirmado"
	Cancelado  Estado = "cancelado"
)

// Linea guarda una COPIA del nombre y del precio al momento de comprar.
// Si mañana el producto cambia de precio, el pedido histórico no debe cambiar.
// También guarda el puntero al producto para poder devolver el stock.
type Linea struct {
	producto       *productos.Producto
	nombre         string
	precioUnitario float64
	cantidad       int
}

// Nombre del producto al momento de la compra.
func (l Linea) Nombre() string { return l.nombre }

// PrecioUnitario al momento de la compra.
func (l Linea) PrecioUnitario() float64 { return l.precioUnitario }

// Cantidad comprada.
func (l Linea) Cantidad() int { return l.cantidad }

// Subtotal de la línea.
func (l Linea) Subtotal() float64 { return l.precioUnitario * float64(l.cantidad) }

// Pedido es una compra confirmada. El método de envío se guarda como la
// interfaz MetodoEnvio: el pedido no sabe ni le importa cuál es.
type Pedido struct {
	id      int
	cliente *clientes.Cliente
	lineas  []Linea
	envio   MetodoEnvio
	estado  Estado
}

// ID del pedido.
func (p *Pedido) ID() int { return p.id }

// Cliente que hizo el pedido.
func (p *Pedido) Cliente() *clientes.Cliente { return p.cliente }

// Envio elegido.
func (p *Pedido) Envio() MetodoEnvio { return p.envio }

// Estado actual.
func (p *Pedido) Estado() Estado { return p.estado }

// Lineas devuelve una copia de las líneas (encapsulación).
func (p *Pedido) Lineas() []Linea {
	copia := make([]Linea, len(p.lineas))
	copy(copia, p.lineas)
	return copia
}

// Subtotal de los productos, sin envío.
func (p *Pedido) Subtotal() float64 {
	total := 0.0
	for _, l := range p.lineas {
		total += l.Subtotal()
	}
	return total
}

// CostoEnvio delega en la interfaz: aquí ocurre el polimorfismo.
func (p *Pedido) CostoEnvio() float64 { return p.envio.Costo(p.Subtotal()) }

// Total a pagar: productos + envío.
func (p *Pedido) Total() float64 { return p.Subtotal() + p.CostoEnvio() }

// String implementa fmt.Stringer.
func (p *Pedido) String() string {
	return fmt.Sprintf("Pedido %d | %s | %d producto(s) | envio %s $%.2f | total $%.2f | %s",
		p.id, p.cliente.Nombre(), len(p.lineas), p.envio.Nombre(), p.CostoEnvio(), p.Total(), p.estado)
}

// cancelar es interna: solo el Gestor puede cancelar, para que el historial
// y el stock se mantengan coherentes.
func (p *Pedido) cancelar() error {
	if p.estado == Cancelado {
		return ErrYaCancelado
	}
	for _, l := range p.lineas {
		if err := l.producto.DevolverStock(l.cantidad); err != nil {
			return err
		}
	}
	p.estado = Cancelado
	return nil
}

// Gestor lleva el historial de pedidos en un SLICE: el orden cronológico
// importa (el pedido 1 va antes que el 2) y siempre se listan todos.
type Gestor struct {
	pedidos     []*Pedido
	siguienteID int
}

// NewGestor crea un gestor sin pedidos.
func NewGestor() *Gestor {
	return &Gestor{siguienteID: 1}
}

// Confirmar convierte el carrito en un pedido con el método de envío elegido.
//
// Se hace en DOS PASADAS a propósito:
//  1. revisar que haya stock de TODOS los productos;
//  2. recién entonces descontar.
//
// Si se descontara mientras se revisa y el tercer producto no tuviera stock,
// los dos primeros ya habrían perdido unidades sin que exista pedido.
func (g *Gestor) Confirmar(c *carrito.Carrito, envio MetodoEnvio) (*Pedido, error) {
	if c.EstaVacio() {
		return nil, ErrCarritoVacio
	}
	if envio == nil {
		return nil, ErrSinEnvio
	}
	items := c.Items()

	// Pasada 1: solo validar.
	for _, it := range items {
		if !it.Producto().HayStock(it.Cantidad()) {
			return nil, fmt.Errorf("%w: %s (hay %d, se piden %d)", productos.ErrSinStock,
				it.Producto().Nombre(), it.Producto().Stock(), it.Cantidad())
		}
	}

	// Pasada 2: descontar y copiar los datos de cada línea.
	lineas := make([]Linea, 0, len(items))
	for _, it := range items {
		p := it.Producto()
		if err := p.DescontarStock(it.Cantidad()); err != nil {
			return nil, err // no debería pasar después de la pasada 1
		}
		lineas = append(lineas, Linea{producto: p, nombre: p.Nombre(),
			precioUnitario: p.Precio(), cantidad: it.Cantidad()})
	}

	pedido := &Pedido{id: g.siguienteID, cliente: c.Cliente(), lineas: lineas,
		envio: envio, estado: Confirmado}
	g.pedidos = append(g.pedidos, pedido)
	g.siguienteID++
	c.Vaciar()
	return pedido, nil
}

// Buscar devuelve el pedido con ese ID.
func (g *Gestor) Buscar(id int) (*Pedido, error) {
	for _, p := range g.pedidos {
		if p.id == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("%w: ID %d", ErrNoExiste, id)
}

// Cancelar anula un pedido confirmado y devuelve el stock de sus productos.
func (g *Gestor) Cancelar(id int) error {
	p, err := g.Buscar(id)
	if err != nil {
		return err
	}
	return p.cancelar()
}

// Listar devuelve una copia del historial en orden cronológico.
func (g *Gestor) Listar() []*Pedido {
	copia := make([]*Pedido, len(g.pedidos))
	copy(copia, g.pedidos)
	return copia
}

// TotalVendido suma solo los pedidos confirmados (los cancelados no cuentan).
func (g *Gestor) TotalVendido() float64 {
	total := 0.0
	for _, p := range g.pedidos {
		if p.estado == Confirmado {
			total += p.Total()
		}
	}
	return total
}
