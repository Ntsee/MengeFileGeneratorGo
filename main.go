package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Params struct {

	Name			string
	MainXML			string
	MainDirectory	string
	MainImage		string
	WallImage		string
	OutputDirectory	string
	HasView			bool

	Width			int
	Height			int
}

var behaviorFlag = flag.String("behavior", "", "specifies the path to the Behavior PNG file")
var wallFlag = flag.String("wall", "", "specifies the path to the Wall PNG file")
var outputFlag = flag.String("output", "", "specifies the path to the output directory")
var viewFlag = flag.Bool("view", false, "specifies if menge should run the visualizer")

func init() {

	flag.StringVar(behaviorFlag, "b", "", "specifies the path to the Behavior PNG file")
	flag.StringVar(wallFlag, "w", "", "specifies the path to the Wall PNG file")
	flag.StringVar(outputFlag, "o", "", "specifies the path to the output directory")
	flag.BoolVar(viewFlag, "v", false, "specifies if menge should run the visualizer")
}

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Too few arguments: missing path to the Scenario XML File")
		return
	}

	scenarioParams := new(Params)
	scenarioParams.MainXML = args[0]
	scenarioParams.Name = strings.Replace(filepath.Base(args[0]), ".xml", "", 1)


	scenarioParams.MainDirectory = args[0][:strings.LastIndex(args[0], filepath.Base(args[0]))]
	scenarioParams.MainImage = strings.Replace(args[0], ".xml", ".png", 1)
	scenarioParams.WallImage = strings.Replace(args[0], ".xml", "Walls.png", 1)
	scenarioParams.OutputDirectory = scenarioParams.Name
	scenarioParams.HasView = *viewFlag

	if *behaviorFlag != "" {
		scenarioParams.MainImage = *behaviorFlag
	}

	if *wallFlag != "" {
		scenarioParams.WallImage = *wallFlag
	}

	if *outputFlag != "" {
		scenarioParams.OutputDirectory = *outputFlag
	}

	wallRegions := new(RegionParams)
	err := ReadMap(scenarioParams, wallRegions, 0x0000)
	if err != nil {
		fmt.Println(err)
		return
	}

	walkableRegions := new(RegionParams)
	err = ReadMap(scenarioParams, walkableRegions, 0xFFFF)
	if err != nil {
		fmt.Println(err)
		return
	}

	colorDictionary, err := ReadColorDictionary(scenarioParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	scenarioXML, err := ReadScenarioXML(scenarioParams.MainXML)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err = os.Stat(scenarioParams.OutputDirectory); os.IsNotExist(err) {
		os.Mkdir(scenarioParams.OutputDirectory, 0777)
	}

	behaviorXML := CreateBehaviorXML(scenarioParams, scenarioXML, colorDictionary)
	err = WriteBehaviorXML(behaviorXML, fmt.Sprintf("%s/%sB.xml", scenarioParams.OutputDirectory, scenarioParams.Name))
	if err != nil {
		fmt.Println(err)
		return
	}

	sceneXML, err := CreateSceneXML(scenarioParams, wallRegions, scenarioXML)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = WriteSceneXML(sceneXML, fmt.Sprintf("%s/%sS.xml", scenarioParams.OutputDirectory, scenarioParams.Name))
	if err != nil {
		fmt.Println(err)
		return
	}

	vertices, edges := CreateGraph(walkableRegions)
	err = WriteRoadMapTXT(scenarioParams, vertices, edges, fmt.Sprintf("%s/%s.txt", scenarioParams.OutputDirectory, scenarioParams.Name))
	if err != nil {
		fmt.Println(err)
		return
	}

	viewerXML, err := CreateViewXML(scenarioParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	if scenarioParams.HasView {
		err = WriteViewXML(viewerXML, fmt.Sprintf("%s/%sV.xml", scenarioParams.OutputDirectory, scenarioParams.Name))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	projectXML := CreateProjectXML(scenarioParams)
	err = WriteProjectXML(projectXML, fmt.Sprintf("%s/%s.xml", scenarioParams.OutputDirectory, scenarioParams.Name))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Use the following command to run the scenario in menge:")
	fmt.Printf("./menge -d 1000 -p %s/%s.xml", scenarioParams.OutputDirectory, scenarioParams.Name)
}
