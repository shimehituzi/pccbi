package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/shimehituzi/pccbi/internal/bitmap"
)

func FyneLoop(fbm [2]bitmap.FyneBitMap) {
	myApp := app.New()
	w := myApp.NewWindow("BitMap")

	dim := fbm[0].GetLength()
	f := 0.0
	frame := binding.BindFloat(&f)
	l := 0
	labeling := binding.BindInt(&l)
	scale := 1.5
	size := fyne.NewSize(float32(dim.D2)*float32(scale), float32(dim.D1)*float32(scale))

	labeledRaster := canvas.NewRaster(func(w, h int) image.Image {
		return fbm[0].GetImage(int(f), l)
	})
	labeledRaster.ScaleMode = canvas.ImageScalePixels
	labeledRaster.Resize(size)

	bitmapRaster := canvas.NewRaster(func(w, h int) image.Image {
		return fbm[1].GetImage(int(f), l)
	})
	bitmapRaster.ScaleMode = canvas.ImageScalePixels
	bitmapRaster.Resize(size)

	labelLength := fbm[0].GetLabelLength(int(f))
	labelingOptions := []string{"All"}
	for i := 1; i <= labelLength; i++ {
		labelingOptions = append(labelingOptions, fmt.Sprint(i))
	}
	labelingRadio := widget.NewRadioGroup(labelingOptions, func(s string) {
		var (
			i   int
			err error
		)
		if s != "All" {
			i, err = strconv.Atoi(s)
			if err != nil {
				fyne.LogError("Failed to convert string", err)
			}
		} else {
			i = 0
		}
		if err = labeling.Set(i); err != nil {
			fyne.LogError("Failed to set binding value", err)
		}
		labeledRaster.Refresh()
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

		labelLength := fbm[0].GetLabelLength(int(f))
		labelingOptions := []string{"All"}
		for i := 1; i <= labelLength; i++ {
			labelingOptions = append(labelingOptions, fmt.Sprint(i))
		}
		labelingRadio.Options = labelingOptions
		labelingRadio.Selected = "All"
		if err := labeling.Set(0); err != nil {
			fyne.LogError("Failed to set binding value", err)
		}
		labelingRadio.Refresh()

		labeledRaster.Refresh()
		bitmapRaster.Refresh()
	}

	frameLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(frame, "frame: %0.0f"))

	labeledRasterContent := container.New(layout.NewGridWrapLayout(size), labeledRaster)
	line := canvas.NewLine(color.Opaque)
	line.Position1 = fyne.NewPos(0, 0)
	line.Position2 = fyne.NewPos(0, float32(dim.D1)*float32(scale))
	bitmapRasterContent := container.New(layout.NewGridWrapLayout(size), bitmapRaster)
	rasterContent := container.New(
		layout.NewHBoxLayout(),
		labeledRasterContent,
		line,
		bitmapRasterContent,
	)

	content := container.New(
		layout.NewVBoxLayout(),
		rasterContent, layout.NewSpacer(),
		frameSlider, layout.NewSpacer(),
		frameLabel, layout.NewSpacer(),
		labelingRadio, layout.NewSpacer(),
	)
	w.SetContent(content)
	w.ShowAndRun()
}
