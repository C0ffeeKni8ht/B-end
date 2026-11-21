package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type element struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CategoryID string `json:"categoryId"`
	Syntax     string `json:"syntax"`
	Example    string `json:"example"`
}

type parameter struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Units string `json:"units"`
}

type elementParameter struct {
	ID          int    `json:"id"`
	ElementID   string `json:"elementId"`
	ParameterID string `json:"parameterId"`
	Value       string `json:"value"`
}

func main() {
	base := flag.String("url", "http://localhost:8000/api", "базова адреса API")
	flag.Parse()

	client := &http.Client{Timeout: 5 * time.Second}

	var elements []element
	var parameters []parameter
	var links []elementParameter
	fetch(client, *base+"/elements", &elements)
	fetch(client, *base+"/parameters", &parameters)
	fetch(client, *base+"/element-parameters", &links)

	fmt.Println("=== Elements ===")
	for _, e := range elements {
		fmt.Printf("#%d %s [категорія %s] — %s | %s\n", e.ID, e.Name, e.CategoryID, e.Syntax, e.Example)
	}
	fmt.Println("\n=== Parameters ===")
	for _, p := range parameters {
		if p.Units == "" {
			fmt.Printf("#%d %s\n", p.ID, p.Name)
			continue
		}
		fmt.Printf("#%d %s (%s)\n", p.ID, p.Name, p.Units)
	}

	elementName := map[string]string{}
	for _, e := range elements {
		elementName[strconv.Itoa(e.ID)] = e.Name
	}
	parameterByID := map[string]parameter{}
	for _, p := range parameters {
		parameterByID[strconv.Itoa(p.ID)] = p
	}
	fmt.Println("\n=== ElementParameters (зіставлено з Elements і Parameters) ===")
	for _, l := range links {
		p := parameterByID[l.ParameterID]
		fmt.Printf("#%d %s → %s = %s %s\n", l.ID, elementName[l.ElementID], p.Name, l.Value, p.Units)
	}

	for _, id := range []int{2, 99} {
		body, status, err := get(client, fmt.Sprintf("%s/elements?id=%d", *base, id))
		if err != nil {
			fmt.Println("Помилка з'єднання:", err)
			os.Exit(1)
		}
		fmt.Printf("\nGET /elements?id=%d → код %d\n%s", id, status, body)
	}
}

func fetch(client *http.Client, url string, result any) {
	body, status, err := get(client, url)
	if err != nil {
		fmt.Println("Помилка з'єднання:", err)
		fmt.Println("Запустіть сервер: go run ./cmd/api")
		os.Exit(1)
	}
	if status != http.StatusOK {
		fmt.Printf("Неочікувана відповідь від %s (код %d): %s\n", url, status, body)
		os.Exit(1)
	}
	if err := json.Unmarshal(body, result); err != nil {
		fmt.Printf("Некоректний JSON від %s: %v\n", url, err)
		os.Exit(1)
	}
}

func get(client *http.Client, url string) ([]byte, int, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}
