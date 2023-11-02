package service

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGetOrgFromFile(t *testing.T) {
	csvTest := `organizacion, usuario, rol
org1,jperez,admin
org1,jperez,superadmin
org1,asosa,writer
org2,jperez,admin
org2,rrodriguez,writer
org2,rrodriguez,editor
`

	err := createTempCSV("data/test.csv", csvTest)
	if err != nil {
		t.Errorf("error creando CSV: %v", err)
		return
	}
	defer deleteTempCSV("data/test.csv")

	orgsEsperado := []Organization{
		{
			Organization: "org1",
			Users: []User{
				{
					Username: "jperez",
					Roles:    []string{"admin", "superadmin"},
				},
				{
					Username: "asosa",
					Roles:    []string{"writer"},
				},
			},
		},
		{
			Organization: "org2",
			Users: []User{
				{
					Username: "jperez",
					Roles:    []string{"admin"},
				},
				{
					Username: "rrodriguez",
					Roles:    []string{"writer", "editor"},
				},
			},
		},
	}

	s := NewService()
	orgs, err := s.GetOrgFromFile("test.csv")
	if err != nil {
		t.Errorf("error en GetOrgFromFile: %v", err)
		return
	}

	if !assert.ObjectsAreEqual(orgsEsperado, orgs) {
		t.Errorf("las organizaciones obtenidas son diferentes de las esperadas\n [obtenidas]: %v \n [esperadas]: %v \n", orgs, orgsEsperado)
	}
}

func createTempCSV(filePath, data string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0770); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func deleteTempCSV(filePath string) {
	_ = os.Remove(filePath)
	_ = os.RemoveAll(filepath.Dir(filePath))
}
