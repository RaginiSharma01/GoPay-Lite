package razorpay

import (
	"errors"
	"os"

	"github.com/razorpay/razorpay-go"
)

var Client *razorpay.Client

func Init() error {
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")

	if keyID == "" || keySecret == "" {
		return errors.New("Razorpay keys not configured")
	}

	Client = razorpay.NewClient(keyID, keySecret)
	return nil
}
