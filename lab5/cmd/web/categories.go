package main

import (
	"net/http"
	"strconv"
	"strings"

	"lab5/models"
)

func (s *server) parseElement(r *http.Request) (models.Data, string) {
	name := strings.TrimSpace(r.PostFormValue("name"))
	syntax := strings.TrimSpace(r.PostFormValue("syntax"))
	example := strings.TrimSpace(r.PostFormValue("example"))
	categoryID, _ := strconv.Atoi(r.PostFormValue("categoryId"))

	for _, c := range []struct{ label, v string }{{"Назва", name}, {"Синтаксис", syntax}, {"Приклад", example}} {
		if msg := checkText(c.label, c.v, true); msg != "" {
			return nil, msg
		}
	}
	if s.app.Categories.FindByID(categoryID) == nil {
		return nil, "Оберіть існуючу категорію."
	}
	return models.Data{
		"name": name, "syntax": syntax, "example": example, "categoryId": categoryID,
	}, ""
}

func (s *server) elementOps() entityOps {
	return entityOps{
		name: "elements",
		add: func(r *http.Request) (string, error) {
			data, msg := s.parseElement(r)
			if msg != "" {
				return msg, nil
			}
			_, err := s.app.Elements.Add(data)
			return "", err
		},
		update: func(r *http.Request, id int) (bool, string, error) {
			if s.app.Elements.Get(id) == nil {
				return false, "", nil
			}
			data, msg := s.parseElement(r)
			if msg != "" {
				return true, msg, nil
			}
			data["id"] = id
			found, err := s.app.Elements.UpdateByID(data)
			return found, "", err
		},
		remove: func(id int) (bool, string, error) {
			found, err := s.app.DeleteElement(id)
			return found, "", err
		},
		load: func(id int) (map[string]string, bool) {
			d := s.app.Elements.FindByID(id)
			return dataValues(d), d != nil
		},
	}
}
