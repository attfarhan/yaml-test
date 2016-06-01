package main

import (
	"fmt"
	// "log"
	// "github.com/shurcooL/go-goon"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/attfarhan/yaml"
	"os"

	original "github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"path/filepath"
	"sourcegraph.com/sourcegraph/srclib/graph"
	"sourcegraph.com/sourcegraph/srclib/unit"
)

var data = `language: go

go:
    - 1.4
    - 1.5
    - 1.6
    - tip

go_import_path: gopkg.in/yaml.v2`

var data2 = `language: node_js
node_js:
  - 0.8
  - 0.10
  - 0.11`

type T struct {
	value  []string
	line   []int
	column []int
}

func main() {
	t := &T{}
	f := original.Unmarshal([]byte(data2), t)
	fmt.Println("ORIGINAL :", f)
	var x []*yaml.Node
	p := yaml.NewParser([]byte(data))
	node := p.Parse()
	tokenList := yaml.Explore(node, x)
	getLineAndColumn(tokenList, data, t)
	for i, _ := range t.value {
		start, end, value := findOffsets(data, t.line[i], t.column[i], t.value[i])
		fmt.Println("column: ", t.column[i], "line: ", t.line[i], "start: ", start, "end: ", end, "value: ", value)
	}
	Execute()
}

func getLineAndColumn(tokenList []*yaml.Node, fileString string, out *T) {
	for _, token := range tokenList {
		out.value = append(out.value, token.Value)
		out.line = append(out.line, token.Line)
		out.column = append(out.column, token.Column)
		// a, b := findOffsets(data, token.Line, token.Column, token.Value)
		fmt.Println("token: ", token, "line: ", token.Line, "column: ", token.Column)
		getLineAndColumn(token.Children, data, out)
	}
}
func findOffsets(fileText string, line, column int, token string) (start, end int, value string) {

	// we count our current line and column position.
	currentCol := 0
	currentLine := 0
	for offset, ch := range fileText {
		// see if we found where we wanted to go to.
		if currentLine == line && currentCol == column {
			end = offset + len([]byte(token))
			return offset, end, token
		}

		// line break - increment the line counter and reset the column.
		if ch == '\n' {
			currentLine++
			currentCol = 0
		} else {
			currentCol++
		}
	}
	return -1, -1, token // not found.
}

func Execute() error {
	inputBytes := []byte(`{"Name":"yaml-test","Type":"yaml","Files":[".travis.yml"]}`)
	var units unit.SourceUnits
	json.NewDecoder(bytes.NewReader(inputBytes)).Decode(&units)
	var u *unit.SourceUnit
	json.NewDecoder(bytes.NewReader(inputBytes)).Decode(&u)
	units = unit.SourceUnits{u}

	os.Stdin.Close()
	if len(units) == 0 {
		log.Fatal("input contains no source unit data.")
	}
	out, err := Graph(units)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		return err
	}
	return nil
}

func Graph(units unit.SourceUnits) (*graph.Output, error) {
	if len(units) > 1 {
		return nil, errors.New("unexpected multiple units")
	}
	u := units[0]
	fmt.Println(units)

	// out is a graph.Output struct with a Ref field of pointers to graph.Ref
	out := graph.Output{Refs: []*graph.Ref{}}
	fmt.Println("Graph5")
	fmt.Println("Out:", out)

	// Decode source unit
	// Get files
	// Iterate over files, parse YAML
	// For each token, get the byte ranges, token string, and add to Refs

	for _, currentFile := range u.Files {

		f, err := ioutil.ReadFile(currentFile)
		if err != nil {
			log.Printf("failed to read a source unit file: %s", err)
			continue
		}
		file := string(f)
		t := &T{}
		var x []*yaml.Node
		p := yaml.NewParser([]byte(file))
		node := p.Parse()
		tokenList := yaml.Explore(node, x)
		getLineAndColumn(tokenList, file, t)
		for i, _ := range t.value {
			fmt.Println("iterating thru values")
			start, end, value := findOffsets(file, t.line[i], t.column[i], t.value[i])
			out.Refs = append(out.Refs, &graph.Ref{
				DefUnitType: "URL",
				DefUnit:     "Circle",
				DefPath:     filepath.ToSlash(currentFile) + string(start),
				Unit:        value,
				File:        filepath.ToSlash(currentFile),
				Start:       uint32(start),
				End:         uint32(end),
			})
		}
	}
	return &out, nil
}

// refs = append(refs, &graph.Ref{
// 	DefUnitType: "URL",
// 	DefUnit:     "MDN",
// 	DefPath:     mdnDefPath(d.Property),
// 	Unit:        u.Name,
// 	File:        filepath.ToSlash(filePath),
// 	Start:       uint32(s),
// 	End:         uint32(e),
