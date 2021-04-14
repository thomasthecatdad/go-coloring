package testHarness

import (
	g "../graphs"
	"bufio"
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
func parseCheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func contains(strs []string, query string) bool {
	for _, v := range strs {
		if v == query {
			return true
		}
	}
	return false
}

func nodeMatch(nList []g.Node, nNameMap map[string]*g.Node, nNeighborNameMap map[string][]string) []g.Node {
	var nNeighborMap map[string][]*g.Node
	nNeighborMap = make(map[string][]*g.Node)

	//Construct proper nNeighborMap
	for k, v := range nNeighborNameMap {
		neighborPointers := make([]*g.Node, 0)
		for _, neighborName := range v {
			//Check for directed edges
			n2Neighbors, ok := nNeighborNameMap[neighborName]
			if !ok || !contains(n2Neighbors, k) {
				log.Fatalf( "DIRECTED EDGE, neighbor node %s lacks reverse pointer to %s", neighborName, k)
			}

			//Retrieve neighbor pointer
			neighbor, ok := nNameMap[neighborName]
			if !ok {
				log.Fatalf( "Parsing error, neighbor node pointer not found for %s with neighbor %s", k, neighborName)
			}
			neighborPointers = append(neighborPointers, neighbor)
		}

		nNeighborMap[k] = neighborPointers
	}

	//Retains original node order, rebuilds nodeList
	var newNodeList []g.Node
	for _, node := range nList {
		node.Neighbors = nNeighborMap[node.Name]
		newNodeList = append(newNodeList, node)
	}
	return newNodeList
}

func ParseFile(fileName string) g.Graph {
	f, err := os.Open(fileName)
	parseCheck(err)

	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	var n = scanner.Text()
	scanner.Scan()
	var d = scanner.Text()
	scanner.Scan()
	deg, err := strconv.Atoi(scanner.Text())

	nodeList := make([]g.Node, 0)

	var nodeNameMap map[string]*g.Node
	nodeNameMap = make(map[string]*g.Node)

	var nodeNeighborNameMap map[string][]string
	nodeNeighborNameMap = make(map[string][]string)

	for scanner.Scan() {
		splitted1 := strings.Split(scanner.Text(), ":")

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

		newNode := g.Node {Name: nodeName}
		nodeNameMap[nodeName] = &newNode
		nodeNeighborNameMap[nodeName] = neighborNames
		nodeList = append(nodeList, newNode)
	}

	refinedNodeList := nodeMatch(nodeList, nodeNameMap, nodeNeighborNameMap)

	return g.Graph{
		Name: n,
		Description: d,
		MaxDegree: deg,
		Nodes: refinedNodeList,
	}
}