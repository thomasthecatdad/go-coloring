package reductions

import (
	g "github.com/thomaseb191/go-coloring/graphs"
	"log"
)

var AllAlgIds = []int{0} //TODO: ADD ADDITIONAL IDS

func RunReduction(gr g.Graph, id int, poolSize int, debug int) (g.Graph, string) {
	var outGraph g.Graph
	var algoName string

	switch id {
	case 0:
		outGraph = RunNaive(gr, poolSize, debug)
		algoName = "Naive"
	//TODO: ADD ADDITIONAL ALGORITHMS


	default:
		log.Fatal("No such algorithm found")
	}
	return outGraph, algoName
}