package quac

import (
	"image"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Manual Scan In...",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	imd.Color = colornames.Red

	pic, err := loadPicture("apple.jpeg")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())

	var (
		boxes        [][]pixel.Vec
		camPos       = pixel.ZV
		camSpeed     = 100.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		rotation     = 0.0
		flipped      = false
	)

	for !win.Closed() {

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Min.Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			pos := win.MousePosition().Scaled(1 / camZoom).Add(camPos)
			if flipped {
				pos.X = win.Bounds().Max.X - pos.X
				pos.Y = win.Bounds().Max.Y - pos.Y
			}
			boxes = append(boxes, []pixel.Vec{pos})
		}
		if win.JustReleased(pixelgl.MouseButtonLeft) {
			pos := win.MousePosition().Scaled(1 / camZoom).Add(camPos)
			if flipped {
				pos.X = win.Bounds().Max.X - pos.X
				pos.Y = win.Bounds().Max.Y - pos.Y
			}
			boxes[len(boxes)-1] = append(boxes[len(boxes)-1], pos)
		}

		if win.JustPressed(pixelgl.KeyR) {
			rotation += math.Pi
			flipped = !flipped
		}

		if !flipped {
			imd.SetMatrix(cam)
		} else {
			imd.SetMatrix(cam.Rotated(win.Bounds().Center(), rotation))
		}

		imd.Clear()
		for _, box := range boxes {
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

		if win.JustPressed(pixelgl.KeyLeft) {
			camPos.X += camSpeed
		}
		if win.JustPressed(pixelgl.KeyRight) {
			camPos.X -= camSpeed
		}
		if win.JustPressed(pixelgl.KeyDown) {
			camPos.Y += camSpeed
		}
		if win.JustPressed(pixelgl.KeyUp) {
			camPos.Y -= camSpeed
		}
		if win.JustPressed(pixelgl.KeyEqual) {
			camZoom *= 1.2
		}
		if win.JustPressed(pixelgl.KeyMinus) {
			camZoom /= 1.2
		}

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		win.Clear(colornames.Aliceblue)
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Rotated(win.Bounds().Center(), rotation))
		imd.Draw(win)
		win.Update()
	}
}

func ScanManual() {
	pixelgl.Run(run)
}
