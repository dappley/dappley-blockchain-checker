package helper

import(
	"errors"
)

//Checks the flag arguments. If the argument is a default value, then return an error message.
func CheckFlags(email string, password string, test string, main string) (err error) {
	switch {
	case email == "default_email":
		err = errors.New("Error: Email is missing!")
	case password == "default_password":
		err = errors.New("Error: Password is missing!")
	case test == "default_test.txt":
		err = errors.New("Error: Test server test result is missing!")
	case main == "default_main.txt":
		err = errors.New("Error: Main server test result is missing!")
	default:
		err = nil
	}
	return err
}