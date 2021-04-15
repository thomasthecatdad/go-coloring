package graphs

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Generates all of the edges in JSON format for the charts API.
// Takes in one of our graphs and converts into one of theirs.
func generateEdges(gr *Graph) []opts.GraphLink {
	// TODO : Generate all of the edges into JSON for the API.
	links := make([]opts.GraphLink, 0)

	return links
}

// Generates all of the nodes in JSON fomrat for the charts API.
// Takes in one of our graphs and converts into one of theirs.
func generateNodes(gr *Graph) []opts.GraphNode {
	// TODO : Generate all of the nodes into JSON for the API.
	return nil
}

// Generates a graph which can then be converted to HTML.
func generateGraph(gr *Graph, graphNum int) *charts.Graph{
	graph := charts.NewGraph()
	title := fmt.Sprintf("Graph %[1]", graphNum)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
	)
	graph.AddSeries("graph", generateNodes(gr), generateEdges(gr)).
		SetSeriesOptions(
			charts.WithGraphChartOpts(
				opts.GraphChart{
					Force: &opts.GraphForce{Repulsion: 8000},
					Layout: "force",
				}),
		)
	return graph
}

// TODO: Write method to convert an arbitrary number of graphs to HTML.