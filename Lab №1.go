package main

import (
	"fmt"
	"math"
	"net/http"
)

func calculate() (float64, bool, float64) {

	b := 5.0
	x := 2.0

	numerator := b*math.Exp(-1.8*b+x) +
		math.Pow(math.Atan(x/2), 2)

	exponent := -1.8 * b

	var minusOnePower float64
	if int(math.Abs(exponent))%2 == 0 {
		minusOnePower = 1
	} else {
		minusOnePower = -1
	}

	denominator := minusOnePower +
		math.Sqrt(math.Abs(math.Log(x)+math.Log10(b)))

	argument := numerator / denominator
	S := math.Log(argument) / math.Log(b)

	x2 := 3
	y2 := 1
	z2 := 4

	result := x2+y2 <= z2 && z2 < 2*x2

	beta := 1.85

	p := math.Atan(beta) - math.Pow(math.Cos(beta), 2)

	p = math.Round(p*100) / 100

	q := math.Log(math.Abs(p+3)) +
		math.Sqrt(p*p-1)

	return S, result, q
}

func handler(w http.ResponseWriter, r *http.Request) {

	S, result, q := calculate()

	fmt.Fprintf(w, "Task 1:\n")
	fmt.Fprintf(w, "result S = %.6f\n", S)

	fmt.Fprintf(w, "Task 2:\n")
	fmt.Fprintf(w, "bool (%t)\n", result)

	fmt.Fprintf(w, "Task 3:\n")
	fmt.Fprintf(w, "result q = %.6f\n", q)
}

func main() {

	S, result, q := calculate()

	fmt.Printf("Task 1:\n")
	fmt.Printf("result S = %.6f\n", S)

	fmt.Printf("Task 2:\n")
	fmt.Printf("bool (%t)\n", result)

	fmt.Printf("Task 3:\n")
	fmt.Printf("result q = %.6f\n", q)

	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}
