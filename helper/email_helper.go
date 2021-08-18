package helper

import (
	"net/mail"
)

//Returns true when one of the server test result returns failcase.
func ContainsFailure(test_results []bool) (final_result bool) {
	for _, test_result := range test_results {
		final_result = final_result || test_result
	}
	return
}

//Checks whether the input email address is a valid address or not.
func Valid_email(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}