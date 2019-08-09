package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edges map[Tuple]bool
type Vertices map [Tuple]int

func isAdjacent(p1 Tuple, p2 Tuple) bool {

	if p1.X + 1 == p2.X || p1.X - 1 == p2.X {
		return p1.Y == p2.Y
	} else if p1.Y + 1 == p2.Y || p1.Y - 1 ==p2.Y {
		return p1.X == p2.X
	} else {
		return false
	}
}

func addIfAbsent(edge Tuple, edges Edges) {
	vid, other_vid := edge.X, edge.Y
	if vid == other_vid {
		return
	}
	if _, exists := edges[Tuple{vid, other_vid}]; exists {
		return
	}
	if _, exists := edges[Tuple{other_vid, vid}]; exists {
		return
	}

	edges[edge] = true
}

func graphSquares(squares []RoutingSquare, vertices Vertices, edges Edges) {
	for square_index := range squares {
		s := squares[square_index]

		x1, x2 := s.X1, s.X2
		y1, y2 := s.Y1, s.Y2
		corners := []Tuple {
			{x1, y2},
			{x1, y1},
			{x2, y1},
			{x2, y2},
		}

		// Add new vertices
		for i := range corners {
			corner := corners[i]
			if _, exists := vertices[corner]; !exists {
				vertices[corner] = len(vertices)
			}
		}

		// Add routers
		for i := range s.Routers {
			router := s.Routers[i]
			if _, exists := vertices[router]; !exists {
				vertices[router] = len(vertices)
			}
		}

		// Add edges to the corners
		for i := range corners {
			vid := vertices[corners[i]]
			for j:= i + 1; j< len(corners); j++ {
				next_vid := vertices[corners[j]]
				addIfAbsent(Tuple{vid, next_vid}, edges)
			}
		}

		// Add edges connecting routers to corners
		for i := range s.Routers {
			vid := vertices[s.Routers[i]]
			for j := range corners {
				corner_vid := vertices[corners[j]]
				addIfAbsent(Tuple{vid, corner_vid}, edges)
			}
		}

		// Add edges connecting internal routers
		if len(s.Routers) > 1 {
			for i := range s.Routers {
				vid := vertices[s.Routers[i]]
				for j := i + 1; j< len(s.Routers); j++ {
					other_vid := vertices[s.Routers[j]]
					addIfAbsent(Tuple{vid, other_vid}, edges)
				}
			}
		}
	}
}

func connectGraphSquares(squares []RoutingSquare, vertices Vertices, edges Edges, borderDict map[int][]int) {

	// Test
	for i := range squares {
		square := squares[i]
		fmt.Printf("Square ID: %d\nTop Left: (%d, %d) Bottom Right: (%d, %d)\n", square.Id_num, square.X1, square.Y1, square.X2, square.Y2)
		for j := range square.Routers {
			router := square.Routers[j]
			fmt.Printf("- Router: (%d, %d)\n", router.X, router.Y)
		}
		fmt.Println()
	}
	// End Test

	for sid := range borderDict {
		square := squares[sid]

		for i := range borderDict[sid] {
			adjacent_square := squares[borderDict[sid][i]]
			for j := range square.Routers {
				router := square.Routers[j]
				vid := vertices[router]
				for k := range adjacent_square.Routers {
					adjacent_router := adjacent_square.Routers[k]
					other_vid := vertices[adjacent_router]
					if isAdjacent(router, adjacent_router) {
						addIfAbsent(Tuple{vid, other_vid}, edges)
					}
				}
			}
		}
	}
}

func CreateGraph(regionParams *RegionParams) (Vertices, Edges) {

	vertices := Vertices{}
	edges := Edges{}

	graphSquares(regionParams.SquareList, vertices, edges)
	connectGraphSquares(regionParams.SquareList, vertices, edges, regionParams.BorderDict)
	return vertices, edges
}

func WriteRoadMapTXT(params *Params, vertices Vertices, edges Edges, file string) error {

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	// Vertices Length
	_, err = w.WriteString(fmt.Sprintf("%d\n", len(vertices)))
	if err != nil {
		return err
	}

	// Output vertices sorted by ID
	vertexIDMap := make(map[int]Tuple, len(vertices))
	for location, id := range vertices {
		vertexIDMap[id] = location
	}

	for vid := 0; vid<len(vertices); vid++ {
		degree := 0
		for edge := range edges {
			if vid == edge.X || vid == edge.Y {
				degree += 1
			}
		}

		location := vertexIDMap[vid]
		_, err = w.WriteString(fmt.Sprintf("%d %d %d\n", degree, location.X, params.Height - location.Y))
		if err != nil {
			return err
		}
	}

	// Edges Length
	_, err = w.WriteString(fmt.Sprintf("%d\n", len(edges)))
	if err != nil {
		return err
	}

	// Output edges
	for edge := range edges {
		_, err = w.WriteString(fmt.Sprintf("%d %d\n", edge.X, edge.Y))
		if err != nil {
			return err
		}
	}

	w.Flush()
	return nil

}