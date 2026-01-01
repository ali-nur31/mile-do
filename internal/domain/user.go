package domain

type UserInput struct {
	Email    string
	Password string
}

type UserOutput struct {
	AccessToken  string
	RefreshToken string
}
