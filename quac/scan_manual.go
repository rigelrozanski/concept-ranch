package quac

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
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
		return nil, img, err
	}
	return pixel.PictureDataFromImage(img), img, nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Manual Scan In...",
		Bounds: pixel.R(0, 0, 1100, 850),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	pic, img, err := loadPicture("apple.jpeg")
	if err != nil {
		panic(err)
	}

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
		boxes        [][]pixel.Vec
		boxesColours []color.RGBA

		camPos  = pixel.ZV
		camZoom = 1.0
		//camSpeed     = 100.0
		//camZoomSpeed = 1.2
		rotation      = 0.0
		flipped       = false
		mouseDragging = false
	)

	for !win.Closed() {

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Min.Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			mouseDragging = true
			pos := win.MousePosition().Scaled(1 / camZoom).Add(camPos)
			if flipped {
				pos.X = win.Bounds().Max.X - pos.X
				pos.Y = win.Bounds().Max.Y - pos.Y
			}
			boxes = append(boxes, []pixel.Vec{pos})
			if !flipped {
				boxesColours = append(boxesColours, colornames.Red)
			} else {
				boxesColours = append(boxesColours, colornames.Blue)
			}
		}
		if win.JustReleased(pixelgl.MouseButtonLeft) {
			mouseDragging = false
			pos := win.MousePosition().Scaled(1 / camZoom).Add(camPos)
			if flipped {
				pos.X = win.Bounds().Max.X - pos.X
				pos.Y = win.Bounds().Max.Y - pos.Y
			}
			boxes[len(boxes)-1] = append(boxes[len(boxes)-1], pos)
		}

		if win.JustPressed(pixelgl.KeyR) && !mouseDragging {
			rotation += math.Pi
			flipped = !flipped
		}

		if !flipped {
			imd.SetMatrix(cam)
		} else {
			imd.SetMatrix(cam.Rotated(win.Bounds().Center(), rotation))
		}

		imd.Clear()
		for i, box := range boxes {
			imd.Color = boxesColours[i]
			if len(box) == 2 {
				imd.Push(box[0], box[1])
				imd.Rectangle(5)
			}
			if len(box) == 1 {
				pos := win.MousePosition().Scaled(1 / camZoom).Add(camPos)
				if flipped {
					pos.X = win.Bounds().Max.X - pos.X
					pos.Y = win.Bounds().Max.Y - pos.Y
				}
				imd.Push(box[0], pos)
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

		if win.JustPressed(pixelgl.KeyEnter) {

			// capture box areas
			smn := boxes[0][0].Add(imgOffset).Scaled(1 / imgScale)
			smx := boxes[0][1].Add(imgOffset).Scaled(1 / imgScale)

			imgH := img.Bounds().Dy()
			simg := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(int(smn.X), imgH-int(smn.Y), int(smx.X), imgH-int(smx.Y)))

			// save file
			f, err := os.Create("outimage.png")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			err = png.Encode(f, simg)
			if err != nil {
				panic(err)
			}

			//win.Destroy()
			//return
		}

		win.Clear(colornames.Aliceblue)
		doc.Draw(win, pixel.IM.Moved(win.Bounds().Center()).
			Rotated(win.Bounds().Center(), rotation).
			Scaled(win.Bounds().Center(), imgScale),
		)
		imd.Draw(win)
		win.Update()
	}
}

func ScanManual() {
	pixelgl.Run(run)
}
