package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type BehaviorXML struct {
	XMLName		xml.Name		`xml:"BFSM"`
	GoalSets	[]BGoalSet		`xml:"GoalSet,omitempty"`
	States		[]BState		`xml:"State,omitempty"`
	Transitions	[]BTransition	`xml:"Transition,omitempty"`
}

type BGoalSet struct {
	XMLName		xml.Name	`xml:"GoalSet"`
	ID			int			`xml:"id,attr"`
	Goals		[]Goal		`xml:"Goal,omitempty"`
}

type Goal struct {
	XMLName		xml.Name	`xml:"Goal"`
	ID			int			`xml:"id,attr"`
	Type		string		`xml:"type,attr"`
	Capacity	int			`xml:"capacity,attr"`
	Weight		float32		`xml:"weight,attr"`
	Radius		float32		`xml:"radius,attr"`
	X			int			`xml:"x,attr"`
	Y			int			`xml:"y,attr"`
}

type BState struct {
	XMLName			xml.Name		`xml:"State"`
	Name			string			`xml:"name,attr"`
	Weight			float32			`xml:"weight,attr"`
	Final			int				`xml:"final,attr"`
	GoalSelector	*GoalSelector	`xml:"GoalSelector,omitempty"`
	VelComponent	*VelComponent	`xml:"VelComponent,omitempty"`
	Action			*Action			`xml:"Action,omitempty"`
}

type GoalSelector struct {
	XMLName		xml.Name	`xml:"GoalSelector"`
	Persistent	int			`xml:"persistent,attr,omitempty"`
	Type		string		`xml:"type,attr"`
	GoalSet		int			`xml:"goal_set,attr"`
}

type VelComponent struct {
	XMLName		xml.Name	`xml:"VelComponent,omitempty"`
	Type		string		`xml:"type,attr,omitempty"`
	FileName	string		`xml:"file_name,attr,omitempty"`
}

type Action struct {
	XMLName 		xml.Name	`xml:"Action"`
	Type			string		`xml:"type,attr"`
	X				int			`xml:"x_value,attr"`
	Y				int			`xml:"y_value,attr"`
	ExitReset		bool		`xml:"exit_reset,attr"`
	Distribution	string		`xml:"dist,attr"`
}

type BTransition struct {
	XMLName		xml.Name	`xml:"Transition"`
	From		string		`xml:"from,attr,omitempty"`
	To			string		`xml:"to,attr,omitempty"`
	Condition	*Condition	`xml:"Condition,omitempty"`
	Target		*Target		`xml:"Target,omitempty"`
}

type Condition struct {
	XMLName			xml.Name	`xml:"Condition,atrr,omitempty"`
	Type			string		`xml:"type,attr,omitempty"`
	Distance		float32		`xml:"distance,attr,omitempty"`
	PerAgent		int			`xml:"per_agent,attr"`
	Distribution	string		`xml:"dist,attr,omitempty"`
	Min				int			`xml:"min,attr"`
	Max				int			`xml:"max,attr"`
}

type Target struct {
	XMLName		xml.Name	`xml:"Target"`
	Type		string		`xml:"type,attr,omitempty"`
	States		[]BState	`xml:"State,omitempty"`
}

func ReadBehaviorXML(path string) (BehaviorXML, error) {

	var behaviorXML BehaviorXML
	xmlFile, err := os.Open(path)
	if err != nil {
		return behaviorXML, err
	}
	defer xmlFile.Close()
	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return behaviorXML, err
	}

	err = xml.Unmarshal(bytes, &behaviorXML)
	return behaviorXML, nil
}

func WriteBehaviorXML(behaviorXML BehaviorXML, file string) error {

	bytes, err := xml.MarshalIndent(&behaviorXML, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

func createGoalSets(params *Params, behaviorXML *BehaviorXML, scenarioXML ScenarioXML, colorDictionary map[RGB][]Tuple) {

	totalGoalSets := 0
	for groupID := range scenarioXML.Groups {
		group := scenarioXML.Groups[groupID]

		for goalSetID := range group.GoalSets {
			goalSet := group.GoalSets[goalSetID]
			goalColor := RGB{goalSet.Color.R, goalSet.Color.G, goalSet.Color.B}

			bGoalSet := BGoalSet{}
			bGoalSet.ID = goalSetID + totalGoalSets
			for goalID := range colorDictionary[goalColor] {
				bGoalSet.Goals = append(bGoalSet.Goals, Goal {
					ID: goalID,
					Type: "circle",
					Weight: 1,
					Radius: 0.5,
					Capacity: goalSet.Capacity,
					X: colorDictionary[goalColor][goalID].X,
					Y: int(math.Max(0, float64(params.Height - colorDictionary[goalColor][goalID].Y - 1))),
				})
			}
			behaviorXML.GoalSets = append(behaviorXML.GoalSets, bGoalSet)

		}

		totalGoalSets += len(group.GoalSets)
	}

}

func createStates(params *Params, behaviorXML *BehaviorXML, scenarioXML ScenarioXML, colorDictionary map[RGB][]Tuple) {

	totalGoalSets := 0
	for groupID := range scenarioXML.Groups {
		group := scenarioXML.Groups[groupID]
		behaviorXML.States = append(behaviorXML.States, BState{
			Name: fmt.Sprintf("Start_%d", groupID),
			Final: 0,
			GoalSelector: &GoalSelector{ Type: "identity", Persistent: 0},
			VelComponent: &VelComponent{Type: "zero"},
		})

		behaviorXML.States = append(behaviorXML.States, BState{
			Name: fmt.Sprintf("Start_Wait_%d", groupID),
			Final: 0,
			GoalSelector: &GoalSelector{ Type: "identity", Persistent: 0},
			VelComponent: &VelComponent{Type: "zero"},
		})

		spawnColor := RGB{group.Spawn.Color.R, group.Spawn.Color.G, group.Spawn.Color.B}
		for spawnID := range colorDictionary[spawnColor] {
			spawn := colorDictionary[spawnColor][spawnID]
			behaviorXML.States = append(behaviorXML.States, BState{
				Name:			fmt.Sprintf("Spawn_%d_%d", groupID, spawnID),
				Final:			0,
				GoalSelector:	&GoalSelector{Type: "identity"},
				VelComponent:	&VelComponent{Type: "zero"},
				Action:			&Action{
					Type: "teleport",
					ExitReset: true,
					Distribution: "c",
					X: spawn.X,
					Y: Max(0, params.Height - spawn.Y - 1),
				},
			})
		}

		for goalSetID := range group.GoalSets {
			behaviorXML.States = append(behaviorXML.States, BState {
				Name: fmt.Sprintf("Travel_%d_%d", groupID, goalSetID),
				GoalSelector: &GoalSelector {
					Type: "random",
					Persistent: 1,
					GoalSet: goalSetID + totalGoalSets,
				},
				VelComponent: &VelComponent {
					Type: "road_map",
					FileName: fmt.Sprintf("%s.txt", params.Name),
				},
			})

			behaviorXML.States = append(behaviorXML.States, BState {
				Name: fmt.Sprintf("Arrive_%d_%d", groupID, goalSetID),
				GoalSelector: &GoalSelector {
					Type: "identity",
					Persistent: 0,
				},
				VelComponent: &VelComponent {
					Type: "goal",
				},
			})

			behaviorXML.States = append(behaviorXML.States, BState {
				Name: fmt.Sprintf("Wait_%d_%d", groupID, goalSetID),
				GoalSelector: &GoalSelector {
					Type: "identity",
					Persistent: 0,
				},
				VelComponent: &VelComponent {
					Type: "goal",
				},
			})
		}

		totalGoalSets += len(group.GoalSets)
	}
}

func createTransitions(params *Params, behaviorXML *BehaviorXML, scenarioXML ScenarioXML, colorDictionary map[RGB][]Tuple) {

	for groupID := range  scenarioXML.Groups {
		group := scenarioXML.Groups[groupID]

		// Start -> Wait
		behaviorXML.Transitions = append(behaviorXML.Transitions, BTransition{
			From:		fmt.Sprintf("Start_%d", groupID),
			To:			fmt.Sprintf("Start_Wait_%d", groupID),
			Condition:	&Condition {
				Type: "timer",
				Distribution: "u",
				PerAgent: 1,
				Min: group.Spawn.Min,
				Max: group.Spawn.Max,
			},
		})

		// Wait -> Spawn
		waitSpawnTransition := BTransition{
			From:		fmt.Sprintf("Start_Wait_%d", groupID),
			Condition:	&Condition{Type: "auto"},
			Target:		&Target{Type: "prob"},
		}
		spawnColor := RGB{group.Spawn.Color.R, group.Spawn.Color.G, group.Spawn.Color.B}
		for spawnID := range colorDictionary[spawnColor] {
			waitSpawnTransition.Target.States = append(waitSpawnTransition.Target.States, BState{
				Name:	fmt.Sprintf("Spawn_%d_%d", groupID, spawnID),
				Weight: 1,
			})
		}
		behaviorXML.Transitions = append(behaviorXML.Transitions, waitSpawnTransition)

		// Spawn -> Travel
		spawnTravelTransition := BTransition{
			Condition:	&Condition{Type: "auto"},
			Target:		&Target{ Type: "prob"},
		}
		for i := range waitSpawnTransition.Target.States {
			spawnTravelTransition.From += fmt.Sprintf("%s,", waitSpawnTransition.Target.States[i].Name)
		}
		spawnTravelTransition.From = strings.TrimSuffix(spawnTravelTransition.From, ",")
		for transitionID := range group.Spawn.Transitions {
			transition := group.Spawn.Transitions[transitionID]
			spawnTravelTransition.Target.States = append(spawnTravelTransition.Target.States, BState{
				Name:	fmt.Sprintf("Travel_%d_%d", groupID, transition.To),
				Weight:	transition.Chance,
			})
		}
		behaviorXML.Transitions = append(behaviorXML.Transitions, spawnTravelTransition)

		for goalSetID := range group.GoalSets {
			goalSet := group.GoalSets[goalSetID]

			// Travel -> Arrive
			behaviorXML.Transitions = append(behaviorXML.Transitions, BTransition{
				From:	fmt.Sprintf("Travel_%d_%d", groupID, goalSetID),
				To:		fmt.Sprintf("Arrive_%d_%d", groupID, goalSetID),
				Condition:	&Condition{
					Type: "goal_reached",
					Distance: 1,
				},
			})

			// Arrive -> Wait
			behaviorXML.Transitions = append(behaviorXML.Transitions, BTransition{
				From:	fmt.Sprintf("Arrive_%d_%d", groupID, goalSetID),
				To:		fmt.Sprintf("Wait_%d_%d", groupID, goalSetID),
				Condition:	&Condition{
					Type: "timer",
					PerAgent: 1,
					Distribution: "u",
					Min: goalSet.Min,
					Max: goalSet.Max,
				},
			})

			// Wait -> Next Goal
			waitNextGoalTransition := BTransition{
				From:	fmt.Sprintf("Wait_%d_%d", groupID, goalSetID),
				Condition:	&Condition{Type: "auto"},
				Target:		&Target{Type: "prob"},
			}
			for i := range goalSet.Transitions {
				transition := goalSet.Transitions[i]
				if goalSetID == transition.To {
					waitNextGoalTransition.Target.States = append(waitNextGoalTransition.Target.States, BState{
						Name:	fmt.Sprintf("Arrive_%d_%d", groupID, goalSetID),
						Weight:	transition.Chance,
					})
				} else {
					waitNextGoalTransition.Target.States = append(waitNextGoalTransition.Target.States, BState{
						Name:	fmt.Sprintf("Travel_%d_%d", groupID, goalSetID),
						Weight:	transition.Chance,
					})
				}
			}
			behaviorXML.Transitions = append(behaviorXML.Transitions, waitNextGoalTransition)
		}
	}
}

func CreateBehaviorXML(params *Params, scenarioXML ScenarioXML, colorDictionary map[RGB][]Tuple) BehaviorXML {

	behaviorXML := BehaviorXML{}
	createGoalSets(params, &behaviorXML, scenarioXML, colorDictionary)
	createStates(params, &behaviorXML, scenarioXML, colorDictionary)
	createTransitions(params, &behaviorXML, scenarioXML, colorDictionary)
	return behaviorXML
}