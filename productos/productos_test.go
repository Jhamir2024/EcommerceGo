package productos

import (
	"errors"
	"testing"
)

func TestNewProductoValida(t *testing.T) {
	if _, err := NewProducto(1, "", 10, 1); !errors.Is(err, ErrNombreVacio) {
		t.Errorf("nombre vacio: se esperaba ErrNombreVacio, llego %v", err)
	}
	if _, err := NewProducto(1, "Mouse", -1, 1); !errors.Is(err, ErrPrecio) {
		t.Errorf("precio negativo: se esperaba ErrPrecio, llego %v", err)
	}
	if _, err := NewProducto(1, "Mouse", 10, -2); !errors.Is(err, ErrStockNegativo) {
		t.Errorf("stock negativo: se esperaba ErrStockNegativo, llego %v", err)
	}
}

func TestSetPrecioNoAceptaInvalido(t *testing.T) {
	p, _ := NewProducto(1, "Mouse", 10, 1)
	if err := p.SetPrecio(0); !errors.Is(err, ErrPrecio) {
		t.Fatalf("se esperaba ErrPrecio, llego %v", err)
	}
	if p.Precio() != 10 {
		t.Errorf("el precio no debia cambiar, quedo en %.2f", p.Precio())
	}
}

func TestDescontarStock(t *testing.T) {
	p, _ := NewProducto(1, "Teclado", 45.90, 5)
	if err := p.DescontarStock(10); !errors.Is(err, ErrSinStock) {
		t.Fatalf("se esperaba ErrSinStock, llego %v", err)
	}
	if err := p.DescontarStock(2); err != nil {
		t.Fatal(err)
	}
	if p.Stock() != 3 {
		t.Errorf("stock esperado 3, quedo %d", p.Stock())
	}
}

func TestCatalogoBuscarYEliminar(t *testing.T) {
	c := NewCatalogo()
	p, _ := c.Agregar("Mouse", 12.5, 20)
	if _, err := c.Buscar(p.ID()); err != nil {
		t.Fatal(err)
	}
	if err := c.Eliminar(p.ID()); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Buscar(p.ID()); !errors.Is(err, ErrNoExiste) {
		t.Errorf("se esperaba ErrNoExiste, llego %v", err)
	}
}
