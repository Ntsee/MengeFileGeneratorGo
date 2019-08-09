package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

/*

<View width="640" height="480" z_up="1" >
    <Camera xpos="783.639" ypos="112.585" zpos="200.133" xtgt="783.639" ytgt="112.595" ztgt="0.00312764" far="500" near="0.01" fov="0.0" orthoScale="0.118519" />
    <Light x="1" y="0" z="-1" type="directional" space="camera" diffR="1" diffG="0.8" diffB="0.8"/>
    <Light x="-1" y="0" z="-1" type="directional" space="camera" diffR="0.8" diffG="0.8" diffB="1"/>
    <Light x="0" y="0" z="1" type="directional" space="world" diffR="0.8" diffG="0.8" diffB="0.8"/>
</View>
 */

type ViewXML struct {
	XMLName	xml.Name	`xml:"View"`
	Width	int			`xml:"width,attr"`
	Height	int			`xml:"height,attr"`
	ZUp	int			`xml:"z_up,attr"`
	Camera	Camera		`xml:"Camera"`
	Lights	[]Light		`xml:"Light"`
}

type Camera struct {
	XMLName		xml.Name	`xml:"Camera"`
	X			float32		`xml:"xpos,attr"`
	Y			float32		`xml:"ypos,attr"`
	Z			float32		`xml:"zpos,attr"`
	XTarget		float32		`xml:"xtgt,attr"`
	YTarget		float32		`xml:"ytgt,attr"`
	ZTarget		float32		`xml:"ztgt,attr"`
	Far			float32		`xml:"far,attr"`
	Near		float32		`xml:"near,attr"`
	FOV			float32		`xml:"fov,attr"`
	OrthoScale	float32		`xml:"orthoScale,attr"`
}

type Light struct {
	XMLName	xml.Name	`xml:"Light"`
	X		int			`xml:"x,attr"`
	Y		int			`xml:"y,attr"`
	Z		int			`xml:"z,attr"`
	Type	string		`xml:"type,attr"`
	Space	string		`xml:"space,attr"`
	DiffR	float32		`xml:"diffR,attr"`
	DiffG	float32		`xml:"diffG,attr"`
	DiffB	float32		`xml:"diffB,attr"`
}

func ReadViewXML(file string) (ViewXML, error) {

	var viewXML ViewXML
	xmlFile, err := os.Open(file)
	if err != nil {
		return viewXML, err
	}
	defer xmlFile.Close()
	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return viewXML, err
	}

	err = xml.Unmarshal(bytes, &viewXML)
	return viewXML, nil
}

func WriteViewXML(viewXML ViewXML, file string) error {

	bytes, err := xml.MarshalIndent(&viewXML, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

func CreateViewXML(params *Params) (ViewXML, error) {

	viewXML, err := ReadViewXML("base_viewer.xml")
	if err != nil {
		return viewXML, err
	}

	viewXML.Width = params.Width
	viewXML.Height = params.Height
	viewXML.Camera.X = float32(params.Width / 2 - 50)
	viewXML.Camera.XTarget = viewXML.Camera.X
	viewXML.Camera.Y = float32(params.Height / 2 - 50)
	viewXML.Camera.YTarget = viewXML.Camera.Y + .01
	viewXML.Camera.OrthoScale = .50

	return viewXML, nil
}