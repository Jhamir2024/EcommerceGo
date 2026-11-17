package pedidos

import (
	"errors"
	"testing"

	"ecommercego/carrito"
	"ecommercego/clientes"
	"ecommercego/productos"
)

// La misma llamada Costo() da resultados distintos según la implementación.
func TestMetodosEnvioPolimorfismo(t *testing.T) {
	casos := []struct {
		envio    MetodoEnvio
		subtotal float64
		esperado float64
	}{
		{EnvioEstandar{}, 30, 3.50},
		{EnvioEstandar{}, 50, 0},
		{EnvioExpress{}, 80, 6.00},
		{RetiroEnTienda{}, 10, 0},
	}
	for _, c := range casos {
		if got := c.envio.Costo(c.subtotal); got != c.esperado {
			t.Errorf("%s con $%.2f: esperado $%.2f, llego $%.2f", c.envio.Nombre(), c.subtotal, c.esperado, got)
		}
	}
}

func preparar(t *testing.T) (*carrito.Carrito, *productos.Producto, *productos.Producto) {
	t.Helper()
	cli, _ := clientes.NewCliente(1, "Ana", "ana@correo.com")
	a, _ := productos.NewProducto(1, "Audifonos", 25, 10)
	b, _ := productos.NewProducto(2, "Teclado", 45.90, 1)
	return carrito.New(cli), a, b
}

func TestConfirmarYCancelar(t *testing.T) {
	c, a, _ := preparar(t)
	if err := c.Agregar(a, 2); err != nil {
		t.Fatal(err)
	}
	g := NewGestor()
	p, err := g.Confirmar(c, EnvioExpress{})
	if err != nil {
		t.Fatal(err)
	}
	if p.Total() != 56 || a.Stock() != 8 || !c.EstaVacio() {
		t.Fatalf("total %.2f, stock %d, carrito vacio %v", p.Total(), a.Stock(), c.EstaVacio())
	}
	if err := g.Cancelar(p.ID()); err != nil {
		t.Fatal(err)
	}
	if a.Stock() != 10 || g.TotalVendido() != 0 {
		t.Errorf("tras cancelar: stock %d, vendido %.2f", a.Stock(), g.TotalVendido())
	}
	if err := g.Cancelar(p.ID()); !errors.Is(err, ErrYaCancelado) {
		t.Errorf("se esperaba ErrYaCancelado, llego %v", err)
	}
}

// Si un producto ya no tiene stock, no se descuenta NINGUNO (validación en dos pasadas).
func TestConfirmarSinStockNoDescuentaNada(t *testing.T) {
	c, a, b := preparar(t)
	_ = c.Agregar(a, 2)
	_ = c.Agregar(b, 1)
	_ = b.DescontarStock(1) // otro cliente compró el último teclado
	if _, err := NewGestor().Confirmar(c, EnvioEstandar{}); !errors.Is(err, productos.ErrSinStock) {
		t.Fatalf("se esperaba ErrSinStock, llego %v", err)
	}
	if a.Stock() != 10 {
		t.Errorf("los audifonos no debian descontarse, stock %d", a.Stock())
	}
}

func TestConfirmarCarritoVacio(t *testing.T) {
	c, _, _ := preparar(t)
	if _, err := NewGestor().Confirmar(c, EnvioEstandar{}); !errors.Is(err, ErrCarritoVacio) {
		t.Errorf("se esperaba ErrCarritoVacio, llego %v", err)
	}
}
