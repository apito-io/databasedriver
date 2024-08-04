package badger

import (
	"context"
	"fmt"
	"log"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	"github.com/apito-io/databasedriver/migration"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

type BadgerDriver struct {
	Db *badger.DB
}

func (b *BadgerDriver) SearchResource(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[any], error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) FindOrganizationAdmin(ctx context.Context, orgId string) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) RunMigration() error {

	isMigrated, err := getValue[bool](b.Db, "migration", "migration")
	if err != nil {
		if !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}
	}

	if *isMigrated {
		fmt.Println("data already migrated, skipping ...")
		return nil
	}

	users := migration.GetMigrationUserData()
	project := migration.GetProjectInfo()
	// initial user data

	ctx := context.Background()

	// insert users
	for _, user := range users {
		_, err := b.CreateSystemUser(ctx, user)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// insert project
	_, err = b.CreateProject(ctx, project)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = b.setValue("migration", "migration", true)
	if err != nil {
		return nil
	}

	return nil
}

func GetSystemBadgerDriver(cfg *protobuff.DriverCredentials) (*BadgerDriver, error) {

	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions("./db/system/badger"))
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	return &BadgerDriver{Db: db}, nil
}
