package tests

// DOESN'T WORK

import (
	"context"
    "testing"
    pb "gw-exchanger/internal/proto/exchange"
	"gw-exchanger/cmd"
    "github.com/stretchr/testify/assert"
)

type MockServer struct {
	server
}

func TestGetExchangeRates (t *testing.T) {
	srv := &MockServer{}
	response, err := srv.GetExchangeRates(context.Background(), &pb.Empty{})

	// TESTCODE!!!
	expectedRates := map[string]float32{
        "USD": 1.0,
        "EUR": 0.9,
        "RUB": 90.0,
    }

	assert.Equal(t, expectedRates, response.Rates)
}