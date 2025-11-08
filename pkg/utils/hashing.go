package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes), err

}

func ComparePasswords(hashedPassword, password string) error {
	//returns an error if the comparison fails
	//!= nil means “the error is not nil” → i.e., the password is wrong.
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
