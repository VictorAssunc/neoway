package cnpj

import (
	"strconv"

	"neoway/pkg/lib/regex"
)

const (
	digitLimit           = 10
	validationDigitCount = 11
)

var (
	firstDigitWeightTable  = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	secondDigitWeightTable = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

// Validate checks if a CNPJ is valid.
func Validate(cnpj string) bool {
	if len(cnpj) != 14 || regex.OnlyDigits.MatchString(cnpj) {
		return false
	}

	var sum int
	for i, digit := range cnpj[:12] {
		d, _ := strconv.Atoi(string(digit))
		sum += d * firstDigitWeightTable[i]
	}

	rest := sum % validationDigitCount
	firstDigit := validationDigitCount - rest
	if firstDigit >= digitLimit {
		firstDigit = 0
	}

	if firstDigit != int(cnpj[12]-'0') {
		return false
	}

	sum = 0
	for i, digit := range cnpj[:13] {
		d, _ := strconv.Atoi(string(digit))
		sum += d * secondDigitWeightTable[i]
	}

	rest = sum % validationDigitCount
	secondDigit := validationDigitCount - rest
	if secondDigit >= digitLimit {
		secondDigit = 0
	}

	return secondDigit == int(cnpj[13]-'0')
}
