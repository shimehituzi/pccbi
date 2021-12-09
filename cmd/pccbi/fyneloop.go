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
	"github.com/shimehituzi/pccbi/internal/codec"
)

func fyneLoop(fbm []codec.FyneBitMap) {
	myApp := app.New()
	w := myApp.NewWindow("BitMap")

	dim := fbm[0].GetLength()
	f := 0.0
	frame := binding.BindFloat(&f)
	scale := 2
	size := fyne.NewSize(float32(dim.D2)*float32(scale), float32(dim.D1)*float32(scale))

	l := 0
	labeling := binding.BindInt(&l)

	raster1 := canvas.NewRaster(func(w, h int) image.Image {
		return fbm[0].GetImage(int(f), 0)
	})
	raster1.ScaleMode = canvas.ImageScalePixels
	raster1.Resize(size)
	raster2 := canvas.NewRaster(func(w, h int) image.Image {
		return fbm[1].GetImage(int(f), l)
	})
	raster2.ScaleMode = canvas.ImageScalePixels
	raster2.Resize(size)

	labelLength := fbm[1].GetLabelLength(int(f))
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
		raster2.Refresh()
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

		labelLength := fbm[1].GetLabelLength(int(f))
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

		raster1.Refresh()
		raster2.Refresh()
	}

	frameLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(frame, "frame: %0.0f"))

	raster1Content := container.New(layout.NewGridWrapLayout(size), raster1)
	raster2Content := container.New(layout.NewGridWrapLayout(size), raster2)
	line := canvas.NewLine(color.Opaque)
	line.Position1 = fyne.NewPos(0, 0)
	hbox := container.New(layout.NewHBoxLayout(), raster2Content, line, raster1Content)

	content := container.New(
		layout.NewVBoxLayout(),
		hbox, layout.NewSpacer(),
		frameSlider, layout.NewSpacer(),
		frameLabel, layout.NewSpacer(),
		labelingRadio, layout.NewSpacer(),
	)
	w.SetContent(content)
	w.ShowAndRun()
}
