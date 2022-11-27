package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/brenoassp/api-crud-persistencia-arquivo/domain"
	"github.com/brenoassp/api-crud-persistencia-arquivo/domain/person"
)

func main() {
	personService, err := person.NewService("person.json")
	if err != nil {
		fmt.Println("Error trying to create person service")
		return
	}

	http.HandleFunc("/person/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var person domain.Person
			err := json.NewDecoder(r.Body).Decode(&person)
			if err != nil {
				fmt.Printf("Error trying to decode body. Body should be a json. Error: %s", err.Error())
				http.Error(w, "Error trying to create person", http.StatusBadRequest)
				return
			}
			if person.ID <= 0 {
				http.Error(w, "Error trying to create person. ID should be a positive integer", http.StatusBadRequest)
				return
			}

			// criar pessoa
			err = personService.Create(person)
			if err != nil {
				fmt.Printf("Error trying to create person: %s", err.Error())
				http.Error(w, "Error trying to create person", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		}

		if r.Method == "GET" {
			path := strings.TrimPrefix(r.URL.Path, "/person/")
			if path == "" {
				// list all
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				people := personService.List()
				err := json.NewEncoder(w).Encode(people)
				if err != nil {
					http.Error(w, "Error trying to list people", http.StatusInternalServerError)
					return
				}
			} else {
				// /person/2 list pessoa com id 2
				personID, err := strconv.Atoi(path)
				if err != nil {
					http.Error(w, "Invalid id given. person ID must be an integer", http.StatusBadRequest)
					return
				}
				person, err := personService.GetByID(personID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-type", "application/json")
				err = json.NewEncoder(w).Encode(person)
				if err != nil {
					http.Error(w, "Error trying to encode person as json", http.StatusInternalServerError)
					return
				}
			}
		}

		if r.Method == "PUT" {
			var person domain.Person
			err := json.NewDecoder(r.Body).Decode(&person)
			if err != nil {
				fmt.Printf("Error trying to decode body. Body should be a json. Error: %s", err.Error())
				http.Error(w, "Error trying to create person", http.StatusBadRequest)
				return
			}
			if person.ID <= 0 {
				http.Error(w, "Error trying to create person. ID should be a positive integer", http.StatusBadRequest)
				return
			}

			// atualizar pessoa
			err = personService.Update(person)
			if err != nil {
				fmt.Printf("Error trying to update person: %s", err.Error())
				http.Error(w, "Error trying to update person", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == "DELETE" {
			path := strings.TrimPrefix(r.URL.Path, "/person/")
			if path == "" {
				http.Error(w, "ID must be provided in the url", http.StatusBadRequest)
				return
			} else {
				personID, err := strconv.Atoi(path)
				if err != nil {
					http.Error(w, "Invalid id given. person ID must be an integer", http.StatusBadRequest)
					return
				}
				err = personService.DeleteByID(personID)
				if err != nil {
					fmt.Printf("Error trying to delete person: %s", err.Error())
					http.Error(w, "Error trying to delete person", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
			}
		}
	})

	http.ListenAndServe(":8080", nil)
}
