package jd

import (
	"bytes"
	"encoding/json"

	"github.com/fatih/color"
)

type DiffElement struct {
	Path      []JsonNode
	OldValues []JsonNode
	NewValues []JsonNode
}

func (d DiffElement) Render() string {
	greenColor := color.New(color.FgGreen).Add(color.Bold)
	b := bytes.NewBuffer(nil)
	b.WriteString("@ ")
	b.Write([]byte(jsonArray(d.Path).Json()))
	b.WriteString("\n")
	for _, oldValue := range d.OldValues {
		if !isVoid(oldValue) {
			oldValueJson, err := json.Marshal(oldValue)
			if err != nil {
				panic(err)
			}
			b.WriteString("- ")
			b.Write(oldValueJson)
			b.WriteString("\n")
		}
	}
	for _, newValue := range d.NewValues {
		if !isVoid(newValue) {
			newValueJson, err := json.Marshal(newValue)
			if err != nil {
				panic(err)
			}
			greenColor.Println(
				b.WriteString("+ "),
			)
			//b.WriteString("+ ")
			greenColor.Println(
				b.WriteString(newValueJson),
			)
			//			b.Write(newValueJson)
			b.WriteString("\n")
		}
	}
	return b.String()
}

type Diff []DiffElement

func (d Diff) Render() string {
	b := bytes.NewBuffer(nil)
	for _, element := range d {
		b.WriteString(element.Render())
	}
	return b.String()
}
