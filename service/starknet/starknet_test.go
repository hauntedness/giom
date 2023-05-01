package starknet

import (
	"testing"

	"github.com/hauntedness/giom/internal/log"
)

func TestEstimateFees(t *testing.T) {
	fees, err := EstimateFees()
	if err != nil {
		t.Error(err)
		return
	}
	if len(fees) == 0 {
		t.Errorf("could not get any fee")
		return
	}
	log.Infos("fees", fees)
}
