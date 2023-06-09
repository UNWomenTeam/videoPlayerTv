package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	entry, err := os.ReadDir("./q")
	if err != nil {
		log.Fatal(err)
	}
	var files []int
	for _, fil := range entry {
		fName := fil.Name()
		fNum := strings.TrimSuffix(fName, ".jpeg")
		intNum, err := strconv.Atoi(fNum)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, intNum)
	}
	sort.Ints(files)
	fmt.Println(entry)
}
