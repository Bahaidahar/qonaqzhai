// Package gateway holds PSP adapters. The Mock gateway always reports a
// captured charge — sufficient for diploma demos and for tests.
package gateway

import (
	"context"

	"qonaqzhai-backend/pkg/config"
	"qonaqzhai-backend/services/payment/internal/ports"
)

// Mock always succeeds and returns a synthetic provider reference.
type Mock struct{}

// NewMock constructs a Mock gateway.
func NewMock() Mock { return Mock{} }

// Charge fakes a successful capture.
func (Mock) Charge(_ context.Context, in ports.ChargeInput) (string, error) {
	ref, _ := config.RandomHex(8)
	return "mock-" + ref, nil
}

var _ ports.Gateway = Mock{}
