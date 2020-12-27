package main

import (
	"math/rand"
	"sort"
)

type flake struct {
	x     int
	y     int
	speed int
	style string
}

type snowField []flake

var flakes snowField

func (sf snowField) update(delta float64, drawHeight int) {
	for i := range sf{
		sf[i] = sf[i].update(delta, drawHeight)
	}
}

func (f flake) update(delta float64, drawHeight int) flake{
	yUpdate := int(float64(f.speed) * (delta / 30))

	if yUpdate < 1 {
		yUpdate = 1
	}
	f.y +=yUpdate
	if f.y > drawHeight {
		f.y -= drawHeight
	}

	return f
}

func createFlakes(flakeCount int, drawWidth int, drawHeight int) snowField {
	flakeArray := make(snowField, flakeCount)

	for index := range flakeArray {
		flakeArray[index].x = rand.Intn(drawWidth)
		flakeArray[index].y = rand.Intn(drawHeight) - drawHeight

		speed := rand.Intn(5) +2
		flakeArray[index].speed = speed

		style := "#666666" // white

		switch speed {
		case 1:
			style = "#00FF00" // green
		case 2:
			style = "#888888" //blue
		case 3:
			style = "#BBBBBB" //yellow
		case 4:
			style = "#DDDDDD" // cyan
		case 5:
			style = "#EEEEEE" // magenta
		case 6:
			style = "#FFFFFF" // red
		}



		flakeArray[index].style = style
	}

	sort.Slice(flakeArray, func(i, j int) bool {
		return flakeArray[i].speed < flakeArray[j].speed
	})


	return flakeArray
}

func (sf snowField) adjustFlakes(newCount int) snowField {
	if newCount == len(sf) {
		return sf
	}

	return createFlakes(newCount, canvasDrawWidth, canvasDrawHeight)
}
