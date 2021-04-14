package graphs

import "fmt"

/*
	Structs for Building and Evaluating Graphs
 */

/*
	Useful functions offered by this file:
		- IsSafe: checks for color differences between neighbors, based on https://www.geeksforgeeks.org/m-coloring-problem-backtracking-5/
		- GetNamesFromNodeList: converts a list of Node pointers to a list of string Node names
		- PrintGraph: prints a graph
*/

type Graph struct {
	Name string
	Description string
	MaxDegree int
	Nodes []Node
}

type Node struct {
	Name string
	Color int
	Neighbors []*Node
}

func IsSafe(gr *Graph) bool {
	// TODO
	return false
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
	fmt.Printf("---End of Graph [%s]", gr.Name)
}