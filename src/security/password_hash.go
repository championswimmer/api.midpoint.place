package security

import (
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
)

const HashCostFactor = 10

func HashPassword(password string) string {
	if password == "" {
		applogger.Error("Hashing empty password")
	}
	hashedPassword := lo.Must(bcrypt.GenerateFromPassword([]byte(password), HashCostFactor))

	return string(hashedPassword)
}

func CheckPasswordHash(password, hash string) bool {
	if password == "" || hash == "" {
		applogger.Error("Comparing empty password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
