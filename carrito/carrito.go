// Package carrito guarda lo que cada cliente quiere comprar antes de confirmar el pedido.
package carrito

import (
	"errors"
	"fmt"
	"sort"

	"ecommercego/clientes"
	"ecommercego/productos"
)

// carritos[clienteID][productoID] = unidades que el cliente quiere llevar.
var carritos = map[int]map[int]int{}

// AgregarAlCarrito añade unidades de un producto al carrito de un cliente.
func AgregarAlCarrito(clienteID int, productoID int, unidades int) error {
	if !clientes.Existe(clienteID) {
		return fmt.Errorf("el cliente %d no existe", clienteID)
	}
	_, _, disponible, err := productos.BuscarProducto(productoID)
	if err != nil {
		return err
	}
	if unidades <= 0 {
		return errors.New("las unidades deben ser mayores que 0")
	}
	if carritos[clienteID] == nil {
		carritos[clienteID] = map[int]int{}
	}
	if carritos[clienteID][productoID]+unidades > disponible {
		return fmt.Errorf("solo hay %d unidades disponibles", disponible)
	}
	carritos[clienteID][productoID] += unidades
	return nil
}

// QuitarDelCarrito saca un producto del carrito de un cliente.
func QuitarDelCarrito(clienteID int, productoID int) error {
	if _, ok := carritos[clienteID][productoID]; !ok {
		return errors.New("ese producto no esta en el carrito")
	}
	delete(carritos[clienteID], productoID)
	return nil
}

// VerCarrito devuelve una línea por producto con su subtotal.
func VerCarrito(clienteID int) []string {
	lineas := []string{}
	for _, productoID := range idsOrdenados(clienteID) {
		unidades := carritos[clienteID][productoID]
		nombre, precio, _, err := productos.BuscarProducto(productoID)
		if err != nil {
			lineas = append(lineas, fmt.Sprintf("- producto %d (ya no existe)", productoID))
			continue
		}
		lineas = append(lineas, fmt.Sprintf("- %s x%d = $%.2f", nombre, unidades, precio*float64(unidades)))
	}
	return lineas
}

// CalcularTotal suma los subtotales del carrito (sin envío).
func CalcularTotal(clienteID int) (float64, error) {
	if len(carritos[clienteID]) == 0 {
		return 0, errors.New("el carrito esta vacio")
	}
	total := 0.0
	for productoID, unidades := range carritos[clienteID] {
		_, precio, _, err := productos.BuscarProducto(productoID)
		if err != nil {
			return 0, err
		}
		total += precio * float64(unidades)
	}
	return total, nil
}

// Contenido devuelve una copia del carrito para que otro paquete (pedidos)
// pueda leerlo sin modificar el original.
func Contenido(clienteID int) map[int]int {
	copia := map[int]int{}
	for productoID, unidades := range carritos[clienteID] {
		copia[productoID] = unidades
	}
	return copia
}

// VaciarCarrito deja el carrito del cliente sin productos.
func VaciarCarrito(clienteID int) {
	delete(carritos, clienteID)
}

// idsOrdenados devuelve los IDs de producto del carrito en orden, porque un map
// no garantiza el orden al recorrerlo.
func idsOrdenados(clienteID int) []int {
	ids := []int{}
	for productoID := range carritos[clienteID] {
		ids = append(ids, productoID)
	}
	sort.Ints(ids)
	return ids
}
