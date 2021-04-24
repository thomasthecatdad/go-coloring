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

// Graph is a struct storing metadata and a list of nodes
//		Name: the name of the graph
//		Description: a description for the graph
//		MaxDegree: the maximum degree of the nodes. MaxDegree is guaranteed to be >= the max degree of all nodes via parsing
//		Nodes: an array of Node pointers that this graph owns
type Graph struct {
	Name string
	Description string
	MaxDegree int
	Nodes []*Node
}

// Node is a struct representing an individual vertex in a Graph
//		Name: the name of the node. Must be unique
//		Color: the color of the Node. May be initialized, may be changed during color-reduction
//		Neighbors: an array of Node pointers to a Node's neighbors
type Node struct {
	Name string
	Ind int
	Color int
	Neighbors []*Node
}

// IsSafe returns a bool for whether or not the Graph represents a valid coloring
// Based on https://www.geeksforgeeks.org/m-coloring-problem-backtracking-5/
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

// DeepCopy copies the Graph and all of its Nodes to hand over to another algorithm
func DeepCopy(gr *Graph) Graph {
	var nodeList []*Node

	var nodeNameMap map[string]*Node
	nodeNameMap = make(map[string]*Node)

	var nodeNeighborNameMap map[string][]string
	nodeNeighborNameMap = make(map[string][]string)

	for _, node := range gr.Nodes {
		newNode := Node{Name:node.Name, Ind: node.Ind, Color:node.Color}
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

// GetNamesFromNodeList converts a list of Node pointers to their names
// This function is mainly useful for printing
func GetNamesFromNodeList(neighborNodes []*Node) []string {
	var neighborNames []string
	for _, node := range neighborNodes {
		neighborNames = append(neighborNames, node.Name)
	}
	return neighborNames
}

// PrintGraph prints the metadata for a graph and then all of its nodes in an established format
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

// contains is a helper method for whether or not a string query is within an array
func contains(strs []string, query string) bool {
	for _, v := range strs {
		if v == query {
			return true
		}
	}
	return false
}

// NodeMatch reduces an array of Node pointers and mapping of their string-named neighbors to actual pointers
//		nList: an array of Node pointers that is edited directly to attach its neighbors
//		nNameMap: a map of string Node names to their pointers
//		nNeighborNameMap: a map of a Node's string name to its string-named neighbors
// This is used during copying and parsing
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
	for _, node := range nList {
		node.Neighbors = nNeighborMap[node.Name]
	}
	return nList
}

// RunColorInit sets all of the Node's colors in a Graph to their index in the Graph's Nodes
func RunColorInit(gr *Graph) *Graph {
	for i, k := range gr.Nodes {
		k.Color = i  //TODO: INDEX 0 OR 1
	}
	return gr
}

// intContains is a helper method for whether or not an int query is within an array
func intContains(nums []int, query int) bool {
	for _, v := range nums {
		if v == query {
			return true
		}
	}
	return false
}

// CountColors counts the total number of unique colors within a Graph
func CountColors(gr *Graph) int {
	var allColors []int
	for _, node := range gr.Nodes {
		if !intContains(allColors, node.Color) {
			allColors = append(allColors, node.Color)
		}
	}
	return len(allColors)
}