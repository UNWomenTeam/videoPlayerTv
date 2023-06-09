package api

import (
	"fmt"
	"io/ioutil"
	"log"
)

const PORT = 8080

const fileh264 = "/home/xela/Projects/1tv/go-base/public/demo.h264"

type WsPlug struct {
	buffer []byte
	chunks []int
	total  int
	h20    *H2O
	serve  chan []byte
}

type H2O struct {
	minNaluPerChunk, interval, current, start, end int
	wss                                            struct{}
}

func NewWsPlug(input chan []byte) *WsPlug {
	h20 := &H2O{minNaluPerChunk: 30}

	buffer, err := ioutil.ReadFile(fileh264)
	if err != nil {
		log.Fatal(err)
	}
	chunks := h20.extractChunks(buffer)
	total := len(chunks)
	return &WsPlug{
		buffer: buffer,
		chunks: chunks,
		total:  total,
		h20:    h20,
		serve:  input,
	}
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
	three := func() {
		if value == 0 {
			state = 3
		} else if value == 1 && i < length {
			if lastIndex != 0 {
				unit = buffer[lastIndex : i-state-1]
				if len(unit) == 0 {
					log.Fatal("eerr nil")
				}
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
	}

	length = len(buffer)

	for i < length {
		//fmt.Println(i)

		//fmt.Printf("is i:%d\n", i)

		value = buffer[i] //0
		i++
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
			//fmt.Println("222222222")
			three()
		case 3:
			three()
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

func (h *WsPlug) sendChunk() {
	if h.h20.current >= h.total {
		h.h20.current = 0
		h.h20.start = 0
	}
	if len(h.chunks) < 1 {
		//fmt.Println("Hulled chanks 321_--")
		return
	} else {
		fmt.Println("[]SEND[OK]")
	}
	h.h20.end = h.chunks[h.h20.current]
	h.h20.current++

	var chunk []byte

	chunk = h.buffer[h.h20.start:h.h20.end]
	h.h20.start = h.h20.end
	fmt.Printf("sendet:%d\n", len(chunk))
	h.serve <- chunk
}
