package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/brenoassp/api-crud-persistencia-arquivo/domain"
)

type Service struct {
	dbFilePath string
	people     domain.People
}

func NewService(dbFilePath string) (Service, error) {
	// verifico se o arquivo existe
	_, err := os.Stat(dbFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// se nao existir, crio arquivo vazio
			err = createEmptyFile(dbFilePath)
			if err != nil {
				return Service{}, err
			}
			return Service{
				dbFilePath: dbFilePath,
				people:     domain.People{},
			}, nil
		}
	}

	// se existir, leio o arquivo e atualizo a variavel people do serviço com as pessoas do arquivo
	jsonFile, err := os.Open(dbFilePath)
	if err != nil {
		return Service{}, fmt.Errorf("Error trying to open file that contains all people: %s", err.Error())
	}

	jsonFileContentByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Service{}, fmt.Errorf("Error trying to read file: %s", err.Error())
	}

	var allPeople domain.People
	json.Unmarshal(jsonFileContentByte, &allPeople)

	return Service{
		dbFilePath: dbFilePath,
		people:     allPeople,
	}, nil
}

func createEmptyFile(dbFilePath string) error {
	var people domain.People = domain.People{
		People: []domain.Person{},
	}
	peopleJSON, err := json.Marshal(people)
	if err != nil {
		return fmt.Errorf("Error trying to encode people as JSON?: %s", err.Error())
	}

	err = ioutil.WriteFile(dbFilePath, peopleJSON, 0755)
	if err != nil {
		return fmt.Errorf("Error trying to write to file. Error: %s", err.Error())
	}

	return nil
}

func (s *Service) Create(person domain.Person) error {
	// verificar se a pessoa ja existe, se ja existe entao retorno erro
	if s.exists(person) {
		return fmt.Errorf("Error trying to create person. There is a person with this ID already registered")
	}

	// adiciono a pessoa na slice de pessoas
	s.people.People = append(s.people.People, person)

	// salvo o arquivo
	err := s.saveFile()
	if err != nil {
		return fmt.Errorf("Error trying save file in method Create. Error: %s", err.Error())
	}

	return nil
}

func (s Service) exists(person domain.Person) bool {
	for _, currentPerson := range s.people.People {
		if currentPerson.ID == person.ID {
			return true
		}
	}
	return false
}

func (s Service) saveFile() error {
	allPeopleJSON, err := json.Marshal(s.people)
	if err != nil {
		return fmt.Errorf("Error trying to encode people as json: %s", err.Error())
	}
	return ioutil.WriteFile(s.dbFilePath, allPeopleJSON, 0755)
}

func (s Service) List() domain.People {
	return s.people
}

func (s Service) GetByID(personID int) (domain.Person, error) {
	for _, currentPerson := range s.people.People {
		if currentPerson.ID == personID {
			return currentPerson, nil
		}
	}
	return domain.Person{}, fmt.Errorf("Person not found")
}

func (s *Service) Update(person domain.Person) error {
	var indexToUpdate int = -1
	for index, currentPerson := range s.people.People {
		if currentPerson.ID == person.ID {
			indexToUpdate = index
			break
		}
	}
	if indexToUpdate < 0 {
		return fmt.Errorf("There is no person with the given ID in our database")
	}

	s.people.People[indexToUpdate] = person
	return s.saveFile()
}

func (s *Service) DeleteByID(personID int) error {
	var indexToDelete int = -1
	for index, currentPerson := range s.people.People {
		if currentPerson.ID == personID {
			indexToDelete = index
			break
		}
	}
	if indexToDelete < 0 {
		return fmt.Errorf("There is no person with the given ID in our database")
	}

	s.people.People = append(s.people.People[:indexToDelete], s.people.People[indexToDelete+1:]...)
	return s.saveFile()
}
