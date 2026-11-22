package models

import (
	"fmt"
	"strconv"
)

type Category struct {
	BaseEntity
	name string
}

type CategoryView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewCategory(data Data) *Category {
	return &Category{
		BaseEntity: BaseEntity{id: num(data, "id")},
		name:       str(data, "name"),
	}
}

func (c *Category) Name() string { return c.name }

func (c *Category) GetAsArray() Data {
	return Data{"id": c.id, "name": c.name}
}

func (c *Category) UpdateData(data Data) {
	c.name = str(data, "name")
}

func (c *Category) ToRow() []string {
	return []string{strconv.Itoa(c.id), c.name}
}

func (c *Category) ToArray() CategoryView { return CategoryView{ID: c.id, Name: c.name} }

func (c *Category) DisplayData() {
	fmt.Printf("%d. %s\n", c.id, c.name)
}
