package main

import (
	"fmt"
	g "github.com/thomaseb191/go-coloring/graphs"
	r "github.com/thomaseb191/go-coloring/reductions"
	t "github.com/thomaseb191/go-coloring/testHarness"
	"log"
	"os"
)

// If running from goland, paths should be res/...
// If running from terminal within src/, paths should be ../res/...
// Examples of calls after running 'go build main.go' include
//		- ./main.exe ../testFiles/test01_naive.txt
//		- ./main.exe ../res/Sample01.txt []
//		- ./main.exe ../res/Sample01.txt [] -1
//		- ./main.exe ../res/Sample01.txt [] -1 3
func main() {
	inputArgs := os.Args
	if len(inputArgs) == 1 {
		// Default behavior
		// TODO: CHANGE TO DESIRED DEFAULT BEHAVIOR
		fmt.Printf("\n\n\n")

		tResults := t.RunTest("../res/Sample02.txt", []int{}, -1, 3)
		for i := 0; i < len(tResults); i++ {
			fmt.Printf("Duration of test %s: %d with %d colors\n", tResults[i].Name, tResults[i].DurationMillis.Milliseconds(), tResults[i].NumColors)
			g.PrintGraph(&tResults[i].Output)
		}

	} else if len(inputArgs) == 2 {
		// Read in file with list of tests to run
		testFileName := os.Args[1]
		testDirectives := t.ParseTestFile(testFileName)

		runTestAndPrintResultAndTrends(testDirectives, testFileName)
	} else if len(inputArgs) >= 3 && len(inputArgs) <= 5 {
		// Run a singular test
		td := t.ParseArgsList(os.Args[1:])

		runTestAndPrintResult(td)
	} else {
		log.Fatal("Incorrect arguments specified")
	}
}

// runTestAndPrintResults is a helper method to run a specific test set
func runTestAndPrintResult(td t.TestDirective) {
	tResults := t.RunTest(td.GraphFile, td.Algos, td.PoolSize, td.Debug)
	for _, k := range tResults {
		//TODO: REFINE TEST OUTPUT, MAYBE EXPORT RESULTS TO FILE
		fmt.Printf("Test Name: %s\n", k.Name)
		g.PrintGraph(&k.Output)
	}
	fmt.Printf("\n-------------------------\n")
}

// runTestAndPrintResultAndTrends is a helper method to print results of tests and generate the trend lines
func runTestAndPrintResultAndTrends(tds []t.TestDirective, testFileName string) {
	var tNumNodes [r.NumAlgos][]int
	var tTimeElapsed [r.NumAlgos][]int
	var tNumberColors [r.NumAlgos][]int
	var tIsSafe [r.NumAlgos][]bool

	for _, td := range tds {
		//Run Tests
		testResults := t.RunTest(td.GraphFile, td.Algos, td.PoolSize, td.Debug)

		algos := td.Algos
		if len(algos) == 0 {
			algos = r.AllAlgIds
		}
		//Extract and format data into arrays
		for i, test := range testResults {
			currAlg := algos[i]
			tNumNodes[currAlg] = append(tNumNodes[currAlg], len(test.Output.Nodes))
			tTimeElapsed[currAlg]= append(tTimeElapsed[currAlg], int(test.DurationMillis.Milliseconds()))
			tNumberColors[currAlg] = append(tNumberColors[currAlg], test.NumColors)
			tIsSafe[currAlg] = append(tIsSafe[currAlg], test.IsSafe)

			fmt.Printf("Test Name: %s\n", test.Name)
			fmt.Printf("\tDurationMillis: %d\tNumColors: %d\tIsSafe: %t\n", test.DurationMillis, test.NumColors, test.IsSafe)
		}
	}
	//Format data into DataPoints
	var tResults map[int]g.DataPoint
	tResults = make(map[int]g.DataPoint)

	for _, id := range r.AllAlgIds {
		tResults[id] = g.DataPoint{
			NumNodes:     tNumNodes[id],
			TimeElapsed:  tTimeElapsed[id],
			NumberColors: tNumberColors[id],
			IsSafe:       tIsSafe[id],
		}
	}
	//TODO: SEND TO JSON

	fmt.Printf("\n-------------------------\n")
	g.GenerateHTMLForDataPoints(tResults) //TODO: CHANGE GRAPH NAME
}