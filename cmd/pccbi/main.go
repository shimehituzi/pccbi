package main

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	srcfile := relativePath[4:]
	bms := LordPly(srcfile)
	FyneLoop(bms)
}
