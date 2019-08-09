package main

// # of midpoints to make on adjacent square regions
const MIDPOINT_DIVISOR = 2

type Tuple struct {
	X, Y int
}

type RoutingSquare struct {
	X1, X2, Y1, Y2 int
	Can_cut        bool
	Id_num         int
	Routers        []Tuple
}

type RegionParams struct {
	PointList  [][]bool
	SquareList []RoutingSquare
	PointDict  map[Tuple]bool
	BorderDict map[int][]int
}

func Square_equals(s1, s2 RoutingSquare) bool {
	return (s1.X1 == s2.X1 && s1.X2 == s2.X2 && s1.Y1 == s2.Y1 && s1.Y2 == s2.Y2)
}

func Square_list_remove(s RoutingSquare, r *RegionParams) {
	index := -1
	for i, sq := range r.SquareList {
		if Square_equals(sq, s) {
			index = i
		}
	}
	if index != -1 {
		r.SquareList = r.SquareList[:index+copy(r.SquareList[index:], r.SquareList[index+1:])]
	}
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Side_ratio(sq1, sq2 RoutingSquare) float64 {
	if sq1.Can_cut {
		ret := 1.1
		if sq1.X1 == sq2.X1 || sq1.X2 == sq2.X2 {
			ret = float64(sq1.X2-sq1.X1) / float64(Max((sq2.X2-sq2.X1), 1))
		} else if sq1.Y1 == sq2.Y1 || sq1.Y2 == sq2.Y2 {
			ret = float64(sq1.Y2-sq1.Y1) / float64(Max((sq2.Y2-sq2.Y1), 1))
		}
		if ret > 1 {
			return 0.0
		} else {
			return ret
		}
	} else {
		return 0.0
	}
}

func Area_ratio(sq1, sq2 RoutingSquare) float64 {
	if sq1.Can_cut {
		if sq1.X1 > sq2.X1 && sq1.X2 < sq2.X2 {
			return float64(sq1.X2-sq1.X1) / float64(Max((sq2.X2-sq2.X1), 1))
		} else if sq1.Y1 > sq2.Y1 && sq1.Y2 < sq2.Y2 {
			return float64(sq1.Y2-sq1.Y1) / float64(Max((sq2.Y2-sq2.Y1), 1))
		} else {
			return 0.0
		}
	} else {
		return 0.0
	}
}

func Single_cut(sq1, sq2 RoutingSquare) []RoutingSquare {
	new_squares := make([]RoutingSquare, 0)

	if sq1.X1 == sq2.X1 && sq1.X2 == sq2.X2 {
		if sq1.Y1 < sq2.Y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq2.Y2, true, sq1.Id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq2.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
		}
	} else if sq1.Y1 == sq2.Y1 && sq1.Y2 == sq2.Y2 {
		if sq1.X1 < sq2.X1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
		}
	} else if sq1.X1 == sq2.X1 {
		if sq1.Y1 < sq2.Y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq1.Y1, sq2.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.X2 + 1, sq2.X2, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq2.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.X2 + 1, sq2.X2, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([]Tuple, 0)})
		}
	} else if sq1.X2 == sq2.X2 {
		if sq1.Y1 < sq2.Y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq1.Y1, sq2.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X1 - 1, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq2.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X1 - 1, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([]Tuple, 0)})
		}
	} else if sq1.Y1 == sq2.Y1 {
		if sq1.X1 < sq2.X1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq1.Y2 + 1, sq2.Y2, false, sq2.Id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq1.Y2 + 1, sq2.Y2, false, sq2.Id_num, make([]Tuple, 0)})
		}
	} else if sq1.Y2 == sq2.Y2 {
		if sq1.X1 < sq2.X1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq2.Y1, sq1.Y1 - 1, false, sq2.Id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq2.Y1, sq1.Y1 - 1, false, sq2.Id_num, make([]Tuple, 0)})
		}
	}

	return new_squares
}

func Double_cut(sq1, sq2 RoutingSquare) []RoutingSquare {
	new_squares := make([]RoutingSquare, 0)

	if sq1.X1 > sq2.X1 && sq1.X2 < sq2.X2 {
		if sq1.Y2+1 == sq2.Y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq1.Y1, sq2.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X1 - 1, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.X2 + 1, sq2.X2, sq2.Y1, sq2.Y2, false, -1, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq2.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X1 - 1, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.X2 + 1, sq2.X2, sq2.Y1, sq2.Y2, false, -1, make([]Tuple, 0)})
		}
	} else if sq1.Y1 > sq2.Y1 && sq1.Y2 < sq2.Y2 {
		if sq1.X2+1 == sq2.X1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq2.Y1, sq1.Y1 - 1, false, sq2.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq1.Y2 + 1, sq2.Y2, false, -1, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq2.Y1, sq1.Y1 - 1, false, sq2.Id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq1.Y2 + 1, sq2.Y2, false, -1, make([]Tuple, 0)})
		}
	}

	return new_squares
}

func Rebuild(sq_list []RoutingSquare, r *RegionParams) {

	r.BorderDict = make(map[int][]int, 0)
	for x := 0; x < len(r.SquareList); x++ {
		r.BorderDict[x] = make([]int, 0)
		r.SquareList[x].Routers = make([]Tuple, 0)
	}

	for y := 0; y < len(sq_list); y++ {
		square := sq_list[y]

		for z := y + 1; z < len(sq_list); z++ {
			new_square := sq_list[z]

			if new_square.X1 >= square.X1 && new_square.X2 <= square.X2 {
				if new_square.Y1 == square.Y2+1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cx := float64(new_square.X2 - new_square.X1) * (float64(i) / float64(MIDPOINT_DIVISOR)) + float64(new_square.X1)
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple {int(cx), new_square.Y1})
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple{ int(cx), square.Y2})
					}
				} else if new_square.Y2 == square.Y1-1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cx := float64(new_square.X2 - new_square.X1) * (float64(i) / float64(MIDPOINT_DIVISOR)) + float64(new_square.X1)
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple {int(cx), new_square.Y2})
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple{ int(cx), square.Y1})
					}
				}

			} else if new_square.Y1 >= square.Y1 && new_square.Y2 <= square.Y2 {
				if new_square.X1 == square.X2+1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cy := float64(new_square.Y2 - new_square.Y1) * (float64(i) / float64(MIDPOINT_DIVISOR)) + float64(new_square.Y1)
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple {new_square.X1, int(cy)})
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple{ square.X2, int(cy)})
					}

				} else if new_square.X2 == square.X1-1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cy := float64(new_square.Y2 - new_square.Y1) * (float64(i) / float64(MIDPOINT_DIVISOR)) + float64(new_square.Y1)
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple {new_square.X2, int(cy)})
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple{ square.X1, int(cy)})
					}
				}
			}
			// ?
			if square.X1 >= new_square.X1 && square.X2 <= new_square.X2 {
				if square.Y1 == new_square.Y2+1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cx := float64(square.X2 - square.X1) * (float64(i) / float64(MIDPOINT_DIVISOR)) + float64(square.X1)
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple {int(cx), square.Y1})
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple{int(cx), new_square.Y2})
					}

				} else if square.Y2 == new_square.Y1-1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cx := float64(square.X2 - square.X1) * (float64(i) / float64(MIDPOINT_DIVISOR)) + float64(square.X1)
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple {int(cx), new_square.Y1})
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple{int(cx), square.Y2})
					}
				}


			} else if square.Y1 >= new_square.Y1 && square.Y2 <= new_square.Y2 {
				if square.X1 == new_square.X2+1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cy := float64(square.Y2 - square.Y1) * (float64(i) / float64(MIDPOINT_DIVISOR)) + float64(square.Y1)
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple {square.X1, int(cy)})
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple{new_square.X2, int(cy)})
					}

				} else if square.X2 == new_square.X1-1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cy := float64(square.Y2 - square.Y1) * (float64(i) / float64(MIDPOINT_DIVISOR)) + float64(square.Y1)
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple {new_square.X1, int(cy)})
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple{square.X2, int(cy)})
					}
				}
			}

			if new_square.Y2 == square.Y1 - 1 {
				//Hanging Case 1
				// next_square is on the top of square
				// right side of next_square extends pass right side of square
				if square.X1 <= new_square.X1 && square.X2 < new_square.X2 && square.X2 > new_square.X1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cx := (float64(i) / float64(MIDPOINT_DIVISOR)) * float64(square.X2 - new_square.X1) + float64(new_square.X1)
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple {int(cx), square.Y1})
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple{int(cx), new_square.Y2})
					}
				}

				// Hanging Case 2
				// next_square is on the top of square
				// left side of next square extends pass the left side of square
				if square.X1 >= new_square.X1 && square.X2 > new_square.X2 && square.X1 < new_square.X2 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cx := (float64(i) / float64(MIDPOINT_DIVISOR)) * float64(new_square.X2 - square.X1) + float64(square.X1)
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple {int(cx), square.Y1})
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple{int(cx), new_square.Y2})
					}
				}
			}

			if square.Y2 == new_square.Y1 - 1 {
				// Hanging Case 3
				// square is on top of next_square
				// right side of square extends pass right side of next_square
				if new_square.X1 <= square.X1 && new_square.X2 < square.X2 && new_square.X2 > square.X1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cx := (float64(i) / float64(MIDPOINT_DIVISOR)) * float64(new_square.X2 - square.X1) + float64(square.X1)
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple {int(cx), square.Y2})
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple{int(cx), new_square.Y1})
					}
				}

				// Hanging Case 4
				// square is on top of next_square
				// left side of square extends pass the left side of next_square
				if new_square.X1 >= square.X1 && new_square.X2 > square.X2 && new_square.X1 < square.X2 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cx := (float64(i) / float64(MIDPOINT_DIVISOR)) * float64(square.X2 - new_square.X1) + float64(new_square.X1)
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple {int(cx), square.Y2})
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple{int(cx), new_square.Y1})
					}
				}
			}

			if new_square.X2 == square.X1 - 1 {
				// Hanging Case 5
				// next_square is to the left of square
				// bottom side of square extends pass the bottom side of next square
				if square.Y1 >= new_square.Y1 && square.Y2 > new_square.Y2 && square.Y1 < new_square.Y2 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cy := (float64(i) / float64(MIDPOINT_DIVISOR)) * float64(new_square.Y2 - square.Y1) + float64(square.Y1)
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple {square.X1, int(cy)})
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple{new_square.X2, int(cy)})
					}
				}

				// Hanging Case 6
				// next_square is to the left of square
				// top side of square extends pass the top side of next_square
				if square.Y1 <= new_square.Y1 && square.Y2 < new_square.Y2 && square.Y2 > new_square.Y1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cy := (float64(i) / float64(MIDPOINT_DIVISOR)) * float64(square.Y2 - new_square.Y1) + float64(new_square.Y1)
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple {square.X1, int(cy)})
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple{new_square.X2, int(cy)})
					}
				}
			}

			if square.X2 == new_square.X1 - 1 {
				// Hanging Case 7
				// square is to the left of next_square
				// bottom side of square extends pass bottom side of next_square
				if square.Y1 >= new_square.Y1 && square.Y2 > new_square.Y2 && square.Y1 < new_square.Y2 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cy := (float64(i) / float64(MIDPOINT_DIVISOR)) * float64(new_square.Y2 - square.Y1) + float64(square.Y1)
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple {square.X1, int(cy)})
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple{new_square.X2, int(cy)})
					}
				}

				// Hanging Case 8
				// square is to the left of next_square
				// top side of square extends pass top side of next_square
				if square.Y1 <= new_square.Y1 && square.Y2 < new_square.Y2 && square.Y2 > new_square.Y1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
					for i:=0; i<=MIDPOINT_DIVISOR; i++ {
						cy := (float64(i) / float64(MIDPOINT_DIVISOR)) * float64(square.Y2 - new_square.Y1) + float64(new_square.Y1)
						r.SquareList[y].Routers = append(r.SquareList[y].Routers, Tuple {square.X2, int(cy)})
						r.SquareList[z].Routers = append(r.SquareList[z].Routers, Tuple{new_square.X1, int(cy)})
					}
				}
			}
		}
	}
}
