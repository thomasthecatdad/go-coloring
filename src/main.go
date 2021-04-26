package main

import (
	"encoding/json"
	"fmt"
	g "github.com/thomaseb191/go-coloring/graphs"
	r "github.com/thomaseb191/go-coloring/reductions"
	t "github.com/thomaseb191/go-coloring/testHarness"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

		runTestAndPrintResult(td, td.Debug)
	} else {
		log.Fatal("Incorrect arguments specified")
	}
}

// runTestAndPrintResults is a helper method to run a specific test set
func runTestAndPrintResult(td t.TestDirective, debug int) {
	tResults := t.RunTest(td.GraphFile, td.Algos, td.PoolSize, td.Debug)
	for _, k := range tResults {
		fmt.Printf("Test Name: %s\n", k.Name)
		if debug % 2 == 1 {
			g.PrintGraph(&k.Output)
		}
		fmt.Printf("IsSafe: %t\tNum Colors: %d\n", g.IsSafe(&k.Output), k.NumColors)
	}
	fmt.Printf("\n-------------------------\n")
}

// runTestAndPrintResultAndTrends is a helper method to print results of tests and generate the trend lines and output results to json
func runTestAndPrintResultAndTrends(tds []t.TestDirective, testFileName string) {
	var tTestNames [r.NumAlgos][]string
	var tNumNodes [r.NumAlgos][]int
	var tTimeElapsed [r.NumAlgos][]int
	var tNumberColors [r.NumAlgos][]int
	var tMaxDegree [r.NumAlgos][]int
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
			tTestNames[currAlg] = append(tTestNames[currAlg], test.Name)
			tNumNodes[currAlg] = append(tNumNodes[currAlg], len(test.Output.Nodes))
			tTimeElapsed[currAlg]= append(tTimeElapsed[currAlg], int(test.DurationMillis.Milliseconds()))
			tNumberColors[currAlg] = append(tNumberColors[currAlg], test.NumColors)
			tMaxDegree[currAlg] = append(tMaxDegree[currAlg], test.Output.MaxDegree)
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
			Names:		  tTestNames[id],
			NumNodes:     tNumNodes[id],
			TimeElapsed:  tTimeElapsed[id],
			NumberColors: tNumberColors[id],
			MaxDegree:	  tMaxDegree[id],
			IsSafe:       tIsSafe[id],
		}
	}

	fmt.Printf("\n-------------------------\n")
	testOutName := extractTestName(testFileName)
	g.GenerateHTMLForDataPoints(tResults, testOutName) //TODO: CHANGE GRAPH NAME
	writeJson(tResults, testOutName)
}

// writeJson is a helper method to write an output to a json output file
func writeJson(tResults map[int]g.DataPoint, testFileName string) {
	b, err := json.Marshal(tResults)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("../json/"+testFileName, b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func extractTestName(testFileName string) string {
	outName := testFileName[0:(len(testFileName)-4)] + ".json"
	if strings.Contains(testFileName, "/") {
		tempNameArray := strings.Split(testFileName, "/")
		tempName := tempNameArray[len(tempNameArray)-1]
		outName = tempName[0:(len(tempName)-4)] + ".json"
	} else if strings.Contains(testFileName, "\\") {
		tempNameArray := strings.Split(testFileName, "\\")
		tempName := tempNameArray[len(tempNameArray)-1]
		outName = tempName[0:(len(tempName)-4)] + ".json"
	}
	return outName
}

