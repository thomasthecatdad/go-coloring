package reductions

import (
	"fmt"
	g "github.com/thomaseb191/go-coloring/graphs"
	"log"
	"math"
	"math/rand"
)

// A Forest is defined as a collection of disjoint trees, where each node is a ForestNode
type Forest struct {
	ID    int
	Root  *ForestNode
	Nodes map[int]*ForestNode
}

// A ForestNode is a node in a Forest, including a Pointer to its relevant node in the original Graph, as well as color, tempcolor, and neighbor information
// Forests are directed, meaning that every node except roots will have a Parent
type ForestNode struct {
	Pointer   *g.Node
	Color     int
	TempColor int
	Parent    *ForestNode   //set later for optimization
	Neighbors []*ForestNode //includes Parent
}

// myChannelData is a struct for communicating data between the goroutines of a CV reduction algorithm.
//	Op refers to the current operation state
//	Val is the most important value
//	Extra is the next most important value
//	F is a pointer to a Forest if necessary
//	Threshold is the level to which down shifting should occur
type myChannelData struct {
	Op        int
	Val       int
	Extra     int
	F         *Forest
	Threshold int
}

//global variables to keep track of which color value to set and the number of nodes in the overall graph
var isTemp bool
var numAllNodes int

// CVReduction is based on the Cole-Vishkin color reduction algorithm and is comprised of the following steps
//		Creation of a worker pool of goroutines
//		[Parallel] Decomposition into maxDegree Forests
//		[Parallel] CV Reduction of each forest into 6 colors
//		[Parallel] Down-shifting of 6 colors into 3 colors
//		Unification of the various forests into a MaxDegree+1 coloring
// CVReduction is based on https://www.cs.bgu.ac.il/~elkinm/book.pdf and https://www.mpi-inf.mpg.de/fileadmin/inf/d1/teaching/winter15/tods/ToDS.pdf
// It is described as having O(Delta^2) + logstar(n) runtime. Because of practical Forest Decomposition, however, our algorithm runs in O(Delta^2) + logstar(n) + O(n) time
func CVReduction(gr g.Graph, poolSize int, debug int) g.Graph {
	if debug%2 == 1 {
		fmt.Printf("Starting CV Reduction \n")
	}
	isTemp = true
	numAllNodes = len(gr.Nodes)
	mainChannel := make(chan myChannelData)
	channels := buildWorkers(gr, poolSize, mainChannel, debug)

	if debug%2 == 1 {
		fmt.Printf("\tStarting Forest Decomposition \n")
	}

	forests := forestDecomposition(gr, channels, mainChannel, debug)
	for _, f := range forests {
		for _, node := range f.Nodes {
			bfsForest(node)
		}
		// Cannot be made parallel for bfs
	}

	if debug%2 == 1 {
		fmt.Printf("\tStarting CV to 6 \n")
	}

	for _, f := range forests {
		cvForestTo6(f, channels, mainChannel, debug)
		shiftDown(f, channels, mainChannel, debug)
		//printForest(f)
	}

	if debug%2 == 1 {
		fmt.Printf("\tStarting Forest Unification \n")
	}

	unifyForests2(forests, &gr)
	return gr
}

// forestDecomposition is the leader implementation of Forest Decomposition
func forestDecomposition(gr g.Graph, c []chan myChannelData, mainChan chan myChannelData, debug int) []*Forest {
	op := 0
	numDone := 0
	numChannels := len(c)

	myForests := make([]*Forest, gr.MaxDegree)
	for i := 0; i < gr.MaxDegree; i++ {
		myForests[i] = &Forest{
			ID:    i,
			Nodes: make(map[int]*ForestNode),
		}
	}

	for i, ch := range c {
		startingInd := len(gr.Nodes) / numChannels * i
		endingInd := len(gr.Nodes) / numChannels * (i + 1)
		if i == len(c)-1 {
			endingInd = len(gr.Nodes)
		}
		ch <- myChannelData{
			Op:    op,
			Val:   startingInd,
			Extra: endingInd,
		}
	}

	for numDone < numChannels {
		rec := <-mainChan
		if rec.Op == -1 {
			numDone++
		} else {
			forestToAdd := rec.Op
			parent := gr.Nodes[rec.Val]
			child := gr.Nodes[rec.Extra]
			fPtr := myForests[forestToAdd]

			existingC, ok := fPtr.Nodes[child.Ind]
			if !ok {
				newChild := ForestNode{
					Pointer:   child,
					Color:     child.Color,
					TempColor: child.Color,
					Parent:    nil,
					Neighbors: make([]*ForestNode, 0),
				}
				fPtr.Nodes[child.Ind] = &newChild
				existingC = &newChild
			}

			existingP, ok := fPtr.Nodes[parent.Ind]
			if !ok {
				newParent := ForestNode{
					Pointer:   parent,
					Color:     parent.Color,
					TempColor: parent.Color,
					Parent:    nil,
					Neighbors: []*ForestNode{existingC},
				}
				fPtr.Nodes[parent.Ind] = &newParent
				existingP = &newParent
			}
			existingC.Neighbors = append(existingC.Neighbors, existingP)

			if fPtr.Root == nil {
				fPtr.Root = existingP
			}
		}
	}
	return myForests
}

// forestDecompositionWorker is the worker implementation of Forest Decomposition based on the Panconesi and Rizzi Decomposition
// Each edge that is divided into a forest is reported to the main thread
func forestDecompositionWorker(gr *g.Graph, startingInd int, endingInd int, mainChannel chan myChannelData) {
	for k := startingInd; k < endingInd && k < len(gr.Nodes); k++ {
		currNode := gr.Nodes[k]
		starter := rand.Intn(len(currNode.Neighbors))
		for i, n := range currNode.Neighbors {
			if n.Ind < currNode.Ind {
				forest := (starter + i) % len(currNode.Neighbors)
				mainChannel <- myChannelData{
					Op:    forest,
					Val:   currNode.Ind,
					Extra: n.Ind,
				}
			}
		}
	}
	mainChannel <- myChannelData{
		Op: -1,
	}
}

// cvForestTo6 is the leader implementation of CV for a given Forest
func cvForestTo6(f *Forest, c []chan myChannelData, mainChan chan myChannelData, debug int) {
	op := 1

	numChannels := len(c)

	for i := 0; i < logStar(float64(len(f.Nodes)))+3; i++ {
		numDone := 0
		for k, ch := range c {
			startingInd := k
			step := numChannels
			ch <- myChannelData{
				Op:    op + i%2, //1 if set to TempColor, 2 if set to Color
				Val:   startingInd,
				Extra: step,
				F:     f,
			}
		}
		for numDone < numChannels {
			rec := <-mainChan
			if rec.Op == -1 {
				numDone++
			}
		}
		isTemp = !isTemp
	}
}

//cvForestTo6Worker is the worker implementation of Cole-Vishkin, setting the new color (either Color or TempColor) accordingly
func cvForestTo6Worker(f *Forest, startingInd int, step int, mainChannel chan myChannelData) {
	//change from previous iteration allows for more likelihood that each channel will have valid work to do
	for k := startingInd; k <= numAllNodes; k += step {
		currNode, ok := f.Nodes[k]
		if ok {
			parent := currNode.Parent
			if parent == nil {
				if isTemp {
					currNode.TempColor = calcColorRoot(currNode.Color)
				} else {
					currNode.Color = calcColorRoot(currNode.TempColor)
				}
			} else {
				if isTemp {
					if currNode.Color == parent.Color {
						log.Fatalf("me and parent is temp same! %d, %s, %s, %t\n%d, %d\n", currNode.Color, currNode.Pointer.Name, parent.Pointer.Name, parent.Parent == nil, currNode.TempColor, parent.TempColor)
					}
					currNode.TempColor = calcColor(currNode.Color, parent.Color)
				} else {
					if currNode.TempColor == parent.TempColor {
						log.Fatalf("me and parent is same! %d, %s, %s, %t\n", currNode.TempColor, currNode.Pointer.Name, parent.Pointer.Name, parent.Parent == nil)
					}
					currNode.Color = calcColor(currNode.TempColor, parent.TempColor)
				}
			}
		}
	}
	mainChannel <- myChannelData{
		Op: -1,
	}
}

// shiftDown is the leader implementation of down shifting process to reduce 6-color Forests to 3-color Forests
func shiftDown(f *Forest, c []chan myChannelData, mainChan chan myChannelData, debug int) {
	op := 3

	numChannels := len(c)

	for i := 0; i < 3; i++ {
		numDone := 0
		for k, ch := range c {
			startingInd := k
			step := numChannels
			ch <- myChannelData{
				Op:    op + i%2, //3 if set to TempColor, 4 if set to Color
				Val:   startingInd,
				Extra: step,
				F:     f,
			}
		}
		for numDone < numChannels {
			rec := <-mainChan
			if rec.Op == -1 {
				numDone++
			}
		}

		numDone = 0
		for k, ch := range c {
			startingInd := k
			step := numChannels
			ch <- myChannelData{
				Op:        op + 2 + i%2, //5 if set to TempColor, 6 if set to Color
				Val:       startingInd,
				Extra:     step,
				Threshold: 6 - i,
				F:         f,
			}
		}
		for numDone < numChannels {
			rec := <-mainChan
			if rec.Op == -1 {
				numDone++
			}
		}
		isTemp = !isTemp
	}
	for _, k := range f.Nodes {
		if !isTemp {
			k.Color = k.TempColor
		}
	}
}

// shiftDownWorker is the worker implementation of the first stage of the down shift algorithm
func shiftDownWorker(f *Forest, startingInd int, step int, mainChannel chan myChannelData) {
	//change from previous iteration allows for more likelihood that each channel will have valid work to do
	for k := startingInd; k <= numAllNodes; k += step {
		currNode, ok := f.Nodes[k]
		if ok {
			parent := currNode.Parent
			if parent == nil {
				if isTemp {
					newColor := currNode.Color
					for newColor == currNode.Color {
						newColor = rand.Intn(3)
					}
					currNode.TempColor = newColor
				} else {
					newColor := currNode.TempColor
					for newColor == currNode.TempColor {
						newColor = rand.Intn(3)
					}
					currNode.Color = newColor
				}
			} else {
				if isTemp {
					currNode.TempColor = parent.Color
				} else {
					currNode.Color = parent.TempColor
				}
			}
		}
	}
	mainChannel <- myChannelData{
		Op: -1,
	}
}

// shiftDownWorker is the worker implementation of the second stage of the down shift algorithm
func shiftDownWorkerCleanup(f *Forest, startingInd int, step int, thresh int, mainChannel chan myChannelData) {
	//change from previous iteration allows for more likelihood that each channel will have valid work to do
	for k := startingInd; k <= numAllNodes; k += step {
		currNode, ok := f.Nodes[k]
		if ok {
			if isTemp {
				newColor := calcSafeReduction(currNode, thresh)
				currNode.TempColor = newColor
			} else {
				newColor := calcSafeReduction(currNode, thresh)
				currNode.Color = newColor
			}
		}
	}
	mainChannel <- myChannelData{
		Op: -1,
	}
}

// unifyForests is a deprecated first attempt at unifying separate Forests
func unifyForests(forests []*Forest, gr *g.Graph) {
	for _, k := range gr.Nodes {
		fPtr, ok := forests[0].Nodes[k.Ind]
		if ok {
			k.Color = fPtr.Color
		} else {
			k.Color = -1
		}
	}
	for i := 1; i < len(forests); i++ {
		for _, k := range gr.Nodes {
			fPtr, ok := forests[i].Nodes[k.Ind]
			if ok && k.Color != -1 {
				//k.Color = calcColor(k.Color, fPtr.Color)
				k.Color = k.Color | fPtr.Color
			} else if ok {
				k.Color = fPtr.Color
			} else if k.Color != -1 {
				//keep the same
			} else {
				k.Color = -1
			}
		}
	}
}

// unifyForests2 is the leader implementation (less efficient) to unify Forests
func unifyForests2(forests []*Forest, gr *g.Graph) {
	for _, k := range gr.Nodes {
		k.Color = -1
	}
	for _, k := range gr.Nodes {
		options := makeRange(0, gr.MaxDegree+1)
		//handle colors already set
		for _, i := range k.Neighbors {
			if i.Color != -1 {
				indInOptions := findIndexOf(options, i.Color)
				if indInOptions != -1 {
					options = remove(options, indInOptions)
				}
			}
		}
		//handle colors not set
		for _, j := range forests {
			fPtr, ok := j.Nodes[k.Ind]
			if ok {
				for _, i := range fPtr.Neighbors {
					if i.Pointer.Color == -1 {
						indInOptions := findIndexOf(options, i.Color)
						if indInOptions != -1 {
							options = remove(options, indInOptions)
						}
					}
				}
			} else {
				continue
			}
		}
		k.Color = options[rand.Intn(len(options))]
	}
}

// findIndexOf searches an array for the desired color, returning -1 if not found
func findIndexOf(options []int, color int) int {
	for ind, k := range options {
		if k == color {
			return ind
		}
	}
	return -1
}

// workerWait is the overall manager for workers, governing the division into different subalgorithms
func workerWait(gr *g.Graph, c chan myChannelData, mainChannel chan myChannelData) {
	rec := <-c
	if rec.Op == 0 {
		forestDecompositionWorker(gr, rec.Val, rec.Extra, mainChannel)
	} else {
		log.Fatal("Wrong op received by worker")
	}

	rec = <-c
	for rec.Op < 7 {
		if rec.Op == 1 || rec.Op == 2 {
			cvForestTo6Worker(rec.F, rec.Val, rec.Extra, mainChannel)
		} else if rec.Op == 3 || rec.Op == 4 {
			shiftDownWorker(rec.F, rec.Val, rec.Extra, mainChannel)
		} else {
			shiftDownWorkerCleanup(rec.F, rec.Val, rec.Extra, rec.Threshold, mainChannel)
		}
		rec = <-c
	}

}

// logStar is the logstar function defined in the CV paper https://www.cs.bgu.ac.il/~elkinm/book.pdf
func logStar(n float64) int {
	if n <= 2 {
		return 0
	} else {
		return 1 + logStar(math.Log2(n))
	}
}

// buildWorkers creates a desired number of workers based on input specifications or a default
func buildWorkers(gr g.Graph, poolSize int, mainChannel chan myChannelData, debug int) []chan myChannelData {
	defaultPool := math.Floor(math.Sqrt(float64(len(gr.Nodes)))) //TODO: ADJUST DEFAULT AS NECESSARY
	var numWorkers int
	if poolSize <= 0 {
		numWorkers = int(defaultPool)
	} else {
		numWorkers = int(math.Min(float64(poolSize), defaultPool))
	}

	if debug%2 == 1 {
		fmt.Printf("\tBuilding %d workers for %d nodes \n", numWorkers, len(gr.Nodes))
	}

	var myChannels []chan myChannelData
	for i := 0; i < numWorkers; i++ {
		c := make(chan myChannelData)
		myChannels = append(myChannels, c)
		go workerWait(&gr, c, mainChannel)
	}
	return myChannels
}

// bfsForest runs a BFS on a given node, orienting the tree and setting parents. Cannot safely be run in parallel
func bfsForest(f *ForestNode) {
	if f != nil {
		for _, k := range f.Neighbors {
			if k.Parent == nil && f.Parent != k {
				k.Parent = f
				bfsForest(k)
			}
		}
	}
}

// calcColor is the crux of the CV algorithm
// implementation borrowed from https://www.zhengqunkoo.com:8443/zhengqunkoo/site/src/commit/ebbab6e24911a02c97b380f2e39f06d9c3e83770/worker.js
func calcColor(me int, parent int) int {
	me1 := uint32(me)
	parent1 := uint32(parent)
	if me == parent {
		log.Fatal("me and parent is same! ", me, parent)
	}

	j := getDifferBitIndex(me1, parent1)
	//midx := getBitLength(me1) - j - 1
	midx := 32 - j - 1
	return int((j << 1) | (me1&(1<<midx))>>midx)
}

// getDifferBitIndex returns the bit index from the left at which 2 uint32s vary (Big Endian)
func getDifferBitIndex(x, y uint32) uint32 {
	if x == y {
		log.Fatal("No differing bit")
	} else {
		bl := uint32(32)
		m := uint32(1) << (bl - 1)
		var idx = 0
		for m != 0 {
			if (x & m) != (y & m) {
				return uint32(idx)
			}
			m = m >> 1
			idx++
		}
	}
	log.Fatal("Why is index returned 0?")
	return 0
}

// getBitLength is a function that returns the bit length of a number. Deprecated in favor of uint32 standardization
func getBitLength(x uint32) uint32 {
	// Minimum bit length is 3
	if x < 6 {
		return 3
	} else {
		var length = 0
		for x != 0 {
			x = x >> 1
			length++
		}
		return uint32(length)
	}
}

// calcColorRoot is the same as calcColor but uses an arbitrary index of 30 (0-indexed)
func calcColorRoot(me int) int {
	me1 := uint32(me)
	j := uint32(30)
	//midx := getBitLength(me1) - j - 1
	midx := 32 - j - 1
	return int((j << 1) | (me1&(1<<midx))>>midx)
}

// calcSafeReduction is used to calculate a safe color to set during the down shifting process
func calcSafeReduction(n *ForestNode, thresh int) int {
	if isTemp {
		if n.TempColor < thresh {
			return n.TempColor
		}
		var colorProposal int
		safe := false
		for !safe {
			colorProposal = rand.Intn(3)
			for i, _ := range n.Neighbors {
				if n.Neighbors[i].TempColor == colorProposal {
					i = 0 //reset the progress
					colorProposal = rand.Intn(3)
				}
			}
			safe = true
		}
		return colorProposal
	} else {
		if n.Color < thresh {
			return n.Color
		}
		var colorProposal int
		safe := false
		for !safe {
			colorProposal = rand.Intn(3)
			for i, _ := range n.Neighbors {
				if n.Neighbors[i].Color == colorProposal {
					i = 0 //reset the progress
					colorProposal = rand.Intn(3)
				}
			}
			safe = true
		}
		return colorProposal
	}
}

// printForest prints the metadata for a Forest and then all of its nodes in an established format
func printForest(f *Forest) {
	fmt.Printf("Graph Name: \t\t%d\n", f.ID)
	for _, v := range f.Nodes {
		neighborNames := GetNamesFromNodeList(v.Neighbors)
		parent := v.Parent
		parentName := "nil"
		if parent != nil {
			parentName = v.Parent.Pointer.Name
		}
		fmt.Printf("Node: %s,\tColor: %d,\tNeighbors: %s,\tParent: %s\n", v.Pointer.Name, v.Color, neighborNames, parentName)
	}
	fmt.Printf("---End of Forest [%d].\n", f.ID)
}

// GetNamesFromNodeList converts a list of Node pointers to their names
// This function is mainly useful for printing
func GetNamesFromNodeList(neighborNodes []*ForestNode) []string {
	var neighborNames []string
	for _, node := range neighborNodes {
		neighborNames = append(neighborNames, node.Pointer.Name)
	}
	return neighborNames
}

// makeRange implementation from https://stackoverflow.com/questions/39868029/how-to-generate-a-sequence-of-numbers
func makeRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// remove implementation from https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang#:~:text=Slices%20and%20arrays%20being%200,so%20on%20and%20so%20forth.
func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}
