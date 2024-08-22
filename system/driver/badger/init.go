package badger

import (
	"context"
	"errors"
	"fmt"
	"github.com/apito-io/databasedriver/migration"
	"log"

	"github.com/apito-io/buffers/protobuff"
	badger "github.com/dgraph-io/badger/v3"
)

type SystemBadgerDriver struct {
	Db *badger.DB
}

func GetSystemBadgerDriver(cfg *protobuff.DriverCredentials) (*SystemBadgerDriver, error) {

	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions("./db/system/badger"))
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	return &SystemBadgerDriver{Db: db}, nil
}

func (b *SystemBadgerDriver) RunMigration(ctx context.Context) error {

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

	// insert users
	for _, user := range users {
		_, err := b.createSystemUser(ctx, user)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// insert project
	_, err = b.createProject(ctx, project)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = b.setValue("migration", "migration", true)
	if err != nil {
		return nil
	}

	return nil
}
