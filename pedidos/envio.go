package pedidos

// MetodoEnvio es la INTERFAZ del sistema: describe qué debe saber hacer
// cualquier forma de entregar un pedido, sin decir cómo.
//
// Pedido guarda un MetodoEnvio y solo llama a Costo(); no sabe si es estándar,
// express o retiro en tienda. Eso es polimorfismo: el mismo llamado se comporta
// distinto según el tipo concreto. Para agregar un envío nuevo (por ejemplo,
// "Courier internacional") basta con crear otro struct con estos tres
// métodos; Pedido y Gestor no cambian.
//
// En Go no se escribe "implements": un tipo cumple la interfaz solo por tener
// los métodos con la misma firma.
type MetodoEnvio interface {
	Nombre() string
	Costo(subtotal float64) float64
	DiasEntrega() int
}

// Las tres implementaciones son structs vacíos: no necesitan datos propios,
// solo comportamiento. Usan receptor de valor porque no modifican nada.

// EnvioEstandar es gratis desde montoEnvioGratis; si no, cuesta costoEstandar.
type EnvioEstandar struct{}

// Constantes de la regla de negocio del envío estándar (definida en el AA1).
const (
	montoEnvioGratis = 50.0
	costoEstandar    = 3.50
	costoExpress     = 6.00
)

// Nombre del método.
func (EnvioEstandar) Nombre() string { return "Estandar" }

// Costo aplica la regla: gratis si el subtotal llega al monto mínimo.
func (EnvioEstandar) Costo(subtotal float64) float64 {
	if subtotal >= montoEnvioGratis {
		return 0
	}
	return costoEstandar
}

// DiasEntrega estimados.
func (EnvioEstandar) DiasEntrega() int { return 3 }

// EnvioExpress tiene un costo fijo y llega al día siguiente.
type EnvioExpress struct{}

// Nombre del método.
func (EnvioExpress) Nombre() string { return "Express" }

// Costo fijo, sin importar el subtotal.
func (EnvioExpress) Costo(float64) float64 { return costoExpress }

// DiasEntrega estimados.
func (EnvioExpress) DiasEntrega() int { return 1 }

// RetiroEnTienda no cobra envío: el cliente retira el pedido.
type RetiroEnTienda struct{}

// Nombre del método.
func (RetiroEnTienda) Nombre() string { return "Retiro en tienda" }

// Costo siempre 0.
func (RetiroEnTienda) Costo(float64) float64 { return 0 }

// DiasEntrega: disponible el mismo día.
func (RetiroEnTienda) DiasEntrega() int { return 0 }

// MetodosDisponibles devuelve todas las implementaciones en un slice del tipo
// de la interfaz: el menú las recorre sin saber cuál es cuál.
func MetodosDisponibles() []MetodoEnvio {
	return []MetodoEnvio{EnvioEstandar{}, EnvioExpress{}, RetiroEnTienda{}}
}
