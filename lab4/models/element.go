package models

import (
	"fmt"
	"strconv"
)

type Element struct {
	BaseEntity
	name       string
	syntax     string
	example    string
	categoryID int
}

type ElementView struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CategoryID int    `json:"categoryId,string"`
	Syntax     string `json:"syntax"`
	Example    string `json:"example"`
	Category   string `json:"-"`
}

func NewElement(data Data) *Element {
	return &Element{
		BaseEntity: BaseEntity{id: num(data, "id")},
		name:       str(data, "name"),
		syntax:     str(data, "syntax"),
		example:    str(data, "example"),
		categoryID: num(data, "categoryId"),
	}
}

func (e *Element) GetAsArray() Data {
	return Data{
		"id":         e.id,
		"name":       e.name,
		"syntax":     e.syntax,
		"example":    e.example,
		"categoryId": e.categoryID,
	}
}

func (e *Element) UpdateData(data Data) {
	e.name = str(data, "name")
	e.syntax = str(data, "syntax")
	e.example = str(data, "example")
	e.categoryID = num(data, "categoryId")
}

func (e *Element) ToRow() []string {
	return []string{
		strconv.Itoa(e.id), e.name, e.syntax, e.example, strconv.Itoa(e.categoryID),
	}
}

func (e *Element) ToArray(categories *CategoryRepository) ElementView {
	categoryName := "—"
	if cat := categories.FindByID(e.categoryID); cat != nil {
		categoryName = str(cat, "name")
	}
	return ElementView{
		ID:         e.id,
		Name:       e.name,
		Syntax:     e.syntax,
		Example:    e.example,
		CategoryID: e.categoryID,
		Category:   categoryName,
	}
}

func (e *Element) DisplayData(categories *CategoryRepository) {
	v := e.ToArray(categories)
	fmt.Printf("%d. %s. Синтаксис: %s. Приклад: %s, Категорія: %s\n",
		v.ID, v.Name, v.Syntax, v.Example, v.Category)
}
