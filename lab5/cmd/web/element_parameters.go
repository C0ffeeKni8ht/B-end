package main

import (
	"net/http"
	"strconv"
	"strings"

	"lab5/models"
)

func (s *server) parseLink(r *http.Request) (models.Data, string) {
	elementID, _ := strconv.Atoi(r.PostFormValue("elementId"))
	parameterID, _ := strconv.Atoi(r.PostFormValue("parameterId"))
	value := strings.TrimSpace(r.PostFormValue("value"))

	if s.app.Elements.Get(elementID) == nil {
		return nil, "Оберіть існуючий елемент."
	}
	if s.app.Parameters.Get(parameterID) == nil {
		return nil, "Оберіть існуючий параметр."
	}
	if msg := checkText("Значення", value, true); msg != "" {
		return nil, msg
	}
	return models.Data{"elementId": elementID, "parameterId": parameterID, "value": value}, ""
}

func (s *server) linkOps() entityOps {
	return entityOps{
		name: "element-parameters",
		add: func(r *http.Request) (string, error) {
			data, msg := s.parseLink(r)
			if msg != "" {
				return msg, nil
			}
			_, err := s.app.ElementParameters.Add(data)
			return "", err
		},
		update: func(r *http.Request, id int) (bool, string, error) {
			if s.app.ElementParameters.Get(id) == nil {
				return false, "", nil
			}
			data, msg := s.parseLink(r)
			if msg != "" {
				return true, msg, nil
			}
			data["id"] = id
			found, err := s.app.ElementParameters.UpdateByID(data)
			return found, "", err
		},
		remove: func(id int) (bool, string, error) {
			found, err := s.app.ElementParameters.DeleteByID(id)
			return found, "", err
		},
		load: func(id int) (map[string]string, bool) {
			d := s.app.ElementParameters.FindByID(id)
			return dataValues(d), d != nil
		},
	}
}
