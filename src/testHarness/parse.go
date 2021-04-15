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
 */

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
			//TODO: TRIM NEIGHBORNAMES FOR WHITESPACE
		}

		_, ok := nodeNameMap[nodeName]
		if ok {
			log.Fatalf("Node %s duplicate definition", nodeName)
		}
		if len(neighborNames) > deg {
			log.Fatalf("Node %s has greater than %d degree", nodeName, deg)
		}

		newNode := g.Node {Name: nodeName}
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