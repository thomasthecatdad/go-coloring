package graphs

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
	"time"
)

type DataPoint struct {
	NumNodes []int
	TimeElapsed []int
	NumberColors []int
	IsSafe []bool
}

var (
	// algoMap is a map that maps algorithm number to its name.
	algoMap = map[int]string{
		0 : "Naive",
		1 : "KW",
		2 : "Cole-Vishkin",
		3 : "Linial Replacement",
	}
)

// generateLineData is a method that generates data points for the line graph.
func generateLineData(data []int) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i]})
	}
	return items
}

// generateLineChart is a method that generates a new line chart based on the the map of algorithms to DataPoint objects.
// In this case we are generating a line chart that graphs Runtime on the Y axis and NumNodes on the X axis.
func generateLineChart(data map[int]DataPoint) *charts.Line{
	lineGraph := charts.NewLine()
	lineGraph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "RunTime Analysis for the different algorithms.",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Cost time(ns)",
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Number of Nodes",
		}),
	)

	// In this case we are assuming the Number of Nodes is what changes in each new run of the test,
	// we are assuming MaxDegree is the same for all runs of the tests.
	lineGraph.SetXAxis(data[0].NumNodes)

	for algoNum, dataPoint := range data {
		lineGraph.AddSeries(algoMap[algoNum], generateLineData(dataPoint.TimeElapsed),
			charts.WithLabelOpts(opts.Label{Show: true, Position: "bottom"}))
			charts.WithTitleOpts(opts.Title{Title: algoMap[algoNum]})
	}

	lineGraph.SetSeriesOptions(
		charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
			Name: "Average",
			Type: "average",
		}),
		charts.WithLineChartOpts(opts.LineChart{
			Smooth: true,
		}),
		charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
			Label: &opts.Label{
				Show:      true,
				Formatter: "{a}: {b}",
			},
		}),
	)

	return lineGraph
}

// GenerateHTMLForDataPoints is a
// Method that converts an arbitrary number of dataPoints to HTML visualisations.
func GenerateHTMLForDataPoints(data map[int]DataPoint) {
	fmt.Printf("Generating html...\n")
	page := components.NewPage()
	page.AddCharts(
		generateLineChart(data),
	)
	now := time.Now()
	path := fmt.Sprintf("../html/%d-%d-%d.html", now.Second(), now.Minute(), now.Hour())
	f, err := os.Create(path)
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))
	fmt.Printf("Done generating html.\n")
}