package domain

import (
	"net/http"
	"strings"
	"time"
)

// MaxPhotoSize is the maximum accepted photo upload size in bytes.
const MaxPhotoSize int64 = 5 * 1024 * 1024

// AllowedPhotoMIMEs is the whitelist of accepted image content types.
// SVG is excluded on purpose — it can carry XSS payloads.
var AllowedPhotoMIMEs = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
	"image/gif":  {},
}

// Photo is a vendor profile image.
type Photo struct {
	ID        string    `json:"id"`
	VendorID  string    `json:"vendorId"`
	MIME      string    `json:"mime"`
	Size      int64     `json:"size"`
	Data      []byte    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

// DetectPhotoMIME sniffs the MIME from the file bytes and returns it only when
// it is one of AllowedPhotoMIMEs. The Content-Type header is untrusted — the
// client controls it. Returns empty string when the bytes do not match any
// allowed type.
func DetectPhotoMIME(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	sniffed := http.DetectContentType(data)
	// DetectContentType returns e.g. "image/jpeg; charset=utf-8" sometimes —
	// strip parameters.
	if idx := strings.Index(sniffed, ";"); idx > 0 {
		sniffed = strings.TrimSpace(sniffed[:idx])
	}
	if _, ok := AllowedPhotoMIMEs[sniffed]; !ok {
		return ""
	}
	return sniffed
}
