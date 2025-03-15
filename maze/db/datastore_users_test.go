package db

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestCreateUser tests the functionality of createUser
func TestCreateUser(t *testing.T) {
	Convey("TestCreateUser: Given the UserInfor when creating a user with", t, func() {
		Convey("values that already exist in the database, a value that implements an error"+
			" interface should be returned", func() {
			user := &UserInfor{TapooID: "FbnnuznkFAN"}
			err := user.createUser("ea49be59-b553-430c-a706-7860dcb3ea12")

			So(err, ShouldNotBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "Duplicate entry 'ea49be59-b553-430c-a706-7860dcb3ea12'")
		})

		Convey("values that have no invalid characters, a nil error value should"+
			" be returned", func() {
			user := &UserInfor{TapooID: "9a9a-7a9a5808e086", Email: "test@naihub.com"}
			err := user.createUser("f538ab54-1692-41bf-9a9a-7a98808e086d")

			So(err, ShouldBeNil)

			data, err := user.getUser()

			So(err, ShouldBeNil)
			So(data.TapooID, ShouldEqual, "9a9a-7a9a5808e086")
			So(data.Email, ShouldEqual, "test@naihub.com")
		})
	})
}

// TestGetUser tests the functionality of getUser
func TestGetUser(t *testing.T) {
	errFunc := func(user *UserInfor, errMsg string) {
		data, err := user.getUser()

		So(err, ShouldNotBeNil)
		So(data, ShouldResemble, new(UserInfoResponse))
		So(err, ShouldImplement, (*error)(nil))
		So(err.Error(), ShouldContainSubstring, errMsg)
	}

	Convey("TestGetUser: Given the UserInfor when fetching a user with", t, func() {
		Convey("the tapoo id provided that does not exist, a value that implements an error"+
			" interface should be returned", func() {
			user := &UserInfor{Level: 23, TapooID: "fake_sample_id"}

			errFunc(user, "sql: no rows in result set")
		})

		Convey("closed db connections, a value that implements an error interface should"+
			" be returned", func() {
			copyOfDb := cloneDb()
			db.Close()
			user := &UserInfor{Level: 2, TapooID: "fake_sample_id"}
			data, err := user.getUser()

			db = copyOfDb

			So(db.Ping(), ShouldBeNil)

			So(err, ShouldNotBeNil)
			So(data, ShouldBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "sql: database is closed")
		})

		Convey("the values used are properly escaped and the tapoo id exists in the db, "+
			"a nil error value should be returned", func() {
			user := &UserInfor{Level: 18, TapooID: "GzlWAL0mP"}

			data, err := user.getUser()

			So(err, ShouldBeNil)
			So(data.CreatedAt, ShouldHappenBefore, time.Now())
			So(data.UpdateAt, ShouldHappenBefore, time.Now())
			So(data.Email, ShouldEqual, "ckumaar0@tripod.com")
			So(data.TapooID, ShouldEqual, "GzlWAL0mP")
		})
	})
}

// TestGetOrCreateUser tests the functionality of GetOrCreateUser
func TestGetOrCreateUser(t *testing.T) {
	errFunc := func(user *UserInfor, errMsg string) {
		data, err := user.GetOrCreateUser()

		So(err, ShouldNotBeNil)
		So(data, ShouldBeNil)
		So(err, ShouldImplement, (*error)(nil))
		So(err.Error(), ShouldContainSubstring, errMsg)
	}

	Convey("TestGetOrCreateUser: Given the UserInfor when fetching or creating"+
		" a user with", t, func() {
		Convey("the empty tapoo ID, a value that implements an error interface"+
			" should be returned", func() {
			user := &UserInfor{TapooID: ""}

			errFunc(user, "invalid Tapoo ID found : '(empty)'")
		})

		Convey("the tapoo ID longer that 64 characters, a value that implements an error"+
			"interface should be returned", func() {
			user := &UserInfor{TapooID: "2af80406-5be2-4569-afba-b14e861-2af81406-5b42-893d-afba-6734grywu"}

			errFunc(user, "invalid Tapoo ID found : '2af80406-5... (Too long)'")
		})

		Convey("the email longer than 64 characters, a value that implements an error"+
			" interface should be returned", func() {
			user := &UserInfor{TapooID: "SWfddew34",
				Email: "2af80406-5be2-4569-afba-b14e861-2af81406-5b42-893d-a3f@niahub.com"}

			errFunc(user, "invalid Email found : '2af80406-5... (Too long)'")
		})

		Convey("the db connection that is invalid, a value that implements an error"+
			" interface should be returned", func() {
			copyOfDb := cloneDb()
			user := &UserInfor{TapooID: "SWfddew34"}

			db.Close()
			data, err := user.GetOrCreateUser()
			db = copyOfDb

			So(db.Ping(), ShouldBeNil)

			So(err, ShouldNotBeNil)
			So(data, ShouldBeNil)
			So(err, ShouldImplement, (*error)(nil))
			So(err.Error(), ShouldContainSubstring, "sql: database is closed")
		})

		Convey("the correct user infor used and the database connection is not invalid"+
			" the error value returned should be a nil value", func() {
			user := &UserInfor{TapooID: "FANVZWeOq2p"}
			data, err := user.GetOrCreateUser()

			So(err, ShouldBeNil)
			So(data.Email, ShouldEqual, "test.user@naihub.com")
			So(data.TapooID, ShouldEqual, "FANVZWeOq2p")
			So(data.CreatedAt, ShouldHappenBefore, time.Now())
			So(data.UpdateAt, ShouldHappenBefore, time.Now())
		})
	})
}

// TestUpdateUser tests the functionality of UpdateUser
func TestUpdateUser(t *testing.T) {
	errFunc := func(user *UserInfor, errMsg string) {
		err := user.UpdateUser()

		So(err, ShouldNotBeNil)
		So(err, ShouldImplement, (*error)(nil))
		So(err.Error(), ShouldContainSubstring, errMsg)
	}

	Convey("TestUpdateUser: Given the UserInfor while updating the user with", t, func() {
		Convey("the empty tapoo ID, a value that implements the error interface is"+
			" returned", func() {
			user := &UserInfor{TapooID: "", Email: "sample_user@naihub.com"}

			errFunc(user, "invalid Tapoo ID found : '(empty)'")
		})

		Convey("the tapoo ID having more that 64 characters, a value that implements"+
			" the error interface is returned", func() {
			user := &UserInfor{Email: "sample_user@naihub.com",
				TapooID: "2af80406-5be2-4569-afba-b14e861-2af81406-5b42-893d-afba-6734grywu"}

			errFunc(user, "invalid Tapoo ID found : '2af80406-5... (Too long)'")
		})

		Convey("the empty email, a value that implements the error interface is"+
			" returned", func() {
			user := &UserInfor{TapooID: "f80406-5be2", Email: ""}

			errFunc(user, "invalid Email found : '(empty)'")
		})

		Convey("the email having more that 64 characters, a value that implements"+
			" the error interface is returned", func() {
			user := &UserInfor{TapooID: "SWfddew34",
				Email: "2af80406-5be2-4569-afba-b14e861-2af81406-5b42-893d-a3f@niahub.com"}

			errFunc(user, "invalid Email found : '2af80406-5... (Too long)'")
		})

		Convey("the correct values used, a nil value error should be returned", func() {
			user := &UserInfor{TapooID: "Vf2TqN5MB", Email: "sample_user@naihub.com"}
			err := user.UpdateUser()

			So(err, ShouldBeNil)

			data, err := user.getUser()

			So(err, ShouldBeNil)
			So(data.Email, ShouldEqual, "sample_user@naihub.com")
			So(data.TapooID, ShouldEqual, "Vf2TqN5MB")
		})
	})
}

// TestExecPrepStmts tests the functionality of execPrepStmts
func TestExecPrepStmts(t *testing.T) {
	errFunc := func(err error, errMsg string) {
		So(err, ShouldNotBeNil)
		So(err, ShouldImplement, (*error)(nil))
		So(err.Error(), ShouldContainSubstring, errMsg)
	}

	Convey("TestExecPrepStmts: Given a query and its other metadata with", t, func() {
		Convey("closed database connection, a value that implements an error interface should be returned", func() {
			copyOfDb := cloneDb()
			db.Close()

			rows, _, err := execPrepStmts(multiRows, "SELECT * FROM users;", "")
			db = copyOfDb

			So(rows, ShouldBeNil)

			errFunc(err, "sql: database is closed")

			// Re-establish a connection if necessary
			So(db.Ping(), ShouldBeNil)
		})

		Convey("singleRow queryType found no resultSet data match, a value that "+
			"implements the error interface should be returned", func() {
			_, row, err := execPrepStmts(singleRow, "SELECT email FROM users WHERE id = ?;", "VZW7274Oq2p")

			So(row, ShouldNotBeNil)
			So(err, ShouldBeNil)

			err = row.Scan(nil)

			errFunc(err, "sql: no rows in result set")
		})

		Convey("query missing some arguments, a value that "+
			"implements the error interface should be returned", func() {
			_, _, err := execPrepStmts(noReturnVal,
				"UPDATE scores SET high_scores = ? WHERE game_level = ? and user_id = ?;", "1000", "12")

			errFunc(err, "sql: expected 3 arguments, got 2")
		})

		Convey("query missing the only argument, a value that "+
			"implements the error interface should be returned", func() {
			_, row, err := execPrepStmts(singleRow, "SELECT email FROM users WHERE id = ?")

			So(row, ShouldNotBeNil)
			So(err, ShouldBeNil)

			err = row.Scan(nil)

			errFunc(err, "Error 1064: You have an error in your SQL syntax;")
		})

		Convey("query having extra arguments, a value that "+
			"implements the error interface should be returned", func() {
			rows, _, err := execPrepStmts(multiRows, "SELECT email FROM users WHERE id LIKE ?;", "V", "S")

			So(rows, ShouldBeNil)

			errFunc(err, "sql: expected 1 arguments, got 2")
		})

		Convey("queryType that is non existent, a value that "+
			"implements the error interface should be returned", func() {
			_, _, err := execPrepStmts(5, "SELECT email FROM users")

			errFunc(err, "invalid queryType found : '5'")
		})

		Convey("noReturnVal queryType having the correct values, should return a nil error value", func() {
			_, _, err := execPrepStmts(noReturnVal,
				"UPDATE scores SET high_scores = ? WHERE game_level = ? and user_id = ?;", "1000", "12", "VZWeOq2p")

			So(err, ShouldBeNil)

			user := &UserInfor{Level: 12, TapooID: "VZWeOq2p"}
			data, err := user.getLevelScore()

			So(err, ShouldBeNil)
			So(data.HighScores, ShouldEqual, 1000)
		})

		Convey("singleRow queryType having the correct values, should return the fetched data and a nil error value", func() {
			d := UserInfoResponse{}
			_, row, err := execPrepStmts(singleRow, "SELECT email FROM users WHERE id = ?;", "VZWeOq2p")
			So(err, ShouldBeNil)

			err = row.Scan(&d.Email)

			So(err, ShouldBeNil)
			So(d.Email, ShouldEqual, "asainsberry4@amazon.com")
		})

		Convey("multiRows queryType having the correct values, should return the fetched data and a nil value error", func() {
			rows, _, err := execPrepStmts(multiRows, "SELECT email FROM users LIMIT 5;")

			So(err, ShouldBeNil)

			count := 0

			for rows.Next() {
				count++
			}

			So(count, ShouldEqual, 5)
		})
	})
}

// cloneDb makes a deep copy of the database connection
// that is used exclusively for testing.
func cloneDb() *sql.DB {
	x := reflect.ValueOf(db)
	copy := &sql.DB{}
	starX := x.Elem()
	y := reflect.New(starX.Type())
	starY := y.Elem()
	starY.Set(starX)
	reflect.ValueOf(copy).Elem().Set(y.Elem())
	return copy
}
