package domain

import (
	"strings"
	"time"
)

// MaxPhotoSize is the maximum accepted photo upload size in bytes.
const MaxPhotoSize int64 = 5 * 1024 * 1024

// Photo is a vendor profile image.
type Photo struct {
	ID        string    `json:"id"`
	VendorID  string    `json:"vendorId"`
	MIME      string    `json:"mime"`
	Size      int64     `json:"size"`
	Data      []byte    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

// ValidPhotoMIME reports whether the MIME type denotes an image.
func ValidPhotoMIME(m string) bool { return strings.HasPrefix(m, "image/") }
