package _const

// Support for all the driver is not ready yet.
// This is just a placeholder for all the driver.
// Community can add more driver here.
const (
	// sql
	MySQLDriver       string = "mysql"      // 100% support [Project | System]
	MariaDBDriver     string = "mariadb"    // 100% support [Project | System]
	SQLiteDriver      string = "sqlite"     // 100% support [Project | System]
	PostgresSQLDriver string = "postgresql" // 100% support [Project | System]
	SQLServerDriver   string = "sqlServer"
	OracleDriver      string = "oracle"

	// nosql
	ArangoDriver    string = "arangodb"
	MongoDBDriver   string = "mongodb"
	CouchbaseDriver string = "couchbase"
	CassandraDriver string = "cassandra"
	Neo4JDriver     string = "neo4j"
	CouchDBDriver   string = "couchdb"

	// Cloud Based
	DynamoDB  string = "dynamoDB"
	FireStore string = "firestore"

	// EmbeddedDB Embedded Database
	EmbeddedDB string = "embedded"

	// KeyValue database
	RedisDriver  string = "redis"
	BadgerDriver string = "badger"
	BoltDriver   string = "bolt"

	MemoryDb string = "memory" // usually its sync.Map{}
)
