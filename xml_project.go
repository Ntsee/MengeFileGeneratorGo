package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type ProjectXML struct {
	XMLName		xml.Name	`xml:"Project"`
	Behavior	string		`xml:"behavior,attr"`
	Scene		string		`xml:"scene,attr"`
	View		string		`xml:"view,attr,omitempty"`
	DumpPath	string		`xml:"dumpPath,attr"`
	Model		string		`xml:"model,attr"`
}

func WriteProjectXML(projectXML ProjectXML, file string) error {

	bytes, err := xml.MarshalIndent(&projectXML, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

func CreateProjectXML(scenarioParams *Params) ProjectXML {

	projectXML := ProjectXML{
		Behavior:	fmt.Sprintf("%sB.xml", scenarioParams.Name),
		Scene:		fmt.Sprintf("%sS.xml", scenarioParams.Name),
		DumpPath:	fmt.Sprintf("images/%s", scenarioParams.Name),
		Model:		"orca",
	}

	if scenarioParams.HasView {
		projectXML.View = fmt.Sprintf("%sV.xml", scenarioParams.Name)
	}

	return projectXML
}