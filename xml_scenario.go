package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

/* Struct Definitions for Scenario XML */
type ScenarioXML struct {
	XMLName		xml.Name	`xml:"Data"`
	Groups		[]Group		`xml:"Group"`
}

type Group struct {
	XMLName		xml.Name	`xml:"Group"`
	Spawn		Spawn		`xml:"Spawn"`
	GoalSets	[]GoalSet	`xml:"GoalSet"`
	Speed		int			`xml:"speed,attr"`
	Amount		int			`xml:"amount,attr"`
}

type Spawn struct {
	XMLName		xml.Name		`xml:"Spawn"`
	Min			int				`xml:"min,attr"`
	Max			int				`xml:"max,attr"`
	Color		*Color			`xml:"Color"`
	Transitions	[]Transition	`xml:"Transition"`
}

type GoalSet struct {
	XMLName 	xml.Name		`xml:"GoalSet"`
	Capacity 	int				`xml:"capacity,attr"`
	Min 		int 			`xml:"min,attr"`
	Max 		int 			`xml:"max,attr"`
	Color 		*Color 			`xml:"Color"`
	Transitions	[]Transition	`xml:"Transition"`

}

type Color struct {
	XMLName xml.Name `xml:"Color"`
	R 		int  	 `xml:"r,attr"`
	G 		int  	 `xml:"g,attr"`
	B 		int  	 `xml:"b,attr"`
}

type Transition struct {
	XMLName xml.Name	`xml:"Transition"`
	To		int			`xml:"to,attr"`
	Chance 	float32		`xml:"chance,attr"`
}

func ReadScenarioXML(file string) (ScenarioXML, error) {

	var scenarioXML ScenarioXML
	xmlFile, err := os.Open(file)
	if err != nil {
		return scenarioXML, err
	}
	defer xmlFile.Close()
	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return scenarioXML, err
	}

	err = xml.Unmarshal(bytes, &scenarioXML)
	return scenarioXML, nil
}

func WriteScenarioXML(scenarioXML ScenarioXML, file string) error {

	bytes, err := xml.MarshalIndent(&scenarioXML, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}