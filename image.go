package main

import (
	"image/png"
	"os"
)

type RGB struct {
	R	int
	G	int
	B	int
}

func ReadColorDictionary(scenarioParams *Params) (map[RGB][]Tuple, error) {

	colorDictionary := make(map[RGB][]Tuple)
	imgFile, err := os.Open(scenarioParams.MainImage)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	imgData, err := png.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	for y:= 0; y<imgData.Bounds().Max.Y; y++ {
		for x:=0; x<imgData.Bounds().Max.X; x++ {
			r, g, b, _ := imgData.At(x, y).RGBA()
			rgb := RGB{int((r >> 8) & r), int((g >> 8) & g), int((b >> 8) & b)}
			colorDictionary[rgb] = append(colorDictionary[rgb], Tuple{x, y})
		}
	}

	return colorDictionary, nil
}

func ReadMap(p *Params, r *RegionParams, wallColor uint32) error {

	imgFile, err := os.Open(p.WallImage)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	imgData, err := png.Decode(imgFile)
	if err != nil {
		return err
	}

	p.Width = imgData.Bounds().Max.X
	p.Height = imgData.Bounds().Max.Y
	r.PointDict = make(map[Tuple]bool)
	r.SquareList = make([]RoutingSquare, 0)
	r.BorderDict = make(map[int][]int)

	for x := 0; x < p.Width; x++ {
		r.PointList = append(r.PointList, make([]bool, p.Height))
	}

	for x := 0; x < p.Width; x++ {
		for y := 0; y < p.Height; y++ {
			pixel := imgData.At(x, y)
			red, _, _, _ := pixel.RGBA()
			if red == wallColor {
				r.PointList[x][y] = true
				r.PointDict[Tuple{x, y}] = true

			} else {
				r.PointDict[Tuple{x, y}] = false
				temp := make([]int, 2)
				temp[0] = x
				temp[1] = y
			}
		}
	}

	GenerateRouting(p, r)
	return nil
}

func GenerateRouting(p *Params, r *RegionParams) {

	id_counter := 0
	done := false
	for !done {
		top_left := Tuple{-1, -1}
		for x, y := 0, 0; y < p.Height; {

			if r.PointList[x][y] {
				top_left = Tuple{x, y}
				break
			}

			if (x + 1) == p.Width {
				y += 1
			}

			x = (x + 1) % p.Width
		}

		if (top_left == Tuple{-1, -1}) {
			done = true
			break
		}

		// Increase X until collide
		temp := Tuple{top_left.X, top_left.Y}
		for r.PointDict[Tuple{temp.X + 1, temp.Y}] {
			temp.X += 1
		}

		collide := false

		// Increase Y until collide
		y_test := Tuple{top_left.X, top_left.Y}

		for !collide {
			y_test.Y += 1
			if y_test.Y >= p.Height {
				collide = true
				break
			}

			for x_val := top_left.X; x_val <= temp.X; x_val++ {
				if !r.PointDict[Tuple{x_val, y_test.Y}] {
					collide = true
				}
			}
		}

		bottom_right := Tuple{temp.X, y_test.Y - 1}
		new_square := RoutingSquare{top_left.X, bottom_right.X, top_left.Y, bottom_right.Y, true, id_counter, make([]Tuple, 0)}
		id_counter++
		RemoveRoutingSquare(new_square, r)
		r.SquareList = append(r.SquareList, new_square)
	}

	length := len(r.SquareList)
	for x := 0; x < length; x++ {
		r.BorderDict[x] = make([]int, 0)
		r.SquareList[x].Routers = make([]Tuple, 0)
	}

	for y, _ := range r.SquareList {
		square := r.SquareList[y]

		for z := y + 1; z < len(r.SquareList); z++ {
			new_square := r.SquareList[z]

			if new_square.X1 >= square.X1 && new_square.X2 <= square.X2 {
				if new_square.Y1 == square.Y2+1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)

				} else if new_square.Y2 == square.Y1-1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}
			} else if new_square.Y1 >= square.Y1 && new_square.Y2 <= square.Y2 {
				if new_square.X1 == square.X2+1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				} else if new_square.X2 == square.X1-1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}
			} else if square.X1 >= new_square.X1 && square.X2 <= new_square.X2 {
				if square.Y1 == new_square.Y2+1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				} else if square.Y2 == new_square.Y1-1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}
			} else if square.Y1 >= new_square.Y1 && square.Y2 <= new_square.Y2 {
				if square.X1 == new_square.X2+1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)

				} else if square.X2 == new_square.X1-1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}
			} else if new_square.Y2 == square.Y1 - 1 {
				//Hanging Case 1
				//next_square is on the top of square
				// right side of next_square extends pass right side of square
				if square.X1 <= new_square.X1 && square.X2 < new_square.X2 && square.X2 > new_square.X1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}

				// Hanging Case 2
				// next_square is on the top of square
				// left side of next square extends pass the left side of square
				if square.X1 >= new_square.X1 && square.X2 > new_square.X2 && square.X1 < new_square.X2 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}
			} else if square.Y2 == new_square.Y1 - 1{
				// Hanging Case 3
				// square is on top of next_square
				// right side of square extends pass right side of next_square
				if new_square.X1 <= square.X1 && new_square.X2 < square.X2 && new_square.X2 > square.X1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}

				// Hanging Case 4
				// square is on top of next_square
				// left side of square extends pass the left side of next_square
				if new_square.X1 >= square.X1 && new_square.X2 > square.X2 && new_square.X1 < square.X2 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}
			} else if new_square.X2 == square.X1 - 1 {
				// Hanging Case 5
				// next_square is to the left of square
				// bottom side of square extends pass the bottom side of next square
				if square.Y1 >= new_square.Y1 && square.Y2 > new_square.Y2 && square.Y1 < new_square.Y2 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}

				// Hanging Case 6
				// next_square is to the left of square
				// top side of square extends pass the top side of next_square
				if square.Y1 <= new_square.Y1 && square.Y2 < new_square.Y2 && square.Y2 > new_square.Y1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}
			} else if square.X2 == new_square.X1 - 1 {
				// Hanging Case 7
				// square is to the left of next_square
				// bottom side of square extends pass bottom side of next_square
				if square.Y1 >= new_square.Y1 && square.Y2 > new_square.Y2 && square.Y1 < new_square.Y2 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}

				// Hanging Case 8
				// square is to the left of next_square
				// top side of square extends pass top side of next_square
				if square.Y1 <= new_square.Y1 && square.Y2 < new_square.Y2 && square.Y2 > new_square.Y1 {
					r.BorderDict[y] = append(r.BorderDict[y], z)
					r.BorderDict[z] = append(r.BorderDict[z], y)
				}
			}
		}
	}


	//fmt.Println(r.BorderDict)

	//Cutting takes place in this loop
	for true {
		rebuilt := false

		for i := 0; i < len(r.SquareList) && !rebuilt; i++ {

			for _, n := range r.BorderDict[i] {

				s_rat := Side_ratio(r.SquareList[i], r.SquareList[n])
				if s_rat > 0.6 {
					new_squares := Single_cut(r.SquareList[i], r.SquareList[n])

					s1 := r.SquareList[n]
					s2 := r.SquareList[i]

					Square_list_remove(s1, r)
					Square_list_remove(s2, r)

					r.SquareList = append(r.SquareList, new_squares...)

					Rebuild(r.SquareList, r)

					rebuilt = true

					break
				}

				a_rat := Area_ratio(r.SquareList[i], r.SquareList[n])
				if a_rat > 0.6 {
					new_squares := Double_cut(r.SquareList[i], r.SquareList[n])

					s1 := r.SquareList[n]
					s2 := r.SquareList[i]

					Square_list_remove(s1, r)
					Square_list_remove(s2, r)

					new_squares[2].Id_num = len(r.SquareList)

					r.SquareList = append(r.SquareList, new_squares...)

					Rebuild(r.SquareList, r)
					rebuilt = true
					break
				}
			}
		}

		if !rebuilt {
			Rebuild(r.SquareList, r)
			break
		}
	}
}

func RemoveRoutingSquare(sq RoutingSquare, r *RegionParams) {

	for i := sq.X1; i <= sq.X2; i++ {
		for j := sq.Y1; j <= sq.Y2; j++ {
			r.PointDict[Tuple{i, j}] = false
			r.PointList[i][j] = false
		}
	}
}
