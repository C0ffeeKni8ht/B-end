package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"lab5/models"
)

func main() {
	dataDir := flag.String("data", "data", "каталог з CSV-файлами")
	inplace := flag.Bool("inplace", false, "змінювати справжні файли (а не копію)")
	flag.Parse()

	workDir := *dataDir
	if !*inplace {
		tmp, err := os.MkdirTemp("", "lab5-data-")
		check(err)
		defer os.RemoveAll(tmp)
		check(copyCSV(*dataDir, tmp))
		workDir = tmp
		fmt.Println("(демо працює з копією даних; початкові файли не змінюються)")
	}

	app, err := models.NewApp(workDir)
	check(err)

	title("Початковий стан CSV-файлів")
	show(workDir, "Element.csv", "Parameters.csv", "ElementParameters.csv")

	title("READ: дані, зчитані з CSV у пам'ять")
	fmt.Println("Elements:")
	app.Elements.DisplayAll(app.Categories)
	fmt.Println("Parameters:")
	app.Parameters.DisplayAll()
	fmt.Println("ElementParameters:")
	app.ElementParameters.DisplayAll(app.Elements, app.Parameters)

	title("CREATE: новий параметр, новий елемент PUSH і значення параметра для нього")
	paramID, err := app.Parameters.Add(models.Data{"name": "Розмір стека", "units": "байт"})
	check(err)
	elementID, err := app.Elements.Add(models.Data{
		"name": "PUSH", "syntax": "PUSH джерело", "example": "push ebp", "categoryId": 1,
	})
	check(err)
	linkID, err := app.ElementParameters.Add(models.Data{
		"elementId": elementID, "parameterId": paramID, "value": "4",
	})
	check(err)
	fmt.Printf("Створено: параметр id=%d, елемент id=%d, зв'язок id=%d\n", paramID, elementID, linkID)
	show(workDir, "Parameters.csv", "Element.csv", "ElementParameters.csv")

	title("UPDATE: елемент #1 (MOV → XCHG) і значення параметра #7 (6 → 7)")
	fmt.Println("ДО оновлення:")
	show(workDir, "Element.csv", "ElementParameters.csv")
	found, err := app.Elements.UpdateByID(models.Data{
		"id": 1, "name": "XCHG", "syntax": "XCHG операнд1, операнд2",
		"example": "xchg eax, ebx", "categoryId": 1,
	})
	check(err)
	fmt.Println("Результат updateById (елемент #1):", found)
	found, err = app.ElementParameters.UpdateByID(models.Data{
		"id": 7, "elementId": 5, "parameterId": 4, "value": "7",
	})
	check(err)
	fmt.Println("Результат updateById (зв'язок #7):", found)
	found, err = app.Elements.UpdateByID(models.Data{"id": 99, "name": "X"})
	check(err)
	fmt.Println("Результат updateById (елемент #99):", found)
	fmt.Println("ПІСЛЯ оновлення:")
	show(workDir, "Element.csv", "ElementParameters.csv")

	title("DELETE: елемент #3 разом з його значеннями параметрів")
	fmt.Println("ElementParameters.csv ДО видалення:")
	show(workDir, "ElementParameters.csv")
	found, err = app.DeleteElement(3)
	check(err)
	fmt.Println("Результат DeleteElement(3):", found)
	found, err = app.DeleteElement(3)
	check(err)
	fmt.Println("Результат DeleteElement(3) повторно:", found)
	fmt.Println("Element.csv та ElementParameters.csv ПІСЛЯ видалення:")
	show(workDir, "Element.csv", "ElementParameters.csv")

	title("DELETE: параметр #3 «Змінює прапорці» (зникає в усіх елементів)")
	found, err = app.DeleteParameter(3)
	check(err)
	fmt.Println("Результат DeleteParameter(3):", found)
	show(workDir, "Parameters.csv", "ElementParameters.csv")

	title("ПЕРЕВІРКА: перезапуск програми (дані читаються з файлів заново)")
	app2, err := models.NewApp(workDir)
	check(err)
	app2.Elements.DisplayAll(app2.Categories)
	app2.ElementParameters.DisplayAll(app2.Elements, app2.Parameters)
}

func title(s string) { fmt.Printf("\n=== %s ===\n", s) }

func show(workDir string, names ...string) {
	for _, name := range names {
		b, err := os.ReadFile(filepath.Join(workDir, name))
		check(err)
		fmt.Printf("--- %s ---\n%s", name, b)
	}
}

func copyCSV(src, dst string) error {
	files, err := filepath.Glob(filepath.Join(src, "*.csv"))
	if err != nil {
		return err
	}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dst, filepath.Base(f)), b, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
