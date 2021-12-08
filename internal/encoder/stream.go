package encoder

func NewStream(cb contourBuffer) *Stream {
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

	stream := Stream{
		StartPoints: startPoints,
		Flags:       flags,
		Codes:       codes,
	}

	return &stream
}
