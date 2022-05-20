package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	typeTOTP        = "TOTP"
	algorithmSHA1   = "SHA1"
	algorithmSHA256 = "SHA256"
	algorithmSHA512 = "SHA512"
	algorithmMD5    = "MD5"
	digitsSix       = 6
	digitsEight     = 8
)

func generateOTP(entry Entry, t time.Time) (string, int64, error) {
	algorithm, err := parseAlgorithm(entry.Algorithm)
	if err != nil {
		return "", 0, fmt.Errorf("parseAlgorithm: %w", err)
	}

	digits, err := parseDigits(entry.Digits)
	if err != nil {
		return "", 0, fmt.Errorf("parseDigits: %w", err)
	}

	switch strings.ToUpper(entry.Type) {
	case typeTOTP:
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
	case algorithmSHA1:
		return otp.AlgorithmSHA1, nil
	case algorithmSHA256:
		return otp.AlgorithmSHA256, nil
	case algorithmSHA512:
		return otp.AlgorithmSHA512, nil
	case algorithmMD5:
		return otp.AlgorithmMD5, nil
	default:
		return 0, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func parseDigits(digits uint8) (otp.Digits, error) {
	switch digits {
	case digitsSix:
		return otp.DigitsSix, nil
	case digitsEight:
		return otp.DigitsEight, nil
	default:
		return 0, fmt.Errorf("unsupported digits: %d", digits)
	}
}
