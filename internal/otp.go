package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateOTP(entry Entry, t time.Time) (string, int64, error) {
	algorithm, err := parseAlgorithm(entry.Algorithm)
	if err != nil {
		return "", 0, fmt.Errorf("parseAlgorithm: %w", err)
	}

	digits, err := parseDigits(entry.Digits)
	if err != nil {
		return "", 0, fmt.Errorf("parseDigits: %w", err)
	}

	switch strings.ToUpper(entry.Type) {
	case "TOTP":
		code, err := totp.GenerateCodeCustom(entry.Secret, t, totp.ValidateOpts{
			Algorithm: algorithm,
			Period:    uint(entry.Period),
			Digits:    digits,
		})
		if err != nil {
			return "", 0, fmt.Errorf("something went wrong while generating totp: %w", err)
		}
		period := int64(entry.Period)
		return code, period - (t.Unix() % period), nil
	default:
		return "", 0, fmt.Errorf("unsupported entry type: %s", entry.Type)
	}
}

func parseAlgorithm(algorithm string) (otp.Algorithm, error) {
	switch strings.ToUpper(algorithm) {
	case "SHA1":
		return otp.AlgorithmSHA1, nil
	case "SHA256":
		return otp.AlgorithmSHA256, nil
	case "SHA512":
		return otp.AlgorithmSHA512, nil
	case "MD5":
		return otp.AlgorithmMD5, nil
	default:
		return 0, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func parseDigits(digits uint8) (otp.Digits, error) {
	switch digits {
	case 6:
		return otp.DigitsSix, nil
	case 8:
		return otp.DigitsEight, nil
	default:
		return 0, fmt.Errorf("unsupported digits: %d", digits)
	}
}
