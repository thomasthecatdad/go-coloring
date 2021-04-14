package testHarness

import (
	//r "../reductions"
	g "../graphs"
)

type TestData struct {
	Name string
	DurationMillis float64
	Output g.Graph
}

//TODO: IMPLEMENT ACTUAL TEST HARNESS