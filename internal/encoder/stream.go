package encoder

func NewStream(voxel *voxel, cb contourBuffer) *Stream {
	startPoints := []int{}
	flags := []byte{}
	codes := []byte{}

	for f := range cb {
		for _, contour := range cb[f] {
			for i, v := range contour {
				startPoints = append(startPoints, f, v.start.y, v.start.x)
				codes = append(codes, v.code...)
				codes = append(codes, 8)
				if i == 0 {
					flags = append(flags, 0)
				} else {
					flags = append(flags, 1)
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
		len(startPoints),
		len(flags),
		len(codes),
	}

	stream := Stream{
		Header:      header,
		StartPoints: startPoints,
		Flags:       flags,
		Codes:       codes,
	}

	return &stream
}
