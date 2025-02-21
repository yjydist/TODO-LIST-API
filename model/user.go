package model

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// TableName
// it is used to define the table name for the User model
func (u *User) TableName() string {
	return "users"
}

// UserRegister
// it is used to register a new user
func (u *User) Save() error{
	// Register a new user
	// DB.Create(u)
	return nil
}

func FindUserByUsername(username string) (User, error) {
	user := User{}
	return user, nil
}