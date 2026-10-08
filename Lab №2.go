package main

import (
	"fmt"
	"math"
	"net/http"
)

func task1() float64 {
	x := 2.0

	var y float64

	if x < 0 {
		y = math.Exp(x) + 1
	} else if x > 0 && x <= 3 {
		y = x*x - math.Log(x)
	} else if x >= 3 {
		y = math.Sqrt(x - 3)
	}

	return y
}

func task2() ([]float64, []float64, []float64, []float64) {

	var xA []float64
	var yA []float64

	for x := 0.5; x <= 1.5; x += 0.2 {
		y := math.Sin(x) / (2*x + 3)

		xA = append(xA, x)
		yA = append(yA, y)
	}

	var xB []float64
	var yB []float64

	x := 1.0
	for i := 0; i < 5; i++ {
		y := math.Sin(x) / (2*x + 3)

		xB = append(xB, x)
		yB = append(yB, y)

		x += 0.1
	}

	return xA, yA, xB, yB
}

func task3() (float64, float64) {

	k := 2
	n := 7

	var S float64 = 0

	for i := k; i <= n; i++ {
		S += float64(i*i+3) / float64(i*i*i-1)
	}

	m := 3
	l := 9

	var P float64 = 1

	for j := m; j <= l; j++ {
		P *= float64(j-1) / float64(j*j+4)
	}

	return S, P
}

func task4() int {

	E := [6]float64{3.2, -4.1, 7.5, -2.0, 0.0, 5.3}

	count := 0

	for _, value := range E {
		if value > 1.0 {
			count++
		}
	}

	return count
}

func handler(w http.ResponseWriter, r *http.Request) {

	y := task1()

	fmt.Fprintf(w, "Task 1:\n")
	fmt.Fprintf(w, "result y = %.6f\n", y)

	xA, yA, xB, yB := task2()

	fmt.Fprintf(w, "\nTask 2:\n")

	fmt.Fprintf(w, "a) x, y:\n")
	for i := 0; i < len(xA); i++ {
		fmt.Fprintf(w, "x = %.1f, y = %.6f\n", xA[i], yA[i])
	}

	fmt.Fprintf(w, "b) x, y:\n")
	for i := 0; i < len(xB); i++ {
		fmt.Fprintf(w, "x = %.1f, y = %.6f\n", xB[i], yB[i])
	}

	S, P := task3()

	fmt.Fprintf(w, "\nTask 3:\n")
	fmt.Fprintf(w, "result S = %.6f\n", S)
	fmt.Fprintf(w, "result P = %.10f\n", P)

	count := task4()

	fmt.Fprintf(w, "\nTask 4:\n")
	fmt.Fprintf(w, "result = %d\n", count)
}

func main() {

	y := task1()

	fmt.Println("Task 1:")
	fmt.Printf("result y = %.6f\n", y)

	xA, yA, xB, yB := task2()

	fmt.Println("\nTask 2:")
	fmt.Println("a)")

	for i := 0; i < len(xA); i++ {
		fmt.Printf("x = %.1f, y = %.6f\n", xA[i], yA[i])
	}

	fmt.Println("b)")

	for i := 0; i < len(xB); i++ {
		fmt.Printf("x = %.1f, y = %.6f\n", xB[i], yB[i])
	}

	S, P := task3()

	fmt.Println("\nTask 3:")
	fmt.Printf("result S = %.6f\n", S)
	fmt.Printf("result P = %.10f\n", P)

	count := task4()

	fmt.Println("\nTask 4:")
	fmt.Printf("result = %d\n", count)

	http.HandleFunc("/", handler)

	fmt.Println("\nСайт запущено: http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Помилка запуску сервера:", err)
	}
}
