package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	// noReturnVal indicates the sql query being executed should not
	// return any value. Return value not expected.
	noReturnVal int = iota

	// singleRow indicates the sql query bieng executed should only
	// return a single row of the expected result set.
	singleRow

	// multiRows indicates the sql query being executed should return
	// multiple rows of the expected result set.
	multiRows
)

// UserInfoResponse defines the expected response from users.
type UserInfoResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	TapooID   string    `json:"id"`
	UpdateAt  time.Time `json:"updated_at"`
}

// createUser creates a new user using the tapoo ID provided.
func (u *UserInfor) createUser(uuid string) error {
	query := `INSERT INTO users (uuid, id, email) VALUES (?, ?, ?);`

	_, _, err := execPrepStmts(noReturnVal, query, uuid, u.TapooID, u.Email)
	return err
}

// getUser checks if the tapoo ID provided exists in users.
func (u *UserInfor) getUser() (*UserInfoResponse, error) {
	query := `SELECT id, email, created_at, updated_at FROM users WHERE id = ?;`

	var d UserInfoResponse

	_, row, err := execPrepStmts(singleRow, query, u.TapooID)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&d.TapooID, &d.Email, &d.CreatedAt, &d.UpdateAt)

	return &d, err
}

// GetOrCreateUser creates the new user with tapoo ID provided if the
// it does not exists. Email used can be empty or not.
func (u *UserInfor) GetOrCreateUser() (*UserInfoResponse, error) {
	switch {
	case len(u.TapooID) == 0:
		return nil, fmt.Errorf(invalidData, "Tapoo ID", u.TapooID+"(empty)")

	case len(u.TapooID) > 64:
		return nil, fmt.Errorf(invalidData, "Tapoo ID", u.TapooID[:10]+"... (Too long)")

	case len(u.Email) > 64:
		return nil, fmt.Errorf(invalidData, "Email", u.Email[:10]+"... (Too long)")
	}

	u4, err := uuid.NewV4()
	if err != nil {
		return nil, errGenUUID
	}

	err = u.createUser(u4.String())

	switch {
	case strings.Contains(err.Error(), "Duplicate entry"):
	default:
		return nil, err
	}

	return u.getUser()
}

// UpdateUser should update the tapoo user information.
// While updating a user, the email should not be empty otherwise
// an error will be returned.
func (u *UserInfor) UpdateUser() error {
	switch {
	case len(u.TapooID) == 0:
		return fmt.Errorf(invalidData, "Tapoo ID", u.TapooID+"(empty)")

	case len(u.TapooID) > 64:
		return fmt.Errorf(invalidData, "Tapoo ID", u.TapooID[:10]+"... (Too long)")

	case len(u.Email) == 0:
		return fmt.Errorf(invalidData, "Email", u.Email+"(empty)")

	case len(u.Email) > 64:
		return fmt.Errorf(invalidData, "Email", u.Email[:10]+"... (Too long)")
	}

	query := `UPDATE users SET email = ? WHERE id = ?;`

	_, _, err := execPrepStmts(noReturnVal, query, u.Email, u.TapooID)
	return err
}

// execPrepStmts executes the Prepared statement for the sql queries.
func execPrepStmts(queryType int, sqlQuery string, val ...interface{}) (*sql.Rows, *sql.Row, error) {
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		return nil, nil, err
	}

	defer stmt.Close()

	switch queryType {
	case noReturnVal:
		_, err = db.Exec(sqlQuery, val...)
		return nil, nil, err

	case singleRow:
		row := db.QueryRow(sqlQuery, val...)
		return nil, row, nil

	case multiRows:
		rows, err := db.Query(sqlQuery, val...)
		return rows, nil, err

	default:
		return nil, nil,
			fmt.Errorf(invalidData, "queryType", queryType)
	}
}
