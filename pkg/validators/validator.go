package validators

import "regexp"

func ValidateMobileNumber(mobileNumber string) bool {
	re := regexp.MustCompile(`^(?:98|\+98|0098|0)?9[0-9]{9}$`)
	return re.MatchString(mobileNumber)
}
