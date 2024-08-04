package system

import (
	"errors"

	"github.com/apito-io/buffers/interfaces"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	_const "github.com/apito-io/databasedriver"
	"github.com/apito-io/databasedriver/system/driver/badger"
	"github.com/apito-io/databasedriver/system/driver/sql"
)

func GetSystemDriver(engineConfig *protobuff.DriverCredentials, conf *shared.CommonDatabaseConfig) (interfaces.SystemDBInterface, error) {

	if engineConfig == nil {
		return nil, errors.New("please select a system database engine")
	}

	if engineConfig.Engine == "" {
		engineConfig = &protobuff.DriverCredentials{
			Engine:   _const.PostgresSQLDriver,
			Host:     conf.Host,
			Port:     conf.Port,
			User:     conf.User,
			Password: conf.Password,
			Database: conf.Database,
		}
	}

	var db interfaces.SystemDBInterface
	var err error

	switch engineConfig.Engine {
	case _const.PostgresSQLDriver, _const.MySQLDriver:
		db, err = sql.GetSystemSQLDriver(engineConfig)
		if err != nil {
			return nil, err
		}
		break
	case _const.EmbeddedDB:
		db, err = badger.GetSystemBadgerDriver(engineConfig)
	default: // default set embedded database
		db, err = badger.GetSystemBadgerDriver(engineConfig)
	}

	// run db migration for all the db
	if db != nil {
		err = db.RunMigration()
	}

	if err != nil {
		return nil, err
	}
	return db, nil
}
