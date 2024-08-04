package sql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/apito-io/buffers/protobuff"
	_const "github.com/apito-io/databasedriver"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type SqlDriver struct {
	ORM              *bun.DB
	DriverCredential *protobuff.DriverCredentials
}

func GetSQLDriver(driverCredentials *protobuff.DriverCredentials) (*SqlDriver, error) {

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

	return &SqlDriver{ORM: orm, DriverCredential: driverCredentials}, nil
}

type Meta struct {
	ID    string `bun:",pk,type:varchar(36),notnull" json:"id"`
	DocID string `bun:",notnull" json:"doc_id"`

	CreatedAt time.Time `bun:",notnull,default:current_date" json:"created_at"`
	UpdatedAt time.Time `bun:",notnull,default:current_date" json:"updated_at"`

	CreatedBy string `bun:",notnull" json:"created_by"`
	UpdatedBy string `bun:",notnull" json:"updated_by"`
	Status    string `bun:",notnull" json:"status"`
}

func (S *SqlDriver) RunMigration(ctx context.Context, projectId string) error {
	_, err := S.ORM.NewCreateTable().
		IfNotExists().
		Model((*Meta)(nil)).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
