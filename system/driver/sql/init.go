package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"os"

	"github.com/apito-io/buffers/protobuff"
	_const "github.com/apito-io/databasedriver"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type SystemSQLDriver struct {
	ORM              *bun.DB
	DriverCredential *protobuff.DriverCredentials
}

func GetSystemSQLDriver(driverCredentials *protobuff.DriverCredentials) (*SystemSQLDriver, error) {

	var orm *bun.DB

	switch driverCredentials.Engine {
	case _const.MySQLDriver:
		sqldb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			driverCredentials.User, driverCredentials.Password,
			driverCredentials.Host, driverCredentials.Database,
			driverCredentials.Port,
		))
		if err != nil {
			return nil, err
		}
		orm = bun.NewDB(sqldb, mysqldialect.New())
	case _const.SQLServerDriver:
		sqldb, err := sql.Open("sqlserver", fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			driverCredentials.User, driverCredentials.Password,
			driverCredentials.Host, driverCredentials.Port,
			driverCredentials.Database,
		))
		if err != nil {
			return nil, err
		}
		orm = bun.NewDB(sqldb, mssqldialect.New())
	case _const.SQLiteDriver:

		// Replace with your SQLite database file path
		dbPath := fmt.Sprintf("./%s.sqlite", driverCredentials.Database)

		// Check if the database file exists
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			// Create the database file
			file, err := os.Create(dbPath)
			if err != nil {
				return nil, err
			}
			file.Close()
			fmt.Printf("Database %s created\n", dbPath)
		}

		// Create a new database connection
		//sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
		sqldb, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, err
		}
		//defer sqldb.Close()
		orm = bun.NewDB(sqldb, sqlitedialect.New())
	case _const.PostgresSQLDriver:
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			driverCredentials.User, driverCredentials.Password,
			driverCredentials.Host, driverCredentials.Port,
			driverCredentials.Database,
		)
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		orm = bun.NewDB(sqldb, pgdialect.New())
		if err := orm.Ping(); err != nil {
			return nil, err
		}
	default:
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			driverCredentials.User, driverCredentials.Password,
			driverCredentials.Host, driverCredentials.Port,
			driverCredentials.Database,
		)
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		orm = bun.NewDB(sqldb, pgdialect.New())
		if err := orm.Ping(); err != nil {
			return nil, err
		}
	}

	return &SystemSQLDriver{ORM: orm, DriverCredential: driverCredentials}, nil
}

func (p *SystemSQLDriver) RunMigration(ctx context.Context) error {

	models := []interface{}{
		(*protobuff.SystemUser)(nil),
		(*protobuff.Project)(nil),
		(*protobuff.ProjectSchema)(nil),
		(*protobuff.PluginDetails)(nil),
		(*protobuff.APIToken)(nil),
		(*protobuff.DriverCredentials)(nil),
		(*protobuff.SystemMessage)(nil),
		(*protobuff.ModelType)(nil),
		(*protobuff.CloudFunction)(nil),
	}

	for _, model := range models {
		_, err := p.ORM.NewCreateTable().IfNotExists().Model(model).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
