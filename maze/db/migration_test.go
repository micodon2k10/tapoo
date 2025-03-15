package db

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestCreateDbConnection tests the functionality of createDbConnection
func TestCreateDbConnection(t *testing.T) {
	copyOfConfig := config

	Convey("TestCreateDbConnection: Given the mysql database configuration", t, func() {
		Convey("with empty values, a value that implements an error interface should be returned ", func() {
			config = new(dbConfig)
			err := createDbConnection()

			// Restore db configuration
			config = copyOfConfig

			So(err, ShouldNotEqual, nil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, `sql: unknown driver "" (forgotten import?)`)
		})

		Convey("with values that are not empty but are incorrect, a value that implements an error "+
			"interface should be returned when a ping is made", func() {
			config = &dbConfig{
				DbHost:         "localhost:3306",
				DbName:         "test_db",
				DbUserName:     "test_user_tapoo",
				DbUserPassword: "fake_password",
				Driver:         "mysql",
			}
			err := createDbConnection()

			// Restore db configuration
			config = copyOfConfig

			So(err, ShouldNotEqual, nil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "Access denied for user 'test_user_tapoo'@")
		})

		Convey("with the correct values a nil error value should be returned", func() {
			err := createDbConnection()

			So(err, ShouldEqual, nil)
		})
	})
}

// TestCheckTablesExist tests the functionality of checkTablesExit
func TestCheckTablesExist(t *testing.T) {
	var (
		copyOfConfig = config

		dropTable = func() error {
			_, err := db.Query("DROP TABLE IF EXISTS scores, users")
			return err
		}
	)

	Convey("TestCheckTablesExist: Given the database configuration ", t, func() {
		Convey("that triggers an error message other than 'no rows in result set' "+
			"that error should be returned", func() {
			config = new(dbConfig)
			err := checkTablesExit()

			config = copyOfConfig

			So(err, ShouldNotBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "Table 'users' already exists")
		})

		Convey("that triggers 'no rows in result set' the missing table(s) should be created ", func() {
			err := dropTable()

			So(err, ShouldBeNil)

			err = checkTablesExit()

			So(err, ShouldBeNil)
		})
	})
}

// TestGetEnvVars tests the functionality of getEnvVars
func TestGetEnvVars(t *testing.T) {
	copyOfConfig := config
	config = new(dbConfig)

	// resetEnvVars set the environment variables with the provided variables
	resetEnvVars := func(key, value string) {
		os.Setenv(key, value)
	}

	// unsetEnvVars removes the db configurations environment variables
	unsetEnvVars := func() {
		os.Unsetenv("TAPOO_DB_NAME")
		os.Unsetenv("TAPOO_DB_USER_NAME")
		os.Unsetenv("TAPOO_DB_USER_PASSWORD")
		os.Unsetenv("TAPOO_DB_HOST")
	}

	unsetEnvVars()

	Convey("TestGetEnvVars: Given the db configuration environment variables", t, func() {
		Convey("with only TAPOO_DB_NAME a value that implements an error interface should be returned", func() {
			err := getEnvVars()

			So(err, ShouldNotBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "envVars: TAPOO_DB_NAME")

			So(config.DbName, ShouldBeEmpty)
			So(config.DbUserName, ShouldBeEmpty)
			So(config.DbUserPassword, ShouldBeEmpty)
			So(config.DbHost, ShouldBeEmpty)
			So(config.Driver, ShouldBeEmpty)
		})

		Convey("with TAPOO_DB_USER_NAME added, a value that implements an error interface should be returned", func() {
			resetEnvVars("TAPOO_DB_NAME", copyOfConfig.DbName)

			err := getEnvVars()

			So(err, ShouldNotBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "envVars: TAPOO_DB_USER_NAME")

			So(config.DbName, ShouldNotBeEmpty)
			So(config.DbUserName, ShouldBeEmpty)
			So(config.DbUserPassword, ShouldBeEmpty)
			So(config.DbHost, ShouldBeEmpty)
			So(config.Driver, ShouldBeEmpty)
		})

		Convey("with TAPOO_DB_USER_PASSWORD added, a value that implements an error interface should be returned", func() {
			resetEnvVars("TAPOO_DB_USER_NAME", copyOfConfig.DbUserName)

			err := getEnvVars()

			So(err, ShouldNotBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "envVars: TAPOO_DB_USER_PASSWORD")

			So(config.DbName, ShouldNotBeEmpty)
			So(config.DbUserName, ShouldNotBeEmpty)
			So(config.DbUserPassword, ShouldBeEmpty)
			So(config.DbHost, ShouldBeEmpty)
			So(config.Driver, ShouldBeEmpty)
		})
		Convey("with TAPOO_DB_HOST added, a value that implements an error interface should be returned", func() {
			resetEnvVars("TAPOO_DB_USER_PASSWORD", copyOfConfig.DbUserPassword)

			err := getEnvVars()

			So(err, ShouldNotBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "envVars: TAPOO_DB_HOST")

			So(config.DbName, ShouldNotBeEmpty)
			So(config.DbUserName, ShouldNotBeEmpty)
			So(config.DbUserPassword, ShouldNotBeEmpty)
			So(config.DbHost, ShouldBeEmpty)
			So(config.Driver, ShouldBeEmpty)
		})

		Convey("with all environment variables set, no error that should be thrown", func() {
			resetEnvVars("TAPOO_DB_HOST", copyOfConfig.DbHost)
			err := getEnvVars()

			So(err, ShouldBeNil)

			So(config.DbName, ShouldNotBeEmpty)
			So(config.DbUserName, ShouldNotBeEmpty)
			So(config.DbUserPassword, ShouldNotBeEmpty)
			So(config.DbHost, ShouldNotBeEmpty)
			So(config.Driver, ShouldNotBeEmpty)
		})
	})
}
