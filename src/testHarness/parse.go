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



func ParseFile(fileName string, colorInit bool) g.Graph {
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

	var nodeList []*g.Node

	var nodeNameMap map[string]*g.Node
	nodeNameMap = make(map[string]*g.Node)

	var nodeNeighborNameMap map[string][]string
	nodeNeighborNameMap = make(map[string][]string)

	counter := 0 //TODO: INDEX 0 OR 1

	for scanner.Scan() {
		splitted1 := strings.Split(scanner.Text(), ":")

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

	refinedNodeList := g.NodeMatch(nodeList, nodeNameMap, nodeNeighborNameMap)

	return g.Graph{
		Name: n,
		Description: d,
		MaxDegree: deg,
		Nodes: refinedNodeList,
	}
}