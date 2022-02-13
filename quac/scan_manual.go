package quac

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"

	"github.com/disintegration/imaging"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	cmn "github.com/rigelrozanski/common"
	"golang.org/x/image/colornames"
)

func loadPicture(path string) (pic pixel.Picture, img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, img, err
	}
	defer file.Close()

	img, _, err = image.Decode(file)
	if err != nil {
		return nil, img, fmt.Errorf("error decoding: %v; error: %v", path, err)
	}
	return pixel.PictureDataFromImage(img), img, nil
}

func concatImage(img1, img2 image.Image) image.Image {

	maxX := img1.Bounds().Dx()
	maxX2 := img2.Bounds().Dx()
	if maxX2 > maxX {
		maxX = maxX2
	}

	y1 := img1.Bounds().Dy()
	y2 := y1 + img1.Bounds().Dy()

	concatRect := image.Rectangle{image.Point{0, 0},
		image.Point{maxX, y2}}
	img2Rect := img2.Bounds().Add(image.Point{0, y1})

	concatImg := image.NewRGBA(concatRect)
	draw.Draw(concatImg, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(concatImg, img2Rect, img2, image.Point{0, 0}, draw.Src)
	return concatImg
}

var (
	scanimgFilepath string
	quit            = false
)

const (
	rotated0deg   byte = 0x00
	rotated90deg  byte = 0x01
	rotated180deg byte = 0x02
	rotated270deg byte = 0x03
)

func getRotatedPos(win *pixelgl.Window, pos pixel.Vec, rotation byte) pixel.Vec {
	pos2 := pixel.Vec{}
	switch rotation {
	case rotated0deg:
		pos2.X = pos.X
		pos2.Y = pos.Y
	case rotated90deg:
		diff := math.Abs(win.Bounds().Max.Y-win.Bounds().Max.X) / 2
		pos2.X = diff + pos.Y
		pos2.Y = win.Bounds().Max.X - diff - pos.X
	case rotated180deg:
		pos2.X = win.Bounds().Max.X - pos.X
		pos2.Y = win.Bounds().Max.Y - pos.Y
	case rotated270deg:
		diff := math.Abs(win.Bounds().Max.Y-win.Bounds().Max.X) / 2
		pos2.X = win.Bounds().Max.Y + diff - pos.Y
		pos2.Y = -diff + pos.X
	}
	return pos2
}

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Manual Scan In... q=quit,r=rotate,mouse-drag=create-box",
		Bounds: pixel.R(0, 0, 1100, 850),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	if scanimgFilepath == "" {
		// XXX todo new logic which selects the default scan path from config
	}

	pic, img, err := loadPicture(scanimgFilepath)
	if err != nil {
		panic(err)
	}

	// ensure scan dir
	scanDir := path.Join(QuDir, "working_scan")
	_ = os.Mkdir(scanDir, os.ModePerm)

	//picRatio := pic.Bounds().Resized(pic.Bounds().Center(), pixel.V(1024, 768)) // resizes around the center
	doc := pixel.NewSprite(pic, pic.Bounds())
	imgScale := (1100.0) / pic.Bounds().W()
	imgScale2 := (850.0) / pic.Bounds().H()
	if imgScale2 < imgScale {
		imgScale = imgScale2
	}
	widthOffset := (pic.Bounds().W()*imgScale - 1100.0) / 2
	heightOffset := (pic.Bounds().H()*imgScale - 850.0) / 2
	imgOffset := pixel.V(widthOffset, heightOffset)

	var (
		// boxes of the []pixel.Vex indexes 0 & 1 are rectangle coordinates of the box
		boxes                      [][]pixel.Vec
		boxesRotation              []byte
		boxesIsConnectedToPrevious []bool

		camPos  = pixel.ZV
		camZoom = 1.0
		//camSpeed     = 100.0
		//camZoomSpeed = 1.2
		camRotation = 0.0

		rotation byte = rotated0deg
		//flipped          = false

		mouseDragging = false
	)

	for !win.Closed() {

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Min.Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			mouseDragging = true
			pos := win.MousePosition().Scaled(1 / camZoom).Add(camPos)
			//if flipped {
			//    pos.X = win.Bounds().Max.X - pos.X
			//    pos.Y = win.Bounds().Max.Y - pos.Y
			//}
			pos2 := getRotatedPos(win, pos, rotation)
			boxes = append(boxes, []pixel.Vec{pos2})
			boxesRotation = append(boxesRotation, rotation)
		}
		if win.JustReleased(pixelgl.MouseButtonLeft) {
			mouseDragging = false
			pos := win.MousePosition().Scaled(1 / camZoom).Add(camPos)
			//if flipped {
			//    pos.X = win.Bounds().Max.X - pos.X
			//    pos.Y = win.Bounds().Max.Y - pos.Y
			//}

			pos2 := getRotatedPos(win, pos, rotation)
			boxes[len(boxes)-1] = append(boxes[len(boxes)-1], pos2)
			isConnected := win.Pressed(pixelgl.KeyC)
			boxesIsConnectedToPrevious = append(boxesIsConnectedToPrevious, isConnected)
		}

		if win.JustPressed(pixelgl.KeyR) && !mouseDragging {
			//camRotation += math.Pi / 2
			//flipped = !flipped
			switch rotation {
			case rotated0deg: // DONT ASK ME WHY, we need this order to go clockwise
				camRotation = 3 * math.Pi / 2
				rotation = rotated270deg
			case rotated270deg:
				camRotation = math.Pi
				rotation = rotated180deg
			case rotated180deg:
				camRotation = math.Pi / 2
				rotation = rotated90deg
			case rotated90deg:
				camRotation = 0
				rotation = rotated0deg
			}
		}
		if win.JustPressed(pixelgl.KeyQ) && !mouseDragging {
			quit = true
			err = os.RemoveAll(scanDir)
			if err != nil {
				log.Fatal(err)
			}
			win.Destroy()
			return
		}

		//if !flipped {
		//    imd.SetMatrix(cam)
		//} else {
		imd.SetMatrix(cam.Rotated(win.Bounds().Center(), camRotation))
		//}

		imd.Clear()
		for i, box := range boxes {
			switch boxesRotation[i] {
			case rotated0deg:
				imd.Color = colornames.Blue
			case rotated90deg:
				imd.Color = colornames.Green
			case rotated180deg:
				imd.Color = colornames.Red
			case rotated270deg:
				imd.Color = colornames.Orange
			}
			if len(box) == 2 {
				imd.Push(box[0], box[1])
				imd.Rectangle(5)
				if boxesIsConnectedToPrevious[i] && i > 0 {
					imd.Push(
						pixel.Rect{boxes[i-1][0], boxes[i-1][1]}.Center(),
						pixel.Rect{box[0], box[1]}.Center(),
					)
					imd.Line(3)
				}
			}
			if len(box) == 1 {
				pos := win.MousePosition().Scaled(1 / camZoom).Add(camPos)
				//if flipped {
				//    pos.X = win.Bounds().Max.X - pos.X
				//    pos.Y = win.Bounds().Max.Y - pos.Y
				//}

				pos2 := getRotatedPos(win, pos, rotation)
				imd.Push(box[0], pos2)
				imd.Rectangle(5)
			}
		}

		// TODO zoom and navigation
		//if win.JustPressed(pixelgl.KeyLeft) {
		//camPos.X += camSpeed
		//}
		//if win.JustPressed(pixelgl.KeyRight) {
		//camPos.X -= camSpeed
		//}
		//if win.JustPressed(pixelgl.KeyDown) {
		//camPos.Y += camSpeed
		//}
		//if win.JustPressed(pixelgl.KeyUp) {
		//camPos.Y -= camSpeed
		//}
		//if win.JustPressed(pixelgl.KeyEqual) {
		//camZoom *= 1.2
		//}
		//if win.JustPressed(pixelgl.KeyMinus) {
		//camZoom /= 1.2
		//}
		//camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		// undo
		if win.JustPressed(pixelgl.KeyU) {
			boxes = boxes[:len(boxes)-1]
			boxesRotation = boxesRotation[:len(boxesRotation)-1]
		}

		if win.JustPressed(pixelgl.KeyEnter) {

			// create images
			imgs := []image.Image{}
			imgsIsConnectedToPrevious := []bool{}
			for i, box := range boxes {

				// capture box areas
				smn := box[0].Add(imgOffset).Scaled(1 / imgScale)
				smx := box[1].Add(imgOffset).Scaled(1 / imgScale)

				imgH := img.Bounds().Dy()
				simg := img.(interface {
					SubImage(r image.Rectangle) image.Image
				}).SubImage(image.Rect(int(smn.X), imgH-int(smn.Y), int(smx.X), imgH-int(smx.Y)))

				// skip if just a click
				if simg.Bounds().Dy() < 10 || simg.Bounds().Dx() < 10 {
					continue
				}

				//if boxesIsFlipped[i] {
				//simg = imaging.Rotate(simg, 180, color.RGBA{0, 0, 0, 1})
				//}

				switch boxesRotation[i] {
				case rotated0deg:
				case rotated90deg:
					simg = imaging.Rotate(simg, 90, color.RGBA{0, 0, 0, 1})
				case rotated180deg:
					simg = imaging.Rotate(simg, 180, color.RGBA{0, 0, 0, 1})
				case rotated270deg:
					simg = imaging.Rotate(simg, 270, color.RGBA{0, 0, 0, 1})
				}

				imgs = append(imgs, simg)
				imgsIsConnectedToPrevious = append(
					imgsIsConnectedToPrevious, boxesIsConnectedToPrevious[i])
			}

			// save files
			for i, img := range imgs {
				if img == nil {
					continue
				}
				f, err := os.Create(path.Join(scanDir, fmt.Sprintf("outimage_%v.png", i)))
				if err != nil {
					panic(err)
				}
				defer f.Close()
				err = png.Encode(f, img)
				if err != nil {
					panic(err)
				}
			}
			// reload images
			reimgs := []image.Image{}
			for i := range imgs {
				_, img, err := loadPicture(path.Join(scanDir, fmt.Sprintf("outimage_%v.png", i)))
				if err != nil {
					panic(err)
				}
				reimgs = append(reimgs, img)
			}

			imgConcats := []image.Image{}
			for i, img := range reimgs {
				if !imgsIsConnectedToPrevious[i] {
					imgConcats = append(imgConcats, img)
					continue
				}
				if i == 0 && imgsIsConnectedToPrevious[i] {
					panic("cannot be first img connected to previous")
				}
				imgConcats[len(imgConcats)-1] = concatImage(imgConcats[len(imgConcats)-1], img)
			}

			// save final images as ideas
			for i, img := range imgConcats {
				if img == nil {
					continue
				}
				filepath := path.Join(scanDir, fmt.Sprintf("outimageCON_%v.png", i))
				f, err := os.Create(filepath)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				err = png.Encode(f, img)
				if err != nil {
					panic(err)
				}

				idea := NewIdeaFromFile("UNTAGGED", filepath)
				err = cmn.Copy(filepath, idea.Path())
				if err != nil {
					log.Fatal(err)
				}
				PrependLast(idea.Id)
				IncrementID()

				fmt.Println("Added the following idea:")
				View(idea.Path())

			}

			err = os.RemoveAll(scanDir)
			if err != nil {
				log.Fatal(err)
			}
			win.Destroy()
			return
		}

		win.Clear(colornames.Aliceblue)
		doc.Draw(win, pixel.IM.Moved(win.Bounds().Center()).
			Rotated(win.Bounds().Center(), camRotation).
			Scaled(win.Bounds().Center(), imgScale),
		)
		imd.Draw(win)
		win.Update()
	}
}

func ScanManual(pathToImageOrDir string) {

	if pathToImageOrDir == "" && len(DefaultScanDir) > 0 {
		pathToImageOrDir = DefaultScanDir
	}

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
			if !file.IsDir() && file.Name()[0] != '.' {
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

	for _, imgFile := range imgFiles {
		if quit {
			break
		}
		scanimgFilepath = imgFile
		pixelgl.Run(run)
		if quit {
			break
		}
		if DeleteWhenScanning {
			_ = os.Remove(imgFile)
		}
	}
}
