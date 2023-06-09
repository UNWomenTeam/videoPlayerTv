package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deepch/vdk/av"
	"github.com/gen2brain/x264-go"
	"github.com/kbinani/screenshot"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const WorkW = "./0_olog/8_wrk_st42_jpgGen_w"

// const WorkQ = "./8_wrk_st42_jpgGen_q"
const WorkQ = "./0_olog/6_q"

const replaceBytesFromJpeg = true

//const replaceBytesFromJpeg = false

func StreamFrames(ch chan *av.Packet) {
	//file, err := os.Create("screen.264.mp4")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	//	os.Exit(1)
	//}
	// open folder pictures [START]
	dir := "/home/xela/Projects/1tv/videos/1685523600_3g"
	chInputJpeg := make(chan image.Image)
	entry, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	// подготовим порядок
	var files []int
	fileNames := make(map[int]string)
	for _, fil := range entry {
		fName := fil.Name()
		fNum := strings.TrimSuffix(fName, ".jpeg")
		fNum = strings.Replace(fNum, ".", "", -1)
		intNum, err := strconv.Atoi(fNum)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, intNum)
		fileNames[intNum] = fName
	}
	sort.Ints(files)
	//порядок файлов джепег готов, откроем файл, отправим в канал
	go func() {
		for _, val := range files {
			fName := filepath.Join(dir, fileNames[val])
			img, err := getImageFromFilePath(fName)
			//dat, err := ioutil.ReadFile(fName)
			if err != nil {
				log.Fatal(err)
			}
			chInputJpeg <- img
		}
	}()
	// сделаем конвертор av

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	flushStatus := make(chan chan bool) //
	go func(b2 *bytes.Buffer) {
		dirBolvank := WorkQ
		entry, err := os.ReadDir(dirBolvank)
		if err != nil {
			log.Fatal(err)
		}
		// подготовим порядок
		var filesBolvank []int
		for _, fil := range entry {
			fNameBolvank := fil.Name()
			intNumBolvank, err := strconv.Atoi(fNameBolvank)
			if err != nil {
				log.Fatal(err)
			}
			filesBolvank = append(filesBolvank, intNumBolvank)
		}
		sort.Ints(filesBolvank)

		irStart := 501
		//irStart := 369
		//irStart := 523
		irEnd := 1000
		//irStart := 1
		//irEnd := 48
		globDur := time.Now()
		//fmt.Println(globDur.Format("00:00:00.000"))
		fmt.Println(globDur)
		cur := 0
		for {
			rer := <-flushStatus
			frm := b2.Bytes()
			// Save to file
			//filePth := "./st42.mp4"
			//file, err := os.OpenFile(filePth, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			//file.Write(frm)
			//file.Close()
			b2.Reset()
			// теперь можно трансформировать кадр в АВ
			//откроем болванку

			fBolvankRaw := cur + irStart
			fmt.Printf("len:%d|file:%d\n", len(frm), fBolvankRaw)
			if fBolvankRaw == irEnd {
				cur = 0
				fmt.Println("Pause 10s on repid")
				time.Sleep(10 * time.Second)
			}
			cur++
			fBolvank := fmt.Sprintf("%d", fBolvankRaw)
			curBolvankFile := filepath.Join(dirBolvank, fBolvank)
			fhBolvank, err := ioutil.ReadFile(curBolvankFile)
			if err != nil {
				log.Fatal(err)
			}
			var containerAV av.Packet
			err = json.Unmarshal(fhBolvank, &containerAV)
			if err != nil {
				log.Fatal(err)
			}
			// залогируем таймауты
			//wrkdir := "./0_olog"
			//fileNameLog := "timerates.log"
			//file, err := os.OpenFile(filepath.Join(wrkdir, fileNameLog), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			//text := fmt.Sprintf("\n[%d][%d]=>Dur:%d|gDur:%d||Time:%d|Idx:%d|Key:%v", cur, fBolvankRaw, containerAV.Duration*time.Millisecond, time.Now().Sub(globDur).Milliseconds(),
			//	containerAV.Time, containerAV.Idx, containerAV.IsKeyFrame)
			//_, err = file.Write([]byte(text))
			//if err != nil {
			//	log.Fatal(err)
			//}
			//err = file.Close()
			//if err != nil {
			//	log.Fatal(err)
			//}

			// заменим байты //83744
			if replaceBytesFromJpeg {
				containerAV.Data = frm
			}
			//containerAV.Data = frm

			//containerAV.Time = 0
			//containerAV.Duration = 0
			//containerAV.Idx = 0
			//containerAV.IsKeyFrame = false
			//if cur == 1 {
			//	containerAV.IsKeyFrame = true
			//}

			//containerAV := av.Packet{
			//	IsKeyFrame:      false,
			//	Idx:             0,
			//	CompositionTime: 0,
			//	Time:            0,
			//	Duration:        0,
			//	Data:            frm,
			//}
			ch <- &containerAV
			rer <- true //оповестим, что можно писать новые байты в буфер
		}
	}(&b)
	// open folder pictures [END]

	//bounds := screenshot.GetDisplayBounds(0)

	opts := &x264.Options{
		Width:     1920,
		Height:    1080,
		FrameRate: 25,
		Tune:      "zerolatency",
		Preset:    "ultrafast",
		Profile:   "high",
		//Profile:   "baseline",
		LogLevel: x264.LogError,
	}

	enc, err := x264.NewEncoder(foo, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	defer enc.Close()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Second / time.Duration(25))

	start := time.Now()
	frame := 0
	cnt := 0
	for range ticker.C {
		select {
		case <-s: // os.exit
			enc.Flush()
			os.Exit(0)
		default:
			frame++
			log.Printf("frame: %v", frame)
			//откроем изображение
			//dc := gg.NewContext(1920, 1080)
			// подготовимся к рисованию
			//сделаем запрос джепега
			imageFrame := <-chInputJpeg
			//reader := bytes.NewReader(rawJpeg)
			//imageFrame, _, err := image.Decode(reader)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//img, err := screenshot.CaptureRect(bounds)
			//if err != nil {
			//	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			//	continue
			//}

			err = enc.Encode(imageFrame)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			log.Printf("t: %v", time.Since(start))
			start = time.Now()
			cnt++
			err = enc.Flush()
			if err != nil {
				log.Fatal(err)
			}
			if cnt > 15 {
				cnt = 0
				rerTo := make(chan bool)
				flushStatus <- rerTo
				<-rerTo
			}

			//if cnt > 100 {
			//	//enc.Flush()
			//
			//	err = file.Close()
			//	if err != nil {
			//		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			//		os.Exit(1)
			//	}
			//
			//	os.Exit(0)
			//}
		}
	}
}
func StreamScreen(ch chan *av.Packet) {
	//file, err := os.Create("screen.264.mp4")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	//	os.Exit(1)
	//}
	// open folder pictures [START]
	dir := "/home/xela/Projects/1tv/videos/1685523600_3g"
	chInputJpeg := make(chan image.Image)
	entry, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	// подготовим порядок
	var files []int
	fileNames := make(map[int]string)
	for _, fil := range entry {
		fName := fil.Name()
		fNum := strings.TrimSuffix(fName, ".jpeg")
		fNum = strings.Replace(fNum, ".", "", -1)
		intNum, err := strconv.Atoi(fNum)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, intNum)
		fileNames[intNum] = fName
	}
	sort.Ints(files)
	//порядок файлов джепег готов, откроем файл, отправим в канал
	go func() {
		for _, val := range files {
			fName := filepath.Join(dir, fileNames[val])
			img, err := getImageFromFilePath(fName)
			//dat, err := ioutil.ReadFile(fName)
			if err != nil {
				log.Fatal(err)
			}
			chInputJpeg <- img
		}
	}()
	// сделаем конвертор av

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	flushStatus := make(chan chan bool) //
	go func(b2 *bytes.Buffer) {
		dirBolvank := "./6_q"
		entry, err := os.ReadDir(dirBolvank)
		if err != nil {
			log.Fatal(err)
		}
		// подготовим порядок
		var filesBolvank []int
		for _, fil := range entry {
			fNameBolvank := fil.Name()
			intNumBolvank, err := strconv.Atoi(fNameBolvank)
			if err != nil {
				log.Fatal(err)
			}
			filesBolvank = append(filesBolvank, intNumBolvank)
		}
		sort.Ints(filesBolvank)

		irStart := 369
		irEnd := 1000
		cur := 0
		for {
			rer := <-flushStatus
			frm := b2.Bytes()
			b2.Reset()
			// теперь можно трансформировать кадр в АВ
			//откроем болванку
			//fmt.Printf("len:%d\n", len(frm))
			fBolvankRaw := cur + irStart
			if fBolvankRaw == irEnd {
				cur = 0
			}
			cur++
			fBolvank := fmt.Sprintf("%d", fBolvankRaw)
			curBolvankFile := filepath.Join(dirBolvank, fBolvank)
			fhBolvank, err := ioutil.ReadFile(curBolvankFile)
			if err != nil {
				log.Fatal(err)
			}
			var containerAV av.Packet
			err = json.Unmarshal(fhBolvank, &containerAV)
			if err != nil {
				log.Fatal(err)
			}
			// заменим байты //83744
			containerAV.Data = frm

			//containerAV := av.Packet{
			//	IsKeyFrame:      false,
			//	Idx:             0,
			//	CompositionTime: 0,
			//	Time:            0,
			//	Duration:        0,
			//	Data:            frm,
			//}
			ch <- &containerAV
			rer <- true //оповестим, что можно писать новые байты в буфер
		}
	}(&b)
	// open folder pictures [END]

	bounds := screenshot.GetDisplayBounds(0)

	opts := &x264.Options{
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
		FrameRate: 25,
		Tune:      "zerolatency",
		Preset:    "ultrafast",
		Profile:   "baseline",
		LogLevel:  x264.LogError,
	}

	enc, err := x264.NewEncoder(foo, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	defer enc.Close()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Second / time.Duration(25))

	start := time.Now()
	frame := 0
	cnt := 0
	for range ticker.C {
		select {
		case <-s: // os.exit
			enc.Flush()
			os.Exit(0)
		default:
			frame++
			log.Printf("frame: %v", frame)
			//откроем изображение
			//dc := gg.NewContext(1920, 1080)
			// подготовимся к рисованию
			//сделаем запрос джепега
			//imageFrame := <-chInputJpeg
			//reader := bytes.NewReader(rawJpeg)
			//imageFrame, _, err := image.Decode(reader)
			//if err != nil {
			//	log.Fatal(err)
			//}
			imageFrame, err := screenshot.CaptureRect(bounds)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				continue
			}

			err = enc.Encode(imageFrame)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			log.Printf("t: %v", time.Since(start))
			start = time.Now()
			cnt++
			err = enc.Flush()
			if err != nil {
				log.Fatal(err)
			}
			rerTo := make(chan bool)
			flushStatus <- rerTo
			<-rerTo
			//if cnt > 100 {
			//	//enc.Flush()
			//
			//	err = file.Close()
			//	if err != nil {
			//		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			//		os.Exit(1)
			//	}
			//
			//	os.Exit(0)
			//}
		}
	}
}
func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}
