package main

import (
	"fmt"
	// "log"
	"github.com/go-yaml/yaml"
	// "github.com/shurcooL/go-goon"
)

var data = `language: go

go:
    - 1.4
    - 1.5
    - 1.6
    - tip

go_import_path: gopkg.in/yaml.v2`

func main() {
	// t := T{}
	var x []*yaml.Node
	p := yaml.NewParser([]byte(data))
	node := p.Parse()
	tokenList := yaml.Explore(node, x)
	getLineAndColumn(tokenList, data)
}

func getLineAndColumn(tokenList []*yaml.Node, fileString string) {
	for _, token := range tokenList {
		fmt.Println("line: ", token.Line)
		fmt.Println("column: ", token.Column)
		fmt.Println("value:", token.Value)
		fmt.Println(findOffsets(data, token.Line, token.Column, token.Value))
		getLineAndColumn(token.Children, data)
	}
}
func findOffsets(fileText string, line, column int, token string) (start, end int) {

	// we count our current line and column position.
	currentCol := 0
	currentLine := 0
	for offset, ch := range fileText {
		// see if we found where we wanted to go to.
		if currentLine == line && currentCol == column {
			end = offset + len([]byte(token))
			return offset, end
		}

		// line break - increment the line counter and reset the column.
		if ch == '\n' {
			currentLine++
			currentCol = 0
		} else {
			currentCol++
		}
	}
	return -1, -1 // not found.
}

		refs = append(refs, &graph.Ref{
			DefUnitType: "URL",
			DefUnit:     "MDN",
			DefPath:     mdnDefPath(d.Property),
			Unit:        u.Name,
			File:        filepath.ToSlash(filePath),
			Start:       uint32(s),
			End:         uint32(e),


