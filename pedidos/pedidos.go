// Package pedidos convierte el carrito en una compra: descuenta stock, calcula
// el envío, guarda el pedido y permite cancelarlo y ver el total vendido.
package pedidos

import (
	"errors"
	"fmt"

	"ecommercego/carrito"
	"ecommercego/clientes"
	"ecommercego/productos"
)

const (
	montoEnvioGratis = 50.0 // desde este subtotal el envío es gratis
	costoEnvio       = 3.50
)

// Datos de los pedidos, guardados por ID de pedido. En el AA2 serán un struct Pedido.
var (
	clienteDe = map[int]int{}         // pedidoID -> clienteID
	totales   = map[int]float64{}     // pedidoID -> total pagado (con envío)
	estados   = map[int]string{}      // pedidoID -> "confirmado" o "cancelado"
	detalles  = map[int]map[int]int{} // pedidoID -> productoID -> unidades
	siguiente = 1
)

// ConfirmarPedido compra todo lo que hay en el carrito del cliente.
// Devuelve el ID del pedido y el total a pagar (incluye envío).
func ConfirmarPedido(clienteID int) (int, float64, error) {
	if !clientes.Existe(clienteID) {
		return 0, 0, fmt.Errorf("el cliente %d no existe", clienteID)
	}
	subtotal, err := carrito.CalcularTotal(clienteID)
	if err != nil {
		return 0, 0, err
	}
	items := carrito.Contenido(clienteID)

	// Primero se revisa el stock de TODOS los productos y recién después se
	// descuenta: así un pedido nunca queda a medias.
	for productoID, unidades := range items {
		nombre, _, disponible, err := productos.BuscarProducto(productoID)
		if err != nil {
			return 0, 0, err
		}
		if unidades > disponible {
			return 0, 0, fmt.Errorf("no hay stock suficiente de %s", nombre)
		}
	}
	for productoID, unidades := range items {
		productos.DescontarStock(productoID, unidades) // ya se validó arriba
	}

	id := siguiente
	clienteDe[id] = clienteID
	totales[id] = subtotal + calcularEnvio(subtotal)
	estados[id] = "confirmado"
	detalles[id] = items
	siguiente++

	carrito.VaciarCarrito(clienteID)
	return id, totales[id], nil
}

// CancelarPedido marca el pedido como cancelado y devuelve el stock.
func CancelarPedido(pedidoID int) error {
	estado, ok := estados[pedidoID]
	if !ok {
		return fmt.Errorf("el pedido %d no existe", pedidoID)
	}
	if estado == "cancelado" {
		return errors.New("el pedido ya estaba cancelado")
	}
	for productoID, unidades := range detalles[pedidoID] {
		// Si el producto fue eliminado del catálogo no hay stock que devolver.
		productos.DevolverStock(productoID, unidades)
	}
	estados[pedidoID] = "cancelado"
	return nil
}

// ListarPedidos devuelve una línea por pedido, ordenadas por ID.
func ListarPedidos() []string {
	lista := []string{}
	for id := 1; id < siguiente; id++ {
		nombre, _, err := clientes.BuscarCliente(clienteDe[id])
		if err != nil {
			nombre = "cliente desconocido"
		}
		lista = append(lista, fmt.Sprintf("Pedido %d - %s - $%.2f - %s", id, nombre, totales[id], estados[id]))
	}
	return lista
}

// TotalVendido suma el dinero de los pedidos confirmados (no cuenta los cancelados).
func TotalVendido() float64 {
	total := 0.0
	for id, estado := range estados {
		if estado == "confirmado" {
			total += totales[id]
		}
	}
	return total
}

// calcularEnvio aplica la regla de la tienda: envío gratis desde $50.
func calcularEnvio(subtotal float64) float64 {
	if subtotal >= montoEnvioGratis {
		return 0
	}
	return costoEnvio
}
