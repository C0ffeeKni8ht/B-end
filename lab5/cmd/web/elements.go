package main

import (
	"net/http"
	"strings"

	"lab5/models"
)

func parseCategory(r *http.Request) (models.Data, string) {
	name := strings.TrimSpace(r.PostFormValue("name"))
	if msg := checkText("Назва категорії", name, true); msg != "" {
		return nil, msg
	}
	return models.Data{"name": name}, ""
}

func (s *server) categoryOps() entityOps {
	return entityOps{
		name: "categories",
		add: func(r *http.Request) (string, error) {
			data, msg := parseCategory(r)
			if msg != "" {
				return msg, nil
			}
			_, err := s.app.Categories.Add(data)
			return "", err
		},
		update: func(r *http.Request, id int) (bool, string, error) {
			if s.app.Categories.FindByID(id) == nil {
				return false, "", nil
			}
			data, msg := parseCategory(r)
			if msg != "" {
				return true, msg, nil
			}
			data["id"] = id
			found, err := s.app.Categories.UpdateByID(data)
			return found, "", err
		},
		remove: func(id int) (bool, string, error) {
			if s.app.Categories.FindByID(id) == nil {
				return false, "", nil
			}
			for _, e := range s.app.ElementViews() {
				if e.CategoryID == id {
					return true, "Не можна видалити категорію, у якій є елементи.", nil
				}
			}
			found, err := s.app.Categories.DeleteByID(id)
			return found, "", err
		},
		load: func(id int) (map[string]string, bool) {
			d := s.app.Categories.FindByID(id)
			return dataValues(d), d != nil
		},
	}
}
