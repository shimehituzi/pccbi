package main

import (
	"fmt"
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
	frame := binding.BindFloat(&f)
	l := 0
	// labeling := binding.BindInt(&l)
	scale := 3
	size := fyne.NewSize(float32(dim.D2)*float32(scale), float32(dim.D1)*float32(scale))

	raster := canvas.NewRaster(func(w, h int) image.Image {
		return fbm.GetImage(int(f), l)
	})
	raster.ScaleMode = canvas.ImageScalePixels
	raster.Resize(size)

	labelLength := fbm.GetLabelLength(int(f))
	labelingOptions := []string{"All"}
	for i := 1; i <= labelLength; i++ {
		labelingOptions = append(labelingOptions, fmt.Sprint(i))
	}
	labelingRadio := widget.NewRadioGroup(labelingOptions, func(s string) {

	})
	labelingRadio.Required = true
	labelingRadio.Selected = "All"
	labelingRadio.Horizontal = true

	frameSlider := widget.NewSliderWithData(0, float64(dim.D0-1), frame)
	frameSlider.OnChanged = func(f float64) {
		err := frame.Set(f)
		if err != nil {
			fyne.LogError("Failed to set binding value", err)
		}
		raster.Refresh()
		labelLength := fbm.GetLabelLength(int(f))
		labelingOptions := []string{"All"}
		for i := 1; i <= labelLength; i++ {
			labelingOptions = append(labelingOptions, fmt.Sprint(i))
		}
		labelingRadio.Options = labelingOptions
		labelingRadio.Refresh()
	}

	frameLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(frame, "frame: %0.0f"))

	bitmapContent := container.New(layout.NewGridWrapLayout(size), raster)
	content := container.New(
		layout.NewVBoxLayout(),
		bitmapContent, layout.NewSpacer(),
		frameSlider, layout.NewSpacer(),
		frameLabel, layout.NewSpacer(),
		labelingRadio, layout.NewSpacer(),
	)
	w.SetContent(content)
	w.ShowAndRun()
}
