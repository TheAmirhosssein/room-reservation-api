package validators

import (
	"regexp"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
)

func ValidateMobileNumber(mobileNumber string) bool {
	re := regexp.MustCompile(`^(?:98|\+98|0098|0)?9[0-9]{9}$`)
	return re.MatchString(mobileNumber)
}

func IsRoleValid(role string) bool {
	switch role {
	case entity.AdminRole, entity.SupportRole, entity.UserRole:
		return true
	default:
		return false
	}
}
