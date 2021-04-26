package reductions

import (
	"fmt"
	g "github.com/thomaseb191/go-coloring/graphs"
)
// runNaiveGoRoutine is a helper function that runs the naive algorithm as a goroutine and
// sends the result back through a channel.
func runNaiveGoRoutine(gr g.Graph, poolSize int, debug int, c chan g.Graph) {
	c <- RunNaive(gr, poolSize, debug)
}

// convertBinsToGraph is a helper method that converts color "bins" into graphs.
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

// combineColors is a helper method that runs the naive algorithms on smaller bins and combines them
// into a new set of smaller bins.
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

func combineColorsSequential(bins [][]*g.Node, gr g.Graph) [][]*g.Node{
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
	return newBins
}

func checkIfNodeInColorSet(colorSet []*g.Node, neighbors []*g.Node) bool {
	hasAny := false
	for _, k := range neighbors {
		for _, node := range colorSet {
			if k.Name == node.Name {
				hasAny = true
				break
			}
		}

	}
	return hasAny
}

func combineColorsWithoutNaive(bins [][]*g.Node, gr g.Graph, c chan [][]*g.Node) {
	maxDegree := gr.MaxDegree
	fmt.Printf("Number of colors in bins: %d\n", len(bins))
	for k := maxDegree + 1; k < len(bins); k++ {
		for j := 0; j < len(bins[k]); j++ {
			for color := 0; color < maxDegree + 1; color++ {
				hasAny := checkIfNodeInColorSet(bins[color], gr.Nodes[j].Neighbors)
				if ! hasAny {
					bins[color] = append(bins[color], gr.Nodes[j])
					break
				}
			}
		}
	}
	if len(bins) < maxDegree + 1 {
		c <- bins
	} else {
		c <- bins[:maxDegree + 1]
	}
}

// kwReduction is the main method that runs the KW algorithm.
func kwReduction(gr g.Graph, poolSize int, debug int) g.Graph {
	fmt.Printf("Starting KW Reduction \n")
	degree := gr.MaxDegree
	startIndexes := make([]int, 0)
	size := len(gr.Nodes)
	c := make(chan g.Graph)
	// If we can't split the graph into bins,
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

	//for i := 0; i < len(startIndexes); i++ {
	//	currStart := startIndexes[i]
	//	var nextStart int
	//	if i + 1 != len(startIndexes) {
	//		nextStart = startIndexes[i + 1]
	//	} else {
	//		nextStart = size
	//	}
	//	grCopy := g.Graph {
	//		Name: gr.Name,
	//		Description: gr.Description,
	//		MaxDegree: gr.MaxDegree,
	//		Nodes: gr.Nodes[currStart:nextStart],
	//	}
	//	go runNaiveGoRoutine(grCopy, poolSize, debug, c)
	//}
	// colorBins := make([][]*g.Node, 0)

	//for i := 0; i < len(startIndexes); i++{
	//	graph := <- c
	//	numColors := g.CountColors(&graph)
	//	newBins := make([][]*g.Node, numColors)
	//	colorToIndex := make(map[int]int)
	//	latestIndex := 0
	//	for _, node := range graph.Nodes {
	//		if _, ok := colorToIndex[node.Color]; ! ok {
	//			colorToIndex[node.Color] = latestIndex
	//			latestIndex++
	//		}
	//		newBins[colorToIndex[node.Color]] = append(newBins[colorToIndex[node.Color]], node)
	//	}
	//	colorBins = append(colorBins, newBins...)
	//}
	//close(c)

	numColors := g.CountColors(&gr)
	colorBins := make([][]*g.Node, numColors)
	colorToIndex := make(map[int]int)
	latestIndex := 0
	for _, node := range gr.Nodes {
		if _, ok := colorToIndex[node.Color]; ! ok {
			colorToIndex[node.Color] = latestIndex
			latestIndex++
		}
		colorBins[colorToIndex[node.Color]] = append(colorBins[colorToIndex[node.Color]], node)
	}
	// colorBins = append(colorBins, newBins...)


	for len(colorBins) > degree+1 {
		fmt.Printf("Number of bins: %d\n", len(colorBins))
		printBins(colorBins)
		d := make(chan [][]*g.Node)
		binIndexes := make([]int, 0)
		colors := len(colorBins)

		for x := 0; x < colors; x++ {
			if x%(2*(degree+1)) == 0 {
				binIndexes = append(binIndexes, x)
			}
		}
		tempBins := make([][]*g.Node, 0)

		for i := 0; i < len(binIndexes); i++ {
			currStart := binIndexes[i]
			var nextStart int
			if i+1 != len(binIndexes) {
				nextStart = binIndexes[i+1]
			} else {
				nextStart = len(colorBins)
			}
			go combineColorsWithoutNaive(colorBins[currStart:nextStart], gr, d)
		}
		for i := 0; i < len(binIndexes); i++ {
			bins := <-d
			tempBins = append(tempBins, bins...)
		}
		close(d)
		//for i := 0; i < len(binIndexes); i++ {
		//	currStart := binIndexes[i]
		//		var nextStart int
		//		if i + 1 != len(binIndexes) {
		//			nextStart = binIndexes[i + 1]
		//		} else {
		//			nextStart = len(colorBins)
		//		}
		//		bins := combineColorsSequential(colorBins[currStart:nextStart], gr)
		//		tempBins = append(tempBins, bins...)
		//}
		colorBins = tempBins
		tempBins = make([][]*g.Node, 0)
	}
	return convertBinsToGraph(colorBins, gr)
}

func printBins(colorBins [][]*g.Node) {
	for i, b := range colorBins {
		fmt.Printf("Bin #%d\t", i)
		for _, n := range b {
			fmt.Printf("%s ", n.Name)
		}
		fmt.Printf("\n")
	}
}


