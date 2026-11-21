package models

import (
	"fmt"
	"strconv"
)

type Parameters struct {
	BaseEntity
	name  string
	units string
}

type ParametersView struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Units string `json:"units"`
}

func NewParameters(data Data) *Parameters {
	return &Parameters{
		BaseEntity: BaseEntity{id: num(data, "id")},
		name:       str(data, "name"),
		units:      str(data, "units"),
	}
}

func (p *Parameters) GetAsArray() Data {
	return Data{"id": p.id, "name": p.name, "units": p.units}
}

func (p *Parameters) UpdateData(data Data) {
	p.name = str(data, "name")
	p.units = str(data, "units")
}

func (p *Parameters) ToRow() []string {
	return []string{strconv.Itoa(p.id), p.name, p.units}
}

func (p *Parameters) ToArray() ParametersView {
	return ParametersView{ID: p.id, Name: p.name, Units: p.units}
}

func (p *Parameters) DisplayData() {
	if p.units == "" {
		fmt.Printf("%d. %s\n", p.id, p.name)
		return
	}
	fmt.Printf("%d. %s (%s)\n", p.id, p.name, p.units)
}
