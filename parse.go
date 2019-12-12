package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"sort"

	"github.com/llgcode/draw2d/draw2dimg"
)

func getColor(x int, y int, img image.Image) float64 {
	r, g, b, _ := img.At(x, y).RGBA()
	smallR := float64(r>>12) * 16
	smallG := float64(g>>12) * 16
	smallB := float64(b>>12) * 16

	return (smallR + smallG + smallB) / 3.0 / 255
}

func getColorGreenRed(x int, y int, img image.Image) float64 {
	r, g, _, _ := img.At(x, y).RGBA()
	smallR := float64(r >> 8)
	smallG := float64(g >> 8)

	return (smallG - smallR) / 255
}

func getPart(img image.Image, x int, y int) [5][5]float64 {
	return [5][5]float64{
		[5]float64{getColor(x-2, y-2, img), getColor(x-1, y-2, img), getColor(x, y-2, img), getColor(x+1, y-2, img), getColor(x+2, y-2, img)},
		[5]float64{getColor(x-2, y-1, img), getColor(x-1, y-1, img), getColor(x, y-1, img), getColor(x+1, y-1, img), getColor(x+2, y-1, img)},
		[5]float64{getColor(x-2, y, img), getColor(x-1, y, img), getColor(x, y, img), getColor(x+1, y, img), getColor(x+2, y, img)},
		[5]float64{getColor(x-2, y+1, img), getColor(x-1, y+1, img), getColor(x, y+1, img), getColor(x+1, y+1, img), getColor(x+2, y+1, img)},
		[5]float64{getColor(x-2, y+2, img), getColor(x-1, y+2, img), getColor(x, y+2, img), getColor(x+1, y+2, img), getColor(x+2, y+2, img)},
	}
}
func getPartGreenRed(img image.Image, x int, y int) [5][5]float64 {
	return [5][5]float64{
		[5]float64{getColorGreenRed(x-2, y-2, img), getColorGreenRed(x-1, y-2, img), getColorGreenRed(x, y-2, img), getColorGreenRed(x+1, y-2, img), getColorGreenRed(x+2, y-2, img)},
		[5]float64{getColorGreenRed(x-2, y-1, img), getColorGreenRed(x-1, y-1, img), getColorGreenRed(x, y-1, img), getColorGreenRed(x+1, y-1, img), getColorGreenRed(x+2, y-1, img)},
		[5]float64{getColorGreenRed(x-2, y, img), getColorGreenRed(x-1, y, img), getColorGreenRed(x, y, img), getColorGreenRed(x+1, y, img), getColorGreenRed(x+2, y, img)},
		[5]float64{getColorGreenRed(x-2, y+1, img), getColorGreenRed(x-1, y+1, img), getColorGreenRed(x, y+1, img), getColorGreenRed(x+1, y+1, img), getColorGreenRed(x+2, y+1, img)},
		[5]float64{getColorGreenRed(x-2, y+2, img), getColorGreenRed(x-1, y+2, img), getColorGreenRed(x, y+2, img), getColorGreenRed(x+1, y+2, img), getColorGreenRed(x+2, y+2, img)},
	}
}

func getMedian(block [5][5]float64) float64 {
	var flatten = make([]float64, 25)

	for i := 0; i < len(block); i++ {
		for j := 0; j < len(block[i]); j++ {
			flatten[i*len(block)+j] = block[i][j]
		}
	}

	sort.Float64s(flatten)

	return flatten[12]
}

func calc(part [5][5]float64, coefs [5][5]float64) float64 {
	sum := float64(0)
	for i := 0; i < len(part); i++ {
		for j := 0; j < len(part[i]); j++ {
			sum += (part[i][j] - part[2][2]) * coefs[i][j]
		}
	}
	return 1 / (1 + math.Exp(-sum))

}
func calcComparison(part [5][5]float64, coefs [5][5]float64, pattern [5][5]float64) float64 {
	sum := float64(0)
	for i := 0; i < len(part); i++ {
		for j := 0; j < len(part[i]); j++ {
			sum += (part[i][j] - pattern[i][j]) * coefs[i][j]
		}
	}
	return 1 / (1 + math.Exp(-sum))
}

func probColor(probability float64) color.Color {
	return color.Gray{uint8(probability * 255)}
}

func reduce(img *image.RGBA) image.Image {
	bounds := img.Bounds()
	// 	0  0  0  0  0
	// 	0  0  0  0  0
	// 	0  0  X  0  0
	// 	0  0  0  0  0
	// 	0  0  0  0  0
	newImage := image.NewRGBA(bounds)

	for x := bounds.Min.X + 2; x < bounds.Max.X-2; x++ {
		for y := bounds.Min.Y + 2; y < bounds.Max.Y-2; y++ {
			part := getPartGreenRed(img, x, y)

			probability := calcComparison(part, [5][5]float64{
				[5]float64{1, 1, 1, 1, 1},
				[5]float64{1, 1, 1, 1, 1},
				[5]float64{1, 1, 1, 1, 1},
				[5]float64{1, 1, 1, 1, 1},
				[5]float64{1, 1, 1, 1, 1},
			}, [5][5]float64{
				[5]float64{1, 1, 1, 1, 1},
				[5]float64{0, 0, 0, 0, 0},
				[5]float64{0, 0, 0, 0, 0},
				[5]float64{0, 0, 0, 0, 0},
				[5]float64{0, 0, 0, 0, 0},
			})

			newImage.Set(x, y, probColor(probability))
		}
	}

	return newImage
}

func getIntervals(p image.Image) [][4]int {
	b := p.Bounds()
	res := make([][4]int, 0)
	for x := b.Min.X + 2; x < b.Max.X-2; x++ {
		for y := b.Min.Y + 2; y < b.Max.Y-2; y++ {
			leftProb := getColor(x-1, y, p)
			topProb := getColor(x, y+1, p)
			currentProb := getColor(x, y, p)
			if currentProb > 0.90 && leftProb < 0.5 && topProb < 0.5 {
				startX := x
				startY := y
				endX := startX + 1

				for getColor(endX, startY, p) > 0.8 {
					endX++
				}

				if (endX - startX) < 15 {
					continue
				}

				res = append(res, [4]int{startX, startY, endX, startY})
			}
		}
	}

	return res
}

func getIntervalsLengths(intervals [][4]int) []int {
	res := make([]int, 0)

	for _, interval := range intervals {
		res = append(res, interval[2]-interval[0])
	}

	return res
}

func getCellCenters(intervals [][4]int, img image.Image) [][2]int {
	points := make([][2]int, 0)

	maxIntervalIndex := 0
	step := intervals[0][2] - intervals[0][0]

	for index, interval := range intervals[1:] {
		length := interval[2] - interval[0]
		if length > step {
			step = length
			maxIntervalIndex = index
		}
	}

	maxInterval := intervals[maxIntervalIndex]
	maxIntervalCenter := [2]int{maxInterval[0] + (step >> 1), maxInterval[1] + (step >> 1)}
	centerColor := getColor(maxIntervalCenter[0], maxIntervalCenter[1], img)

	leftTopCorner := [2]int{maxIntervalCenter[0] % step, maxIntervalCenter[1] % step}

	bounds := img.Bounds()

	for x := leftTopCorner[1]; x < bounds.Max.X; x += step {
		for y := leftTopCorner[0]; y < bounds.Max.Y; y += step {
			color := getColor(x, y, img)
			if math.Abs(color-centerColor) < 0.1 {
				points = append(points, [2]int{x, y})
			}
		}
	}

	return points
}

func getResultImage(img image.Image) image.Image {
	bounds := img.Bounds()

	newImage := image.NewRGBA(bounds)

	// Set color for each pixel.
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y>>1; y++ {
			if x > bounds.Min.X+2 && x < bounds.Max.X-2 && y < bounds.Max.Y-2 && y > bounds.Min.Y+2 {
				part := getPart(img, x, y)

				median := getMedian(part)

				newImage.Set(x, y, color.Gray{uint8(median * 255)})
			} else {
				newImage.Set(x, y, img.At(x, y))
			}
		}
	}

	for x := bounds.Min.X + 2; x < bounds.Max.X-2; x++ {
		for y := bounds.Min.Y + 2; y < bounds.Max.Y-2; y++ {
			part := getPart(img, x, y)
			// topRight := [5][5]float64{
			// 	[5]float64{0.0, -31.0, 4.0, 2.0, 1.0},
			// 	[5]float64{0.0, -31.0, 8.0, 4.0, 2.0},
			// 	[5]float64{0.0, -31.0, 32.0, 16.0, 8.0},
			// 	[5]float64{0.0, -31.0, -31.0, -31.0, -31.0},
			// 	[5]float64{0.0, 0.0, 0.0, 0.0, 0.0},
			// } highlights letters

			rvalue := calc(part, [5][5]float64{
				[5]float64{-3, -3, 1, 2, 3},
				[5]float64{-3, -3, 1, 2, 3},
				[5]float64{-3, -3, 1, 2, 3},
				[5]float64{-3, -3, 1, 2, 3},
				[5]float64{-3, -3, 1, 2, 3},
			})
			gvalue := calc(part, [5][5]float64{
				[5]float64{-3, -3, -3, -3, -3},
				[5]float64{-3, -3, -3, -3, -3},
				[5]float64{1, 1, 1, 1, 1},
				[5]float64{2, 2, 2, 2, 2},
				[5]float64{3, 3, 3, 3, 3},
			})

			r := uint8(rvalue * 255)
			g := uint8(gvalue * 255)
			b := uint8(0)
			const threshold = 250
			if r > threshold {
				r = 255
			} else {
				r = 0
			}
			if g > threshold {
				g = 255
			} else {
				g = 0
			}

			newImage.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	probabilities := reduce(newImage)

	intervals := getIntervals(probabilities)

	dots := getCellCenters(intervals, img)

	fmt.Println(dots)

	resImage := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			resImage.Set(x, y, img.At(x, y))
		}
	}

	gc := draw2dimg.NewGraphicContext(resImage)
	gc.SetStrokeColor(color.RGBA{255, 0, 0, 255})
	gc.SetLineWidth(3)
	gc.BeginPath()
	for _, center := range dots {
		gc.MoveTo(float64(center[0]-5), float64(center[1]))
		gc.LineTo(float64(center[0]+5), float64(center[1]))
		gc.MoveTo(float64(center[0]), float64(center[1]-5))
		gc.LineTo(float64(center[0]), float64(center[1]+5))
	}
	gc.Stroke()
	gc.Close()
	fmt.Println(getIntervalsLengths(intervals))
	return resImage
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
