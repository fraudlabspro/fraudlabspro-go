package fraudlabspro

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	PaymentMethodCreditCard     string = "CREDITCARD"
	PaymentMethodPayPal         string = "PAYPAL"
	PaymentMethodCashOnDelivery string = "COD"
	PaymentMethodBankDeposit    string = "BANKDEPOSIT"
	PaymentMethodGiftCard       string = "GIFTCARD"
	PaymentMethodCrypto         string = "CRYPTO"
	PaymentMethodWireTransfer   string = "WIRED"
	PaymentMethodOthers         string = "OTHERS"
)

const (
	OrderActionApprove         string = "APPROVE"
	OrderActionReject          string = "REJECT"
	OrderActionRejectBlacklist string = "REJECT_BLACKLIST"
)

type ErrorObject struct {
	Error struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	} `json:"error"`
}

func doHash(s string) string {
	// Step 1: initial prefix
	hash := "fraudlabspro_" + s

	// Step 2: 65536 iterations of sha1
	for i := 0; i < 65536; i++ {
		h := sha1.Sum([]byte("fraudlabspro_" + hash))
		hash = hex.EncodeToString(h[:])
	}

	// Step 3: final sha256
	h2 := sha256.Sum256([]byte(hash))
	return hex.EncodeToString(h2[:])
}

func stripNonDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func formatTo2Decimals(numStr string) (string, error) {
	// Step 1: parse string to float64
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return "", err
	}

	// Step 2: format to 2 decimal places
	return fmt.Sprintf("%.2f", num), nil
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
