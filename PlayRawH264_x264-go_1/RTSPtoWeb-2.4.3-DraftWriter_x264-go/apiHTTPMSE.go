package main

import (
	"encoding/json"
	"fmt"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/mp4f"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//const codakPth = "./4_wrk_w/1_codecs.json"

//const codakPth = "./8_wrk_st42_jpgGen_w/1_codecs.json"

const codakPth = "./0_olog/6_w/1_codecs.json"
const saveCodecsToFile = false
const saveFragmeToFile = false

var cadr500part int
var wrtr500part int
var wrtr1part int
var wrtr1bath int
var inpun int
var inpkd int
var repid bool

// HTTPAPIServerStreamMSE func
func HTTPAPIServerStreamMSE(c *gin.Context) {
	repid = true //из файлов

	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		return
	}
	ch2 := make(chan *av.Packet)
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_mse",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamMSE",
	})

	defer func() {
		err = conn.Close()
		requestLogger.WithFields(logrus.Fields{
			"call": "Close",
		}).Errorln(err)
		log.Println("Client Full Exit")
	}()
	fmt.Println("=====>Hello ws! HTTPAPIServerStreamMSE _UpgradeHTTP FUNC_")
	if !repid {
		if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
			requestLogger.WithFields(logrus.Fields{
				"call": "StreamChannelExist",
			}).Errorln(ErrorStreamNotFound.Error())
			return
		}
		fmt.Println("=====>Hello ws! 1. RemoteAuthorization _UpgradeHTTP FUNC_")
		if !RemoteAuthorization("WS", c.Param("uuid"), c.Param("channel"), c.Query("token"), c.ClientIP()) {
			requestLogger.WithFields(logrus.Fields{
				"call": "RemoteAuthorization",
			}).Errorln(ErrorStreamUnauthorized.Error())
			return
		}
		fmt.Println("=====>Hello ws! 2. _UpgradeHTTP FUNC_")
		Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
		fmt.Println("=====>Hello ws! 3. _UpgradeHTTP FUNC_")
		err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "SetWriteDeadline",
			}).Errorln(err.Error())
			return
		}
		fmt.Println("=====>Hello ws! 4. _UpgradeHTTP FUNC_")
		cid, ch, _, err := Storage.ClientAdd(c.Param("uuid"), c.Param("channel"), MSE)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "ClientAdd",
			}).Errorln(err.Error())
			return
		}
		fmt.Println("=====>Hello ws! 5. _UpgradeHTTP FUNC_")
		defer func() {
			fmt.Println("=====>Hello ws! 5. _EXIT DEFER_")
			Storage.ClientDelete(c.Param("uuid"), cid, c.Param("channel"))
		}()
		codecs, err := Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "StreamCodecs",
			}).Errorln(err.Error())
			return
		}
		// save codecs in file _START_
		if saveCodecsToFile {
			cdkDat, err := json.Marshal(codecs)
			if err != nil {
				log.Fatal(err)
			}
			inpkd++
			fileSave := fmt.Sprintf("./w/%d_codecs.json", inpkd)
			f, err := os.Create(fileSave)
			if err != nil {
				log.Fatal(err)
			}
			f.Write(cdkDat)
			f.Close()

			fileSave2 := fmt.Sprintf("./w/tips_%d_codecs.json", inpkd)
			f2, err := os.Create(fileSave2)
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range codecs {
				f2.Write([]byte(fmt.Sprintf("{tip:\"%s\"},\n", v.Type())))
			}
			f2.Close()
		}

		// save codecs in file _END_
		fmt.Println("=====>Hello ws! 6. _UpgradeHTTP FUNC_")
		muxerMSE := mp4f.NewMuxer(nil)
		fmt.Println("=====>Hello ws! 7. _UpgradeHTTP FUNC_")
		err = muxerMSE.WriteHeader(codecs)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "WriteHeader",
			}).Errorln(err.Error())
			return
		}
		fmt.Println("=====>Hello ws! 8. _UpgradeHTTP FUNC_")
		meta, init := muxerMSE.GetInit(codecs)
		fmt.Println("=====>Hello ws! 9. _UpgradeHTTP FUNC_")
		err = wsutil.WriteServerMessage(conn, ws.OpBinary, append([]byte{9}, meta...))
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "Send",
			}).Errorln(err.Error())
			return
		}
		fmt.Println("=====>Hello ws! 10. _UpgradeHTTP FUNC_")
		err = wsutil.WriteServerMessage(conn, ws.OpBinary, init)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "Send",
			}).Errorln(err.Error())
			return
		}

		fmt.Println("=====>Hello ws!videoStart+++ ")
		var videoStart bool
		controlExit := make(chan bool, 10)
		noClient := time.NewTimer(10 * time.Second)
		go func() {
			defer func() {
				fmt.Println("=====>Hello ws!controlExit <- true+++ ")
				controlExit <- true
				fmt.Println("=====>Hello ws!controlExit <- true+++ OK")
			}()
			for { // отладочные, диагностические функции
				header, _, err := wsutil.NextReader(conn, ws.StateServerSide)
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "Receive",
					}).Errorln(err.Error())
					return
				}
				//fmt.Println("=====>Hello ws!Receive <-  wsutil.NextReader+++ ")
				switch header.OpCode {
				case ws.OpPong:
					noClient.Reset(10 * time.Second)
					//fmt.Println("=====>Hello ws!Receive <-  wsutil.NextReader+++OpPong ")
				case ws.OpClose:
					fmt.Println("=====>Hello ws!Receive <-  wsutil.NextReader---OpClose ")
					return
				}
			}
		}()
		noVideo := time.NewTimer(10 * time.Second)
		pingTicker := time.NewTicker(500 * time.Millisecond)
		defer pingTicker.Stop()
		defer log.Println("client exit")

		var logxtm int
		//var position int
		if saveFragmeToFile {
			var files []int
			entry, err := os.ReadDir("./q")
			if err != nil {
				log.Fatal(err)
			}
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
		}

		fmt.Println("=====>Hello ws!Receive <-  for---START ")
		for {
			select {

			case <-pingTicker.C:
				//fmt.Println("=====>Hello ws!Receive <-  pingTicker.C ")
				err = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
				if err != nil {
					return
				}
				buf, err := ws.CompileFrame(ws.NewPingFrame(nil))
				if err != nil {
					return
				}
				_, err = conn.Write(buf)
				if err != nil {
					return
				}
			case <-controlExit:
				fmt.Println("=====>Hello ws!Receive <-  controlExit ")
				requestLogger.WithFields(logrus.Fields{
					"call": "controlExit",
				}).Errorln("Client Reader Exit")
				return
			case <-noClient.C:
				fmt.Println("=====>Hello ws!Receive <-  noClient.C ")
				requestLogger.WithFields(logrus.Fields{
					"call": "ErrorClientOffline",
				}).Errorln("Client OffLine Exit")
				return
			case <-noVideo.C:
				fmt.Println("=====>Hello ws!Receive <-  noVideo.C ")
				requestLogger.WithFields(logrus.Fields{
					"call": "ErrorStreamNoVideo",
				}).Errorln(ErrorStreamNoVideo.Error())
				return
			case pck := <-ch:
				logxtm++
				if logxtm > 500 {
					fmt.Println("=====>Hello ws!Receive <-ch logxtm _500кадров._ ")
					logxtm = 0
				}
				if saveFragmeToFile {
					inpun++
					f, err := os.Create(fmt.Sprintf("./q/%d", inpun))
					if err != nil {
						log.Fatal(err)
					}
					inpJson, err := json.Marshal(pck)
					if err != nil {
						log.Fatal(err)
					}
					f.Write(inpJson)
					f.Close()
				}

				//fName := fmt.Sprintf("./q/%d.jpeg", files[position])
				//position++
				//dat, err := ioutil.ReadFile(fName)
				//if err != nil {
				//	log.Fatal(err)
				//}
				//pck.Data = dat

				//fmt.Println("New frame!")
				if pck.IsKeyFrame {
					noVideo.Reset(10 * time.Second)
					videoStart = true
				}
				if !videoStart {
					continue
				}
				ready, buf, err := muxerMSE.WritePacket(*pck, false)
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "WritePacket",
					}).Errorln(err.Error())
					return
				}
				if ready {
					err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
					if err != nil {
						requestLogger.WithFields(logrus.Fields{
							"call": "SetWriteDeadline",
						}).Errorln(err.Error())
						return
					}
					//err = websocket.Message.Send(ws, buf)
					err = wsutil.WriteServerMessage(conn, ws.OpBinary, buf)
					if err != nil {
						requestLogger.WithFields(logrus.Fields{
							"call": "Send",
						}).Errorln(err.Error())
						return
					}
				}
			}
		}
	}
	// START MODERN
	//cdkFile := fmt.Sprintf("./6_w/tips_1_codecs.json")
	cdkFile := fmt.Sprintf(codakPth)
	kdkDat, err := ioutil.ReadFile(cdkFile)
	if err != nil {
		log.Fatal(err)
	}
	//var codecsR []Codak
	//var codecs []av.CodecData
	//err = json.Unmarshal(kdkDat, &codecsR)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, v := range codecsR {
	//	nv := CnvtTip(&v)
	//	codecs = append(codecs, nv)
	//}

	var codecsR []h264parser.CodecData
	var codecs []av.CodecData
	err = json.Unmarshal(kdkDat, &codecsR)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range codecsR {
		nv := CnvtTipH264parser(v)
		codecs = append(codecs, nv)
	}

	fmt.Println("open codek [END]")
	// start sort files frames

	fmt.Println("sort files frames [END] ")

	fmt.Println("=====>Hello ws! 6. _UpgradeHTTP FUNC_")
	muxerMSE := mp4f.NewMuxer(nil)
	fmt.Println("=====>Hello ws! 7. _UpgradeHTTP FUNC_")
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	fmt.Println("=====>Hello ws! 8. _UpgradeHTTP FUNC_")
	meta, init := muxerMSE.GetInit(codecs)
	fmt.Println("=====>Hello ws! 9. _UpgradeHTTP FUNC_")
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, append([]byte{9}, meta...))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "Send",
		}).Errorln(err.Error())
		return
	}
	fmt.Println("=====>Hello ws! 10. _UpgradeHTTP FUNC_")
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, init)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "Send",
		}).Errorln(err.Error())
		return
	}

	fmt.Println("=====>Hello ws!videoStart+++ ")

	var videoStart bool
	controlExit := make(chan bool, 10)
	noClient := time.NewTimer(10 * time.Second)
	go func() {
		defer func() {
			fmt.Println("=====>Hello ws!controlExit <- true+++ ")
			controlExit <- true
			fmt.Println("=====>Hello ws!controlExit <- true+++ OK")
		}()
		for { // отладочные, диагностические функции
			header, _, err := wsutil.NextReader(conn, ws.StateServerSide)
			if err != nil {
				requestLogger.WithFields(logrus.Fields{
					"call": "Receive",
				}).Errorln(err.Error())
				return
			}
			//fmt.Println("=====>Hello ws!Receive <-  wsutil.NextReader+++ ")
			switch header.OpCode {
			case ws.OpPong:
				noClient.Reset(10 * time.Second)
				//fmt.Println("=====>Hello ws!Receive <-  wsutil.NextReader+++OpPong ")
			case ws.OpClose:
				fmt.Println("=====>Hello ws!Receive <-  wsutil.NextReader---OpClose ")
				return
			}
		}
	}()
	noVideo := time.NewTimer(10 * time.Second)
	pingTicker := time.NewTicker(500 * time.Millisecond)
	defer pingTicker.Stop()
	defer log.Println("client exit")
	logxtm := 0
	logwrt := 0
	fmt.Println("=====>Hello ws!videoStart+++[FOR WAIT] ")
	if repid {
		go StreamFrames(ch2)
	}

	//chFarame := make(chan *av.Packet)
	//go Repid(chFarame)
	for {
		select {
		case <-pingTicker.C:
			//fmt.Println("=====>Hello ws!Receive <-  pingTicker.C ")
			err = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if err != nil {
				return
			}
			buf, err := ws.CompileFrame(ws.NewPingFrame(nil))
			if err != nil {
				return
			}
			_, err = conn.Write(buf)
			if err != nil {
				return
			}
		case <-controlExit:
			fmt.Println("=====>Hello ws!Receive <-  controlExit ")
			requestLogger.WithFields(logrus.Fields{
				"call": "controlExit",
			}).Errorln("Client Reader Exit")
			return
		case <-noClient.C:
			fmt.Println("=====>Hello ws!Receive <-  noClient.C ")
			requestLogger.WithFields(logrus.Fields{
				"call": "ErrorClientOffline",
			}).Errorln("Client OffLine Exit")
			return
		case <-noVideo.C:
			fmt.Println("=====>Hello ws!Receive <-  noVideo.C ")
			requestLogger.WithFields(logrus.Fields{
				"call": "ErrorStreamNoVideo",
			}).Errorln(ErrorStreamNoVideo.Error())
			return
		case pck := <-ch2:
			//case pck := <-chFarame:
			//fmt.Println("=====>Hello ws!logxtm+++ ")
			logxtm++
			if logxtm > 100 {
				cadr500part++
				fmt.Printf("=====>Hello2 ws!Receive <-ch logxtm _500кадров.[%d]_ \n", cadr500part)
				logxtm = 0
			}
			if pck.IsKeyFrame {
				noVideo.Reset(10 * time.Second)
				videoStart = true
			}
			if !videoStart {
				continue
			}
			ready, buf, err := muxerMSE.WritePacket(*pck, false)
			if err != nil {
				requestLogger.WithFields(logrus.Fields{
					"call": "WritePacket",
				}).Errorln(err.Error())
				return
			}
			if ready {
				//fmt.Printf("=====>_ready.[%d]_ \n", wrtr1part)
				wrtr1part++
				wrtr1bath++
				if wrtr1bath == 10 {
					fmt.Printf("=====>_ready.[%d]_ \n", wrtr1part)
					wrtr1bath = 0
				}
				err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "SetWriteDeadline",
					}).Errorln(err.Error())
					return
				}
				//err = websocket.Message.Send(ws, buf)
				logwrt++
				if logwrt > 100 {
					wrtr500part++
					fmt.Printf("=====>_500кадров.[%d]_ \n", wrtr500part)
					logwrt = 0
				}
				err = wsutil.WriteServerMessage(conn, ws.OpBinary, buf)
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "Send",
					}).Errorln(err.Error())
					return
				}
			}
		}
	}

	fmt.Println("=====>Hello ws!Receive <-  for---END ")
}

type Codak struct {
	tip string `json:"tip"`
}

const avCodecTypeMagic = 233333

func (c *Codak) Type() av.CodecType {
	// h264
	fmt.Printf("Codec:%s|to:av.H264\n", c.tip)
	return av.H264
}

func CnvtTip(baz av.CodecData) av.CodecData {
	return baz
	// f is of type *foo
}

func CnvtTipH264parser(baz h264parser.CodecData) av.CodecData {
	return baz
	// f is of type *foo
}
