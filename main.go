package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FIXatdl struct {
	XMLName    xml.Name   `xml:"fixatdl"`
	Strategies []Strategy `xml:"Strategies>Strategy"`
}

type Strategy struct {
	Name           string         `xml:"name,attr"`
	StrategyLayout StrategyLayout `xml:"StrategyLayout"`
}

type StrategyLayout struct {
	Panel Panel `xml:"Panel"`
}

type Panel struct {
	Orientation string    `xml:"orientation,attr"`
	Elements    []Element `xml:",any"`
}

type Element struct {
	XMLName   xml.Name
	Field     string     `xml:"field,attr"`
	Label     string     `xml:"label,attr"`
	Text      string     `xml:"text,attr"`
	InitValue string     `xml:"initValue,attr"`
	MinValue  string     `xml:"minValue,attr"`
	MaxValue  string     `xml:"maxValue,attr"`
	Step      string     `xml:"step,attr"`
	EnumPairs []EnumPair `xml:"EnumPair"`
}

type EnumPair struct {
	WireValue   string `xml:"wireValue,attr"`
	EnumID      string `xml:"enumID,attr"`
	DisplayName string `xml:"displayName,attr"`
}

func main() {
	xmlFile, err := os.ReadFile("sample.xml")
	if err != nil {
		log.Fatal(err)
	}

	var fixatdl FIXatdl
	err = xml.Unmarshal(xmlFile, &fixatdl)
	if err != nil {
		log.Fatal(err)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Algo Panel")

	//takes first strat now
	strategy := fixatdl.Strategies[0]

	form := &widget.Form{
		Items: []*widget.FormItem{},
	}

	for _, element := range strategy.StrategyLayout.Panel.Elements {
		switch element.XMLName.Local {
		case "Label":
			label := widget.NewLabel(element.Text)
			form.Append("", label)
		case "DropDownList", "EditableDropDownList":
			options := make([]string, len(element.EnumPairs))
			for i, pair := range element.EnumPairs {
				options[i] = pair.DisplayName
			}
			dropdown := widget.NewSelect(options, func(value string) {
				log.Printf("%s selected: %s", element.Label, value)
			})
			form.Append(element.Label, dropdown)
		case "Spinner", "SingleSpinner":
			entry := widget.NewEntry()
			entry.SetText(element.InitValue)
			form.Append(element.Label, entry)
		case "CheckBox":
			check := widget.NewCheck(element.Label, func(value bool) {
				log.Printf("%s checked: %v", element.Label, value)
			})
			form.Append("", check)
		default:
			log.Printf("Unsupported element type: %s", element.XMLName.Local)
		}
	}

	submitButton := widget.NewButton("Submit", func() {
		fmt.Println("Order submitted")
	})

	content := container.NewVBox(
		form,
		submitButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 400))
	myWindow.ShowAndRun()
}
