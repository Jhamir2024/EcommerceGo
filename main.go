// Sistema de Gestión de E-commerce · Programación Orientada a Objetos (Go) · UIDE
// main solo muestra el menú y pide datos; la lógica vive en los structs de cada paquete.
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

// Tienda reúne todas las partes del sistema. En el AA1 eran variables sueltas
// de cada paquete; ahora son campos de un struct, así podría haber dos tiendas
// independientes (útil para las pruebas y para los servicios web de la ECD).
type Tienda struct {
	catalogo *productos.Catalogo
	clientes *clientes.Registro
	carritos map[int]*carrito.Carrito // un carrito por ID de cliente
	pedidos  *pedidos.Gestor
}

// NewTienda crea una tienda vacía con todo inicializado.
func NewTienda() *Tienda {
	return &Tienda{
		catalogo: productos.NewCatalogo(),
		clientes: clientes.NewRegistro(),
		carritos: make(map[int]*carrito.Carrito),
		pedidos:  pedidos.NewGestor(),
	}
}

// carritoDe devuelve el carrito del cliente; si aún no tiene uno, lo crea.
// Primero comprueba que el cliente exista, para no crear carritos huérfanos.
func (t *Tienda) carritoDe(clienteID int) (*carrito.Carrito, error) {
	cli, err := t.clientes.Buscar(clienteID)
	if err != nil {
		return nil, err
	}
	c, ok := t.carritos[clienteID]
	if !ok {
		c = carrito.New(cli)
		t.carritos[clienteID] = c
	}
	return c, nil
}

func main() {
	t := NewTienda()
	t.cargarDatosDemo()

	for {
		mostrarMenu()
		opcion, err := entrada.LeerEntero("Elige: ")
		if err != nil {
			fmt.Println("Escribe un numero del menu")
			continue
		}

		switch opcion {
		case 1:
			t.opcionAgregarProducto()
		case 2:
			t.opcionVerCatalogo()
		case 3:
			t.opcionBuscarProducto()
		case 4:
			t.opcionCambiarPrecio()
		case 5:
			t.opcionEliminarProducto()
		case 6:
			t.opcionRegistrarCliente()
		case 7:
			t.opcionVerClientes()
		case 8:
			t.opcionAgregarAlCarrito()
		case 9:
			t.opcionVerCarrito()
		case 10:
			t.opcionQuitarDelCarrito()
		case 11:
			t.opcionConfirmarPedido()
		case 12:
			t.opcionVerPedidos()
		case 13:
			t.opcionCancelarPedido()
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

// mostrarResultado imprime el error si lo hay; si no, el mensaje de éxito.
// Centraliza el manejo: ningún error del sistema se ignora en el menú.
func mostrarResultado(err error, exito string) {
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(exito)
}

// leerNumero pide un entero y avisa si lo escrito no es un número.
// Devuelve ok=false para que la opción se cancele sin seguir pidiendo datos.
func leerNumero(mensaje string) (int, bool) {
	n, err := entrada.LeerEntero(mensaje)
	if err != nil {
		fmt.Println("Error: debes escribir un numero entero")
		return 0, false
	}
	return n, true
}

// ---------------- PRODUCTOS ----------------

func (t *Tienda) opcionAgregarProducto() {
	nombre := entrada.LeerTexto("Nombre: ")
	precio, err := entrada.LeerDecimal("Precio: ")
	if err != nil {
		fmt.Println("Error: el precio debe ser un numero")
		return
	}
	stock, ok := leerNumero("Stock: ")
	if !ok {
		return
	}
	p, err := t.catalogo.Agregar(nombre, precio, stock)
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	fmt.Println("Producto agregado:", p) // usa Producto.String()
}

func (t *Tienda) opcionVerCatalogo() {
	fmt.Println("--- CATALOGO ---")
	lista := t.catalogo.Listar()
	if len(lista) == 0 {
		fmt.Println("(no hay productos)")
	}
	for _, p := range lista {
		fmt.Println(p)
	}
}

func (t *Tienda) opcionBuscarProducto() {
	id, ok := leerNumero("ID del producto: ")
	if !ok {
		return
	}
	p, err := t.catalogo.Buscar(id)
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	fmt.Println(p)
}

func (t *Tienda) opcionCambiarPrecio() {
	id, ok := leerNumero("ID del producto: ")
	if !ok {
		return
	}
	p, err := t.catalogo.Buscar(id)
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	precio, err := entrada.LeerDecimal("Nuevo precio: ")
	if err != nil {
		fmt.Println("Error: el precio debe ser un numero")
		return
	}
	// El setter valida: un precio de 0 o negativo devuelve error y no se aplica.
	mostrarResultado(p.SetPrecio(precio), "Precio actualizado: "+p.String())
}

func (t *Tienda) opcionEliminarProducto() {
	id, ok := leerNumero("ID del producto: ")
	if !ok {
		return
	}
	mostrarResultado(t.catalogo.Eliminar(id), "Producto eliminado")
}

// ---------------- CLIENTES ----------------

func (t *Tienda) opcionRegistrarCliente() {
	nombre := entrada.LeerTexto("Nombre: ")
	correo := entrada.LeerTexto("Correo: ")
	c, err := t.clientes.Registrar(nombre, correo)
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	fmt.Println("Cliente registrado:", c)
}

func (t *Tienda) opcionVerClientes() {
	fmt.Println("--- CLIENTES ---")
	for _, c := range t.clientes.Listar() {
		fmt.Println(c)
	}
}

// ---------------- CARRITO ----------------

func (t *Tienda) opcionAgregarAlCarrito() {
	clienteID, ok := leerNumero("ID del cliente: ")
	if !ok {
		return
	}
	c, err := t.carritoDe(clienteID)
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	productoID, ok := leerNumero("ID del producto: ")
	if !ok {
		return
	}
	p, err := t.catalogo.Buscar(productoID)
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	unidades, ok := leerNumero("Unidades: ")
	if !ok {
		return
	}
	mostrarResultado(c.Agregar(p, unidades), "Agregado al carrito")
}

func (t *Tienda) opcionVerCarrito() {
	clienteID, ok := leerNumero("ID del cliente: ")
	if !ok {
		return
	}
	c, err := t.carritoDe(clienteID)
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	fmt.Println("--- CARRITO de", c.Cliente().Nombre(), "---")
	if c.EstaVacio() {
		fmt.Println("(vacio)")
		return
	}
	for _, it := range c.Items() {
		fmt.Printf("#%d %-22s x%d  $%.2f\n", it.Producto().ID(), it.Producto().Nombre(),
			it.Cantidad(), it.Subtotal())
	}
	fmt.Printf("Subtotal: $%.2f\n", c.Subtotal())
}

func (t *Tienda) opcionQuitarDelCarrito() {
	clienteID, ok := leerNumero("ID del cliente: ")
	if !ok {
		return
	}
	c, err := t.carritoDe(clienteID)
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	productoID, ok := leerNumero("ID del producto: ")
	if !ok {
		return
	}
	mostrarResultado(c.Quitar(productoID), "Producto quitado del carrito")
}

// ---------------- PEDIDOS ----------------

func (t *Tienda) opcionConfirmarPedido() {
	clienteID, ok := leerNumero("ID del cliente: ")
	if !ok {
		return
	}
	c, err := t.carritoDe(clienteID)
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	if c.EstaVacio() {
		mostrarResultado(pedidos.ErrCarritoVacio, "")
		return
	}

	// Polimorfismo: se recorre un slice de la interfaz MetodoEnvio y cada
	// implementación calcula su propio costo para el mismo subtotal.
	metodos := pedidos.MetodosDisponibles()
	fmt.Printf("Subtotal: $%.2f. Elige el envio:\n", c.Subtotal())
	for i, m := range metodos {
		fmt.Printf("  %d) %-17s $%.2f  (%d dia(s))\n", i+1, m.Nombre(), m.Costo(c.Subtotal()), m.DiasEntrega())
	}
	n, ok := leerNumero("Envio: ")
	if !ok {
		return
	}
	if n < 1 || n > len(metodos) {
		fmt.Println("Error: ese metodo de envio no existe")
		return
	}

	p, err := t.pedidos.Confirmar(c, metodos[n-1])
	if err != nil {
		mostrarResultado(err, "")
		return
	}
	fmt.Printf("Pedido %d confirmado. Total a pagar: $%.2f\n", p.ID(), p.Total())
}

func (t *Tienda) opcionVerPedidos() {
	fmt.Println("--- PEDIDOS ---")
	lista := t.pedidos.Listar()
	if len(lista) == 0 {
		fmt.Println("(no hay pedidos)")
	}
	for _, p := range lista {
		fmt.Println(p)
	}
	fmt.Printf("Total vendido: $%.2f\n", t.pedidos.TotalVendido())
}

func (t *Tienda) opcionCancelarPedido() {
	id, ok := leerNumero("ID del pedido: ")
	if !ok {
		return
	}
	mostrarResultado(t.pedidos.Cancelar(id), "Pedido cancelado y stock devuelto")
}

// cargarDatosDemo agrega datos de ejemplo para probar el sistema enseguida.
// Si alguno fallara, se avisa en lugar de ignorar el error.
func (t *Tienda) cargarDatosDemo() {
	demo := []struct {
		nombre string
		precio float64
		stock  int
	}{
		{"Audifonos Bluetooth", 25.00, 10},
		{"Mouse inalambrico", 12.50, 20},
		{"Teclado mecanico", 45.90, 5},
	}
	for _, d := range demo {
		if _, err := t.catalogo.Agregar(d.nombre, d.precio, d.stock); err != nil {
			fmt.Println("Aviso: no se cargo", d.nombre, "->", err)
		}
	}
	if _, err := t.clientes.Registrar("Ana Torres", "ana@correo.com"); err != nil {
		fmt.Println("Aviso: no se cargo el cliente demo ->", err)
	}
}
