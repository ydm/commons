package commons

import "fmt"

func PrintArrayF64(xs []float64, precision int) {
	fmt.Print("[")

	prefix := ""
	for _, x := range xs {
		fmt.Printf("%s%.*f", prefix, precision, x)
		prefix = ", "
	}

	fmt.Println("]")
}
