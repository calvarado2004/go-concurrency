package data

type UserInterface interface {
	GetAll() ([]*User, error)
	GetByEmail(email string) (*User, error)
	GetOne(id int) (*User, error)
	Insert(user User) (int, error)
	ResetPassword(password string) error
	PasswordMatches(plainText string) (bool, error)
	Update(user User) error
	DeleteByID(id int) error
}

type PlanInterface interface {
	GetAll() ([]*Plan, error)
	GetOne(id int) (*Plan, error)
	SubscribeUserToPlan(user User, plan Plan) error
	AmountForDisplay() string
}
