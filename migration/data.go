package migration

import (
	"fmt"

	"github.com/apito-io/buffers/protobuff"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetProjectInfo() *protobuff.Project {
	/*envMap, err := godotenv.Read(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}*/

	return &protobuff.Project{
		XKey:        "fitness_app_jh478",
		Id:          "fitness_app_jh478",
		Name:        "Fitness App",
		Description: "A Fitness Tracker App",
		Locals:      []string{"en"},
	}
}

func GetMigrationUserData() []*protobuff.SystemUser {

	hash, err := bcrypt.GenerateFromPassword([]byte("#ApitoRocks#"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
	}

	return []*protobuff.SystemUser{
		{
			Id:               uuid.New().String(),
			FirstName:        "System Admin",
			Email:            "admin@apito.io",
			Username:         "admin",
			Secret:           string(hash),
			CurrentProjectId: "fitness_app_jh478",
		},
	}
}
