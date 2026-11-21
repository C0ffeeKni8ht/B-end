package models

import (
	"fmt"
	"strconv"
)

type ElementParameters struct {
	BaseEntity
	elementID   int
	parameterID int
	value       string
}

type ElementParametersView struct {
	ID          int    `json:"id"`
	ElementID   int    `json:"elementId,string"`
	ParameterID int    `json:"parameterId,string"`
	Value       string `json:"value"`
}

func NewElementParameters(data Data) *ElementParameters {
	return &ElementParameters{
		BaseEntity:  BaseEntity{id: num(data, "id")},
		elementID:   num(data, "elementId"),
		parameterID: num(data, "parameterId"),
		value:       str(data, "value"),
	}
}

func (ep *ElementParameters) GetAsArray() Data {
	return Data{
		"id":          ep.id,
		"elementId":   ep.elementID,
		"parameterId": ep.parameterID,
		"value":       ep.value,
	}
}

func (ep *ElementParameters) UpdateData(data Data) {
	ep.elementID = num(data, "elementId")
	ep.parameterID = num(data, "parameterId")
	ep.value = str(data, "value")
}

func (ep *ElementParameters) ToRow() []string {
	return []string{
		strconv.Itoa(ep.id), strconv.Itoa(ep.elementID),
		strconv.Itoa(ep.parameterID), ep.value,
	}
}

func (ep *ElementParameters) ToArray() ElementParametersView {
	return ElementParametersView{
		ID: ep.id, ElementID: ep.elementID,
		ParameterID: ep.parameterID, Value: ep.value,
	}
}

func (ep *ElementParameters) DisplayData(elements *ElementRepository, parameters *ParametersRepository) {
	elementName, parameterName, units := "—", "—", ""
	if e := elements.Get(ep.elementID); e != nil {
		elementName = e.name
	}
	if p := parameters.Get(ep.parameterID); p != nil {
		parameterName = p.name
		units = p.units
	}
	fmt.Printf("%d. %s → %s = %s %s\n", ep.id, elementName, parameterName, ep.value, units)
}
