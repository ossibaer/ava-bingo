package main

import (
	"math"
	"math/rand"
	"strings"
)

type Cell struct {
	value   string
	crossed bool
}

type Bingo struct {
	Board [5][5]Cell
}

func (b Bingo) String() string {
	maxCellWidth := 16
	rows := make([]string, 7)
	for i, row := range b.Board {
		upperHalf := make([]string, 7)
		lowerHalf := make([]string, 7)
		for x := 0; x < 5; x++ {
			item := row[x].value
			var upper, lower string
			if len(item) > maxCellWidth {
				splitIndex := maxCellWidth
				for ; item[splitIndex] != ' '; splitIndex-- {
				}
				upper = item[:splitIndex]
				lower = item[splitIndex+1:]
			} else {
				upper = item
			}
			fill := " "
			if row[x].crossed {
				upper, lower = "", ""
				fill = "#"
			}
			upperWidth, lowerWidth := len(upper), len(lower)
			uL := (maxCellWidth - upperWidth) / 2
			uR := maxCellWidth - upperWidth - uL
			lL := (maxCellWidth - lowerWidth) / 2
			lR := maxCellWidth - lowerWidth - lL
			upperHalf[x+1] = strings.Repeat(fill, uL) + upper + strings.Repeat(fill, uR)
			lowerHalf[x+1] = strings.Repeat(fill, lL) + lower + strings.Repeat(fill, lR)

		}
		rows[i+1] = strings.Join(upperHalf, " | ") + "\n" + strings.Join(lowerHalf, " | ")
	}
	rowSep := "|" + strings.Repeat(strings.Repeat("-", maxCellWidth+2)+"|", 5)
	return strings.Join(rows, "\n "+rowSep+" \n")
}

func (b *Bingo) cross(v string) int {
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			if strings.EqualFold(b.Board[x][y].value, v) {
				b.Board[x][y].crossed = true
				return b.isBingo(x, y)
			}
		}
	}
	return 0
}

func (b *Bingo) isBingo(x, y int) int {
	bingoCount := 0
	rowBingo, columnBingo, leftUpRightDownBingo, leftDownRightUpBingo := true, true, true, true
	for i := 0; i < 5; i++ {
		if !b.Board[x][i].crossed {
			rowBingo = false
		}
		if !b.Board[i][y].crossed {
			columnBingo = false
		}
		if !b.Board[i][i].crossed {
			leftUpRightDownBingo = false
		}
		if !b.Board[i][4-i].crossed {
			leftDownRightUpBingo = false
		}
	}
	if rowBingo {
		bingoCount++
	}
	if columnBingo {
		bingoCount++
	}
	if x == y || int(math.Abs(float64(x-4))) == y {
		if leftUpRightDownBingo {
			bingoCount++
		}
		if leftDownRightUpBingo {
			bingoCount++
		}
	}
	return bingoCount
}

func NewBingo() Bingo {
	possibleFields := []string{
		"Final",
		"Basi",
		"KC",
		"Dancing",
		"Piniata",
		"Robo Captain",
		"Leash",
		"Contact",
		"Dive",
		"DPS hat Aggro",
		"Full Wipe",
		"Fail Air-Compress",
		"Green Chest",
		"Blue Chest",
		"Purple Chest",
		"Gold Chest",
		"Zweite Ebene",
		"Jemand pullt ausversehen",
		"DPS stirbt beim Pre-Attacken",
		"Frost friert Falschen ein",
		"Silence zu late/misst",
		"Core-Role zieht Aggro beim Pull",
		"Falsche Swaps beim Boss",
		"Swaps vergessen",
		"GA nicht up",
	}

	bingo := Bingo{}

	for i := 25; i > 0; i-- {
		index := rand.Intn(i)
		bingo.Board[(i-1)/5][(i-1)%5].value = possibleFields[index]
		possibleFields = append(possibleFields[:index], possibleFields[index+1:]...)
	}
	return bingo
}
