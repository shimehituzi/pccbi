package codec

func NewFrames(voxel Voxel, header *Header) frames {
	lv, numLabels := NewLabels(voxel, header)

	frames := make(frames, header.Length[0])
	for i := range frames {
		frames[i] = make(frame, numLabels[i])
		for j := range frames[i] {
			frames[i][j] = make(segment, header.Length[1])
			for k := range frames[i][j] {
				frames[i][j][k] = make([]byte, header.Length[2])
			}
		}
	}

	for f := range lv {
		for y := range lv[f] {
			for x, label := range lv[f][y] {
				if label != 0 {
					l := label - 1
					frames[f][l][y][x] = 1
				}
			}
		}
	}

	return frames
}

func NewContourBuffer(voxel Voxel, header *Header) contourBuffer {
	frames := NewFrames(voxel, header)

	cb := make(contourBuffer, len(frames))
	for i := range cb {
		cb[i] = make([]contour, len(frames[i]))
	}

	for f, frame := range frames {
		for l, img := range frame {
			cb[f][l] = newContour(bitmap(img))
		}
	}

	return cb
}

func NewFyneContour(cb contourBuffer, header *Header) labeledVoxel {
	fc := make(labeledVoxel, header.Length[0])
	for f := range fc {
		fc[f] = make(label, header.Length[1])
		for y := range fc[f] {
			fc[f][y] = make([]int, header.Length[2])
		}
	}

	for f := range cb {
		for l := range cb[f] {
			for c, cc := range cb[f][l] {
				if c == 0 {
					for _, point := range cc.points {
						fc[f][point.y][point.x] = l*2 + 1
					}
				} else {
					for _, point := range cc.points {
						fc[f][point.y][point.x] = l*2 + 2
					}
				}
			}
		}
	}

	return fc
}
