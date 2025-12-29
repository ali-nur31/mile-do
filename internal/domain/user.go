package domain

type UserInput struct {
	Email    string
	Password string
}

type UserOutput struct {
	Token string
}
