package reductions

import (
	"fmt"
	g "github.com/thomaseb191/go-coloring/graphs"
	"math"
	"sync"
)

func lineal(gr g.Graph, poolSize int, debug int) g.Graph {
	fmt.Println("Starting Lineal Color Reduction")

	for _, n := range gr.Nodes {
		n.Color = n.Color + 1
	}

	totalColors := len(gr.Nodes)
	iter := 0
	p := gr.MaxDegree
	q := 2*p*p

	g.PrintGraph(&gr)

	fmt.Printf("p = %d, q = %d\n", p, q)

	for totalColors > p*p {
		fmt.Printf("Starting Iteration %d\nTotal colors: %d\n", iter, totalColors)

		numBins := totalColors / q
		if numBins == 0 {
			numBins = 1
		}
		subGraphs := make([][]*g.Node, numBins)
		colormap := make(map[*g.Node]int)

		for _, n := range gr.Nodes {
			j := 1
			if totalColors >= 1 {
				j = int(math.Min(math.Ceil(float64(n.Color)/float64(q)), float64(numBins)))
			}

			j = j - 1
			color := n.Color - j*q
			n.Color = color

			subGraphs[j] = append(subGraphs[j], n)
			colormap[n] = j
		}

		var wg sync.WaitGroup

		wg.Add(len(subGraphs))

		fmt.Println(wg)

		var wg2 sync.WaitGroup
		wg2.Add(len(subGraphs))

		for _, sg := range subGraphs {
			go func() {
				defer wg2.Done()
				procedureRefine(sg, p, &wg)
			}()
		}

		wg2.Wait()

		g.PrintGraph(&gr)

		for key, val := range colormap {
			key.Color = key.Color + val*p*p
		}



		totalColors = int(math.Max(math.Floor(float64(totalColors)/float64(q)),1.0)) * p*p

		fmt.Printf("Total colors2: %d\n", totalColors)

		g.PrintGraph(&gr)

		iter++
	}


	return gr
}

func procedureRefine(gr []*g.Node, p int, wg *sync.WaitGroup) {

	var tempmap = make(map[*g.Node]int)

	for _, n := range gr {
		var lt = make([]int, p+1)
		var gt = make([]int, p+1)

		for _, neighbor := range n.Neighbors {
			if neighbor.Color > p {
				continue
			}
			if neighbor.Color < n.Color {
				temp := lt[neighbor.Color]
				lt[neighbor.Color] = temp + 1
			} else {
				temp := gt[neighbor.Color]
				gt[neighbor.Color] = temp + 1
			}
		}

		var minltuses = math.MaxInt64
		var minltcolor = math.MaxInt64

		if len(lt) == 0 {
			minltcolor = 1
		} else {
			for key, val := range lt {
				if key == 0 {
					continue
				}
				if val < minltuses {
					minltuses = val
					minltcolor = key
				}
			}
		}

		var mingtuses = math.MaxInt64
		var mingtcolor = math.MaxInt64

		if len(gt) == 0 {
			mingtcolor = 1
		} else {
			for key, val := range gt {
				if val < mingtuses {
					mingtuses = val
					mingtcolor = key
				}
			}
		}

		tempmap[n] = (mingtcolor-1)*p + minltcolor
	}

	wg.Done()

	fmt.Println("WG DONE")
	fmt.Println(wg)

	wg.Wait()

	fmt.Println("Done waiting")

	for key, val := range tempmap {
		key.Color = val
	}

}