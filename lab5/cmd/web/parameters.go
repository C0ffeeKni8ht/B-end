package main

import (
	"net/http"
	"strings"

	"lab5/models"
)

func parseParameter(r *http.Request) (models.Data, string) {
	name := strings.TrimSpace(r.PostFormValue("name"))
	units := strings.TrimSpace(r.PostFormValue("units"))
	if msg := checkText("Назва параметра", name, true); msg != "" {
		return nil, msg
	}
	if msg := checkText("Одиниці виміру", units, false); msg != "" {
		return nil, msg
	}
	return models.Data{"name": name, "units": units}, ""
}

func (s *server) parameterOps() entityOps {
	return entityOps{
		name: "parameters",
		add: func(r *http.Request) (string, error) {
			data, msg := parseParameter(r)
			if msg != "" {
				return msg, nil
			}
			_, err := s.app.Parameters.Add(data)
			return "", err
		},
		update: func(r *http.Request, id int) (bool, string, error) {
			if s.app.Parameters.Get(id) == nil {
				return false, "", nil
			}
			data, msg := parseParameter(r)
			if msg != "" {
				return true, msg, nil
			}
			data["id"] = id
			found, err := s.app.Parameters.UpdateByID(data)
			return found, "", err
		},
		remove: func(id int) (bool, string, error) {
			found, err := s.app.DeleteParameter(id)
			return found, "", err
		},
		load: func(id int) (map[string]string, bool) {
			d := s.app.Parameters.FindByID(id)
			return dataValues(d), d != nil
		},
	}
}
