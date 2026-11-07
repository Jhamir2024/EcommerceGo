// Sistema de Gestión de E-commerce · Programación Orientada a Objetos (Go) · UIDE
// main solo muestra el menú y pide datos; la lógica vive en los paquetes.
// Ejecutar desde esta carpeta: go run .
package main

import (
	"fmt"

	"ecommercego/carrito"
	"ecommercego/clientes"
	"ecommercego/entrada"
	"ecommercego/pedidos"
	"ecommercego/productos"
)

func main() {
	cargarDatosDemo()

	for {
		mostrarMenu()
		opcion, err := entrada.LeerEntero("Elige: ")
		if err != nil {
			fmt.Println("Escribe un numero del menu")
			continue
		}

		switch opcion {
		case 1:
			opcionAgregarProducto()
		case 2:
			imprimir("CATALOGO", productos.ListarProductos())
		case 3:
			opcionBuscarProducto()
		case 4:
			opcionActualizarPrecio()
		case 5:
			opcionEliminarProducto()
		case 6:
			opcionRegistrarCliente()
		case 7:
			imprimir("CLIENTES", clientes.ListarClientes())
		case 8:
			opcionAgregarAlCarrito()
		case 9:
			opcionVerCarrito()
		case 10:
			opcionQuitarDelCarrito()
		case 11:
			opcionConfirmarPedido()
		case 12:
			imprimir("PEDIDOS", pedidos.ListarPedidos())
			fmt.Printf("Total vendido: $%.2f\n", pedidos.TotalVendido())
		case 13:
			opcionCancelarPedido()
		case 0:
			fmt.Println("Gracias, hasta pronto")
			return
		default:
			fmt.Println("Esa opcion no existe")
		}
	}
}

func mostrarMenu() {
	fmt.Println("\n========== TIENDA GO ==========")
	fmt.Println(" PRODUCTOS   1) Agregar   2) Ver catalogo   3) Buscar")
	fmt.Println("             4) Cambiar precio   5) Eliminar")
	fmt.Println(" CLIENTES    6) Registrar   7) Ver clientes")
	fmt.Println(" CARRITO     8) Agregar   9) Ver carrito   10) Quitar")
	fmt.Println(" PEDIDOS    11) Confirmar  12) Ver pedidos  13) Cancelar")
	fmt.Println("             0) Salir")
}

// imprimir muestra un título y una lista de líneas (o un aviso si está vacía).
func imprimir(titulo string, lineas []string) {
	fmt.Println("---", titulo, "---")
	if len(lineas) == 0 {
		fmt.Println("(no hay datos todavia)")
	}
	for _, l := range lineas {
		fmt.Println(l)
	}
}

// mostrarResultado imprime el error si lo hay; si no, el mensaje de éxito.
func mostrarResultado(err error, exito string) {
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(exito)
}

func opcionAgregarProducto() {
	nombre := entrada.LeerTexto("Nombre: ")
	precio, err1 := entrada.LeerDecimal("Precio: ")
	stock, err2 := entrada.LeerEntero("Stock: ")
	if err1 != nil || err2 != nil {
		fmt.Println("Error: precio o stock no son numeros")
		return
	}
	id, err := productos.AgregarProducto(nombre, precio, stock)
	mostrarResultado(err, fmt.Sprintf("Producto agregado con ID %d", id))
}

func opcionBuscarProducto() {
	id, _ := entrada.LeerEntero("ID del producto: ")
	nombre, precio, stock, err := productos.BuscarProducto(id)
	mostrarResultado(err, fmt.Sprintf("%s - $%.2f - stock %d", nombre, precio, stock))
}

func opcionActualizarPrecio() {
	id, _ := entrada.LeerEntero("ID del producto: ")
	precio, err := entrada.LeerDecimal("Nuevo precio: ")
	if err != nil {
		fmt.Println("Error: el precio no es un numero")
		return
	}
	mostrarResultado(productos.ActualizarPrecio(id, precio), "Precio actualizado")
}

func opcionEliminarProducto() {
	id, _ := entrada.LeerEntero("ID del producto: ")
	mostrarResultado(productos.EliminarProducto(id), "Producto eliminado")
}

func opcionRegistrarCliente() {
	nombre := entrada.LeerTexto("Nombre: ")
	correo := entrada.LeerTexto("Correo: ")
	id, err := clientes.RegistrarCliente(nombre, correo)
	mostrarResultado(err, fmt.Sprintf("Cliente registrado con ID %d", id))
}

func opcionAgregarAlCarrito() {
	cliente, _ := entrada.LeerEntero("ID del cliente: ")
	producto, _ := entrada.LeerEntero("ID del producto: ")
	unidades, _ := entrada.LeerEntero("Unidades: ")
	mostrarResultado(carrito.AgregarAlCarrito(cliente, producto, unidades), "Agregado al carrito")
}

func opcionVerCarrito() {
	cliente, _ := entrada.LeerEntero("ID del cliente: ")
	imprimir("CARRITO", carrito.VerCarrito(cliente))
	if total, err := carrito.CalcularTotal(cliente); err == nil {
		fmt.Printf("Subtotal: $%.2f (envio gratis desde $50)\n", total)
	}
}

func opcionQuitarDelCarrito() {
	cliente, _ := entrada.LeerEntero("ID del cliente: ")
	producto, _ := entrada.LeerEntero("ID del producto: ")
	mostrarResultado(carrito.QuitarDelCarrito(cliente, producto), "Producto quitado del carrito")
}

func opcionConfirmarPedido() {
	cliente, _ := entrada.LeerEntero("ID del cliente: ")
	id, total, err := pedidos.ConfirmarPedido(cliente)
	mostrarResultado(err, fmt.Sprintf("Pedido %d confirmado. Total a pagar: $%.2f", id, total))
}

func opcionCancelarPedido() {
	id, _ := entrada.LeerEntero("ID del pedido: ")
	mostrarResultado(pedidos.CancelarPedido(id), "Pedido cancelado y stock devuelto")
}

// cargarDatosDemo agrega algunos datos para probar el sistema sin escribir todo a mano.
func cargarDatosDemo() {
	productos.AgregarProducto("Audifonos Bluetooth", 25.00, 10)
	productos.AgregarProducto("Mouse inalambrico", 12.50, 20)
	productos.AgregarProducto("Teclado mecanico", 45.90, 5)
	clientes.RegistrarCliente("Ana Torres", "ana@correo.com")
}
