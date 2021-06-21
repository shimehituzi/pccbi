package labeling

type ChainCode []byte

type Countour struct {
	ChainCode ChainCode
	Label     int
}

type Segment struct {
	Countour []Countour
	Label    int
}

type Segments []Segment
