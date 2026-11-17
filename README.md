# Tienda Go · Sistema de Gestión de E-commerce

Proyecto del curso **Programación Orientada a Objetos (Go)** · UIDE · 2026
Autor: **Jhamir Rivera**
Entrega actual: **Aprendizaje Autónomo 2 · Desarrollo del sistema** (structs, encapsulación, errores e interfaces)
Entrega anterior: Aprendizaje Autónomo 1 · Planeación (funciones y paquetes), ver `docs/`.

## Objetivo

Permitir que una tienda en línea pequeña gestione su **catálogo de productos**, sus **clientes**, el **carrito de compras** y los **pedidos**, con control de stock, métodos de envío y total de ventas, desde un programa escrito en Go.

## Cómo ejecutarlo

Requisito: Go 1.22 o superior.

```bash
go build ./...   # comprueba que todo compila
go test ./...    # ejecuta las pruebas
go run .         # abre el menú
```

El programa arranca con 3 productos y 1 cliente de ejemplo (Ana Torres, ID 1).

## Estructura

```
EcommerceGo/                   módulo: ecommercego
├── main.go                    struct Tienda + menú (solo pide datos y muestra resultados)
├── productos/productos.go     struct Producto y struct Catalogo
├── clientes/clientes.go       struct Cliente y struct Registro
├── carrito/carrito.go         struct Item y struct Carrito
├── pedidos/envio.go           interfaz MetodoEnvio y sus 3 implementaciones
├── pedidos/pedidos.go         struct Linea, struct Pedido y struct Gestor
├── entrada/entrada.go         lectura del teclado
├── */*_test.go                pruebas de validaciones, errores e interfaz
└── docs/                      documentos del AA1 y AA2
```

## Del AA1 al AA2: de funciones sueltas a structs con métodos

| Funcionalidad (AA1) | Función en el AA1 | Ahora en el AA2 |
|---|---|---|
| 1 Agregar producto | `productos.AgregarProducto` | `Catalogo.Agregar` → usa `NewProducto` |
| 2 Ver catálogo | `productos.ListarProductos` | `Catalogo.Listar` + `Producto.String` |
| 3 Buscar producto | `productos.BuscarProducto` | `Catalogo.Buscar` |
| 4 Cambiar precio | `productos.ActualizarPrecio` | `Producto.SetPrecio` (valida) |
| 5 Eliminar producto | `productos.EliminarProducto` | `Catalogo.Eliminar` |
| 6 Registrar cliente | `clientes.RegistrarCliente` | `Registro.Registrar` → usa `NewCliente` |
| 7 Ver clientes | `clientes.ListarClientes` | `Registro.Listar` |
| 8 Agregar al carrito | `carrito.AgregarAlCarrito` | `Carrito.Agregar` |
| 9 Ver carrito | `carrito.VerCarrito` | `Carrito.Items` + `Carrito.Subtotal` |
| 10 Quitar del carrito | `carrito.QuitarDelCarrito` | `Carrito.Quitar` |
| 11 Confirmar pedido | `pedidos.ConfirmarPedido` | `Gestor.Confirmar(carrito, MetodoEnvio)` |
| 12 Ver pedidos | `pedidos.ListarPedidos` | `Gestor.Listar` + `Gestor.TotalVendido` |
| 13 Cancelar pedido | `pedidos.CancelarPedido` | `Gestor.Cancelar` |

Las **13 funcionalidades** del AA1 funcionan, y se agregó la elección del **método de envío**.

## Decisiones de diseño

### 1. ¿Por qué estas entidades?

- **Producto, Cliente, Pedido**: son los sustantivos del negocio; cada uno tiene datos y reglas propias (un producto no puede tener precio negativo, un cliente necesita un correo válido, un pedido no se cancela dos veces).
- **Carrito e Item**: el carrito existe antes del pedido y cambia mucho; separarlo evita pedidos "a medio armar".
- **Linea** (de pedido) es distinta de **Item** (de carrito): la línea **copia el nombre y el precio** al momento de comprar. Si mañana el producto sube de precio, el pedido histórico no cambia.
- **Catalogo, Registro, Gestor** son los "dueños" de cada colección. En el AA1 eran maps sueltos a nivel de paquete; ahora son campos de un struct, y `Tienda` los agrupa. Así se pueden crear tiendas independientes (por ejemplo, en las pruebas).

### 2. Receptor de valor o puntero

| Receptor | Cuándo | Ejemplos |
|---|---|---|
| **Valor** `(p Producto)` | El método solo **lee**; trabajar con una copia no afecta | `Nombre()`, `Precio()`, `HayStock()`, `String()`, `Costo()` |
| **Puntero** `(p *Producto)` | El método **modifica** el original | `SetPrecio()`, `DescontarStock()`, `Carrito.Agregar()`, `Catalogo.Eliminar()` |

Los constructores devuelven **punteros** (`*Producto`) porque el catálogo, el carrito y el pedido deben compartir **el mismo** producto: cuando se vende, todos ven el stock actualizado. Si se guardaran copias, cada uno tendría un stock distinto.

### 3. Slice o map

| Estructura | Dónde | Por qué |
|---|---|---|
| `map[int]*Producto` | `Catalogo` | La operación más usada es **buscar por ID** (carrito, pedidos, menú). |
| `map[int]*Cliente` | `Registro` | Igual: se busca por ID al armar carritos. |
| `map[int]*Carrito` | `Tienda` | Un carrito por cliente, se accede por el ID del cliente. |
| `[]Item` | `Carrito` | **Importa el orden** en que se agregaron, y son pocas líneas. |
| `[]Linea` | `Pedido` | El detalle se muestra en orden. |
| `[]*Pedido` | `Gestor` | Historial **cronológico**; siempre se listan todos. |

Como los maps no tienen orden fijo, `Listar()` los pasa a un slice y los ordena por ID.

### 4. Encapsulación

- Todos los campos de los structs están en **minúscula** (no exportados): desde `main` no se puede escribir `p.precio = -5`.
- Se leen con **getters** al estilo Go (`Precio()`, no `GetPrecio()`).
- Se cambian con **setters** que usan receptor puntero y **validan** (`SetPrecio`, `SetNombre`, `SetCorreo`). Si el valor no es válido, devuelven un error y el dato no cambia.
- `Carrito.Items()`, `Pedido.Lineas()` y `Gestor.Listar()` devuelven **copias** del slice, para que nadie modifique la colección por fuera de los métodos.
- `Pedido.cancelar()` es interno: solo `Gestor.Cancelar` puede cancelar, así el stock y el historial siempre quedan coherentes.

### 5. Manejo de errores

- Toda función que puede fallar devuelve `error` como último valor, y quien la llama lo revisa (en el menú, `mostrarResultado`).
- Cada paquete declara **errores base** (`ErrNoExiste`, `ErrPrecio`, `ErrSinStock`, `ErrCarritoVacio`…) que se pueden comparar con `errors.Is`.
- Con `fmt.Errorf("%w: ...")` se agrega detalle sin perder el error base: el usuario ve *"no hay stock suficiente: solo hay 5 de Teclado mecanico"* y el código sigue reconociendo `ErrSinStock`.
- `Gestor.Confirmar` valida en **dos pasadas** (primero revisa el stock de todo, después descuenta) para no dejar productos descontados si uno falla.

### 6. Interfaz: `MetodoEnvio`

```go
type MetodoEnvio interface {
    Nombre() string
    Costo(subtotal float64) float64
    DiasEntrega() int
}
```

| Implementación | Costo | Entrega |
|---|---|---|
| `EnvioEstandar` | Gratis desde $50; si no, $3.50 (regla del AA1) | 3 días |
| `EnvioExpress` | $6.00 fijo | 1 día |
| `RetiroEnTienda` | $0 | mismo día |

`Pedido` guarda un `MetodoEnvio` y solo llama a `Costo()`: no sabe cuál es. Eso es **polimorfismo**. Para agregar un envío nuevo basta con crear otro struct con esos tres métodos; `Pedido` y `Gestor` no cambian.

También `Producto`, `Cliente` y `Pedido` implementan la interfaz estándar **`fmt.Stringer`** (método `String()`), por eso `fmt.Println(p)` imprime un texto legible.

## Paquetes externos

No se usan paquetes de terceros: solo la biblioteca estándar (`fmt`, `errors`, `strings`, `strconv`, `sort`, `bufio`, `os`, `testing`).

## Próxima etapa (ECD)

Cada funcionalidad se expondrá como servicio web con JSON (`net/http`), protegido con `sync.Mutex`, con más pruebas.

| # | Funcionalidad | Servicio previsto |
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
