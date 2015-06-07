package user

type User struct {
	Id           string
	Name         string
	passwordHash string
	authToken    string
}
