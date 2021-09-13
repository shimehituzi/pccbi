package processing

func NewStream(voxel *voxel, cb contourBuffer) *Stream {
	numOuter := 0
	numInner := 0

	outerStartPoints := []int{}
	outerCodes := []byte{}
	innerStartPoints := []int{}
	innerCodes := []byte{}

	for f := range cb {
		for _, contour := range cb[f] {
			for i, v := range contour {
				if i == 0 {
					outerStartPoints = append(outerStartPoints, v.start.x, v.start.y)
					numOuter += 2
					outerCodes = append(outerCodes, v.code...)
					outerCodes = append(outerCodes, 8)
				} else {
					innerStartPoints = append(innerStartPoints, v.start.x, v.start.y)
					numInner += 2
					innerCodes = append(innerCodes, v.code...)
					innerCodes = append(innerCodes, 8)
				}
			}
		}
	}

	header := []int{
		voxel.header.axis[0],   //axisX
		voxel.header.axis[1],   //axisY
		voxel.header.axis[2],   //axisZ
		voxel.header.length[0], //frames
		voxel.header.length[1], //height
		voxel.header.length[2], //width
		voxel.header.bias[0],   //biasFrames
		voxel.header.bias[1],   //biasHeight
		voxel.header.bias[2],   //biasWidth
		numOuter,
		numInner,
		len(outerCodes),
		len(innerCodes),
	}

	stream := Stream{
		Header:           header,
		OuterStartPoints: outerStartPoints,
		InnerStartPoints: innerStartPoints,
		OuterCodes:       outerCodes,
		InnerCodes:       innerCodes,
	}

	return &stream
}
