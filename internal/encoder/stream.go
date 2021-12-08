package encoder

func NewStream(cb contourBuffer, voxelHeader *VoxelHeader) (*Stream, StreamHeader) {
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
		voxelHeader.axis[0],   //axisX
		voxelHeader.axis[1],   //axisY
		voxelHeader.axis[2],   //axisZ
		voxelHeader.length[0], //frames
		voxelHeader.length[1], //height
		voxelHeader.length[2], //width
		voxelHeader.bias[0],   //biasFrames
		voxelHeader.bias[1],   //biasHeight
		voxelHeader.bias[2],   //biasWidth
		len(startPoints),
		len(flags),
		len(codes),
	}

	stream := Stream{
		StartPoints: startPoints,
		Flags:       flags,
		Codes:       codes,
	}

	return &stream, header
}
