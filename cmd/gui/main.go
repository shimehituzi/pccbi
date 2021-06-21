package main

func main() {
	srcfile := "./DATABASE/orig/loot/loot/Ply/loot_vox10_1000.ply"
	bms := LordPly(srcfile)
	FyneLoop(bms)
}
