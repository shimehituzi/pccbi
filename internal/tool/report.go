package tool

import (
	"fmt"
	"time"
)

func Report(isLossless bool, bits int, numPoints int, times [6]time.Time) {
	if !isLossless {
		return
	}
	allTime := times[5].Sub(times[0]).Seconds()
	encTime := times[2].Sub(times[1]).Seconds()
	decTime := times[4].Sub(times[3]).Seconds()
	fmt.Println("=====================================")
	fmt.Printf("     all time\t: %10.2f sec.\n", allTime)
	fmt.Printf("  encode time\t: %10.2f sec.\n", encTime)
	fmt.Printf("  decode time\t: %10.2f sec.\n", decTime)
	fmt.Println("-------------------------------------")
	fmt.Printf("   total bits\t: %10d bits\n", bits)
	fmt.Printf("   num points\t: %10d points\n", numPoints)
	fmt.Printf("  coding rate\t: %10.5f b/p\n", float64(bits)/float64(numPoints))
	fmt.Println("=====================================")
}
