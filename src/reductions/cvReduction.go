package reductions

import (
	"fmt"
	g "github.com/thomaseb191/go-coloring/graphs"
	"log"
	"math"
	"math/rand"
)

type Forest struct {
	ID    int
	Root  *ForestNode
	Nodes map[int]*ForestNode
}

type ForestNode struct {
	Pointer   *g.Node
	Color     int
	TempColor int
	Parent    *ForestNode   //set later for optimization
	Neighbors []*ForestNode //includes Parent
}

type myChannelData struct {
	Op        int
	Val       int
	Extra     int
	F         *Forest
	Threshold int
}

var isTemp bool
var numAllNodes int

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
		if i == len(c) - 1 {
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
				if child.Color == 0 {
					fmt.Printf("CHILD COLOR IS 0 FOR %s\n", child.Name)
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
				if parent.Color == 0 {
					fmt.Printf("PARENT COLOR IS 0 FOR %s\n", parent.Name)
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

func cvForestTo6(f *Forest, c []chan myChannelData, mainChan chan myChannelData, debug int) {
	op := 1

	numChannels := len(c)

	for i := 0; i < logStar(float64(len(f.Nodes)))+3; i++ {
		//fmt.Printf("WE'RE DOING THIS FOR %d ITERATIONS.\n", logStar(float64(len(f.Nodes))))
		numDone := 0
		for k, ch := range c {
			//startingInd := len(f.Nodes) / numChannels * i
			//endingInd := len(f.Nodes) / numChannels * (i + 1)
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

func cvForestTo6Worker(f *Forest, startingInd int, step int, mainChannel chan myChannelData) {
	//change from previous iteration allows for more likelihood that each channel will have valid work to do
	for k := startingInd; k <= numAllNodes; k += step {
		currNode, ok := f.Nodes[k]
		if ok {
			parent := currNode.Parent
			if parent == nil {
				if isTemp {
					currNode.TempColor = calcColorRoot(currNode.Color)
					if currNode.TempColor == 0 {
						fmt.Printf("I JUST SET %s (temp) TO 0\n", currNode.Pointer.Name)
					}
					//currNode.TempColor = calcColor(currNode.Color, currNode.Neighbors[0].Color)
					//currNode.TempColor = currNode.Color
				} else {
					currNode.Color = calcColorRoot(currNode.TempColor)
					//currNode.Color = calcColor(currNode.TempColor, currNode.Neighbors[0].TempColor)
					//currNode.Color = currNode.TempColor
					if currNode.Color == 0 {
						fmt.Printf("I JUST SET %s TO 0\n", currNode.Pointer.Name)
					}
				}
			} else {
				if isTemp {
					//fmt.Println("isTemp!", currNode.Pointer.Name, currNode.Color, parent.Pointer.Name, parent.Color)
					if currNode.Color == parent.Color {
						log.Fatalf("me and parent is temp same! %d, %s, %s, %t\n%d, %d\n", currNode.Color, currNode.Pointer.Name, parent.Pointer.Name, parent.Parent == nil, currNode.TempColor, parent.TempColor)
					}
					currNode.TempColor = calcColor(currNode.Color, parent.Color)
					if currNode.TempColor == 0 {
						fmt.Printf("I JUST SET %s (temp) TO 0\n", currNode.Pointer.Name)
					}
					fmt.Printf("%d calculated from %d and %d for %s\n", currNode.TempColor, currNode.Color, parent.Color, currNode.Pointer.Name)
				} else {
					//fmt.Println("is not Temp!", currNode.Pointer.Name, currNode.Color, parent.Pointer.Name, parent.Color)
					if currNode.TempColor == parent.TempColor {
						printForest(f)
						//TODO: FIX RIGHT HERE
						log.Fatalf("me and parent is same! %d, %s, %s, %t\n", currNode.TempColor, currNode.Pointer.Name, parent.Pointer.Name, parent.Parent == nil)
					}
					currNode.Color = calcColor(currNode.TempColor, parent.TempColor)
					if currNode.Color == 0 {
						fmt.Printf("I JUST SET %s TO 0\n", currNode.Pointer.Name)
					}
					fmt.Printf("%d calculated from %d and %d for %s\n", currNode.Color, currNode.TempColor, parent.TempColor, currNode.Pointer.Name)
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

func shiftDownWorker(f *Forest, startingInd int, step int, mainChannel chan myChannelData) {
	//change from previous iteration allows for more likelihood that each channel will have valid work to do
	for k := startingInd; k <= numAllNodes; k += step {
		currNode, ok := f.Nodes[k]
		if ok {
			//fmt.Printf("Examining \t\t %s\n", currNode.Pointer.Name)
			parent := currNode.Parent
			if parent == nil {
				if isTemp {
					newColor := currNode.Color
					for newColor == currNode.Color {
						newColor = rand.Intn(3)
					}
					//fmt.Printf("%s (root) shifted down from %d to %d\n", currNode.Pointer.Name, currNode.Color, newColor)
					currNode.TempColor = newColor
				} else {
					newColor := currNode.TempColor
					for newColor == currNode.TempColor {
						newColor = rand.Intn(3)
					}
					//fmt.Printf("%s (root) shifted down from %d to %d\n", currNode.Pointer.Name, currNode.TempColor, newColor)
					currNode.Color = newColor
				}
			} else {
				if isTemp {
					currNode.TempColor = parent.Color
					//fmt.Printf("%s shifted down from %d to %d\n", currNode.Pointer.Name, currNode.Color, currNode.TempColor)
				} else {
					currNode.Color = parent.TempColor
					//fmt.Printf("%s shifted down from %d to %d\n", currNode.Pointer.Name, currNode.TempColor, currNode.Color)
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
				fmt.Printf("HERE, unifying %d and %d to %d\n", k.Color, fPtr.Color, k.Color&fPtr.Color)
				//k.Color = calcColor(k.Color, fPtr.Color) //TODO: FIX IF NECESSARY
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

//func unifyForests3(forests []*Forest, gr *g.Graph) {
//
//}

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

func findIndexOf(options []int, color int) int {
	for ind, k := range options {
		if k == color {
			return ind
		}
	}
	return -1
}

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
		// TODO: MAKE PARALLEL
	}

	if debug%2 == 1 {
		fmt.Printf("\tStarting CV to 6 \n")
	}

	for _, f := range forests {
		cvForestTo6(f, channels, mainChannel, debug)
		shiftDown(f, channels, mainChannel, debug)
		printForest(f)
	}

	if debug%2 == 1 {
		fmt.Printf("\tStarting Forest Unification \n")
	}

	unifyForests2(forests, &gr)
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

func bfsForest(f *ForestNode) {
	if f != nil {
		for _, k := range f.Neighbors {
			if k.Parent == nil && f.Parent != k {
				//fmt.Printf("\tAssigning %s as parent of %s\n.", f.Pointer.Name, k.Pointer.Name)
				k.Parent = f
				bfsForest(k)
			}
		}
	}
}

// implementation borrowed from https://github.com/BenWiederhake/cole-vishkin/blob/master/cv.cpp
// implementation borrowed from https://www.zhengqunkoo.com:8443/zhengqunkoo/site/src/commit/ebbab6e24911a02c97b380f2e39f06d9c3e83770/worker.js
func calcColor(me int, parent int) int {
	me1 := uint32(me)
	parent1 := uint32(parent)
	if me == parent {
		log.Fatal("me and parent is same! ", me, parent)
	}
	/*
	//i := findIndexFromRight(me, parent)
	//i2 := i
	//if i2 == 0 {
	//	i2 = 1
	//}
	//shamt := int(math.Ceil(math.Log2(float64(i2))))
	//return (((me >> i) & 1) << shamt) + i

	//i := minFromZero(uint(me) ^ uint(parent)) //one indexed?
	//shamt := widthOfI(i)
	//origBit := 1 & (uint(me) >> (i-1))
	////return int(bits.Reverse(i | (orig_bit << shamt)))
	//return int(i | (origBit << shamt))

	//xored := me ^ parent
	//num := trailingZeros(xored)
	//orig_bit := 1 & (me >> num)
	//return orig_bit | (num << 1)
	*/

	j := getDifferBitIndex(me1, parent1)
	//midx := getBitLength(me1) - j - 1
	midx := 32 - j - 1
	fmt.Printf("calc color from %d and %d to %d\n", me, parent, int((j << 1) | (me1&(1<<midx))>>midx))
	fmt.Printf("\tdiffers at index %d\n", j)
	fmt.Printf("\tcalc color from %16b and %16b to %16b\n", me, parent, int((j << 1) | (me1&(1<<midx))>>midx))
	return int((j << 1) | (me1&(1<<midx))>>midx)
}

func getDifferBitIndex(x, y uint32) uint32 {
	if x == y {
		log.Fatal("No differing bit")
	} else {
		// Get the longest bit length of x and y
		//bl := 0
		//if x > y {
		//	bl = getBitLength(x)
		//} else {
		//	bl = getBitLength(y)
		//}
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

func minFromZero(u uint) uint {
	if 1&u == 1 {
		return 1
	}
	return 1 + minFromZero(u>>1)
}

func trailingZeros(me int) int {
	if me == 0 {
		return 0
	} else if me%2 == 1 {
		return 0
	} else {
		return trailingZeros(me>>1) + 1
	}
}

func calcColorRoot(me int) int {
	/*
	//i2 := me
	//if i2 == 0 {
	//	i2 = 1
	//}
	//shamt := int(math.Ceil(math.Log2(float64(i2))))
	//i := shamt
	//return (((me >> i) & 1) << shamt) + i

	//xored := me ^ parent
	//num := trailingZeros(xored)

	//num := 1
	//orig_bit := 1 & (me >> num)
	//return orig_bit | (num << 1)

		//i := minFromZero(uint(me) ^ uint(15)) //one indexed?
		//shamt := widthOfI(i)
		//origBit := 1 & (uint(me) >> (i-1))
		//return int(i | (origBit << shamt))
	*/
	//j := getDifferBitIndex(me, parent)
	//j := rand.Intn(getBitLength(me))  //TODO: FIX ROOT STUFF

	me1 := uint32(me)
	j := uint32(30)
	//midx := getBitLength(me1) - j - 1
	midx := 32 - j - 1
	//fmt.Printf("calc color from %d (root) to %d\n", me, (j << 1) | (me&(1<<midx))>>midx)
	//fmt.Printf("\tdiffers at index %d\n", j)
	//fmt.Printf("\tcalc color from %16b (root) to %16b\n", me, (j << 1) | (me&(1<<midx))>>midx)
	return int((j << 1) | (me1&(1<<midx))>>midx)
}

func findIndexFromRight(me int, parent int) int {
	if me == parent {
		log.Fatalf("Parent and Child same color: %d\n", me, parent)
	}
	if me%2 != parent%2 {
		return 0
	} else {
		return 1 + findIndexFromRight(me>>1, parent>>1)
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
			for i, _ := range n.Neighbors {
				if n.Neighbors[i].TempColor == colorProposal {
					i=0 //reset the progress
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
					i=0 //reset the progress
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