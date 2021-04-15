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
	fmt.Printf("Generating Edges... \n")
	edges := make([]opts.GraphLink, 0)
	for _, x := range gr.Nodes {
		for _, neighbor := range x.Neighbors {
			fmt.Printf("Generating Edge between %s and %s\n", x.Name, neighbor.Name)
			edges = append(edges,
				opts.GraphLink {
					Source: x.Name,
					Target: neighbor.Name,
				})
		}
	}
	return edges
}

// Generates all of the nodes in JSON fomrat for the charts API.
// Takes in one of our graphs and converts into one of theirs.
func generateNodes(gr *Graph) []opts.GraphNode {
	fmt.Printf("Generating Nodes... \n")
	nodes := make([]opts.GraphNode, 0)
	for _, x := range gr.Nodes {
		fmt.Printf("Generating Node %s\n", x.Name)
		nodes = append(nodes,
			opts.GraphNode {
				Name: x.Name,
				Category: x.Color,
			})
	}
	return nodes
}

// generate graph is a Generates a graph which can then be converted to HTML.
func generateGraph(gr *Graph) *charts.Graph{
	categories := make([]*opts.GraphCategory, 0)
	numColors := CountColors(gr)
	for i := 0; i < numColors; i++ {
		categories = append(categories,
			&opts.GraphCategory{
				Name: fmt.Sprintf("%d", i),
				Label: &opts.Label {
					Show: true,
					Position: "right",
				},
			})
	}
	graph := charts.NewGraph()
	title := fmt.Sprintf("Graph %s", gr.Name)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
		charts.WithLegendOpts(
			opts.Legend{
				Show: true,
				Data: categories,
			}),
	)
	nodes := generateNodes(gr)
	edges := generateEdges(gr)
	graph.AddSeries("graph", nodes, edges).
		SetSeriesOptions(
			charts.WithGraphChartOpts(
				opts.GraphChart {
					Force: &opts.GraphForce{Repulsion: 8000},
					Layout: "force",
					Categories: categories,
				}),
			charts.WithLabelOpts(
				opts.Label {
					Show: true,
					Position: "right",
				}),
			charts.WithLineStyleOpts(
				opts.LineStyle {
					Curveness: 0.3,
				}),
		)
	return graph
}

// Method that converts an arbitrary number of graphs to HTML visualisations.
func GenerateHTMLForMany(grs []*Graph, testNum int) {
	page := components.NewPage()
	for _, x := range grs {
		page.AddCharts(
			generateGraph(x),
		)
	}
	path := fmt.Sprintf("html/test%dResults.html", testNum)
	f, err := os.Create(path)
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))
}

// Method that converts one graph to HTML visualized.
func GenerateHTMLForOne(gr *Graph) {
	fmt.Printf("Generating HTML... for one graph\n")
	page := components.NewPage()
	page.AddCharts(
		generateGraph(gr),
	)
	path := fmt.Sprintf("html/%s.html", gr.Name)
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
