package reductions

import (
	g "github.com/thomaseb191/go-coloring/graphs"
	"log"
)

// AllAlgIds - A list of all valid algorithm IDs for when t.RunTest is given an empty array.
var AllAlgIds = []int{0, 1} //TODO: ADD ADDITIONAL IDS

// RunReduction calls the respective color-reducing algorithm for a graph, algorithm id, number of worker pools, and debug setting
// 		gr: a graph that the algorithm will own
// 		id: an ID mapping to an algorithm
// 		poolSize: the number of worker goroutines allowed for parallel algorithms
// 		debug: 0 if just generate output, 1 if allow prints, 2 if just graph and output, 3 if allow graph and prints
func RunReduction(gr g.Graph, id int, poolSize int, debug int) (g.Graph, string) {
	var outGraph g.Graph
	var algoName string

	switch id {
	case 0:
		outGraph = RunNaive(gr, poolSize, debug)
		algoName = "Naive"
	case 1:
		outGraph = kwReduction(gr, poolSize, debug)
		algoName = "Kuhn-Wattenhofer"
	//TODO: ADD ADDITIONAL ALGORITHMS


	default:
		log.Fatalf("No such algorithm found for %d.\n", id)
	}
	return outGraph, algoName
}