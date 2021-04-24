package testHarness

import (
	"bufio"
	g "github.com/thomaseb191/go-coloring/graphs"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
	Useful functions offered by this file:
		- ParseFile: parse a fileName to get a Graph
		- ParseTestFile: parses a fileName to get a list of TestDirectives
		- ConvertStringToIntArray converts an input string and parses it into an int array
 */

// TestDirective is a struct contains a string fileName for a graph and a list of algos and any other settings
//		GraphFile: the fileName to read in to build the graph
//		Algos: an integer array list of IDs for algorithms to run
//		PoolSize: the number of goroutine workers to allow (default -1)
//		Debug: the debug level for printing and displaying test results
type TestDirective struct {
	GraphFile string
	Algos []int
	PoolSize int
	Debug int
}

// Most of parsing reference taken from // Reference from https://gobyexample.com/reading-files

// parseCheck is a helper method to log a fatal error message if an error exists
func parseCheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}


// ParseFile takes a fileName and whether or not colors should be initialized to their index in the node array.
// Errors if file is incorrectly set up, given max degree is too small, or if directed edges are found
func ParseFile(fileName string, colorInit bool) g.Graph {
	//Initialize readers
	f, err := os.Open(fileName)
	parseCheck(err)

	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	//Parse metadata
	scanner.Scan()
	var n = scanner.Text()
	scanner.Scan()
	var d = scanner.Text()
	scanner.Scan()
	deg, err := strconv.Atoi(scanner.Text())

	//Initialize node tracking containers
	var nodeList []*g.Node

	var nodeNameMap map[string]*g.Node
	nodeNameMap = make(map[string]*g.Node)

	var nodeNeighborNameMap map[string][]string
	nodeNeighborNameMap = make(map[string][]string)

	counter := 0 //TODO: INDEX 0 OR 1

	//Read nodes line by line
	for scanner.Scan() {
		splitted1 := strings.Split(strings.ReplaceAll(scanner.Text(), " ", ""), ":")

		nodeName := splitted1[0]
		neighborNames := []string{splitted1[1]}
		if strings.Contains(splitted1[1], ",") {
			neighborNames = strings.Split(splitted1[1], ",")
		}

		_, ok := nodeNameMap[nodeName]
		if ok {
			log.Fatalf("Node %s duplicate definition", nodeName)
		}
		if len(neighborNames) > deg {
			log.Fatalf("Node %s has greater than %d degree", nodeName, deg)
		}

		newNode := g.Node {Name: nodeName, Ind: len(nodeList)}
		if (colorInit) {
			newNode.Color = counter
		}
		nodeNameMap[nodeName] = &newNode
		nodeNeighborNameMap[nodeName] = neighborNames
		nodeList = append(nodeList, &newNode)
		counter++
	}

	//Map string node names to their actual pointers
	refinedNodeList := g.NodeMatch(nodeList, nodeNameMap, nodeNeighborNameMap)

	return g.Graph{
		Name: n,
		Description: d,
		MaxDegree: deg,
		Nodes: refinedNodeList,
	}
}

// ParseTestFile takes a fileName and parses it to create an array of TestDirectives
func ParseTestFile(fileName string) []TestDirective {
	//Initialize readers
	f, err := os.Open(fileName)
	parseCheck(err)

	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var directiveList []TestDirective

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "%") {
			continue
		}
		splitted1 := strings.Split(scanner.Text(), " ")
		directiveList = append(directiveList, ParseArgsList(splitted1))
	}
	return directiveList
}

// ParseArgsList parses an array of Strings to create a TestDirective
func ParseArgsList(argList []string) TestDirective {
	poolSize := -1
	debugLevel := 3
	var err error = nil

	if len(argList) > 2 {
		poolSize, err = strconv.Atoi(argList[2])
		if err != nil {
			log.Fatal("Error parsing poolSize input")
		}
	}
	if len(argList) > 3 {
		debugLevel, err = strconv.Atoi(argList[3])
		if err != nil {
			log.Fatal("Error parsing debug input")
		}
	}
	return TestDirective{
		GraphFile: argList[0],
		Algos: ConvertStringToIntArray(argList[1]),
		PoolSize: poolSize,
		Debug: debugLevel,
	}
}

// ConvertStringToIntArray converts an input string and parses it into an int array
// Accepts [1,2,3]; [1,2,3,; 1,2,3]; and []
// Whitespace will throw an error with two many arguments
func ConvertStringToIntArray(str string) []int {
	removeBrackets1 := strings.ReplaceAll(str, "[", "")
	removeBrackets2 := strings.ReplaceAll(removeBrackets1, "]", "")
	if len(removeBrackets2) == 0 {
		return []int{}
	}

	splitted := strings.Split(removeBrackets2, ",")
	if len(splitted) == 0 {
		conv, err := strconv.Atoi(removeBrackets2)
		if err != nil {
			log.Fatal("Error parsing algo ID inputs")
		}
		return []int{conv}
	}

	var res []int
	for _, val := range splitted {
		conv, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal("Error parsing algo ID inputs")
		}
		res = append(res, conv)
	}
	return res
}