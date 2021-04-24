package reductions

import (
	"fmt"
	g "github.com/thomaseb191/go-coloring/graphs"
	"log"
	"math"
	"math/rand"
)

type Forest struct {
	ID int
	Root *ForestNode
	Nodes map[int]*ForestNode
}

type ForestNode struct {
	Pointer *g.Node
	Color int
	TempColor int
	Parent *ForestNode //set later for optimization
	Neighbors []*ForestNode //includes Parent
}

type myChannelData struct {
	Op int
	Val int
	Extra int
	F *Forest
	Threshold int
}

//// SafeCounter taken from https://tour.golang.org/concurrency/9
//type SafeCounter struct {
//	mu sync.Mutex
//	v  int
//}
//
//func safeInc() int {
//	c.mu.Lock()
//	val := c.v
//	// Lock so only one goroutine at a time can access the map c.v.
//	c.v++
//	c.mu.Unlock()
//	return val
//}
//
//var c SafeCounter
var isTemp bool
var numAllNodes int

func forestDecomposition(gr g.Graph, c []chan myChannelData, mainChan chan myChannelData, debug int) []*Forest {
	op := 0
	numDone := 0
	numChannels := len(c)

	myForests := make([]*Forest, gr.MaxDegree)
	for i := 0; i < gr.MaxDegree; i ++ {
		myForests[i] = &Forest{
			ID: i,
			Nodes: make(map[int]*ForestNode),
		}
	}

	for i, ch := range c {
		startingInd := len(gr.Nodes) / numChannels * i
		endingInd := len(gr.Nodes) / numChannels * (i + 1)
		ch <- myChannelData{
			Op: op,
			Val: startingInd,
			Extra: endingInd,
		}
	}

	for numDone < numChannels {
		rec := <- mainChan
		if rec.Op == -1 {
			numDone ++
		} else {
			forestToAdd := rec.Op
			parent := gr.Nodes[rec.Val]
			child := gr.Nodes[rec.Extra]
			fPtr := myForests[forestToAdd]

			existingC, ok := fPtr.Nodes[child.Ind]
			if !ok {
				newChild := ForestNode{
					Pointer: child,
					Color: child.Color,
					Parent: nil,
					Neighbors: make([]*ForestNode, 0),
				}
				fPtr.Nodes[child.Ind] = &newChild
				existingC = &newChild
			}

			existingP, ok := fPtr.Nodes[parent.Ind]
			if !ok {
				newParent := ForestNode{
					Pointer: parent,
					Color: parent.Color,
					Parent: nil,
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

func forestDecompositionWorker(gr *g.Graph, startingInd int, endingInd int, mainChannel chan myChannelData) {
	for k := startingInd; k < endingInd && k < len(gr.Nodes); k++ {
		currNode := gr.Nodes[k]
		starter := rand.Intn(len(currNode.Neighbors))
		for i, n := range currNode.Neighbors {
			if n.Ind < currNode.Ind {
				forest := (starter + i) % len(currNode.Neighbors)
				mainChannel <- myChannelData{
					Op: forest,
					Val: currNode.Ind,
					Extra: n.Ind,
				}
			}
		}
	}
	mainChannel <- myChannelData{
		Op: -1,
	}
}

func cvForestTo6(f *Forest, c []chan myChannelData, mainChan chan myChannelData, debug int) {
	op := 1

	numChannels := len(c)

	for i := 0; i < logStar(float64(len(f.Nodes))) + 3; i ++ {
		numDone := 0
		for i, ch := range c {
			//startingInd := len(f.Nodes) / numChannels * i
			//endingInd := len(f.Nodes) / numChannels * (i + 1)
			startingInd := i
			step := numChannels
			ch <- myChannelData{
				Op: op + i % 2, //1 if set to TempColor, 2 if set to Color
				Val: startingInd,
				Extra: step,
				F: f,
			}
		}
		for numDone < numChannels {
			rec := <-mainChan
			if rec.Op == -1 {
				numDone ++
			}
		}
		isTemp = !isTemp
	}
}

func cvForestTo6Worker(f *Forest, startingInd int, step int, mainChannel chan myChannelData) {
	//change from previous iteration allows for more likelihood that each channel will have valid work to do
	for k := startingInd; k < len(f.Nodes); k += step {
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
					fmt.Println("isTemp!", currNode.Pointer.Name, currNode.Color, parent.Pointer.Name, parent.Color)
					currNode.TempColor = calcColor(currNode.Color, parent.Color)
				} else {
					fmt.Println("is not Temp!", currNode.Pointer.Name, currNode.Color, parent.Pointer.Name, parent.Color)
					currNode.Color = calcColor(currNode.TempColor, parent.TempColor)
				}
			}
		}
	}
	mainChannel <- myChannelData{
		Op: -1,
	}
}

func shiftDown(f *Forest, c []chan myChannelData, mainChan chan myChannelData, debug int) {
	op := 3

	numChannels := len(c)

	for i := 0; i < 3; i ++ {
		numDone := 0
		for i, ch := range c {
			startingInd := i
			step := numChannels
			ch <- myChannelData{
				Op: op + i % 2, //3 if set to TempColor, 4 if set to Color
				Val: startingInd,
				Extra: step,
				F: f,
			}
		}
		for numDone < numChannels {
			rec := <-mainChan
			if rec.Op == -1 {
				numDone ++
			}
		}

		numDone = 0
		for i, ch := range c {
			startingInd := i
			step := numChannels
			ch <- myChannelData{
				Op: op + 2 + i % 2, //5 if set to TempColor, 6 if set to Color
				Val: startingInd,
				Extra: step,
				Threshold: 6-i,
				F: f,
			}
		}
		for numDone < numChannels {
			rec := <-mainChan
			if rec.Op == -1 {
				numDone ++
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

func shiftDownWorker(f *Forest, startingInd int, step int, mainChannel chan myChannelData) {
	//change from previous iteration allows for more likelihood that each channel will have valid work to do
	for k := startingInd; k < len(f.Nodes); k += step {
		currNode, ok := f.Nodes[k]
		if ok {
			parent := currNode.Parent
			if parent == nil {
				if isTemp {
					newColor := currNode.Color
					for newColor != currNode.Color {
						newColor = rand.Intn(3)
					}
					currNode.TempColor = newColor
				} else {
					newColor := currNode.TempColor
					for newColor != currNode.TempColor {
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

func shiftDownWorkerCleanup(f *Forest, startingInd int, step int, thresh int, mainChannel chan myChannelData) {
	//change from previous iteration allows for more likelihood that each channel will have valid work to do
	for k := startingInd; k < len(f.Nodes); k += step {
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
				fmt.Printf("HERE, unifying %d and %d to %d\n", k.Color, fPtr.Color, k.Color & fPtr.Color)
				//k.Color = calcColor(k.Color, fPtr.Color) //TODO: FIX IF NECESSARY
				k.Color = k.Color & fPtr.Color
			} else if ok {
				k.Color = fPtr.Color
			} else if k.Color != -1 {

			} else {
				k.Color = -1
			}
		}
	}
}



func workerWait(gr *g.Graph, c chan myChannelData, mainChannel chan myChannelData) {
	rec := <- c
	if rec.Op == 0 {
		forestDecompositionWorker(gr, rec.Val, rec.Extra, mainChannel)
	} else {
		log.Fatal("Wrong op received by worker")
	}

	rec = <- c
	for rec.Op < 7 {
		if rec.Op == 1 || rec.Op == 2 {
			cvForestTo6Worker(rec.F, rec.Val, rec.Extra, mainChannel)
		} else if rec.Op == 3 || rec.Op == 4 {
			shiftDownWorker(rec.F, rec.Val, rec.Extra, mainChannel)
		} else {
			shiftDownWorkerCleanup(rec.F, rec.Val, rec.Extra, rec.Threshold, mainChannel)
		}
		rec = <- c
	}



}

func CVReduction(gr g.Graph, poolSize int, debug int) g.Graph {
	if debug % 2 == 1 {
		fmt.Printf("Starting CV Reduction \n")
	}
	isTemp = true
	numAllNodes = len(gr.Nodes)
	mainChannel := make(chan myChannelData)
	channels := buildWorkers(gr, poolSize, mainChannel, debug)

	if debug % 2 == 1 {
		fmt.Printf("\tStarting Forest Decomposition \n")
	}

	forests := forestDecomposition(gr, channels, mainChannel, debug)
	for _, f := range forests {
		bfsForest(f.Root) // TODO: MAKE PARALLEL
	}

	if debug % 2 == 1 {
		fmt.Printf("\tStarting CV to 6 \n")
	}

	for _, f := range forests {
		cvForestTo6(f, channels, mainChannel, debug)
		shiftDown(f, channels, mainChannel, debug)
		printForest(f)
	}

	if debug % 2 == 1 {
		fmt.Printf("\tStarting Forest Unification \n")
	}

	unifyForests(forests, &gr)
	return gr
}

func logStar(n float64) int {
	if n <= 2 {
		return 0
	} else {
		return 1 + logStar(math.Log2(n))
	}
}

func buildWorkers(gr g.Graph, poolSize int, mainChannel chan myChannelData, debug int) []chan myChannelData {
	defaultPool := math.Floor(math.Sqrt(float64(len(gr.Nodes)))) //TODO: DETERMINE DEFAULT
	var numWorkers int
	if poolSize <= 0 {
		numWorkers = int(defaultPool)
	} else {
		numWorkers = int(math.Min(float64(poolSize), defaultPool))
	}

	if debug % 2 == 1 {
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

func bfsForest(f *ForestNode) {
	if f != nil {
		for _, k := range f.Neighbors {
			if k != f.Parent {
				k.Parent = f
				bfsForest(k)
			}
		}
	}
}

// implementation borrowed from https://github.com/BenWiederhake/cole-vishkin/blob/master/cv.cpp
func calcColor(me int, parent int) int {
	fmt.Println("calc color", me, parent)
	//i := findIndexFromRight(me, parent)
	//i2 := i
	//if i2 == 0 {
	//	i2 = 1
	//}
	//shamt := int(math.Ceil(math.Log2(float64(i2))))
	//return (((me >> i) & 1) << shamt) + i

	xored := me ^ parent
	num := trailingZeros(xored)
	orig_bit := 1 & (me >> num)
	return orig_bit | (num << 1)
}

func trailingZeros(me int) int {
	if me == 0 {
		return 0
	} else if me % 2 == 1 {
		return 0
	} else {
		return trailingZeros(me >> 1) + 1
	}
}

func calcColorRoot(me int) int {
	//i2 := me
	//if i2 == 0 {
	//	i2 = 1
	//}
	//shamt := int(math.Ceil(math.Log2(float64(i2))))
	//i := shamt
	//return (((me >> i) & 1) << shamt) + i

	//xored := me ^ parent
	//num := trailingZeros(xored)
	num := 1
	orig_bit := 1 & (me >> num)
	return orig_bit | (num << 1)
}

func findIndexFromRight(me int, parent int) int {
	if me == parent {
		log.Fatalf("Parent and Child same color: %d\n", me, parent)
	}
	if me % 2 != parent % 2 {
		return 0
	} else {
		return 1 + findIndexFromRight(me >> 1, parent >> 1)
	}
}

func calcSafeReduction(n *ForestNode, thresh int) int {
	if isTemp {
		if n.TempColor < thresh {
			return n.TempColor
		}
		var colorProposal int
		safe := false
		for !safe {
			colorProposal = rand.Intn(3)
			for _, k :=  range n.Neighbors {
				if k.TempColor == colorProposal {
					continue
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
			for _, k :=  range n.Neighbors {
				if k.Color == colorProposal {
					continue
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
		fmt.Printf("Node: %s,\tColor: %d,\tNeighbors: %s\n", v.Pointer.Name, v.Color, neighborNames)
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