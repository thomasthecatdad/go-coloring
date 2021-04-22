package reductions

import (
	"fmt"
	g "github.com/thomaseb191/go-coloring/graphs"
)

func runNaiveGoRoutine(gr g.Graph, poolSize int, debug int, c chan g.Graph) {
	c <- RunNaive(gr, poolSize, debug)
}

func convertBinsToGraph(bins [][]*g.Node, original g.Graph) g.Graph {
	nodes := make([]*g.Node, 0)
	for color := 0; color < len(bins); color++ {
		for _, node := range bins[color] {
			nodes = append(nodes,
				&g.Node{
					Name: node.Name,
					Color: color,
					Neighbors: node.Neighbors,
				})
		}
	}

	return g.Graph{
		Name: original.Name,
		Description: "Color Reduced with KW",
		MaxDegree: original.MaxDegree,
		Nodes: nodes,
	}
}

func combineColors(bins [][]*g.Node, gr g.Graph, c chan [][]*g.Node) {
	binsGraph := convertBinsToGraph(bins, gr)
	newGr := RunNaive(binsGraph, -1, 3)
	colorToIndex := make(map[int]int)
	latestIndex := 0
	numColors := g.CountColors(&newGr)
	newBins := make([][]*g.Node, numColors)
	for _, node := range newGr.Nodes {
		if _, ok := colorToIndex[node.Color]; ! ok {
			colorToIndex[node.Color] = latestIndex
			latestIndex++
		}
		newBins[colorToIndex[node.Color]] = append(newBins[colorToIndex[node.Color]], node)
	}
	c <- newBins
}

func kwReduction(gr g.Graph, poolSize int, debug int) g.Graph {
	fmt.Printf("Starting KW Reduction \n")
	degree := gr.MaxDegree
	startIndexes := make([]int, 0)
	size := len(gr.Nodes)
	c := make(chan g.Graph)
	if size < 2 * (degree + 1) {
		gr.Description = "Color Reduced with KW"
		go runNaiveGoRoutine(gr, poolSize, debug, c)
		return <- c
	}
	for x := 0; x < size; x++ {
		if x % (2 * (degree + 1)) == 0 {
			startIndexes = append(startIndexes, x)
		}
	}

	for i := 0; i < len(startIndexes); i++ {
		currStart := startIndexes[i]
		var nextStart int
		if i + 1 != len(startIndexes) {
			nextStart = startIndexes[i + 1]
		} else {
			nextStart = size
		}
		grCopy := g.DeepCopy(&gr)
		grCopy.Nodes = gr.Nodes[currStart:nextStart]
		go runNaiveGoRoutine(grCopy, poolSize, debug, c)
	}
	colorBins := make([][]*g.Node, 0)

	for i := 0; i < len(startIndexes); i++{
		graph := <- c
		numColors := g.CountColors(&graph)
		newBins := make([][]*g.Node, numColors)
		colorToIndex := make(map[int]int)
		latestIndex := 0
		for _, node := range graph.Nodes {
			if _, ok := colorToIndex[node.Color]; ! ok {
				colorToIndex[node.Color] = latestIndex
				latestIndex++
			}
			newBins[colorToIndex[node.Color]] = append(newBins[colorToIndex[node.Color]], node)
		}
		colorBins = append(colorBins, newBins...)
	}
	close(c)

	for len(colorBins) > degree + 1 {
		d := make(chan [][]*g.Node)
		binIndexes := make([]int, 0)
		numColors := len(colorBins)

		for x := 0; x < numColors; x++ {
			if x % (2 * (degree + 1)) == 0 {
				binIndexes = append(binIndexes, x)
			}
		}
		newBins := make([][]*g.Node, 0)

		for i := 0; i < len(binIndexes); i++ {
			currStart := binIndexes[i]
			var nextStart int
			if i + 1 != len(binIndexes) {
				nextStart = binIndexes[i + 1]
			} else {
				nextStart = len(colorBins)
			}
			go combineColors(colorBins[currStart:nextStart], gr, d)
		}
		for i := 0; i < len(binIndexes); i++ {
			bins := <- d
			newBins = append(newBins, bins...)
		}
		close(d)
		colorBins = newBins
		newBins = make([][]*g.Node, 0)
	}
	return convertBinsToGraph(colorBins, gr)
}


