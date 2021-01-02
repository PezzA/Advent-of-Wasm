package main

import (
	"math/rand"
	"sort"
)

type flake struct {
	x            float64
	y            float64
	size         float64
	style        string
	fallTime     float64
	fallDuration float64
}

type snowField []flake

var flakes snowField

func (sf snowField) update(delta float64, drawHeight int) {

	floatHeight := float64(drawHeight)
	for i := range sf {
		sf[i] = sf[i].update(delta, floatHeight)
	}
}

func (f flake) update(delta float64, drawHeight float64) flake {
	f.fallDuration += delta * 2
	f.y = f.fallDuration / f.fallTime
	if f.y > drawHeight {
		f.y -= drawHeight
		f.fallDuration = 0
	}

	return f
}

func createFlakes(flakeCount int, drawWidth int, drawHeight int) snowField {
	flakeArray := make(snowField, flakeCount)

	for index := range flakeArray {
		flakeArray[index].x = float64(rand.Intn(drawWidth))

		speed := float64(rand.Intn(5) + 2)
		fallSpeed := float64(0)
		flakeArray[index].size = speed

		style := "#666666"

		switch speed {
		case 1:
			style = "#00FF00"
		case 2:
			fallSpeed = 12 + rand.Float64()*(2)
			style = "#888888"
		case 3:
			fallSpeed = 10 + rand.Float64()*(2)
			style = "#BBBBBB"
		case 4:
			fallSpeed = 8 + rand.Float64()*(2)
			style = "#DDDDDD"
		case 5:
			fallSpeed = 6 + rand.Float64()*(2)
			style = "#EEEEEE"
		case 6:
			fallSpeed = 4 + rand.Float64()*(2)
			style = "#FFFFFF"
		}
		flakeArray[index].fallTime = (fallSpeed * 1000) / float64(drawHeight)
		flakeArray[index].fallDuration = rand.Float64() * fallSpeed * float64(drawHeight)
		flakeArray[index].style = style
	}

	sort.Slice(flakeArray, func(i, j int) bool {
		return flakeArray[i].size < flakeArray[j].size
	})

	return flakeArray
}

func (sf snowField) adjustFlakes(newCount int, drawWidth int, drawHeight int) snowField {
	if newCount == len(sf) {
		return sf
	}

	return createFlakes(newCount, drawWidth, drawHeight)
}
