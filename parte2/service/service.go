package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type Service interface {
	GetOrgFromFile(fileName string) ([]Organization, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

type User struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type Organization struct {
	Organization string `json:"organization"`
	Users        []User `json:"users"`
}

// GetOrgFromFile lee un archivo CSV y retorna un slice de Organization
func (s *service) GetOrgFromFile(fileName string) ([]Organization, error) {
	csvFilePath := fmt.Sprintf("./data/%s", fileName)

	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error abriendo el archivo: %s", err))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err = reader.Read(); err != nil {
		return nil, errors.New(fmt.Sprintf("error leyendo el encabezado: %s", err))
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error leyendo el archivo: %s", err))
	}

	orgMap := s.getOrgMap(records)

	var result []Organization
	for _, org := range orgMap {
		result = append(result, *org)
	}

	return result, nil

}

// getOrgMap retorna un map con las organizaciones pobladas con sus usuarios
func (s *service) getOrgMap(records [][]string) map[string]*Organization {
	orgMap := make(map[string]*Organization)

	for _, record := range records {
		orgName := record[0]
		userName := record[1]
		role := record[2]

		org, exists := orgMap[orgName]
		if !exists {
			org = &Organization{
				Organization: orgName,
			}
			orgMap[orgName] = org
		}

		user, userIndex := s.findUser(org.Users, userName)
		if user == nil {
			user = &User{
				Username: userName,
				Roles:    []string{role},
			}
			org.Users = append(org.Users, *user)
		} else {
			org.Users[userIndex].Roles = append(org.Users[userIndex].Roles, role)
		}
	}
	return orgMap
}

// findUser retorna el usuario y su indice en el slice si existe, de lo contrario retorna nil y -1
func (s *service) findUser(users []User, username string) (*User, int) {
	for i, user := range users {
		if user.Username == username {
			return &users[i], i
		}
	}
	return nil, -1
}
