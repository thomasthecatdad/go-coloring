package testHarness

/*
	Structs for Building and Evaluating Graphs
 */

type TestData struct {
	name string
	durationMillis float64
	output Graph
}

type Graph struct {
	name string
	description string
	maxDegree int
	nodes []Node
}

type Node struct {
	name string
	color int
	neighbors []*Node
}