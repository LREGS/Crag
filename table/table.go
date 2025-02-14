package table

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// these are here for when we can be bothered to add colour to it
var (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

type ForecastTotals struct {
	HighestTemp  float64
	LowestTemp   float64
	AvgTemp      float64
	AvgWindSpeed float64
	TotalPrecip  float64
	Datapoints   int
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

// we want our max width to be larger but we dont want to rely on max width to align the name section with the bottom section.
// because we are padding based on max width when a string is larger than max width the top section isn't adjusted.
// I guess the easiest way to get round it atm is to check this at the end of making the table and then just pad if required.
// or we make the top row last once we know what the longest column is and just make sure that the columns are within the max lenght

const (
	width = 56
)

// do we want forecast data prep and the table builder to be two seperate things? Would simplyfy the struct and mean we only need strings
// to be passed to the table builder and not the whole forecast
type ForecastTable struct {
	forecast map[string]*ForecastTotals
	windows  [][]time.Time
	name     string
	table    strings.Builder
}

func NewForecastTable(name string, forecast map[string]*ForecastTotals, windows [][]time.Time) *ForecastTable {
	return &ForecastTable{
		forecast: forecast,
		windows:  windows,
		name:     name,
		table:    strings.Builder{},
	}
}

const (
	topLeftCorner      = "┌"
	topRightCorner     = "┐"
	bottomLeftCorner   = "└"
	bottomRightCorner  = "┘"
	horizontalLine     = "─"
	verticalLine       = "│"
	topIntersection    = "┬"
	bottomIntersection = "┴"
	leftIntersection   = "├"
	rightIntersection  = "┤"
)

// longest forecast data row will define width
func (t *ForecastTable) nameBox(width int) string {

	// builds the name box in relation to the width of the largest forecast data row

	nameBox := NewLineBuilder(width)
	nameBox.WriteString(topLeftCorner)
	nameBox.WriteString(strings.Repeat("─", (width - 2)))
	nameBox.WriteString(topRightCorner)
	nameBox.WriteRune('\n')
	nameBox.WriteString(verticalLine)
	nameBox.WriteString(t.name)
	nameBox.WriteString(strings.Repeat(" ", ((width - 2) - len(t.name))))
	nameBox.WriteString(verticalLine)
	nameBox.WriteRune('\n')

	return nameBox.String()
}

// func (t *ForecastTable) columnHeaders(columnNames []string)

func (t *ForecastTable) build() {
	// headers - physical definition of a header

	t.table.Reset()

	// add the name box to the table
	t.table.WriteString(t.nameBox(width))

	// this section builds and populates the headline forecast data for each crag over the days
	columns := NewLineBuilder(width)
	colsMap := make(map[string][]string, len(t.forecast))

	//stylistacially id prefer this column data to spawn into the middle of the table somehow and not be left hand aligned
	// as this will make the white space look a bit less jarring

	columnWidths := make([]int, 4) // we want this to be able to grow as columns grow actually but we know how many columns beforehand just no magic number
	for k, v := range t.forecast {
		cols := []string{
			// why dont we have a date for each row here?!
			fmt.Sprintf(" Temp %d/%d/%d ", int(v.HighestTemp), int(v.LowestTemp), int(v.AvgTemp)),
			fmt.Sprintf(" Total Precip %d ", int(v.TotalPrecip)),
			fmt.Sprintf(" Wind %dmp ↓ ", int(v.AvgWindSpeed)),
			fmt.Sprintf(" 1/2/3 "),
		}
		// I think this is ok or should it be outside of the other for loop
		for i, v := range cols {
			if columnWidths[i] < utf8.RuneCountInString(v) {
				columnWidths[i] = utf8.RuneCountInString(v)
			}
		}
		// key is date
		colsMap[k] = cols
	}

	// top column seperator
	topForecastRow := NewLineBuilder(width)
	botForecastRow := NewLineBuilder(width)
	topForecastRow.WriteString(leftIntersection)
	botForecastRow.WriteString(leftIntersection)
	for i, c := range columnWidths {
		if i != 3 {
			topForecastRow.WriteString(strings.Repeat(horizontalLine, (c)))
			topForecastRow.WriteString(topIntersection)
			botForecastRow.WriteString(strings.Repeat(horizontalLine, (c)))
			botForecastRow.WriteString(bottomIntersection)

		} else {
			topForecastRow.WriteString(strings.Repeat(horizontalLine, c))
			topForecastRow.WriteString(rightIntersection)
			botForecastRow.WriteString(strings.Repeat(horizontalLine, c))
			botForecastRow.WriteString(rightIntersection)

		}

	}
	t.table.WriteString(topForecastRow.String())
	t.newline()

	for _, v := range colsMap {
		columns.WriteString(verticalLine)
		for i, c := range v {
			if i != 3 {
				columns.WriteString(c)
				columns.WriteString(strings.Repeat(" ", (columnWidths[i] - utf8.RuneCountInString(c))))
				columns.WriteString(verticalLine)

			} else {
				columns.WriteString(c)
				columns.WriteString(verticalLine)
				if utf8.RuneCountInString(columns.String()) < width {
					columns.WriteString(strings.Repeat(" ", ((width - utf8.RuneCountInString(columns.String())) - 1)))
				}
				// end of table
				columns.WriteString(verticalLine)

				t.table.WriteString(columns.String())
				t.newline()
				t.table.WriteString(botForecastRow.String())
				t.newline()
				columns.Reset()

			}
		}
	}
}

func (t *ForecastTable) newline() {
	t.table.WriteRune('\n')
}

// if we have changed our params re max width do we still need to create a new builder in this way and is creating a new builder per line
// even that good?

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
//   I thiknk the spacing here is maybe more clear visually but im not sure at all now#?!
// 	|06:00       09:00       12:00 - - 1500 - - 1800 - - 2100 - - 2400 - - 0300 - -
// 	|  *   *   *   *   *   *   *
//		Symbols
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
	// why are we passing and copying the whole forecast to this function just to find the longest name?
	// Are we actually needing to find the longest name more than once per function call?
	l := 0
	for k := range totals {
		if len(k) > l {
			l = len(k)
		}
	}
	return l
}

/* dynamic table column withds

- need to move column seperators on top of data row move with the pipes inside the data row
  but not ever column will be the same width every time
	I think we can maybe create the data string and then split it using the pipe and then
	calculating the length of each element after the split that will tell us where the top pipe goes I think?




*/
