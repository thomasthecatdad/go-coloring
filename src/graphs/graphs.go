package graphs

import (
	"fmt"
	"log"
)

/*
	Structs for Building and Evaluating Graphs
 */

/*
	Useful functions offered by this file:
		- IsSafe: checks for color differences between neighbors, based on https://www.geeksforgeeks.org/m-coloring-problem-backtracking-5/
		- GetNamesFromNodeList: converts a list of Node pointers to a list of string Node names
		- PrintGraph: prints a graph
		- NodeMatch: convert a map of names into map of pointers
*/

type Graph struct {
	Name string
	Description string
	MaxDegree int
	Nodes []*Node
}

type Node struct {
	Name string
	Color int
	Neighbors []*Node
}

func IsSafe(gr *Graph) bool {
	for _, node := range gr.Nodes {
		for _, neighbor := range node.Neighbors {
			if node.Color == neighbor.Color {
				return false
			}
		}
	}
	return true
}

func DeepCopy(gr *Graph) Graph {
	var nodeList []*Node

	var nodeNameMap map[string]*Node
	nodeNameMap = make(map[string]*Node)

	var nodeNeighborNameMap map[string][]string
	nodeNeighborNameMap = make(map[string][]string)

	for _, node := range gr.Nodes {
		newNode := Node{Name:node.Name, Color:node.Color}
		nodeList = append(nodeList, &newNode)
		nodeNameMap[node.Name] = &newNode
		nodeNeighborNameMap[node.Name] = GetNamesFromNodeList(node.Neighbors)
	}

	copiedNodes := NodeMatch(nodeList, nodeNameMap, nodeNeighborNameMap)

	return Graph{
		Name: gr.Name,
		Description: gr.Description,
		MaxDegree: gr.MaxDegree,
		Nodes: copiedNodes,
	}
}

func GetNamesFromNodeList(neighborNodes []*Node) []string {
	var neighborNames []string
	for _, node := range neighborNodes {
		neighborNames = append(neighborNames, node.Name)
	}
	return neighborNames
}

func PrintGraph(gr *Graph) {
	fmt.Printf("Graph Name: \t\t%s\n", gr.Name)
	fmt.Printf("Graph Description: \t%s\n", gr.Description)
	fmt.Printf("Graph Max Degree: \t%d\n", gr.MaxDegree)
	for _, v := range gr.Nodes {
		neighborNames := GetNamesFromNodeList(v.Neighbors)
		fmt.Printf("Node: %s,\tColor: %d,\tNeighbors: %s\n", v.Name, v.Color, neighborNames)
	}
	fmt.Printf("---End of Graph [%s].\n", gr.Name)
}

func contains(strs []string, query string) bool {
	for _, v := range strs {
		if v == query {
			return true
		}
	}
	return false
}

func NodeMatch(nList []*Node, nNameMap map[string]*Node, nNeighborNameMap map[string][]string) []*Node {
	var nNeighborMap map[string][]*Node
	nNeighborMap = make(map[string][]*Node)

	//Construct proper nNeighborMap
	for k, v := range nNeighborNameMap {
		neighborPointers := make([]*Node, 0)
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
	//var newNodeList []*Node
	for _, node := range nList {
		node.Neighbors = nNeighborMap[node.Name]
		//newNodeList = append(newNodeList, &node)
	}
	//return newNodeList
	return nList
}

func RunColorInit(gr *Graph) *Graph {
	for i, k := range gr.Nodes {
		k.Color = i  //TODO: INDEX 0 OR 1
	}
	return gr
}

func intContains(nums []int, query int) bool {
	for _, v := range nums {
		if v == query {
			return true
		}
	}
	return false
}

func CountColors(gr *Graph) int {
	var allColors []int
	for _, node := range gr.Nodes {
		if !intContains(allColors, node.Color) {
			allColors = append(allColors, node.Color)
		}
	}
	return len(allColors)
}