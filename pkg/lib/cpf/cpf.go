package cpf

import (
	"strconv"

	"neoway/pkg/lib/regex"
)

const (
	firstDigitWeightStart  = 10
	secondDigitWeightStart = 11

	digitLimit           = 10
	validationDigitCount = 11
)

var (
	invalidCPFs = map[string]struct{}{
		"00000000000": {},
		"11111111111": {},
		"22222222222": {},
		"33333333333": {},
		"44444444444": {},
		"55555555555": {},
		"66666666666": {},
		"77777777777": {},
		"88888888888": {},
		"99999999999": {},
	}
)

// Validate checks if a CPF is valid.
func Validate(cpf string) bool {
	if _, ok := invalidCPFs[cpf]; ok || len(cpf) != 11 || regex.OnlyDigits.MatchString(cpf) {
		return false
	}

	var sum int
	for i, digit := range cpf[:9] {
		d, _ := strconv.Atoi(string(digit))
		sum += d * (firstDigitWeightStart - i)
	}

	rest := sum % validationDigitCount
	firstDigit := validationDigitCount - rest
	if firstDigit >= digitLimit {
		firstDigit = 0
	}

	if firstDigit != int(cpf[9]-'0') {
		return false
	}

	sum = 0
	for i, digit := range cpf[:10] {
		d, _ := strconv.Atoi(string(digit))
		sum += d * (secondDigitWeightStart - i)
	}

	rest = sum % validationDigitCount
	secondDigit := validationDigitCount - rest
	if secondDigit >= digitLimit {
		secondDigit = 0
	}

	return secondDigit == int(cpf[10]-'0')
}
