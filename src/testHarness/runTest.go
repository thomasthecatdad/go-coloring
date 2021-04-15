package testHarness

import (
	r "github.com/thomaseb191/go-coloring/reductions"
	"fmt"
	//d "../display" //TODO: IMPORT
	g "github.com/thomaseb191/go-coloring/graphs"
	"time"
)

type TestData struct {
	Name string
	DurationMillis time.Duration
	Output g.Graph
	NumColors int
}

func RunTest(fileName string, algos []int, poolSize int, debug int) []TestData {
	//Parse and build the graph. Check the colors
	var testDatas []TestData
	initGraph := ParseFile(fileName, false)
	if debug % 2 == 1 {
		fmt.Printf("Initial IsSafe() for %s without color init: %t\n", initGraph.Name, g.IsSafe(&initGraph))
	}
	g.RunColorInit(&initGraph)

	//Start the time, run an algorithm
	if len(algos) == 0 {
		algos = r.AllAlgIds
	}

	for _, algo := range algos {
		copiedGraph := g.DeepCopy(&initGraph)
		start := time.Now()
		outGraph, algoName := r.RunReduction(copiedGraph, algo, poolSize, debug)

		//Stop the time, check the algorithm
		elapsed := time.Since(start)
		numColors := g.CountColors(&outGraph)

		if debug % 2 == 1 {
			fmt.Printf("Output IsSafe() for %s_%s in %d: %t\n", initGraph.Name, algoName, elapsed.Milliseconds(), g.IsSafe(&initGraph))
			fmt.Printf("\t\tNum Colors: %d\n", numColors)
		}

		newTest := TestData{
			Name: initGraph.Name + "_" + algoName,
			DurationMillis: elapsed,
			Output: outGraph,
			NumColors: numColors,
		}
		testDatas = append(testDatas, newTest)
		//Render if desired
		if debug >= 2 {
			//TODO: RENDER VISUALIZATION
		}
	}
	return testDatas
}