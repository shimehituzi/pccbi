package tool

import "sort"

func Preprocessing(srcPath, sortedPath, etcPath string) int {
	data, header := readPly(srcPath)

	sort.Sort(data)

	writePly(sortedPath, data, header)

	etcData := make(plyData, len(data))
	for i := range etcData {
		etcData[i] = data[i][3:]
	}
	writePly(etcPath, etcData, header)

	return len(data)
}

func (data plyData) Len() int { return len(data) }

func (data plyData) Swap(i, j int) { data[i], data[j] = data[j], data[i] }

func (data plyData) Less(i, j int) bool {
	switch {
	case data[i][0] < data[j][0]:
		return true
	case data[i][0] > data[j][0]:
		return false
	default:
		switch {
		case data[i][1] < data[j][1]:
			return true
		case data[i][1] > data[j][1]:
			return false
		default:
			switch {
			case data[i][2] < data[j][2]:
				return true
			case data[i][2] > data[j][2]:
				return false
			default:
				return false
			}
		}
	}
}
