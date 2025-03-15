package db

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

// TestMain sets up the test environment by loading the mock data.
func TestMain(m *testing.M) {
	withErrorExit := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	withErrorExit(db.Ping())

	// drop the users and the scores tables if they exist
	_, err := db.Query("DROP TABLE IF EXISTS scores, users;")
	withErrorExit(err)

	// recreate the tables
	err = checkTablesExit()
	withErrorExit(err)

	// load mock data users mock data
	loadData := func(filePath, table string) error {
		mysql.RegisterLocalFile(filePath)
		_, err = db.Exec(`LOAD DATA LOCAL INFILE '` + filePath +
			`' INTO TABLE ` + table +
			` FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n' IGNORE 1 LINES;`)
		return err
	}

	withErrorExit(loadData("sample_data_users.csv", "users"))
	withErrorExit(loadData("sample_data_scores.csv", "scores"))

	os.Exit(m.Run())
}

// TestCreateLevelScore tests the functionality of createLevelScore
func TestCreateLevelScore(t *testing.T) {
	Convey("TestCreateLevelScore: Given the UserInfor to create level scores with correct data", t, func() {
		Convey("recreating game_level and user_id combination that already exist should return"+
			"a value that implements an error interface", func() {
			user := &UserInfor{TapooID: "Vf2TqN5MB", Level: 1}
			err := user.createLevelScore("sample_uuid_value")

			So(err, ShouldNotBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "Duplicate entry 'Vf2TqN5MB-1' for key 'user_id'")
		})

		Convey("creating a new game_level and user_id combination should return a nil value", func() {
			user := &UserInfor{TapooID: "06PE0LPzyCL", Level: 20}
			err := user.createLevelScore("sample_uuid_value")

			So(err, ShouldBeNil)

			data, err := user.getLevelScore()

			So(err, ShouldBeNil)
			So(data.TapooID, ShouldEqual, "06PE0LPzyCL")
			So(data.Level, ShouldEqual, 20)
			So(data.HighScores, ShouldEqual, 0)
		})
	})
}

// TestGetLevelScore tests the functionality of getLevelScore
func TestGetLevelScore(t *testing.T) {
	Convey("TestGetLevelScore: Given the UserInfor to get level scores with", t, func() {
		Convey("closed db connection is used, should return a value that implements"+
			" an error interface", func() {
			copyOfDb := cloneDb()
			db.Close()

			user := &UserInfor{TapooID: "VZWeOq2p", Level: 1}
			data, err := user.getLevelScore()

			db = copyOfDb

			So(db.Ping(), ShouldBeNil)

			So(data, ShouldEqual, nil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "sql: database is closed")
		})

		Convey("the tapoo id entry does not exist an error should be returned, "+
			"should return a value that implements an error interface", func() {
			user := &UserInfor{TapooID: "VZW2eOq2p", Level: 1}
			data, err := user.getLevelScore()

			So(data, ShouldResemble, new(LevelScoreResponse))
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "sql: no rows in result set")
		})

		Convey("variables whose tapoo ID entry exists in the db should return a nil value error", func() {
			user := &UserInfor{TapooID: "VZWeOq2p", Level: 1}
			data, err := user.getLevelScore()

			So(err, ShouldBeNil)
			So(data.Email, ShouldEqual, "")
			So(data.CreatedAt, ShouldHappenBefore, time.Now())
			So(data.UpdateAt, ShouldHappenBefore, time.Now())
			So(data.TapooID, ShouldEqual, "VZWeOq2p")
			So(data.Level, ShouldEqual, 1)
			So(data.HighScores, ShouldEqual, 533)
		})
	})
}

// TestGetOrCreateLevelScore tests the functionality of GetOrCreateLevelScore
func TestGetOrCreateLevelScore(t *testing.T) {
	errfunc := func(info *UserInfor, errMsg string) {
		scores, err := info.GetOrCreateLevelScore()

		So(err, ShouldNotBeNil)
		So(err, ShouldImplement, (*error)(nil))
		So(err.Error(), ShouldContainSubstring, errMsg)
		So(scores, ShouldBeNil)
	}

	Convey("TestGetOrCreateLevelScore: Given UserInfo data to get or create level scores with", t, func() {
		Convey("game level less than zero, a value that implements an error interface should be returned", func() {
			errfunc(&UserInfor{TapooID: "VZWeOq2p", Level: -1}, "invalid game level found : '-1'")
		})

		Convey("an empty tapoo ID, a value that implements an error interface should be returned", func() {
			errfunc(&UserInfor{TapooID: "", Level: 1}, "invalid Tapoo ID found : '(empty)'")
		})

		Convey("tapoo ID longer than 64 characters, a value that implements an error"+
			" interface should be returned", func() {
			user := &UserInfor{TapooID: "a6a1-b43d84afa437-d3140030-4a5c-4352-9a8a-8fe4d988502-9a8a-8fe4d9", Level: 3}

			errfunc(user, "invalid Tapoo ID found : 'a6a1-b43d8... (Too long)'")
		})

		Convey("tapoo ID that is not a foerign key in users table, a value that implements an"+
			" error interface in returned", func() {
			errfunc(&UserInfor{TapooID: "VZWe   Oq2p", Level: 3}, "Cannot add or update a child row: a foreign key constraint fails")
		})

		Convey("all variables correctly used and have no invalid characters, "+
			"the error value returned should be nil", func() {
			user := &UserInfor{TapooID: "06PE0LPzyCL", Level: 3}
			scores, err := user.GetOrCreateLevelScore()

			So(err, ShouldBeNil)
			So(scores.HighScores, ShouldEqual, 1203)
			So(scores.TapooID, ShouldEqual, "06PE0LPzyCL")
			So(scores.Email, ShouldEqual, "")
			So(scores.Level, ShouldEqual, 3)
			So(scores.CreatedAt, ShouldHappenBefore, time.Now())
			So(scores.UpdateAt, ShouldHappenBefore, time.Now())
		})
	})
}

// TestGetTopFiveScores tests the functionality of GetTopFiveScores
func TestGetTopFiveScores(t *testing.T) {
	Convey("TestGetTopFiveScores: Given the UserInfor to fetch top five scores with ", t, func() {
		Convey("the game level as a value less than zero, a value that implements "+
			"an error interface should be returned", func() {
			user := &UserInfor{Level: -23, TapooID: "fake_tapoo_id"}
			data, err := user.GetTopFiveScores()

			So(data, ShouldHaveLength, 0)
			So(err, ShouldNotBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "invalid game level found : '-23'")
		})

		Convey("the database connection used bieng invalid, a value that "+
			"implements an error interface should be returned", func() {
			copyOfDb := cloneDb()
			db.Close()

			user := &UserInfor{Level: 3, TapooID: "fake_tapoo_id"}
			data, err := user.GetTopFiveScores()

			db = copyOfDb

			So(db.Ping(), ShouldBeNil)

			So(data, ShouldHaveLength, 0)
			So(err, ShouldNotBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "sql: database is closed")
		})

		Convey("the valid database connection and game level value used, the error value"+
			" returned should be equal to nil", func() {
			user := &UserInfor{Level: 2, TapooID: ""}
			data, err := user.GetTopFiveScores()

			// test data
			topScores := []int{1948, 1653, 1616, 1584, 1027}
			userIDs := []string{"GzlWAL0mP", "Vf2TqN5MB", "FANVZWeOq2p", "FbnnuznkFAN", "Fbn56nuznk"}
			userEmails := []string{"ckumaar0@tripod.com", "sgravell1@europa.eu",
				"test.user@naihub.com", "sclaussen3@cam.ac.uk", "sample.user@niahub.com"}

			So(err, ShouldBeNil)
			So(data, ShouldHaveLength, 5)

			for _, item := range data {
				So(topScores, ShouldContain, item.HighScores)
				So(userIDs, ShouldContain, item.TapooID)
				So(userEmails, ShouldContain, item.Email)
			}
		})
	})
}

// TestUpdateLevelScores tests the functionality of UpdateLevelScores
func TestUpdateLevelScores(t *testing.T) {
	errfunc := func(info *UserInfor, highScores int, errMsg string) {
		err := info.UpdateLevelScore(highScores)

		So(err, ShouldNotBeNil)
		So(err, ShouldImplement, (*error)(nil))
		So(err.Error(), ShouldContainSubstring, errMsg)
	}

	Convey("TestUpdateLevelScores: Given the UserInfor and High Scores with", t, func() {
		Convey("the game level less than zero, a value that implements "+
			"an error interface should be returned", func() {
			user := &UserInfor{Level: -12, TapooID: "82hsgdj"}
			errfunc(user, 3223, "invalid game level found : '-12'")
		})

		Convey("the high Scores less than zero, a value that implements "+
			"an error interface should be returned", func() {
			user := &UserInfor{Level: 23, TapooID: "82hsgdj"}
			errfunc(user, -326, "invalid high scores found : '-326'")
		})

		Convey("the tapoo ID empty, a value that implements "+
			"an error interface should be returned", func() {
			user := &UserInfor{Level: 23, TapooID: ""}
			errfunc(user, 326, "invalid Tapoo ID found : '(empty)'")
		})

		Convey("the tapoo ID longer than 64 charactes, a value that implements "+
			"an error interface should be returned", func() {
			user := &UserInfor{Level: 23,
				TapooID: "a6a1-b43d84afa437-d3140030-4a5c-4352-9a8a-8fe4d988502-9a8a-8fe4d9"}

			errfunc(user, 326, "invalid Tapoo ID found : 'a6a1-b43d8... (Too long)'")
		})

		Convey("the correct values provided, the error value returned should be nil", func() {
			user := &UserInfor{Level: 1, TapooID: "VZWeOq2p"}
			err := user.UpdateLevelScore(1000)

			So(err, ShouldBeNil)

			data, err := user.getLevelScore()

			So(err, ShouldBeNil)
			So(data.HighScores, ShouldEqual, 1000)
			So(data.Level, ShouldEqual, 1)
		})
	})
}
