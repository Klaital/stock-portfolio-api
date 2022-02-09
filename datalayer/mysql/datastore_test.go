package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"testing"
)

type DataStoreSuite struct {
	suite.Suite
	Store          *DataStore
	dockerpool     *dockertest.Pool
	testDbResource *dockertest.Resource
	testdb         *sql.DB
}

func TestDataStoreSuite(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skip datastore tests")
	}

	suite.Run(t, new(DataStoreSuite))
}

func (suite *DataStoreSuite) SetupSuite() {
	// Start up a database to test against
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.WithError(err).Fatal("Failed to init dockertest pool")
	}
	suite.dockerpool = pool

	suite.testDbResource, err = pool.Run("mysql", "5.7", []string{"MYSQL_DATABASE=autotest", "MYSQL_USER=autotester", "MYSQL_PASSWORD=autotest123", "MYSQL_ROOT_PASSWORD=meh"})
	if err != nil {
		log.WithError(err).Fatal("Failed to init testDbResource")
	}
	if err = suite.testDbResource.Expire(5 * 60); err != nil {
		log.WithError(err).Fatal("Failed to set expiration timeout on test DB")
	}
	portNum, err := strconv.ParseInt(suite.testDbResource.GetPort("3306/tcp"), 10, 32)
	if err != nil {
		log.WithField("port", suite.testDbResource.GetPort("3306/tcp")).WithError(err).Fatal("Failed to parse db port")
	}

	if err := pool.Retry(func() error {
		var err error
		suite.testdb, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true", "autotester", "autotest123", "localhost", portNum, "autotest"))
		if err != nil {
			return err
		}
		return suite.testdb.Ping()
	}); err != nil {
		log.WithError(err).Fatal("Could not connect to db")
	}

	// Run the DB schema migrations
	migrateDriver, err := mysql.WithInstance(suite.testdb, &mysql.Config{})
	if err != nil {
		log.WithError(err).Fatal("Failed to create migrations driver")
	}
	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", migrateDriver)
	if err != nil {
		log.WithError(err).Fatal("Failed to init db migrator")
	}

	if err = migrator.Up(); err != nil {
		log.WithError(err).Fatal("Failed to update test DB schema")
	}

	// Initialize the DataStore now that the DB is ready
	suite.Store, err = New(context.Background(), "localhost", "autotester", "autotest123", "autotest", int(portNum), bcrypt.MinCost)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to test DB")
	}
}

func (suite *DataStoreSuite) SetupTest() {
	// truncate all the tables before each test so that test data does not leak
	tablesToTruncate := []string{"positions", "users"}
	for _, table := range tablesToTruncate {
		_, err := suite.testdb.Exec(fmt.Sprintf(`DELETE FROM %s`, table))
		if err != nil {
			log.WithField("table", table).WithError(err).Fatal("Failed to truncate table")
		}
	}
}

func (suite *DataStoreSuite) TearDownSuite() {
	if err := suite.dockerpool.Purge(suite.testDbResource); err != nil {
		log.WithError(err).Fatal("Failed to purge test DB")
	}
}
