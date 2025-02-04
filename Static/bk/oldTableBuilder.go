package nothing

import (
	"fmt"
	"strings"
	"time"
)

// theres got to be a better way to share types between these projects

type ForecastTotals struct {
	HighestTemp  float64
	LowestTemp   float64
	AvgTemp      float64
	AvgWindSpeed float64
	TotalPrecip  float64
	Datapoints   int // wtf is this because ive forgot?
}

// A window is a gap in precipitation of more than 1hour and gives the amount of rain the area has just experienced.

type Window struct {
	Time                       []string
	AvgTemp                    float64
	AvgWindSpeed               float64
	WindDirection              int
	PrecipInLastXAmountofHours int
}

// Met office weather model is updated and released every hour for the next 72 hours.

type ForecastPayload struct {
	LastModelRunTime string
	Totals           map[string]*ForecastTotals
	Windows          [][]time.Time // doesnt make sense im guessing they're in order but windows should just be packaged individually with each repotr
	// days since last rain / longest dry period
}

const (
	width = 80
)

type Headers struct {
	headers []string
}

type TableBuilder struct {
	builder strings.Builder
}

func NewTableBuilder() *TableBuilder {
	return &TableBuilder{builder: strings.Builder{}}
}

func (t *TableBuilder) build(name string, totals map[string]*ForecastTotals, Windows [][]time.Time) string {
	// headers - physical definition of a header

	// why arent we getting direction from the forecast totals as well because im sure its part of the dataset?!

	b := t.builder

	colsMap := make(map[string][]string, len(totals))

	for k, v := range totals {
		cols := []string{
			fmt.Sprintf("Temp %d/%d/%d", int(v.HighestTemp), int(v.LowestTemp), int(v.AvgTemp)),
			fmt.Sprintf("Total Precip %d", int(v.TotalPrecip)),
			fmt.Sprintf("Wind %dmp ↓", int(v.AvgWindSpeed)),
			fmt.Sprintf("1/2/3"),
		}
		colsMap[k] = cols
	}

	b.Reset()

	nameSection := NewLineBuilder(width)

	topRow := NewLineBuilder(width)
	topRow.WriteString("┌")
	topRow.WriteString(strings.Repeat("─", width-2))
	topRow.WriteString("┐")

	nameRow := NewLineBuilder(width)

	nameRow.WriteString("│")

	nameSection.WriteString(topRow.String())
	nameSection.WriteRune('\n')
	nameSection.WriteString(nameRow.String())
	nameSection.WriteRune('\n')

	// generate header section based on length of the columns

	// widestCol := map[string]int

	colNumber := 4
	columns := NewLineBuilder(width)
	for _, v := range colsMap {
		columns.WriteRune('│')

		for i, c := range v {
			if i != colNumber-1 {
				columns.WriteString(c)
				columns.WriteString(" │ ")
			} else {
				columns.WriteString(c)
				columns.WriteString(strings.Repeat(" ", width-(len(columns.String())+2)))

				columns.WriteString(" │ ")
				nameSection.WriteString(columns.String())
				nameSection.WriteRune('\n')
				columns.Reset()

			}

		}
	}

	nameSection.WriteString(columns.String())

	// for k, v := range colsMap {
	// 	_ = k
	// 	// we need to evaluate all the keys to determine the order of the columns and

	// }

	// b.WriteString(nameSection.String())
	// b.WriteString(topRow.String())
	// b.WriteString(nameRow.String())
	b.WriteString(nameSection.String())

	return b.String()

}

func NewLineBuilder(size int) *strings.Builder {
	var b strings.Builder
	b.Grow(size)
	return &b
}

// 	+------------------------+ +------------------------+ +------------------------+
// 	|  	 Stanage, Peak District, info etc
// 	+------------------------+ +------------------------+ +------------------------+
// 	|	Tempv +8 | Total Precip | Avg Wind | Problems 6/7/8 colour coded|
// 	|  			lo/hi/avg| 6													 |
// 	+------------------------+ +------------------------+ +------------------------+
// 	|08/11:
// 	|0600 - - 0900 - - 1200 - - 1500 - - 1800 - - 2100 - - 2400 - - 0300 - -
// 	| symbols
// 	|09/11:
// 	|0600 - - 0900 - - 1200 - - 1500 - - 1800 - - 2100 - - 2400 - - 0300 - -
// 	|Symbols
// 	----------------------------------------------------------------------------------
// }

// The above will be a full blcok for each crag, windows will only appear if there is a weather window.
// each block is going to be made from 3 elements:

// header - gives crag info
// weather chunk - display the weather data
// windows chunk - displays the weather windows data
// I will just try and have a fixed width based on the max length of the current longest name in keys
// and maybe at some point centralise it to the screen

func longestName(totals map[string]*ForecastTotals) int {
	l := 0
	for k := range totals {
		if len(k) > l {
			l = len(k)
		}
	}
	return l
}
