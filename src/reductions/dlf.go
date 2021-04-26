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

func dlf(gr g.Graph, poolSize int, debug int) g.Graph {

	var wg sync.WaitGroup
	wg.Add(len(gr.Nodes))

	incoming := make(map[string][]chan message)
	outgoing := make(map[string][]chan message)
	wglist := []*sync.WaitGroup{}

	for _, node := range gr.Nodes {
		for _, neighbor := range node.Neighbors {
			edge := make(chan message, 1)
			outgoing[node.Name] = append(outgoing[node.Name], edge)
			incoming[neighbor.Name] = append(incoming[neighbor.Name], edge)
		}
		var newwg sync.WaitGroup
		wglist = append(wglist, &newwg)
	}

	println(incoming)
	println(outgoing)

	wglist[0].Add(len(gr.Nodes))
	for _, node := range gr.Nodes {
		fmt.Println(node.Name, incoming[node.Name], outgoing[node.Name])
		node := node
		go func() {
			vertex(node, incoming[node.Name], outgoing[node.Name], gr.MaxDegree, wglist)
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

func vertex(n *g.Node, incoming []chan message, outgoing []chan message, maxDegree int, wg []*sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())
	degree := len(n.Neighbors)

	set := s.NewOrderedSet()

	for i := 0; i <= maxDegree; i++ {
		set.Add(i)
	}

	m := message{
		degree:   degree,
		rndvalue: -1,
		color:    -1,
	}

	iter := 0

	for {
		fmt.Println(n.Name, "round: ", iter)
		m.rndvalue = rand.Float32()
		m.color = set.Values()[0].(int)
		//fmt.Println(outgoing)
		for _, ch := range outgoing {
			ch <- m
		}

		deg := -1
		rnd := float32(-1.0)

		mycolor := true

		var newincoming []chan message

		for _, ch := range incoming {
			incmsg := <-ch

			if incmsg.color == -1 {
				//incoming = append(incoming[:i], incoming[i:]...)
				continue
			}

			newincoming = append(newincoming, ch)

			if incmsg.degree < degree {
				continue
			} else if incmsg.degree == degree && incmsg.rndvalue < m.rndvalue {
				continue
			} else {
				if incmsg.degree < deg {
					continue
				} else if incmsg.degree == deg && incmsg.rndvalue < rnd {
					continue
				}
			}
			mycolor = false
			deg = incmsg.degree
			rnd = incmsg.rndvalue
		}

		incoming = newincoming

		if mycolor {
			n.Color = m.color
			//fmt.Println("Color", m.color, "selected by node ", n.Name, "in round", iter)
			break
		} else {
			set.Remove(m.color)
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
