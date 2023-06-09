package main

import (
	"encoding/json"
	"fmt"
	"github.com/deepch/vdk/av"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

func Repid(chFarame chan *av.Packet) {
	frameDirName := "6_q"
	fldrFiles, err := os.ReadDir(fmt.Sprintf("./%s", frameDirName))
	if err != nil {
		log.Fatal(err)
	}
	var files []int
	for _, fil := range fldrFiles {
		fName := fil.Name()
		//fNum := strings.TrimSuffix(fName, ".jpeg")
		intNum, err := strconv.Atoi(fName)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, intNum)
	}
	sort.Ints(files)
	for {
		for _, rawPacketFname := range files {
			//fmt.Println("=====>Hello ws!files+++[1.READ] ")
			fPackName := fmt.Sprintf("./%s/%d", frameDirName, rawPacketFname)
			datPack, err := ioutil.ReadFile(fPackName)
			if err != nil {
				log.Fatal(err)
			}
			var pack av.Packet
			err = json.Unmarshal(datPack, &pack)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println("=====>Hello ws!files+++[2.WRITE_WAIT_] ")
			chFarame <- &pack
			//fmt.Println("=====>Hello ws!files+++[3.WRITE_OK_] ")
			time.Sleep(100 * time.Microsecond)

		}
	}
}
