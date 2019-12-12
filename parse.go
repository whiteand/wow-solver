package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
)

func getGrid(img image.Image, step int, startX int, startY int) [][]int {
	bounds := img.Bounds()

	res := make([][]int, 0)

	for x := startX; x < bounds.Max.X; x += step {
		newRow := make([]int, 0)
		for y := startY; y < bounds.Max.Y; y += step {
			r, g, b, _ := img.At(x, y).RGBA()
			newRow = append(newRow, int((r>>8)+(g>>8)+(b>>8))/3)
		}
		res = append(res, newRow)
	}

	return res
}

func getK(numbers []int, k int) int {
	if len(numbers) == 1 {
		return numbers[0]
	}
	less := make([]int, 0)
	equal := make([]int, 0)
	greater := make([]int, 0)

	pivot := numbers[len(numbers)>>1]

	for _, v := range numbers {
		if v < pivot {
			less = append(less, v)
			continue
		}
		if v > pivot {
			greater = append(greater, v)
			continue
		}
		equal = append(equal, v)
	}

	if k < len(less) {
		return getK(less, k)
	}

	if k < len(less)+len(equal) {
		return pivot
	}

	return getK(greater, k-len(less)-len(equal))
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isTooSmall(img image.Image, step int) bool {
	bounds := img.Bounds()

	if bounds.Max.X/step > 50 {
		return true
	}
	if bounds.Max.X/step < 2 {
		return false
	}

	grid := getGrid(img, step, step>>1, step>>1)
	sameColors := 0
	for i := 0; i < len(grid)-1; i++ {
		for j := 0; j < len(grid[i])-1; j++ {
			color := grid[i][j]
			rightColor := grid[i][j+1]
			bottomColor := grid[i+1][j]
			rightBottomColor := grid[i+1][j+1]

			dr := Abs(rightColor - color)
			db := Abs(bottomColor - color)
			drb := Abs(rightBottomColor - color)

			if dr < 5 && db < 5 && drb < 5 {
				sameColors++
				if sameColors > 44 {
					return true
				}
			}

		}
	}

	return false
}

func findStep(img image.Image) int {
	bounds := img.Bounds()
	left := 10
	right := bounds.Max.X >> 1

	for left+1 < right {
		middle := (right + left) >> 1
		if isTooSmall(img, middle) {
			left = middle
		} else {
			right = middle
		}
	}

	return right
}

func countSameColors(img image.Image, step int, shiftX int, shiftY int) (int, int) {
	colorsCount := make(map[int]int)

	grid := getGrid(img, step, shiftX, shiftY)

	for _, row := range grid {
		for _, value := range row {
			count, has := colorsCount[value>>1]

			if !has {
				colorsCount[value>>1] = 1
			} else {
				colorsCount[value>>1] = count + 1
			}
		}
	}

	maxCount := 0
	maxColor := 0

	for color, count := range colorsCount {
		fmt.Println("color", color, "count", count)
		if count > maxCount {
			maxCount = count
			maxColor = color
		}
	}

	return maxCount, (maxColor << 1)
}

func findShiftsAndColor(img image.Image, step int) (int, int, int) {
	shiftX := 0
	shiftY := 0
	maxCount, cellColor := countSameColors(img, step, shiftX, shiftY)

	for x := shiftX + 1; x < step; x++ {
		for y := shiftY + 1; y < step; y++ {
			count, color := countSameColors(img, step, x, y)

			if count > maxCount {
				maxCount = count
				cellColor = color
				shiftX = x
				shiftY = y
			}
		}
	}
	return shiftX, shiftY, cellColor
}

func getResultImage(img image.Image) image.Image {
	bounds := img.Bounds()

	newImage := image.NewRGBA(bounds)

	// Set color for each pixel.
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			newImage.Set(x, y, img.At(bounds.Max.X-x, y))
		}
	}

	step := findStep(newImage)

	shiftX, shiftY, cellColor := findShiftsAndColor(newImage, step)

	fmt.Println(cellColor)

	gc := draw2dimg.NewGraphicContext(newImage)
	gc.SetStrokeColor(color.RGBA{uint8(255 - cellColor), uint8(255 - cellColor), uint8(255 - cellColor), 255})
	gc.BeginPath()

	for x := shiftX; x < bounds.Max.X; x += step {
		for y := shiftY; y < bounds.Max.Y; y += step {
			xf := float64(x)
			yf := float64(y)
			gc.MoveTo(xf-5, yf)
			gc.LineTo(xf+5, yf)
			gc.MoveTo(xf, yf-5)
			gc.LineTo(xf, yf+5)
		}
	}

	gc.Close()
	gc.Stroke()

	return newImage
}

func main() {
	reader, err := os.Open("./input.jpg")

	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	img, _, err := image.Decode(reader)

	if err != nil {
		log.Fatal(err)
	}

	resultImage := getResultImage(img)

	err = draw2dimg.SaveToPngFile("./output.png", resultImage)

	if err != nil {
		log.Fatal(err)
	}

}
