package reductions

import (
	g "../graphs"
	"fmt"
)

// RunNaive
/* Implements Naive Color Reduction Alg found here:
https://stanford.edu/~rezab/classes/cme323/S16/projects_reports/bae.pdf
MaxDegree 4 means there are 5 colors nodes can be colored as [0,1,2,3,4]
 */
func RunNaive(gr g.Graph, poolSize int, debug int) g.Graph {
	if debug % 2 == 1 {
		fmt.Printf("Starting reduction for %s algorithm...\n", "Naive")
	}

	size := len(gr.Nodes)

	for i := gr.MaxDegree+1; i < size; i++ {
		color := MinColor(*gr.Nodes[i], gr.MaxDegree)
		if color == -1 {
			fmt.Errorf("MinColor() did not return a valid value")
		}
		gr.Nodes[i].Color = color
	}

	return gr
}

func MinColor(n g.Node, poolSize int) int {
	for i := 0; i < poolSize; i++ {
		contained := false
		for _, neighbor := range n.Neighbors {
			if neighbor.Color == i {
				contained = true
				break
			}
		}
		if contained {
			continue
		}
		return i
	}

	return -1
}
