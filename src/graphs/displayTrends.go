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
	Names []string
	NumNodes []int
	TimeElapsed []int
	NumberColors []int
	MaxDegree []int //Added by Tyler, no implentation on visualization side yet
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
	degIV := IsDegreeOnlyIV(data)

	var lineGraph *charts.Line
	if degIV {
		lineGraph = generateDegreeLineChart(data)
	} else {
		lineGraph = generateNodeLineChart(data)
	}


	return lineGraph
}

func generateNodeLineChart(data map[int]DataPoint) *charts.Line {
	lineGraph := charts.NewLine()
	categories := make([]*opts.GraphCategory, 0)
	numAlgos := len(data)

	for i := 0; i < numAlgos; i++ {
		categories = append(categories,
			&opts.GraphCategory{
				Name: fmt.Sprintf("%s", algoMap[i]),
				Label: &opts.Label{
					Show:     true,
					Position: "right",
				},
			})
	}
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
		charts.WithLegendOpts(opts.Legend{
			Left: "60%",
			Show: true,
			Data: categories,
		}),
	)

	// In this case we are assuming the Number of Nodes is what changes in each new run of the test,
	// we are assuming MaxDegree is the same for all runs of the tests.
	var algInd int
	for i, v := range data {
		if v.NumNodes != nil {
			algInd = i
			break
		}
	}
	lineGraph.SetXAxis(data[algInd].NumNodes)

	for algoNum, dataPoint := range data {
		lineGraph.AddSeries(algoMap[algoNum], generateLineData(dataPoint.TimeElapsed),
			charts.WithLabelOpts(opts.Label{Show: true, Position: "bottom"}))
	}

	lineGraph.SetSeriesOptions(
		charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
			Name: "Average",
			Type: "average",
		}),

		charts.WithLineChartOpts(opts.LineChart{
			Smooth: false,
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

func generateDegreeLineChart(data map[int]DataPoint) *charts.Line {
	lineGraph := charts.NewLine()
	categories := make([]*opts.GraphCategory, 0)
	numAlgos := len(data)

	for i := 0; i < numAlgos; i++ {
		categories = append(categories,
			&opts.GraphCategory{
				Name: fmt.Sprintf("%s", algoMap[i]),
				Label: &opts.Label{
					Show:     true,
					Position: "right",
				},
			})
	}
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
			Name: "Max Degree",
		}),
		charts.WithLegendOpts(opts.Legend{
			Left: "60%",
			Show: true,
			Data: categories,
		}),
	)

	// In this case we are assuming the Number of Nodes is what changes in each new run of the test,
	// we are assuming MaxDegree is the same for all runs of the tests.
	var algInd int
	for i, v := range data {
		if v.NumNodes != nil {
			algInd = i
			break
		}
	}
	lineGraph.SetXAxis(data[algInd].MaxDegree)

	for algoNum, dataPoint := range data {
		lineGraph.AddSeries(algoMap[algoNum], generateLineData(dataPoint.TimeElapsed),
			charts.WithLabelOpts(opts.Label{Show: true, Position: "bottom"}))
	}

	lineGraph.SetSeriesOptions(
		charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
			Name: "Average",
			Type: "average",
		}),

		charts.WithLineChartOpts(opts.LineChart{
			Smooth: false,
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

func IsDegreeOnlyIV(data map[int]DataPoint) bool {
	numNodes := 0
	for _, alg := range data {
		if alg.NumNodes != nil {
			for _, p := range alg.NumNodes {
				if numNodes != 0 && p != numNodes {
					return false
				}
				numNodes = p
			}
		}
	}
	return true
}

// GenerateHTMLForDataPoints is a
// Method that converts an arbitrary number of dataPoints to HTML visualisations.
func GenerateHTMLForDataPoints(data map[int]DataPoint, testFileName string) {
	fmt.Printf("Generating html...\n")
	page := components.NewPage()
	page.AddCharts(
		generateLineChart(data),
	)
	now := time.Now()
	path := fmt.Sprintf("../html/%s-%d-%d-%d.html", testFileName[0:6], now.Hour(), now.Minute(), now.Second())
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
	fmt.Printf("Done generating html.\t%s\n", path)
}