package boltdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/databasedriver/migration"
	bolt "go.etcd.io/bbolt"
	"log"
)

type SystemBoltDBDriver struct {
	DB *bolt.DB
}

func GetSystemBoltDriver(engineConfig *protobuff.DriverCredentials) (*SystemBoltDBDriver, error) {

	// Open the system.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("./db/system.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	return &SystemBoltDBDriver{DB: db}, nil
}

func (b *SystemBoltDBDriver) RunMigration(ctx context.Context) error {
	isMigrated, err := getValue[bool](b.DB, "migration", "migration")
	if err != nil {
		if !errors.Is(bolt.ErrBucketNotFound, err) && err != nil {
			return err
		}
	}

	if isMigrated != nil && *isMigrated {
		fmt.Println("Data already migrated, skipping ...")
		return nil
	}

	// Get migration data
	users := migration.GetMigrationUserData()
	project := migration.GetProjectInfo()

	// Insert users
	for _, user := range users {
		_, err := b.createSystemUser(ctx, user)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// Insert project
	_, err = b.createProject(ctx, project)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Mark migration as complete
	err = b.setValue("migration", "migration", true)
	if err != nil {
		return err
	}

	return nil
}
