package main

import "fmt"

func main() {
	const maxProductos = 3 // limite de productos (reto 2)
	cantidad := 0
	nombre1, nombre2, nombre3 := "", "", ""
	precio1, precio2, precio3 := 0.0, 0.0, 0.0
	stock1, stock2, stock3 := 0, 0, 0
	totalVentas := 0.0

	for {
		fmt.Println("\n=== E-COMMERCE ===")
		fmt.Println("1) Agregar producto")
		fmt.Println("2) Ver catalogo")
		fmt.Println("3) Vender producto")
		fmt.Println("4) Ver logo de la tienda")
		fmt.Println("5) Salir")
		fmt.Print("Elige: ")
		var opcion int
		fmt.Scan(&opcion)

		switch opcion {
		case 1:
			if cantidad == maxProductos {
				fmt.Println("Catalogo lleno: maximo", maxProductos, "productos")
			} else {
				var nombre string
				var precio float64
				var stock int
				fmt.Print("Nombre (una palabra): ")
				fmt.Scan(&nombre)
				fmt.Print("Precio: ")
				fmt.Scan(&precio)
				fmt.Print("Stock: ")
				fmt.Scan(&stock)

				if precio <= 0 || stock < 0 {
					fmt.Println("Precio o stock no valido")
				} else {
					cantidad++
					if cantidad == 1 {
						nombre1, precio1, stock1 = nombre, precio, stock
					} else if cantidad == 2 {
						nombre2, precio2, stock2 = nombre, precio, stock
					} else {
						nombre3, precio3, stock3 = nombre, precio, stock
					}
					fmt.Println("Producto agregado:", nombre)
				}
			}
		case 2:
			if cantidad == 0 {
				fmt.Println("No hay productos todavia")
			} else {
				fmt.Println("Productos:", cantidad, "de", maxProductos)
				if cantidad >= 1 {
					fmt.Printf("1) %s - $%.2f - stock %d\n", nombre1, precio1, stock1)
				}
				if cantidad >= 2 {
					fmt.Printf("2) %s - $%.2f - stock %d\n", nombre2, precio2, stock2)
				}
				if cantidad >= 3 {
					fmt.Printf("3) %s - $%.2f - stock %d\n", nombre3, precio3, stock3)
				}
				fmt.Printf("Total vendido: $%.2f\n", totalVentas)
			}
		case 3:
			var num, unidades int
			fmt.Print("Numero de producto: ")
			fmt.Scan(&num)
			if num < 1 || num > cantidad {
				fmt.Println("Ese producto no existe")
			} else {
				fmt.Print("Unidades: ")
				fmt.Scan(&unidades)

				// p apunta al stock del producto elegido
				var p *int
				var precio float64
				switch num {
				case 1:
					p, precio = &stock1, precio1
				case 2:
					p, precio = &stock2, precio2
				case 3:
					p, precio = &stock3, precio3
				}

				if unidades <= 0 || unidades > *p {
					fmt.Println("Cantidad no valida o sin stock suficiente")
				} else {
					*p -= unidades
					venta := precio * float64(unidades)
					totalVentas += venta
					fmt.Printf("Venta de $%.2f realizada\n", venta)
				}
			}
		case 4:
			// dibujo con for (reto 3)
			for i := 1; i <= 4; i++ {
				for j := 0; j < i; j++ {
					fmt.Print("* ")
				}
				fmt.Println()
			}
			fmt.Println("TIENDA GO")
		case 5:
			fmt.Println("Gracias, hasta pronto")
			return
		default:
			fmt.Println("Opcion no existe") // reto 2
		}
	}
}
