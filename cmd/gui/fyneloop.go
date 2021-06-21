package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"

	"github.com/shimehituzi/pccbi/internal/plyio"
)

func FyneLoop(bms *plyio.BitMaps) {
	myApp := app.New()
	w := myApp.NewWindow("Raster")

	// 7 < index < 999 で DATABASE/pbm 画像の番号と同じになる
	raster := canvas.NewRasterFromImage(bms.Data[606-bms.Bias[0]])
	w.SetContent(raster)
	w.Resize(fyne.NewSize(float32(bms.Length[2]), float32(bms.Length[1])))
	w.ShowAndRun()
}
