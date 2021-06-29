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

func fyneLoop(fbm refactoring.FyneBitMap) {
	myApp := app.New()
	w := myApp.NewWindow("BitMap")

	dim := fbm.GetLength()
	f := 0.0
	frame := binding.BindFloat(&f)
	scale := 3
	size := fyne.NewSize(float32(dim.D2)*float32(scale), float32(dim.D1)*float32(scale))

	bitmapRaster := canvas.NewRaster(func(w, h int) image.Image {
		return fbm.GetImage(int(f))
	})
	bitmapRaster.ScaleMode = canvas.ImageScalePixels
	bitmapRaster.Resize(size)

	frameSlider := widget.NewSliderWithData(0, float64(dim.D0-1), frame)
	frameSlider.OnChanged = func(f float64) {
		err := frame.Set(f)
		if err != nil {
			fyne.LogError("Failed to set binding value", err)
		}

		bitmapRaster.Refresh()
	}

	frameLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(frame, "frame: %0.0f"))

	bitmapRasterContent := container.New(layout.NewGridWrapLayout(size), bitmapRaster)

	content := container.New(
		layout.NewVBoxLayout(),
		bitmapRasterContent, layout.NewSpacer(),
		frameSlider, layout.NewSpacer(),
		frameLabel, layout.NewSpacer(),
	)
	w.SetContent(content)
	w.ShowAndRun()
}
