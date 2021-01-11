package main

import (
	"math/rand"
	"sort"
)

type flake struct {
	x            float64
	y            float64
	size         float64
	drawSize     float64
	style        string
	fallTime     float64
	fallDuration float64
	origIndex    int
}

type snowField []flake

var flakes snowField

func (sf snowField) update(delta float64, drawHeight int, speed float64) {
	floatHeight := float64(drawHeight)
	for i := range sf {
		sf[i] = sf[i].update(delta, floatHeight, speed)
	}
}

func (f flake) update(delta float64, drawHeight float64, speed float64) flake {
	f.fallDuration += delta * speed
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
		flakeArray[index].drawSize = speed * 2
		flakeArray[index].origIndex = index
	}

	return flakeArray
}

func adjustFlakes(newCount int, current snowField, drawWidth int, drawHeight int) snowField {
	if newCount == len(current) {
		return current
	}

	sort.Slice(current, func(i, j int) bool {
		return current[i].origIndex < current[j].origIndex
	})

	if newCount > len(current) {
		newFlakesArray := createFlakes(newCount-len(current), drawWidth, drawHeight)
		current = append(current, newFlakesArray...)
		for i := range current {
			current[i].origIndex = i
		}

		sort.Slice(current, func(i, j int) bool {
			return current[i].style < current[j].style
		})

		return current
	}

	current = current[0:newCount]
	sort.Slice(current, func(i, j int) bool {
		return current[i].style < current[j].style
	})

	return current
}
