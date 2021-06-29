package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shimehituzi/pccbi/internal/processing"
)

func fyneLoop(fbm []processing.FyneBitMap) {
	myApp := app.New()
	w := myApp.NewWindow("BitMap")

	dim := fbm[0].GetLength()
	f := 0.0
	frame := binding.BindFloat(&f)
	scale := 3
	size := fyne.NewSize(float32(dim.D2)*float32(scale), float32(dim.D1)*float32(scale))

	raster := canvas.NewRaster(func(w, h int) image.Image {
		return fbm[0].GetImage(int(f))
	})
	raster.ScaleMode = canvas.ImageScalePixels
	raster.Resize(size)
	raster2 := canvas.NewRaster(func(w, h int) image.Image {
		return fbm[1].GetImage(int(f))
	})
	raster2.ScaleMode = canvas.ImageScalePixels
	raster2.Resize(size)

	frameSlider := widget.NewSliderWithData(0, float64(dim.D0-1), frame)
	frameSlider.OnChanged = func(f float64) {
		err := frame.Set(f)
		if err != nil {
			fyne.LogError("Failed to set binding value", err)
		}

		raster.Refresh()
		raster2.Refresh()
	}

	frameLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(frame, "frame: %0.0f"))

	rasterContent := container.New(layout.NewGridWrapLayout(size), raster)
	raster2Content := container.New(layout.NewGridWrapLayout(size), raster2)
	line := canvas.NewLine(color.Opaque)
	line.Position1 = fyne.NewPos(0, 0)
	line.Position2 = fyne.NewPos(0, float32(dim.D1)*float32(scale))
	hbox := container.New(layout.NewHBoxLayout(), rasterContent, line, raster2Content)

	content := container.New(
		layout.NewVBoxLayout(),
		hbox, layout.NewSpacer(),
		frameSlider, layout.NewSpacer(),
		frameLabel, layout.NewSpacer(),
	)
	w.SetContent(content)
	w.ShowAndRun()
}
