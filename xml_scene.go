package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

type SceneXML struct {
	XMLName 		xml.Name		`xml:"Experiment"`
	Version 		float32			`xml:"version,attr,omitempty"`
	SpatialQuery	*SpatialQuery	`xml:"SpatialQuery,omitempty"`
	Common			*Common			`xml:"Common,omitempty"`
	GCF				*GCF				`xml:"GCF,omitempty"`
	Helbing			*Helbing			`xml:"Helbing,omitempty"`
	Dummy			*Dummy			`xml:"Dummy,omitempty"`
	AgentProfiles	[]AgentProfile 	`xml:"AgentProfile,omitempty"`
	AgentGroups		[]AgentGroup	`xml:"AgentGroup,omitempty"`
	ObstacleSet		*ObstacleSet		`xml:"ObstacleSet,omitempty"`
}

type SpatialQuery struct {
	XMLName			xml.Name	`xml:"SpatialQuery"`
	TestVisibility	bool		`xml:"test_visibility,attr"`
	Type			string		`xml:"type,attr,omitempty"`
}

type Common struct {
	XMLName				xml.Name	`xml:"Common"`
	TimeStep			int			`xml:"time_step,attr,omitempty"`
	Class 				int			`xml:"class,attr,omitempty"`
	PreferredSpeed		int			`xml:"pref_speed,attr,omitempty"`
	MaxAcceleration		int			`xml:"max_accel,attr,omitempty"`
	MaxAngleVelocity	int			`xml:"max_angle_vel,attr,omitempty"`
	MaxNeighbors		int			`xml:"max_neighbors,attr,omitempty"`
	NeighborDistance	int			`xml:"neighbor_dist,attr,omitempty"`
	ObstacleSet			int			`xml:"obstacleSet,attr,omitempty"`
	Priority			float32		`xml:"priority,attr,omitempty"`
	Radius				float32		`xml:"r,attr,omitempty"`
}

type PedVO struct {
	XMLName		xml.Name	`xml:"PedVO"`
	Buffer		float32		`xml:"buffer,attr,omitempty"`
	Factor		float32		`xml:"factor,attr,omitempty"`
	TAU			int			`xml:"tau,attr,omitempty"`
	TAUObstacle	int			`xml:"tauObst,attr,omitempty"`
	TurningBias	float32		`xml:"turningBias,attr,omitempty"`
}

type GCF struct {
	XMLName				xml.Name	`xml:"GCF"`
	AgentForceStrength	float32		`xml:"agent_force_strength,attr"`
	AgentInterpWidth	float32		`xml:"agent_interp_width,attr"`
	MaxAgentDistance	int			`xml:"max_agent_dist,attr"`
	MaxAgentForce		int			`xml:"max_agent_force,attr"`
	ReactionTime		float32		`xml:"reaction_time,attr"`
	MoveScale			float32		`xml:"move_scale,attr"`
	SlowWidth			float32		`xml:"slow_width,attr"`
	StandDepth			float32		`xml:"stand_depth,attr"`
	SwayChange			float32		`xml:"sway_change,attr"`
}

type Helbing struct {
	XMLName			xml.Name	`xml:"Helbing"`
	AgentScale		int			`xml:"agent_scale,attr"`
	BodyForce		int			`xml:"body_force,attr"`
	Friction		int			`xml:"friction,attr"`
	ObstacleScale	int			`xml:"obstacle_scale,attr"`
	ForceDistance	float32		`xml:"force_distance,attr"`
	ReactionTime	float32		`xml:"reaction_time,attr"`
}

type Orca struct {
	XMLName	xml.Name`xml:"ORCA"`
	TAU			int	`xml:"tau,attr"`
	TAUObstacle	int	`xml:"tauObst,attr"`
}

type Dummy struct {
	XMLName	xml.Name`xml:"Dummy"`
	StdDev	float32	`xml:"stddev,attr"`
}

type AgentProfile struct {
	XMLName		xml.Name	`xml:"AgentProfile"`
	Name		string		`xml:"name,attr,omitempty"`
	Inherits	string		`xml:"inherits,attr,omitempty"`
	Common		*Common		`xml:"Common,omitempty"`
	PedVO		*PedVO		`xml:"PedVO,omitempty"`
	Helbing		*Helbing	`xml:"Helbing,omitempty"`
	Orca		*Orca		`xml:"ORCA,omitempty"`
	GCF			*GCF		`xml:"GCF,omitempty"`
}

type AgentGroup struct {
	XMLName			xml.Name		`xml:"AgentGroup"`
	ProfileSelector *ProfileSelector `xml:"ProfileSelector"`
	StateSelector	*StateSelector	`xml:"StateSelector"`
	Generator		*Generator		`xml:"Generator"`
}

type ProfileSelector struct {
	XMLName	xml.Name	`xml:"ProfileSelector"`
	Name	string		`xml:"name,attr,omitempty"`
	Type	string		`xml:"type,attr,omitempty"`
}

type StateSelector struct {
	XMLName	xml.Name	`xml:"StateSelector"`
	Name	string		`xml:"name,attr,omitempty"`
	Type	string		`xml:"type,attr,omitempty"`
}

type Generator struct {
	XMLName		xml.Name	`xml:"Generator"`
	AnchorX		int			`xml:"anchor_x,attr"`
	AnchorY		int			`xml:"anchor_y,attr"`
	CountX		int			`xml:"count_x,attr"`
	CountY		int			`xml:"count_y,attr"`
	OffsetX		int			`xml:"offset_x,attr"`
	OffsetY		int			`xml:"offset_y,attr"`
	Rotation	int			`xml:"rotation,attr"`
	Type		string		`xml:"type,attr"`
}

type ObstacleSet struct {
	XMLName		xml.Name 	`xml:"ObstacleSet"`
	Type		string		`xml:"type,attr"`
	Class		int			`xml:"class,attr"`
	Obstacles	[]Obstacle	`xml:"Obstacle"`
}

type Obstacle struct {
	XMLName		xml.Name	`xml:"Obstacle"`
	Closed		int			`xml:"closed,attr"`
	Vertices	[]Vertex	`xml:"Vertex"`
}

type Vertex struct {
	XMLName	xml.Name	`xml:"Vertex"`
	X		int			`xml:"p_x,attr"`
	Y		int			`xml:"p_y,attr"`
}

func ReadSceneXML(file string) (SceneXML, error) {

	var sceneXML SceneXML
	xmlFile, err := os.Open(file)
	if err != nil {
		return sceneXML, err
	}
	defer xmlFile.Close()
	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return sceneXML, err
	}

	err = xml.Unmarshal(bytes, &sceneXML)
	return sceneXML, nil
}

func WriteSceneXML(sceneXML SceneXML, file string) error {

	bytes, err := xml.MarshalIndent(&sceneXML, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

func CreateSceneXML(params *Params, regionParams *RegionParams, scenarioXML ScenarioXML) (SceneXML, error) {

	sceneXML, err := ReadSceneXML("base_scene.xml")
	if err != nil {
		return sceneXML, err
	}

	// Agent Profiles and Groups
	for groupID := range scenarioXML.Groups {
		group := scenarioXML.Groups[groupID]

		sceneXML.AgentProfiles = append(sceneXML.AgentProfiles, AgentProfile{
			Name:	fmt.Sprintf("group_%d", groupID),
			Inherits:	"base",
			Common:	&Common{
				Class: groupID,
				PreferredSpeed: group.Speed,
			},
		})

		sceneXML.AgentGroups = append(sceneXML.AgentGroups, AgentGroup{
			ProfileSelector:	&ProfileSelector{
				Name:	fmt.Sprintf("group_%d", groupID),
				Type:	"const",
			},
			StateSelector: &StateSelector{
				Name:	fmt.Sprintf("Start_%d", groupID),
				Type:	"const",
			},
			Generator:	&Generator{
				Type: "rect_grid",
				AnchorX: -10,
				AnchorY: -10,
				CountX: group.Amount,
				CountY: 1,
				OffsetX: 1,
				OffsetY: 0,
				Rotation: 0,
			},
		})
	}

	// ObstacleSet
	obstacleSet := ObstacleSet{
		Type:	"explicit",
		Class:	1,
	}

	border := Obstacle{
		Closed: 1,
		Vertices: []Vertex {
			{X: 0, 				Y: params.Height},
			{X: params.Width,	Y: params.Height},
			{X: params.Width,	Y: 0},
			{X: 0,				Y: 0},
		},
	}
	obstacleSet.Obstacles = append(obstacleSet.Obstacles, border)

	for i := range regionParams.SquareList {
		square := regionParams.SquareList[i]
		obstacleSet.Obstacles = append(obstacleSet.Obstacles, Obstacle{
			Closed: 1,
			Vertices: []Vertex {
				{X: square.X1, Y: int(math.Max(0, float64(params.Height - square.Y2 - 1)))},
				{X: square.X2, Y: int(math.Max(0, float64(params.Height - square.Y2 - 1)))},
				{X: square.X2, Y: int(math.Max(0, float64(params.Height - square.Y1 - 1)))},
				{X: square.X1, Y: int(math.Max(0, float64(params.Height - square.Y1 - 1)))},
			},
		})
	}

	sceneXML.ObstacleSet = &obstacleSet
	return sceneXML, nil
}