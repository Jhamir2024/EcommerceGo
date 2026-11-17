// Package productos modela los artículos de la tienda (struct Producto) y el
// catálogo que los guarda (struct Catalogo).
package productos

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Errores que otros paquetes pueden comparar con errors.Is.
// Se declaran una sola vez para no repetir el texto en cada función.
var (
	ErrNoExiste      = errors.New("el producto no existe")
	ErrPrecio        = errors.New("el precio debe ser mayor que 0")
	ErrNombreVacio   = errors.New("el nombre no puede estar vacio")
	ErrStockNegativo = errors.New("el stock no puede ser negativo")
	ErrSinStock      = errors.New("no hay stock suficiente")
	ErrUnidades      = errors.New("las unidades deben ser mayores que 0")
)

// Producto es un artículo del catálogo.
// Sus campos están en minúscula (no exportados): desde otro paquete solo se
// pueden leer con los getters y cambiar con los setters, que validan.
// Así nadie puede dejar un precio negativo o un stock de -5.
type Producto struct {
	id     int
	nombre string
	precio float64
	stock  int
}

// NewProducto crea un producto válido o devuelve el error que lo impide.
// Devuelve un puntero porque el catálogo, el carrito y los pedidos deben
// compartir el MISMO producto: si se vende, todos ven el stock actualizado.
func NewProducto(id int, nombre string, precio float64, stock int) (*Producto, error) {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return nil, ErrNombreVacio
	}
	if precio <= 0 {
		return nil, ErrPrecio
	}
	if stock < 0 {
		return nil, ErrStockNegativo
	}
	return &Producto{id: id, nombre: nombre, precio: precio, stock: stock}, nil
}

// Getters: receptor de valor porque solo leen; trabajar con una copia no afecta.

// ID devuelve el identificador del producto.
func (p Producto) ID() int { return p.id }

// Nombre devuelve el nombre del producto.
func (p Producto) Nombre() string { return p.nombre }

// Precio devuelve el precio unitario.
func (p Producto) Precio() float64 { return p.precio }

// Stock devuelve las unidades disponibles.
func (p Producto) Stock() int { return p.stock }

// HayStock indica si se pueden vender esas unidades.
func (p Producto) HayStock(unidades int) bool { return unidades > 0 && unidades <= p.stock }

// String hace que Producto cumpla la interfaz fmt.Stringer: fmt.Println(p)
// imprime este texto en lugar de los campos internos.
func (p Producto) String() string {
	return fmt.Sprintf("#%d %-22s $%7.2f  stock: %d", p.id, p.nombre, p.precio, p.stock)
}

// Setters y métodos que modifican: receptor puntero, porque deben cambiar el
// producto original y no una copia. Todos validan antes de modificar.

// SetPrecio cambia el precio si el nuevo valor es válido.
func (p *Producto) SetPrecio(nuevo float64) error {
	if nuevo <= 0 {
		return ErrPrecio
	}
	p.precio = nuevo
	return nil
}

// SetNombre cambia el nombre si no está vacío.
func (p *Producto) SetNombre(nuevo string) error {
	nuevo = strings.TrimSpace(nuevo)
	if nuevo == "" {
		return ErrNombreVacio
	}
	p.nombre = nuevo
	return nil
}

// DescontarStock resta unidades al vender. Falla si no alcanzan.
func (p *Producto) DescontarStock(unidades int) error {
	if unidades <= 0 {
		return ErrUnidades
	}
	if unidades > p.stock {
		// %w envuelve el error base: el mensaje lleva el detalle y errors.Is
		// sigue reconociendo ErrSinStock.
		return fmt.Errorf("%w de %s: hay %d, se piden %d", ErrSinStock, p.nombre, p.stock, unidades)
	}
	p.stock -= unidades
	return nil
}

// DevolverStock suma unidades cuando se cancela un pedido.
func (p *Producto) DevolverStock(unidades int) error {
	if unidades <= 0 {
		return ErrUnidades
	}
	p.stock += unidades
	return nil
}

// Catalogo guarda los productos de la tienda.
// Se usa un map por ID porque la operación más frecuente es buscar un
// producto por su ID (carrito, pedidos, menú): con un map es directa.
type Catalogo struct {
	productos   map[int]*Producto
	siguienteID int
}

// NewCatalogo crea un catálogo vacío listo para usar (el map ya inicializado).
func NewCatalogo() *Catalogo {
	return &Catalogo{productos: make(map[int]*Producto), siguienteID: 1}
}

// Agregar crea un producto con el siguiente ID libre y lo guarda.
func (c *Catalogo) Agregar(nombre string, precio float64, stock int) (*Producto, error) {
	p, err := NewProducto(c.siguienteID, nombre, precio, stock)
	if err != nil {
		return nil, err // el ID no se consume si los datos son inválidos
	}
	c.productos[p.id] = p
	c.siguienteID++
	return p, nil
}

// Buscar devuelve el producto con ese ID o ErrNoExiste.
func (c *Catalogo) Buscar(id int) (*Producto, error) {
	p, ok := c.productos[id]
	if !ok {
		return nil, fmt.Errorf("%w: ID %d", ErrNoExiste, id)
	}
	return p, nil
}

// Listar devuelve los productos ordenados por ID.
// Los maps de Go no tienen orden fijo, por eso se pasan a un slice y se ordenan.
func (c *Catalogo) Listar() []*Producto {
	lista := make([]*Producto, 0, len(c.productos))
	for _, p := range c.productos {
		lista = append(lista, p)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].id < lista[j].id })
	return lista
}

// Eliminar quita un producto del catálogo.
func (c *Catalogo) Eliminar(id int) error {
	if _, err := c.Buscar(id); err != nil {
		return err
	}
	delete(c.productos, id)
	return nil
}
