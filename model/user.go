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

// UserLogin
// it is used to login a user
func (u *User) Login() {
	// Login a user
	DB.Where("email = ?", u.Email).First(u)

}

// UserUpdate
// it is used to update a user
func (u *User) Update() {
	// Update a user
	DB.Save(u)
}


// UserDelete