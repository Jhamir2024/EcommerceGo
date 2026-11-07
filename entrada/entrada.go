// Package entrada lee datos del teclado. Usa bufio para poder leer textos con
// espacios (fmt.Scan solo lee una palabra).
package entrada

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var lector = bufio.NewReader(os.Stdin)

// LeerTexto muestra un mensaje y devuelve la línea escrita, sin espacios a los lados.
func LeerTexto(mensaje string) string {
	fmt.Print(mensaje)
	linea, _ := lector.ReadString('\n')
	return strings.TrimSpace(linea)
}

// LeerEntero muestra un mensaje y convierte la respuesta a int.
func LeerEntero(mensaje string) (int, error) {
	return strconv.Atoi(LeerTexto(mensaje))
}

// LeerDecimal muestra un mensaje y convierte la respuesta a float64.
// Acepta coma o punto decimal (12,5 o 12.5).
func LeerDecimal(mensaje string) (float64, error) {
	texto := strings.Replace(LeerTexto(mensaje), ",", ".", 1)
	return strconv.ParseFloat(texto, 64)
}
