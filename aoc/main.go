package main

import (
	"fmt"
	"github.com/pezza/advent-of-code/puzzles"
	"github.com/pezza/advent-of-wasm/wasm"
	"strings"
	"syscall/js"
	"time"
)

func getStarData() string {
	var sb strings.Builder
	for day := 1; day <= 25; day ++ {
		sb.WriteString("<tr>")
		sb.WriteString(fmt.Sprintf("<td><span class=\"day\">%d<span></td>",day))
		for year := 2015; year <= 2020; year++ {
			sb.WriteString("<td>")
			_, err := puzzles.GetPuzzle(day, year)

			if err == nil {
				sb.WriteString(fmt.Sprintf("<span class=\"star completed\">*</span><span class=\"star completed\">*</span>&nbsp<button onClick=\"runPuzzle(%d,%d);\">Run</button>", day, year))
			}else{
				sb.WriteString("<span class=\"star\">*</span><span class=\"star\">*</span>")
			}
			sb.WriteString("</td>")
		}
		sb.WriteString("</tr>")
	}
	return sb.String()
}

func runPuzzle(this js.Value, p []js.Value) interface{} {
	day, year := p[0].Int() ,p[1].Int()

	puzzle, err := puzzles.GetPuzzle(day, year)

	if err != nil {
		fmt.Println(err)
	}

	year, day, name := puzzle.Describe()
	fmt.Printf("Running puzzle: -- %d: Day %d - %v\n", year, day, name)
	start := time.Now()
	res := puzzle.PartOne(puzzle.PuzzleInput(), nil)
	fmt.Printf("Part 1 => %s. Complete in %d milliseconds\n", res ,time.Since(start).Milliseconds())

	start = time.Now()
	res = puzzle.PartTwo(puzzle.PuzzleInput(), nil)
	fmt.Printf("Part 2 => %s. Complete in %d milliseconds\n", res ,time.Since(start).Milliseconds())

	return js.ValueOf(0)
}

func main (){
	done := make(chan bool, 0)
	doc := wasm.NewJsDoc()
	doc.Document.Set("runPuzzle", js.FuncOf(runPuzzle))
	doc.SetInnerHTMLById("starBody", getStarData())
	<-done
}
