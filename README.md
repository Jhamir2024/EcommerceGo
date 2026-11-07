# Tienda Go · Sistema de Gestión de E-commerce

Proyecto del curso **Programación Orientada a Objetos (Go)** · UIDE · 2026
Autor: **Jhamir Rivera**
Entrega actual: **Aprendizaje Autónomo 1 · Planeación del sistema** (funciones y paquetes)

## Objetivo

Permitir que una tienda en línea pequeña gestione su **catálogo de productos**, sus **clientes**, el **carrito de compras** y los **pedidos**, con control de stock y total de ventas, desde un programa escrito en Go.

## Alcance

**Sí hace:**
- Registrar, consultar, cambiar de precio y eliminar productos, con control de stock.
- Registrar y consultar clientes, validando el correo.
- Armar un carrito por cliente con subtotales.
- Confirmar pedidos (descuenta stock y calcula envío: gratis desde $50; si no, $3.50).
- Cancelar pedidos (devuelve el stock) y mostrar el total vendido.

**No hace (fuera del alcance):**
- Pagos reales con tarjeta ni facturación electrónica.
- Envíos ni seguimiento de paquetes.
- Inicio de sesión con contraseña.
- Guardar datos en una base de datos: por ahora todo vive en memoria mientras el programa corre.

## Módulos y paquetes

```
EcommerceGo/                 módulo: ecommercego
├── go.mod
├── main.go                  menú de consola (solo pide datos y llama a los paquetes)
├── productos/productos.go   catálogo y stock
├── clientes/clientes.go     registro de clientes
├── carrito/carrito.go       carrito de compras por cliente
├── pedidos/pedidos.go       confirmar, cancelar y total vendido
├── entrada/entrada.go       lectura del teclado
├── docs/                    PDF de la planeación
└── practicas/semana01/      menú de la semana 1 (versión anterior)
```

Dependencias entre paquetes: `main` usa todos los paquetes; `pedidos` usa `carrito`, `clientes` y `productos`; `carrito` usa `clientes` y `productos`.

### Funciones por módulo

**productos**

| Función | Parámetros | Retorna | Qué hace |
|---|---|---|---|
| `AgregarProducto` | `nombre string, precio float64, unidades int` | `int, error` | Registra un producto y devuelve su ID |
| `ListarProductos` | — | `[]string` | Lista el catálogo |
| `BuscarProducto` | `id int` | `string, float64, int, error` | Devuelve nombre, precio y stock |
| `ActualizarPrecio` | `id int, nuevoPrecio float64` | `error` | Cambia el precio |
| `EliminarProducto` | `id int` | `error` | Quita el producto del catálogo |
| `Existe` | `id int` | `bool` | Indica si el producto existe |
| `DescontarStock` | `id int, unidades int` | `error` | Resta stock al vender |
| `DevolverStock` | `id int, unidades int` | `error` | Suma stock al cancelar |
| `precioValido` *(interna)* | `precio float64` | `bool` | Valida que el precio sea mayor que 0 |

**clientes**

| Función | Parámetros | Retorna | Qué hace |
|---|---|---|---|
| `RegistrarCliente` | `nombre string, correo string` | `int, error` | Registra un cliente y devuelve su ID |
| `BuscarCliente` | `id int` | `string, string, error` | Devuelve nombre y correo |
| `ListarClientes` | — | `[]string` | Lista los clientes |
| `Existe` | `id int` | `bool` | Indica si el cliente existe |
| `correoValido` *(interna)* | `correo string` | `bool` | Revisa que tenga `@` y `.` |
| `correoRegistrado` *(interna)* | `correo string` | `bool` | Evita correos repetidos |

**carrito**

| Función | Parámetros | Retorna | Qué hace |
|---|---|---|---|
| `AgregarAlCarrito` | `clienteID int, productoID int, unidades int` | `error` | Añade unidades validando el stock |
| `QuitarDelCarrito` | `clienteID int, productoID int` | `error` | Saca un producto del carrito |
| `VerCarrito` | `clienteID int` | `[]string` | Lista el carrito con subtotales |
| `CalcularTotal` | `clienteID int` | `float64, error` | Suma el carrito (sin envío) |
| `Contenido` | `clienteID int` | `map[int]int` | Copia del carrito para `pedidos` |
| `VaciarCarrito` | `clienteID int` | — | Deja el carrito vacío |
| `idsOrdenados` *(interna)* | `clienteID int` | `[]int` | Ordena los productos del carrito |

**pedidos**

| Función | Parámetros | Retorna | Qué hace |
|---|---|---|---|
| `ConfirmarPedido` | `clienteID int` | `int, float64, error` | Compra el carrito: descuenta stock y devuelve ID y total |
| `CancelarPedido` | `pedidoID int` | `error` | Cancela y devuelve el stock |
| `ListarPedidos` | — | `[]string` | Lista los pedidos con su estado |
| `TotalVendido` | — | `float64` | Suma los pedidos confirmados |
| `calcularEnvio` *(interna)* | `subtotal float64` | `float64` | Envío gratis desde $50; si no, $3.50 |

**entrada**

| Función | Parámetros | Retorna | Qué hace |
|---|---|---|---|
| `LeerTexto` | `mensaje string` | `string` | Lee una línea (acepta espacios) |
| `LeerEntero` | `mensaje string` | `int, error` | Lee un número entero |
| `LeerDecimal` | `mensaje string` | `float64, error` | Lee un decimal (con coma o punto) |

## Paquetes externos

No se usan paquetes de terceros: todo se resuelve con la biblioteca estándar de Go (`fmt`, `errors`, `strings`, `strconv`, `sort`, `bufio`, `os`). Así el proyecto compila sin `go get`. Si más adelante hace falta uno, se instalará con `go get` y se documentará aquí.

## Cómo ejecutarlo

Requisito: Go 1.22 o superior.

```bash
go build ./...   # comprueba que todo compila
go run .         # abre el menú
```

El programa arranca con 3 productos y 1 cliente de ejemplo para poder probarlo enseguida.

## Funcionalidades (futuros servicios web)

| # | Funcionalidad | Servicio web previsto (ECD) |
|---|---|---|
| 1 | Agregar producto | `POST /productos` |
| 2 | Ver catálogo | `GET /productos` |
| 3 | Buscar producto | `GET /productos/{id}` |
| 4 | Cambiar precio | `PUT /productos/{id}` |
| 5 | Eliminar producto | `DELETE /productos/{id}` |
| 6 | Registrar cliente | `POST /clientes` |
| 7 | Ver clientes | `GET /clientes` |
| 8 | Agregar al carrito | `POST /carrito/{cliente}` |
| 9 | Ver carrito | `GET /carrito/{cliente}` |
| 10 | Quitar del carrito | `DELETE /carrito/{cliente}/{producto}` |
| 11 | Confirmar pedido | `POST /pedidos` |
| 12 | Ver pedidos y total vendido | `GET /pedidos` |
| 13 | Cancelar pedido | `PUT /pedidos/{id}/cancelar` |

## Próximas etapas

- **AA2:** los datos pasarán de maps sueltos a structs (`Producto`, `Cliente`, `Pedido`) con constructores, encapsulación, errores e interfaces (por ejemplo, una interfaz `MetodoEnvio`).
- **ECD:** cada funcionalidad se expondrá como servicio web con JSON, protegido con `sync.Mutex` y con pruebas.
