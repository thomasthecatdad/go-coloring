package reductions

import (
	"fmt"
	s "github.com/goombaio/orderedset"
	g "github.com/thomaseb191/go-coloring/graphs"
	"math/rand"
	"sync"
	"time"
)

type messageShared struct {
	degree int
	rndval float32
	avail  *s.OrderedSet
}

var data = make(map[string]messageShared)

func dlfShared(gr g.Graph, poolSize int, debug int) g.Graph {
	var wg sync.WaitGroup
	wg.Add(len(gr.Nodes))

	var checkpoint1 sync.WaitGroup
	var checkpoint2 sync.WaitGroup
	var checkpoint3 sync.WaitGroup

	checkpoint1.Add(len(gr.Nodes))

	var lock sync.Mutex

	for _, node := range gr.Nodes {
		if debug%2 == 1 {
			fmt.Println(node.Name)
		}
		node := node
		go func() {
			vertexShared(node, gr.MaxDegree, &checkpoint1, &checkpoint2, &checkpoint3, &lock, debug)
			wg.Done()
		}()
	}

	wg.Wait()

	return gr
}

func vertexShared(n *g.Node, maxDegree int, checkpoint1 *sync.WaitGroup, checkpoint2 *sync.WaitGroup, checkpoint3 *sync.WaitGroup, lock *sync.Mutex, debug int) {
	rand.Seed(time.Now().UnixNano())
	degree := len(n.Neighbors)

	set := s.NewOrderedSet()

	for i := 0; i <= maxDegree; i++ {
		set.Add(i)
	}

	m := messageShared{
		degree: degree,
		rndval: -1,
		avail:  set,
	}

	lock.Lock()
	data[n.Name] = m
	lock.Unlock()

	iter := 0

	for {
		if debug%2 == 1 {
			fmt.Println(n.Name, "round:", iter)
		}

		temp := rand.Float32()

		lock.Lock()
		m = data[n.Name]
		m.rndval = temp
		data[n.Name] = m
		lock.Unlock()

		selectedColor := m.avail.Values()[0].(int)

		myColor := true
		deg := m.degree
		rnd := m.rndval

		checkpoint2.Add(1)
		checkpoint1.Done()
		checkpoint1.Wait()

		for _, neighbor := range n.Neighbors {
			incmsg := data[neighbor.Name]
			if debug%2 == 1 {
				fmt.Println(n.Name, incmsg)
			}
			if incmsg.degree == -1 {
				continue
			}

			if !myColor {
				continue
			}

			if incmsg.degree < deg {
				continue
			} else if incmsg.degree == deg && incmsg.rndval < rnd {
				continue
			}

			myColor = false
		}

		checkpoint3.Add(1)
		checkpoint2.Done()
		checkpoint2.Wait()

		if myColor {
			n.Color = selectedColor

			m.degree = -1
			lock.Lock()
			for _, neighbor := range n.Neighbors {
				data[neighbor.Name].avail.Remove(selectedColor)
				if debug%2 == 1 {
					fmt.Println(neighbor.Name, data[neighbor.Name].avail.Values())
				}
			}
			data[n.Name] = m
			lock.Unlock()
			break
		}

		checkpoint1.Add(1)
		checkpoint3.Done()
		checkpoint3.Wait()

		iter++
	}

	checkpoint3.Done()

	return

}
