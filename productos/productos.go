// Package productos gestiona el catálogo de la tienda: registrar productos,
// consultarlos, cambiar su precio, eliminarlos y controlar su stock.
package productos

import (
	"errors"
	"fmt"
)

// Datos del catálogo. En esta etapa (unidad 1) cada dato se guarda en un map
// cuya clave es el ID del producto. En el Aprendizaje Autónomo 2 estos maps se
// reemplazarán por un struct Producto con sus métodos.
var (
	nombres   = map[int]string{}
	precios   = map[int]float64{}
	stock     = map[int]int{}
	siguiente = 1 // ID que recibirá el próximo producto
)

// AgregarProducto registra un producto nuevo y devuelve su ID.
func AgregarProducto(nombre string, precio float64, unidades int) (int, error) {
	if nombre == "" {
		return 0, errors.New("el nombre no puede estar vacio")
	}
	if !precioValido(precio) {
		return 0, errors.New("el precio debe ser mayor que 0")
	}
	if unidades < 0 {
		return 0, errors.New("el stock no puede ser negativo")
	}
	id := siguiente
	nombres[id], precios[id], stock[id] = nombre, precio, unidades
	siguiente++
	return id, nil
}

// ListarProductos devuelve una línea de texto por cada producto del catálogo,
// ordenadas por ID.
func ListarProductos() []string {
	lista := []string{}
	for id := 1; id < siguiente; id++ {
		if Existe(id) {
			lista = append(lista, fmt.Sprintf("%d) %s - $%.2f - stock %d", id, nombres[id], precios[id], stock[id]))
		}
	}
	return lista
}

// BuscarProducto devuelve el nombre, el precio y el stock de un producto.
func BuscarProducto(id int) (string, float64, int, error) {
	if !Existe(id) {
		return "", 0, 0, errNoExiste(id)
	}
	return nombres[id], precios[id], stock[id], nil
}

// ActualizarPrecio cambia el precio de un producto existente.
func ActualizarPrecio(id int, nuevoPrecio float64) error {
	if !Existe(id) {
		return errNoExiste(id)
	}
	if !precioValido(nuevoPrecio) {
		return errors.New("el precio debe ser mayor que 0")
	}
	precios[id] = nuevoPrecio
	return nil
}

// EliminarProducto quita un producto del catálogo.
func EliminarProducto(id int) error {
	if !Existe(id) {
		return errNoExiste(id)
	}
	delete(nombres, id)
	delete(precios, id)
	delete(stock, id)
	return nil
}

// Existe indica si hay un producto con ese ID.
func Existe(id int) bool {
	_, ok := nombres[id]
	return ok
}

// DescontarStock resta unidades del stock (se usa al confirmar un pedido).
func DescontarStock(id int, unidades int) error {
	if !Existe(id) {
		return errNoExiste(id)
	}
	if unidades <= 0 || unidades > stock[id] {
		return fmt.Errorf("stock insuficiente de %s: hay %d", nombres[id], stock[id])
	}
	stock[id] -= unidades
	return nil
}

// DevolverStock suma unidades al stock (se usa al cancelar un pedido).
func DevolverStock(id int, unidades int) error {
	if !Existe(id) {
		return errNoExiste(id)
	}
	stock[id] += unidades
	return nil
}

// precioValido es una auxiliar interna (minúscula): no se puede usar fuera del paquete.
func precioValido(precio float64) bool {
	return precio > 0
}

// errNoExiste arma el mensaje de error para un ID que no está en el catálogo.
func errNoExiste(id int) error {
	return fmt.Errorf("el producto %d no existe", id)
}
