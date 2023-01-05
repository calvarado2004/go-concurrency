package data

import "database/sql"

func TestNew(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: &UserTest{},
		Plan: &PlanTest{},
	}
}
