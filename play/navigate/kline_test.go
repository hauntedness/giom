package main

import (
	"testing"

	"github.com/hauntedness/giom/internal/log"
	"gonum.org/v1/gonum/floats"
)

func TestNormalize(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	f := floats.Norm(data, 1)
	log.Infos("norm", f)
}
