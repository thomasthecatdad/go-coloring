package reductions

import (
	"fmt"
	s "github.com/goombaio/orderedset"
	g "github.com/thomaseb191/go-coloring/graphs"
	"math/rand"
	"sync"
	"time"
)

type message struct {
	degree   int
	rndvalue float32
	color    int
}

var availableColors = make(map[string]*s.OrderedSet)

func dlf(gr g.Graph, poolSize int, debug int) g.Graph {

	var wg sync.WaitGroup
	wg.Add(len(gr.Nodes))

	incoming := make(map[string][]chan message)
	outgoing := make(map[string][]chan message)

	for _, node := range gr.Nodes {
		for _, neighbor := range node.Neighbors {
			edge := make(chan message, 1)
			outgoing[node.Name] = append(outgoing[node.Name], edge)
			incoming[neighbor.Name] = append(incoming[neighbor.Name], edge)
		}
	}

	var wglist = make([]*sync.WaitGroup, len(gr.Nodes))

	for i := 0; i < len(wglist); i++ {
		var newwg sync.WaitGroup
		wglist[i] = &newwg
	}

	if debug % 2 == 1 {
		println(incoming)
		println(outgoing)
	}

	wglist[0].Add(len(gr.Nodes))
	var lock sync.Mutex
	for _, node := range gr.Nodes {
		if debug % 2 == 1 {
			fmt.Println(node.Name, incoming[node.Name], outgoing[node.Name])
		}
		node := node
		go func() {
			vertex(node, incoming[node.Name], outgoing[node.Name], gr.MaxDegree, wglist, &lock, debug)
			wg.Done()
			for {
				for _, ch := range incoming[node.Name] {
					_, ok := <-ch
					if !ok {
						continue
					}
					//fmt.Println(info)
				}
			}
		}()
	}

	wg.Wait()

	return gr
}

func vertex(n *g.Node, incoming []chan message, outgoing []chan message, maxDegree int, wg []*sync.WaitGroup, lock *sync.Mutex, debug int) {
	rand.Seed(time.Now().UnixNano())
	degree := len(n.Neighbors)

	set := s.NewOrderedSet()

	for i := 0; i <= maxDegree; i++ {
		set.Add(i)
	}

	lock.Lock()
	availableColors[n.Name] = set
	lock.Unlock()

	m := message{
		degree:   degree,
		rndvalue: -1,
		color:    -1,
	}

	iter := 0

	for {
		if debug % 2 == 1 {
			fmt.Println(n.Name, "round: ", iter, len(incoming))
		}
		m.rndvalue = rand.Float32()
		m.color = set.Values()[0].(int)
		//fmt.Println(outgoing)
		for _, ch := range outgoing {
			ch <- m
		}

		deg := m.degree
		rnd := m.rndvalue

		mycolor := true

		var newincoming []chan message

		for _, ch := range incoming {
			incmsg := <-ch

			if incmsg.color == -1 {
				//incoming = append(incoming[:i], incoming[i:]...)
				continue
			}

			newincoming = append(newincoming, ch)

			if incmsg.degree < deg {
				continue
			} else if incmsg.degree == deg && incmsg.rndvalue < rnd {
				continue
			}

			mycolor = false
			deg = incmsg.degree
			rnd = incmsg.rndvalue
		}

		incoming = newincoming

		if mycolor {
			n.Color = m.color
			//fmt.Println("Color", m.color, "selected by node ", n.Name, "in round", iter)
			lock.Lock()
			for _, neighbor := range n.Neighbors {
				availableColors[neighbor.Name].Remove(n.Color)
			}
			lock.Unlock()
			break
		}

		wg[iter+1].Add(1)
		wg[iter].Done()
		wg[iter].Wait()

		iter++

	}

	for _, ch := range outgoing {
		ch <- message{color: -1}
		close(ch)
	}

	//fmt.Println(n.Name, "done")

	wg[iter].Done()

	return
}
