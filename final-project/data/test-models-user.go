package data

import (
	"time"
)

type UserTest struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
	IsAdmin   int
	CreatedAt time.Time
	UpdatedAt time.Time
	Plan      *Plan
}

func (u *UserTest) GetAll() ([]*User, error) {

	var users []*User

	user := User{
		ID:        1,
		Email:     "admin@test.com",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "password",
		Active:    1,
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users = append(users, &user)

	return users, nil
}

func (u *UserTest) GetByEmail(email string) (*User, error) {

	plan := Plan{
		ID:                  1,
		PlanAmount:          100,
		PlanName:            "Test Plan",
		PlanAmountFormatted: "$100",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	user := User{
		ID:        1,
		Email:     "admin@test.com",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "password",
		Active:    1,
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Plan:      &plan,
	}

	return &user, nil
}

func (u *UserTest) GetOne(id int) (*User, error) {

	plan := Plan{
		ID:                  1,
		PlanAmount:          100,
		PlanName:            "Test Plan",
		PlanAmountFormatted: "$100",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	user := User{
		ID:        1,
		Email:     "admin@test.com",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "password",
		Active:    1,
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Plan:      &plan,
	}

	return &user, nil
}

func (u *UserTest) Update(user User) error {

	return nil
}

func (u *UserTest) Delete() error {

	return nil
}

func (u *UserTest) DeleteByID(id int) error {

	return nil
}

func (u *UserTest) Insert(user User) (int, error) {

	return 2, nil
}

func (u *UserTest) ResetPassword(password string) error {

	return nil
}

func (u *UserTest) PasswordMatches(plainText string) (bool, error) {

	return true, nil
}
