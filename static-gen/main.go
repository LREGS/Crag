package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	f, err := os.Create("./static/hello.html")
	if err != nil {
		panic(err)
	}

	err = Example(`<pre>
    ┌───────┬───────┬─────────────┬────────────┬─────────┬──────────────┬───────────────┬────────────────────┐
    │ Crag  │ Date  │ HighestTemp │ LowestTemp │ AvgTemp │ AvgWindSpeed │ WindDirection │ TotalPrecipitation │
    ├───────┼───────┼─────────────┼────────────┼─────────┼──────────────┼───────────────┼────────────────────┤
    │ name  │       │             │            │         │              │               │                    │
    └───────┴───────┴─────────────┴────────────┴─────────┴──────────────┴───────────────┴────────────────────┘
    
</pre>`).Render(context.Background(), f)
	if err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}

	fs := http.FileServer((http.Dir("./static")))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":3030", nil))

}

var header = []string{
	"+-----------+",
	"|   CragName|",
}

var firstRow = []string{
	"+-------+-------+---+---+-------+---+---+",
	"|       |       |   |   |   |   |   |   |",
}

var middleRow = []string{
	"+---------------------------------------+",
	"|       |       |   |   |   |   |   |   |",
}

var bottomRow = []string{
	"|       |       |   |   |   |   |   |   |",
	"+-------+-------+---+---+---+---+---+---+",
}

// const cragTable = `

//                             +-----------+
//                             |   CragName|
// +-------+-------+---+---+-------+---+---+
// |       |       |   |   |   |   |   |   |
// +---------------------------------------+
// |       |       |   |   |   |   |   |   |
// |       |       |   |   |   |   |   |   |
// +---------------------------------------+
// |       |       |   |   |   |   |   |   |
// |       |       |   |   |   |   |   |   |
// +---------------------------------------+
// |       |       |   |   |   |   |   |   |
// |       |       |   |   |   |   |   |   |
// +-------+-------+---+---+---+---+---+---+

// `
