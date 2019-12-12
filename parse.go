package main

import (
	"math"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
)

func getRed(x int, y int, img image.Image) int {
	r, _, _, _ := img.At(x, y).RGBA()
	return int(r >> 8)
}

func getGreen(x int, y int, img image.Image) int {
	_, g, _, _ := img.At(x, y).RGBA()
	return int(g >> 8)
}

func getBlue(x int, y int, img image.Image) int {
	_, _, b, _ := img.At(x, y).RGBA()
	return int(b >> 8)
}

const blockWidth = 5
const blockHeight = 5
const blockLength = blockWidth * blockHeight
const stepX = blockWidth >> 1
const stepY = blockHeight >> 1

func getPart(getColor func(int,int,image.Image) int, centerX int, centerY int, img image.Image) (res []int) {
	for x := centerX - stepX; x < centerX + stepX; x++ {
		for y := centerY - stepY; y < centerY + stepY; y++ {
			color := getColor(x, y, img)
			res = append(res, color)
		}
	}
	return
}

func getK(numbers []int, k int) int {
	
	if (len(numbers) == 1) {
		return numbers[0]
	}
	
	less := make([]int, 0)
	equal := make([]int, 0)
	more := make([]int, 0)
	
	pivot := numbers[len(numbers) >> 1]
	

	for _, value := range numbers {
		if value > pivot {
			more = append(more, value)
		} else if value < pivot {
			less = append(less, value)
		} else {
			equal = append(equal, value)
		}
	}

	if k < len(less) {
		return getK(less, k)
	}

	if k < len(less) + len(equal) {
		return pivot
	}

	return getK(more, k - len(less) - len(equal))
	
}

func meanQuadraticDifference(numbers []int) int {
	median := getK(numbers, len(numbers) >> 1)

	sum := 0

	for _, value := range numbers {
		sum += (value - median) * (value - median)
	}

	divisor := 1

	if len(numbers) > 1 {
		divisor = len(numbers) - 1
	}

	return int(math.Round(math.Sqrt(float64(sum) / float64(divisor))))
}


func getResultImage(img image.Image) image.Image {
	bounds := img.Bounds()

	topLeftX := bounds.Min.X
	topLeftY := bounds.Min.Y
	botRightX := bounds.Max.X
	botRightY := bounds.Max.Y

	newImage := image.NewRGBA(bounds)

	max := 0

	for y := topLeftY + stepY; y < botRightY - stepY; y++ {
		for x := topLeftX + stepX; x < botRightX - stepX; x++ {
			redPart := getPart(getRed, x, y, img)
			bluePart := getPart(getBlue, x, y, img)
			greenPart := getPart(getGreen, x, y, img)

			meanQuadraticDifference := meanQuadraticDifference(redPart) + meanQuadraticDifference(greenPart) + meanQuadraticDifference(bluePart)

			if (meanQuadraticDifference > max) {
				max = meanQuadraticDifference
			}
		}
	}
	for y := topLeftY + stepY; y < botRightY - stepY; y++ {
		for x := topLeftX + stepX; x < botRightX - stepX; x++ {
			redPart := getPart(getRed, x, y, img)
			bluePart := getPart(getBlue, x, y, img)
			greenPart := getPart(getGreen, x, y, img)

			meanQuadraticDifference := meanQuadraticDifference(redPart) + meanQuadraticDifference(greenPart) + meanQuadraticDifference(bluePart)
			c := meanQuadraticDifference * 255 / max

			if c < 100 {
				c = 0
			} else  {
				c = 255
			}

			newImage.Set(x, y, color.Gray{uint8(c)})
		}
	}

	fmt.Println(max)

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
