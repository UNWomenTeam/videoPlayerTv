package main

import (
	"io/ioutil"
	"log"
)

const PORT = 8080

const fileh264 = "/home/xela/Projects/1tv/go-base/public/demo.h264"

type H2O struct {
	minNaluPerChunk, interval, current, start, end int
	wss                                            struct{}
}

func main() {
	var (
		buffer []byte
		chunks []int
		total  int
	)

	h20 := &H2O{}

	buffer, err := ioutil.ReadFile(fileh264)
	if err != nil {
		log.Fatal(err)
	}
	chunks = h20.extractChunks(buffer)
	total = len(chunks)

}

func (h *H2O) extractChunks(buffer []byte) []int {
	var (
		i         int
		length    int
		naluCount int
		value     byte
		unit      []byte
		ntype     byte
		state     int
		lastIndex int
		result    []int
	)
	length = len(buffer)

	for i < length {
		//fmt.Println(i)
		i++
		value = buffer[i]
		// finding 3 or 4-byte start codes (00 00 01 OR 00 00 00 01)
		switch state {
		case 0:
			if value == 0 {
				state = 1
			}
			break
		case 1:
			if value == 0 {
				state = 2
			} else {
				state = 0
			}
			break
		case 2:
		case 3:
			if value == 0 {
				state = 3
			} else if value == 1 && i < length {
				if lastIndex == 0 {
					unit = buffer[lastIndex : i-state-1]
					ntype = unit[0] & 0x1f
					naluCount++
				}
				if naluCount >= h.minNaluPerChunk && ntype != 1 && ntype != 5 {
					pVal := lastIndex - state - 1
					result = append(result, pVal)
					naluCount = 0
				}
				state = 0
				lastIndex = i
			} else {
				state = 0
			}
			break
		default:
			break
		}
	}
	if naluCount > 0 {
		result = append(result, lastIndex)
	}
	return result
}
