package main

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shimehituzi/pccbi/internal/refactoring"
)

func fyneLoop(fbm []refactoring.FyneBitMap) {
	myApp := app.New()
	w := myApp.NewWindow("BitMap")

	dim := fbm[0].GetLength()
	f := 0.0
	frame := binding.BindFloat(&f)
	scale := 1.5
	size := fyne.NewSize(float32(dim.D2)*float32(scale), float32(dim.D1)*float32(scale))

	raster := canvas.NewRaster(func(w, h int) image.Image {
		return fbm[0].GetImage(int(f))
	})
	raster.ScaleMode = canvas.ImageScalePixels
	raster.Resize(size)

	frameSlider := widget.NewSliderWithData(0, float64(dim.D0-1), frame)
	frameSlider.OnChanged = func(f float64) {
		err := frame.Set(f)
		if err != nil {
			fyne.LogError("Failed to set binding value", err)
		}

		raster.Refresh()
	}

	frameLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(frame, "frame: %0.0f"))

	raster1Content := container.New(layout.NewGridWrapLayout(size), raster)

	content := container.New(
		layout.NewVBoxLayout(),
		raster1Content, layout.NewSpacer(),
		frameSlider, layout.NewSpacer(),
		frameLabel, layout.NewSpacer(),
	)
	w.SetContent(content)
	w.ShowAndRun()
}
