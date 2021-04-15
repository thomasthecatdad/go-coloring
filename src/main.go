package main

import (
	"fmt"
	g "github.com/thomaseb191/go-coloring/graphs"
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
		fmt.Printf("Duration of test %s: %d with %d colors\n", tResults[0].Name, tResults[0].DurationMillis.Milliseconds(), tResults[0].NumColors)
		g.PrintGraph(&tResults[0].Output)
	} else if len(inputArgs) == 2 {
		// Read in file with list of tests to run
		testFileName := os.Args[1]
		testDirectives := t.ParseTestFile(testFileName)

		for _, td := range testDirectives {
			runTestAndPrintResult(td)
		}
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