package driven

type PasswordService interface {
	HashPassword(string) (string, error)
	ComparePassword(hashedPassword, textPassword string) error
}
