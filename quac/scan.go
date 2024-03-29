package quac

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/common/colour"
	"github.com/rigelrozanski/thranch/quac/idea"
)

// calibration area is 1/2 squared top left corner

// at 200DPI a sharpie line is about 13 pixels wide
// Test 7px & 7px Down a specific colour
// start colour calibration read 20 pixels down

const (
	caliSetStartX  = 30  // calibration set-x start
	caliSetEndX    = 35  // calibration set-x end
	caliSearchMinY = 5   // calibration search y start
	caliSearchMaxY = 100 // calibration search y max
	thick          = 3   // pixels down and across to check for calibration and parsing

	minIdeaDimention = 50 // must be 50 pixels in each direction to be considered an object

	// Allowable variance per R,G,or B, up or down from the
	// calibration variance to be considered that colour
	variance = 20 * 257

	brightnessVariance = 30 * 257
)

var LastScanCalibrationFile string

func Scan(pathToImageOrDir, opTag string) {

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("PNG", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("JPG", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("JPEG", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	fod, err := os.Stat(pathToImageOrDir)
	if err != nil {
		log.Fatal(err)
	}
	isDir := fod.Mode().IsDir()

	var imgFiles []string
	if isDir {
		files, err := ioutil.ReadDir(pathToImageOrDir)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if !file.IsDir() {
				filepath := path.Join(pathToImageOrDir, file.Name())
				imgFiles = append(imgFiles, filepath)
			}
		}
		if len(imgFiles) == 0 {
			log.Fatal("directory is empty")
		}
	} else {
		imgFiles = []string{pathToImageOrDir}
	}

	// retrieve calibration colours directly from ideas
	var caliNN, caliQP, caliHP, caliQT colour.Colour
	tags := []idea.Tag{idea.MustNewTagReg("scan-calibration", "")}
	imagesIdeas := GetAllIdeas().WithImage().WithTags(tags)
	for _, idear := range imagesIdeas {
		avgCol := getAvgColourFromFile(idear.Path())
		switch {
		case idear.HasTag(idea.MustNewTagReg("orientation", "noon")):
			caliNN = avgCol
		case idear.HasTag(idea.MustNewTagReg("orientation", "quarter-past")):
			caliQP = avgCol
		case idear.HasTag(idea.MustNewTagReg("orientation", "half-past")):
			caliHP = avgCol
		case idear.HasTag(idea.MustNewTagReg("orientation", "quarter-to")):
			caliQT = avgCol
		}
	}

	caliColours := colour.NewColours(caliNN, caliQP, caliHP, caliQT)
	PrintCaliColours(caliColours)
	if !(caliColours.AreUnique(variance)) {
		log.Fatal("non-unique calibration colours")
	}

	fmt.Println("confirm calibration colours (Y/N)")
	consoleScanner := bufio.NewScanner(os.Stdin)
	_ = consoleScanner.Scan()
	in := consoleScanner.Text()
	if in != "Y" {
		fmt.Println("okay! exiting")
		os.Exit(1)
	}

	for _, ifn := range imgFiles {

		file, err := os.Open(ifn)
		if err != nil {
			log.Fatal("Error: Image could not be decoded")
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("creating calibration grid for %v...\n", file.Name())
		// get the calibration img
		// 1 == noon
		// 2 == quarterPast
		// 3 == halfPast
		// 4 == quarterTo
		caliGrid := createCalibrationGrid(img, caliColours)

		//fmt.Println("removing image marks...")
		// remove all the marks from the image
		//imgRM := removeMarks(img, caliGrid)

		fmt.Println("extracting subimages...")
		noonResults := extractSubsetImgs(img, caliGrid, 1)
		quarterPastResults := Rotate90(extractSubsetImgs(img, caliGrid, 2))
		halfPastResults := Rotate180(extractSubsetImgs(img, caliGrid, 3))
		quarterToResults := Rotate270(extractSubsetImgs(img, caliGrid, 4))

		// concat
		var results []image.Image
		results = append(results, noonResults...)
		results = append(results, quarterPastResults...)
		results = append(results, halfPastResults...)
		results = append(results, quarterToResults...)

		// ensure scan dir
		scanDir := path.Join(QuDir, "working_scan")
		_ = os.Mkdir(scanDir, os.ModePerm)

		fmt.Println("saving files...")
		var imgPaths []string
		for i, result := range results {
			caliImgPath := path.Join(scanDir, strconv.Itoa(i)+".png")
			imgPaths = append(imgPaths, caliImgPath)
			f, _ := os.Create(caliImgPath)
			png.Encode(f, result)
		}

		commonTag := false
		if opTag != "" {
			commonTag = true
		}

		for _, imgPath := range imgPaths {
			ViewImageNoFilename(imgPath)
			fmt.Println("please enter tags seperated by spaces then press enter:")
			consoleScanner := bufio.NewScanner(os.Stdin)
			_ = consoleScanner.Scan()
			tags := strings.Fields(consoleScanner.Text())

			if len(tags) > 0 && tags[0] == "DEL" {
				continue
			}

			if commonTag {
				tags = append(tags, opTag)
			}

			clumpedTags := strings.Join(tags, ",")

			// save the new idea
			idea := NewIdeaFromFile(clumpedTags, imgPath)
			err := cmn.Copy(imgPath, idea.Path())
			if err != nil {
				log.Fatal(err)
			}
			IncrementID()
		}

		err = os.RemoveAll(scanDir)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func getAvgColourFromFile(filpath string) colour.Colour {
	file, err := os.Open(filpath)
	if err != nil {
		log.Fatal("Error: Image could not be decoded")
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	bounds := img.Bounds()
	return colour.LoadColours(0, bounds.Max.X, 0, bounds.Max.Y, img).AvgColour()
}

func PrintCaliColours(in colour.Colours) {
	if len(in) != 4 {
		panic("bad number of colours to print")
	}
	fmt.Println("\nCalibration Colours:")
	in[0].PrintColour(fmt.Sprintf("Noon \t\t%v\n", in[0].String()))
	in[1].PrintColour(fmt.Sprintf("Quarter-Past \t%v\n", in[1].String()))
	in[2].PrintColour(fmt.Sprintf("Half-Past \t\t%v\n", in[2].String()))
	in[3].PrintColour(fmt.Sprintf("Quarter-To \t\t%v\n", in[3].String()))
}

// TODO integrate into new command
//func getCalibrationColours(img image.Image) (noon, quarterPast, halfPast, quarterTo colour.Colour, err error) {

//noon, outsideY, err := colour.GetCalibrationColour(caliSetStartX, caliSetEndX, caliSearchMinY, caliSearchMaxY, thick, variance, img)
//if err != nil {
//err := errors.New(fmt.Sprintf("Error during calibration of noon: %v", err))
//return noon, quarterPast, halfPast, quarterTo, err
//}
//quarterPast, outsideY, err = colour.GetCalibrationColour(caliSetStartX, caliSetEndX, outsideY, caliSearchMaxY, thick, variance, img)
//if err != nil {
//err = errors.New(fmt.Sprintf("Error during calibration of quarterPast: %v", err))
//return noon, quarterPast, halfPast, quarterTo, err
//}
//halfPast, outsideY, err = colour.GetCalibrationColour(caliSetStartX, caliSetEndX, outsideY, caliSearchMaxY, thick, variance, img)
//if err != nil {
//err = errors.New(fmt.Sprintf("Error during calibration of halfPast: %v", err))
//return noon, quarterPast, halfPast, quarterTo, err
//}
//quarterTo, outsideY, err = colour.GetCalibrationColour(caliSetStartX, caliSetEndX, outsideY, caliSearchMaxY, thick, variance, img)
//if err != nil {
//err = errors.New(fmt.Sprintf("Error during calibration of quarterTo: %v", err))
//return noon, quarterPast, halfPast, quarterTo, err
//}

//return noon, quarterPast, halfPast, quarterTo, nil
//}

// create the grid seperating indicating areas found from the calibration
func createCalibrationGrid(img image.Image, Cali colour.Colours) [][]uint8 {

	bounds := img.Bounds()
	maxXInc, maxYInc := bounds.Max.X, bounds.Max.Y

	caliGrid := make([][]uint8, maxXInc+1)
	for i := 0; i < bounds.Max.X+1; i++ {
		caliGrid[i] = make([]uint8, maxYInc+1)
	}

	for y := 0; y <= maxYInc; y++ {
		for x := 0; x <= maxXInc; x++ {

			c := colour.LoadColour(x, y, img)
			i, _, withinVariance := Cali.NearestColourTo(c, brightnessVariance, variance)
			if !withinVariance {
				continue
			}
			caliGrid[x][y] = uint8(i + 1) // so doesn't conflict with zero default values
		}
	}

	// filter out lone pixels
	for y := 0; y <= maxYInc; y++ {
		for x := 0; x <= maxXInc; x++ {
			if caliGrid[x][y] != 0 {
				leftOk := x-1 >= 0
				rightOk := x+1 <= maxXInc
				upOk := y-1 >= 0
				downOk := y+1 <= maxYInc
				if (leftOk && caliGrid[x-1][y] != 0) ||
					(rightOk && caliGrid[x+1][y] != 0) ||
					(upOk && caliGrid[x][y-1] != 0) ||
					(downOk && caliGrid[x][y+1] != 0) ||
					(leftOk && upOk && caliGrid[x-1][y-1] != 0) ||
					(leftOk && downOk && caliGrid[x-1][y+1] != 0) ||
					(rightOk && upOk && caliGrid[x+1][y-1] != 0) ||
					(rightOk && downOk && caliGrid[x+1][y+1] != 0) {
					continue
				}
				caliGrid[x][y] = 0
			}
		}
	}

	return caliGrid
}

// turns all the pixels white in img where a value in caliGrid is present
func removeMarks(img image.Image, caliGrid [][]uint8) image.Image {

	bounds := img.Bounds()
	maxXInc, maxYInc := bounds.Max.X, bounds.Max.Y

	upLeft := image.Point{0, 0}
	lowRight := image.Point{maxXInc, maxYInc}

	imgOut := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// expand caliGrid by some pixels
	exPix := 5
	caliGridExp := make([][]uint8, bounds.Max.X+1)
	for i := 0; i < bounds.Max.X+1; i++ {
		caliGridExp[i] = make([]uint8, bounds.Max.Y+1)
	}
	for shiftX := -exPix; shiftX <= exPix; shiftX++ {
		for shiftY := -exPix; shiftY <= exPix; shiftY++ {
			for y := 0; y <= maxYInc; y++ {
				for x := 0; x <= maxXInc; x++ {
					if x+shiftX < 0 || x+shiftX > maxXInc || y+shiftY < 0 || y+shiftY > maxYInc {
						continue
					}
					if caliGrid[x+shiftX][y+shiftY] != 0 {
						caliGridExp[x][y] = 1
					}
				}
			}
		}
	}

	for y := 0; y <= maxYInc; y++ {
		for x := 0; x <= maxXInc; x++ {
			if caliGridExp[x][y] == 0 {
				imgOut.Set(x, y, img.At(x, y))
			} else {
				imgOut.Set(x, y, color.White)
			}
		}
	}
	return imgOut
}

func extractSubsetImgs(img image.Image, caliGrid [][]uint8, target uint8) []image.Image {

	var results []image.Image

	bounds := img.Bounds()
	maxXInc, maxYInc := bounds.Max.X, bounds.Max.Y

	// areas with the target image
	// each new area is assigned a new integer
	areas := make([][]uint8, bounds.Max.X+1)
	for i := 0; i < bounds.Max.X+1; i++ {
		areas[i] = make([]uint8, bounds.Max.Y+1)
	}
	areaI := uint8(0)

	// bounds definition, index is the areaI, hence the first record is a dummy
	boundsMinX := []int{0}
	boundsMaxX := []int{0}
	boundsMinY := []int{0}
	boundsMaxY := []int{0}

	// determine all individual areas
	// keep track of the bounds while we're at it
	for y := 0; y <= maxYInc; y++ {
		for x := 0; x <= maxXInc; x++ {
			if areas[x][y] != 0 {
				continue
			}
			if caliGrid[x][y] == target {
				areaI++
				boundsMinX = append(boundsMinX, x)
				boundsMaxX = append(boundsMaxX, x)
				boundsMinY = append(boundsMinY, y)
				boundsMaxY = append(boundsMaxY, y)

				propagate(x, y, maxXInc, maxYInc, areaI, (&areas),
					&(boundsMinX[areaI]), &(boundsMaxX[areaI]),
					&(boundsMinY[areaI]), &(boundsMaxY[areaI]),
					caliGrid, target)
			}
		}
	}

	// save the resulting images
	for i := uint8(1); i <= areaI; i++ {

		// skip if the dimentions are too small
		if (boundsMaxX[i]-boundsMinX[i] < minIdeaDimention) ||
			(boundsMaxY[i]-boundsMinY[i] < minIdeaDimention) {
			continue
		}

		rect := image.Rect(boundsMinX[i], boundsMinY[i], boundsMaxX[i], boundsMaxY[i])
		resImg := imaging.Crop(img, rect)
		results = append(results, resImg)
	}

	return results
}

func propagate(x, y, maxXInc, maxYInc int, areaI uint8, areas *[][]uint8,
	boundsMinX, boundsMaxX, boundsMinY, boundsMaxY *int,
	caliImg [][]uint8, target uint8) {

	if caliImg[x][y] != target {
		return
	}
	if (*areas)[x][y] != 0 {
		return
	}
	(*areas)[x][y] = areaI

	// update bounds
	if x < (*boundsMinX) {
		(*boundsMinX) = x
	}
	if x > (*boundsMaxX) {
		(*boundsMaxX) = x
	}
	if y < (*boundsMinY) {
		(*boundsMinY) = y
	}
	if y > (*boundsMaxY) {
		(*boundsMaxY) = y
	}

	leftOk := x-1 >= 0
	rightOk := x+1 <= maxXInc
	upOk := y-1 >= 0
	downOk := y+1 <= maxYInc
	if leftOk {
		propagate(x-1, y, maxXInc, maxYInc, areaI, areas, boundsMinX,
			boundsMaxX, boundsMinY, boundsMaxY, caliImg, target)
	}
	if leftOk && downOk {
		propagate(x-1, y+1, maxXInc, maxYInc, areaI, areas, boundsMinX,
			boundsMaxX, boundsMinY, boundsMaxY, caliImg, target)
	}
	if leftOk && upOk {
		propagate(x-1, y-1, maxXInc, maxYInc, areaI, areas, boundsMinX,
			boundsMaxX, boundsMinY, boundsMaxY, caliImg, target)
	}
	if rightOk {
		propagate(x+1, y, maxXInc, maxYInc, areaI, areas, boundsMinX,
			boundsMaxX, boundsMinY, boundsMaxY, caliImg, target)
	}
	if rightOk && downOk {
		propagate(x+1, y+1, maxXInc, maxYInc, areaI, areas, boundsMinX,
			boundsMaxX, boundsMinY, boundsMaxY, caliImg, target)
	}
	if rightOk && upOk {
		propagate(x+1, y-1, maxXInc, maxYInc, areaI, areas, boundsMinX,
			boundsMaxX, boundsMinY, boundsMaxY, caliImg, target)
	}
	if downOk {
		propagate(x, y+1, maxXInc, maxYInc, areaI, areas, boundsMinX,
			boundsMaxX, boundsMinY, boundsMaxY, caliImg, target)
	}
	if upOk {
		propagate(x, y-1, maxXInc, maxYInc, areaI, areas, boundsMinX,
			boundsMaxX, boundsMinY, boundsMaxY, caliImg, target)
	}
}

func Rotate90(imgs []image.Image) []image.Image {
	var out []image.Image
	for _, img := range imgs {
		out = append(out, imaging.Rotate90(img))
	}
	return out
}

func Rotate180(imgs []image.Image) []image.Image {
	var out []image.Image
	for _, img := range imgs {
		out = append(out, imaging.Rotate180(img))
	}
	return out
}

func Rotate270(imgs []image.Image) []image.Image {
	var out []image.Image
	for _, img := range imgs {
		out = append(out, imaging.Rotate270(img))
	}
	return out
}
