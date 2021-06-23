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

	"github.com/shimehituzi/pccbi/internal/bitmap"
)

func FyneLoop(fbm bitmap.FyneBitMap) {
	myApp := app.New()
	w := myApp.NewWindow("BitMap")

	dim := fbm.GetLength()
	f := 0.0
	data := binding.BindFloat(&f)
	scale := 3
	size := fyne.NewSize(float32(dim.D2)*float32(scale), float32(dim.D1)*float32(scale))

	raster := canvas.NewRaster(func(w, h int) image.Image {
		return fbm.GetImage(int(f))
	})
	raster.ScaleMode = canvas.ImageScalePixels
	raster.Resize(size)

	slider := widget.NewSliderWithData(0, float64(dim.D0-1), data)
	slider.OnChanged = func(f float64) {
		err := data.Set(f)
		if err != nil {
			fyne.LogError("Failed to set binding value", err)
		}
		raster.Refresh()
	}

	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "frame: %0.0f"))

	bitmapContent := container.New(layout.NewGridWrapLayout(size), raster)
	content := container.New(layout.NewVBoxLayout(), bitmapContent, slider, layout.NewSpacer(), label)
	w.SetContent(content)
	w.ShowAndRun()
}
