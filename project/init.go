package project

import (
	"errors"

	"github.com/apito-io/buffers/interfaces"
	"github.com/apito-io/buffers/protobuff"
	_const "github.com/apito-io/databasedriver"
	"github.com/apito-io/databasedriver/project/driver/firestore"
	"github.com/apito-io/databasedriver/project/driver/sql"
)

func GetProjectDriver(engineConfig *protobuff.DriverCredentials) (interfaces.ProjectDBInterface, error) {
	var db interfaces.ProjectDBInterface
	var err error
	switch engineConfig.Engine {
	case _const.SQLiteDriver, _const.MySQLDriver, _const.PostgresSQLDriver, _const.SQLServerDriver:
		db, err = sql.GetSQLDriver(engineConfig)
	case _const.FireStore:
		db, err = firestore.GetFirestoreDriver(engineConfig)
	default:
		return nil, errors.New("unsupported database driver passed")
	}
	if err != nil {
		return nil, err
	}
	return db, nil
}
