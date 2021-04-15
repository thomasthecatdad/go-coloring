package graphs

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
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
func generateGraph(gr *Graph) *charts.Graph{
	graph := charts.NewGraph()
	title := fmt.Sprintf("Graph %s", gr.Name)
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

// Method that converts an arbitrary number of graphs to HTML visualisations.
func generateHTML(grs []*Graph, testNum int) {
	page := components.NewPage()
	for _, x := range grs {
		page.AddCharts(
			generateGraph(x),
		)
	}
	path := fmt.Sprintf("res/html/test%dResults.html", testNum)
	f, err := os.Create(path)
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))
}

// Method that converts one graph to HTML visualized.
func generateHTMLForOne(gr *Graph) {
	page := components.NewPage()
	page.AddCharts(
		generateGraph(gr),
	)
	path := fmt.Sprintf("res/html/%s", gr.Name)
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

}
