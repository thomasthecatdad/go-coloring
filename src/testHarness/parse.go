package testHarness

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
	Useful functions offered by this file:
		- ParseFile: parse a fileName to get a Graph
		- GetNamesFromNodeList: converts a list of Node pointers to a list of string Node names
		- PrintGraph: prints a graph
 */

// Most of parsing reference taken from // Reference from https://gobyexample.com/reading-files
func parseCheck(e error) {
	if (e != nil) {
		log.Fatal(e)
	}
}

func nodeMatch(nList []Node, nNameMap map[string]*Node, nNeighborNameMap map[string][]string) []Node {
	var nNeighborMap map[string][]*Node
	nNeighborMap = make(map[string][]*Node)

	//construct proper nNeighborMap
	for k, v := range nNeighborNameMap {
		neighborPointers := make([]*Node, 0)
		for _, neighborName := range v {
			neighbor, ok := nNameMap[neighborName]
			if (!ok) {
				log.Fatalf( "Parsing error, neighbor node pointer not found for %s with neighbor %s", k, neighborName)
			}
			neighborPointers = append(neighborPointers, neighbor)
		}

		nNeighborMap[k] = neighborPointers
	}

	var newNodeList []Node
	for _, node := range nList {
		node.neighbors = nNeighborMap[node.name]
		newNodeList = append(newNodeList, node)
	}
	return newNodeList
}

func ParseFile(fileName string) Graph {
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

	nodeList := make([]Node, 0)

	var nodeNameMap map[string]*Node
	nodeNameMap = make(map[string]*Node)

	var nodeNeighborNameMap map[string][]string
	nodeNeighborNameMap = make(map[string][]string)

	for scanner.Scan() {
		splitted1 := strings.Split(scanner.Text(), ":")

		nodeName := splitted1[0]
		neighborNames := []string{splitted1[1]}
		if (strings.Contains(splitted1[1], ",")) {
			neighborNames = strings.Split(splitted1[1], ",")
		}

		_, ok := nodeNameMap[nodeName]
		if (ok) {
			log.Fatal("Node %s duplicate definition", nodeName)
		}
		if (len(neighborNames) > deg) {
			log.Fatal("Node %s has greater than %d degree", nodeName, deg)
		}

		newNode := Node {name: nodeName}
		nodeNameMap[nodeName] = &newNode
		nodeNeighborNameMap[nodeName] = neighborNames
		nodeList = append(nodeList, newNode)
	}

	refinedNodeList := nodeMatch(nodeList, nodeNameMap, nodeNeighborNameMap)

	return Graph{
		name: n,
		description: d,
		maxDegree: deg,
		nodes: refinedNodeList,
	}
}

func GetNamesFromNodeList(neighborNodes []*Node) []string {
	var neighborNames []string
	for _, node := range neighborNodes {
		neighborNames = append(neighborNames, node.name)
	}
	return neighborNames
}

func PrintGraph(g *Graph) {
	fmt.Printf("Graph Name: \t\t%s\n", g.name)
	fmt.Printf("Graph Description: \t%s\n", g.description)
	fmt.Printf("Graph Max Degree: \t%d\n", g.maxDegree)
	for _, v := range g.nodes {
		neighborNames := GetNamesFromNodeList(v.neighbors)
		fmt.Printf("Node: %s,\tColor: %d,\tNeighbors: %s\n", v.name, v.color, neighborNames)
	}
	fmt.Printf("---End of Graph [%s]", g.name)
}