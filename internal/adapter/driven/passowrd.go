package driven

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordService struct {
	cost int // Cost factor for bcrypt (default is 10)
}

// NewBcryptPasswordService creates a new instance of BcryptPasswordService with a specified cost factor.
func NewBcryptPasswordService(cost int) *BcryptPasswordService {
	return &BcryptPasswordService{cost: cost}
}

// HashPassword hashes the provided password using bcrypt.
func (b *BcryptPasswordService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (b *BcryptPasswordService) ComparePassword(hashedPassword, textPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(textPassword))
}
