package lib

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/gookit/color"
)

// calibration area is 1/2 squared top left corner

const (
	caliSetStartX  = 30        // calibration set-x start
	caliSetEndX    = 35        // calibration set-x end
	caliSearchMinY = 2         // calibration search y start
	caliSearchMaxY = 100       // calibration search y max
	thick          = 3         // pixels down and across to check for calibration and parsing
	whiteMinRGBSum = 600 * 257 // minimum sum of RGB values to be considered white

	// Allowable variance per R,G,or B, up or down from the
	// calibration variance to be considered that colour
	variance = 20 * 257
)

// Pixel struct example
type Colour struct {
	R uint32
	G uint32
	B uint32
}

// NewColour creates a new Colour object
func NewColour(r, g, b uint32) Colour {
	return Colour{
		R: r,
		G: g,
		B: b,
	}
}

type Colours []Colour

func NewColours(cs ...Colour) Colours {
	return cs
}

func RGBAtoColour(r, g, b, a uint32) Colour {
	return Colour{r, g, b}
}

func LoadColoursDown(x, y, down int, img image.Image) Colours {
	var cs Colours
	for i := 0; i < down; i++ {
		cs = append(cs, RGBAtoColour(img.At(x, y+i).RGBA()))
	}
	return cs
}

func LoadColours(x1, x2, y1, y2 int, img image.Image) Colours {
	var cs Colours
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			cs = append(cs, RGBAtoColour(img.At(x, y).RGBA()))
		}
	}
	return cs
}

func AvgColour(x1, x2, y1, y2 int, img image.Image) Colour {
	var R, G, B, i uint32 = 0, 0, 0, 0
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			c := RGBAtoColour(img.At(x, y).RGBA())
			R += c.R
			G += c.G
			B += c.B
			i++
		}
	}
	return NewColour(R/i, G/i, B/i)
}

func (c Colour) IsWhite() bool {
	return c.R+c.B+c.G >= whiteMinRGBSum
}

func (c Colour) String() string {
	return fmt.Sprintf("R: %v,\tG: %v,\tB: %v", c.R/257, c.G/257, c.B/257)
}

func (c Colour) PrintColour(rightHandText string) {
	clib := color.RGB(uint8(c.R/257), uint8(c.G/257), uint8(c.B/257), true) // bg color
	clib.Print("  ")
	fmt.Printf(" %v", rightHandText)
}

func (c Colour) WithinVariance(target Colour) bool {
	if !(c.R+variance >= target.R && c.R-variance <= target.R) {
		return false
	}
	if !(c.G+variance >= target.G && c.G-variance <= target.G) {
		return false
	}
	if !(c.B+variance >= target.B && c.B-variance <= target.B) {
		return false
	}
	return true
}

func (cs Colours) AnyWhite() bool {
	for _, c := range cs {
		if c.R+c.B+c.G >= whiteMinRGBSum {
			return true
		}
	}
	return false
}

func (cs Colours) AvgColour() Colour {
	var R, G, B uint32 = 0, 0, 0
	for _, c := range cs {
		R += c.R
		G += c.G
		B += c.B
	}
	ln := uint32(len(cs))
	return NewColour(R/ln, G/ln, B/ln)
}

func (cs Colours) AllWithinVariance(target Colour) bool {
	for _, c := range cs {
		if !c.WithinVariance(target) {
			return false
		}
	}
	return true
}

// checks that no colours contained within the set have overlapping variance
// TODO make more efficient
func (cs Colours) AreUnique() bool {
	for i, c := range cs {
		for j, c2 := range cs {
			if i == j {
				continue
			}
			if c.WithinVariance(c2) {
				return false
			}
		}
	}
	return true
}

// at 200DPI a sharpie line is about 13 pixels wide
// Test 7px & 7px Down a specific colour
// start colour calibration read 20 pixels down

// Standard Order
// Red - Upright
// Green - Upside-down
// Blue - Read from the Left
// Purple - Read from Right

func Scan(pathToImage string) {

	// You can register another format here
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	file, err := os.Open(pathToImage)
	if err != nil {
		log.Fatal("Error: Image could not be decoded")
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	//bounds := img.Bounds()
	//width, height := bounds.Max.X, bounds.Max.Y

	// determine calibration
	caliNN, caliQP, caliHP, caliQT := getCalibrationColours(img)

	fmt.Println("\nCalibration Colours:\n")
	caliNN.PrintColour(fmt.Sprintf("Noon \t\t%v\n", caliNN.String()))
	caliQP.PrintColour(fmt.Sprintf("Quarter-Past \t%v\n", caliQP.String()))
	caliHP.PrintColour(fmt.Sprintf("Half-Past \t\t%v\n", caliHP.String()))
	caliQT.PrintColour(fmt.Sprintf("Quarter-To \t\t%v\n", caliQT.String()))

	if !(NewColours(caliNN, caliQP, caliHP, caliQT).AreUnique()) {
		log.Fatal("non-unique colours")
	}
}

func getCalibrationColours(img image.Image) (noon, quarterPast, halfPast, quarterTo Colour) {

	noon, outsideY, err := getCalibrationColour(caliSetStartX, caliSetEndX, caliSearchMinY, img)
	if err != nil {
		log.Fatalf("Error during calibration of noon: %v", err)
	}
	quarterPast, outsideY, err = getCalibrationColour(caliSetStartX, caliSetEndX, outsideY, img)
	if err != nil {
		log.Fatalf("Error during calibration of quarterPast: %v", err)
	}
	halfPast, outsideY, err = getCalibrationColour(caliSetStartX, caliSetEndX, outsideY, img)
	if err != nil {
		log.Fatalf("Error during calibration of halfPast: %v", err)
	}
	quarterTo, outsideY, err = getCalibrationColour(caliSetStartX, caliSetEndX, outsideY, img)
	if err != nil {
		log.Fatalf("Error during calibration of quarterTo: %v", err)
	}

	return noon, quarterPast, halfPast, quarterTo
}

//func getCalibrationCoordinates(img image.Image) (caliStartX, caliStartY int) {

//// first determine the set y
//for x := calibrationMinX; x <= calibrationMaxX; x++ {
//for y := calibrationMinY; y <= calibrationMaxY-4; y++ {
//cs := LoadColoursDown(x, y, calibrationThick, img)
////fmt.Printf("debug cs: %v\n", cs)
//if cs.AnyWhite() {
//continue
//}
//if cs.AllWithinVariance(cs.AvgColour()) {
//return x, y
//}
//}
//}

//log.Fatal("could not determine calibration start coordinates")
//return 0, 0
//}

// outsideX represents the first x coordinate outside the colour
func getCalibrationColour(setStartX, setEndX, searchStartY int, img image.Image) (caliColour Colour, outsideY int, err error) {

	caliStartY, caliEndY := 0, 0

	for y := searchStartY; y <= caliSearchMaxY-thick; y++ {
		var cs Colours
		if caliStartY == 0 {
			cs = LoadColours(setStartX, setEndX, y, y+thick, img)
		} else {
			cs = LoadColours(setStartX, setEndX, caliStartY, y+thick, img)
		}

		if cs.AnyWhite() {
			continue
		}
		awv := cs.AllWithinVariance(cs.AvgColour())
		switch {
		case awv && caliStartY == 0:
			caliStartY = y
			caliEndY = y + thick
		case awv && caliStartY != 0:
			caliEndY = y + thick
		case !awv && caliStartY != 0:
			caliEndY = y + thick - 1 // one less must have been the real end then
			caliColour = AvgColour(caliSetStartX, caliSetEndX, caliStartY, caliEndY, img)
			return caliColour, caliEndY + 1, nil
		}
	}

	return caliColour, 0, errors.New("could not determine calibration colour")
}
