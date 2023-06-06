// Package api configures an http server for administration and application resources.
package api

import (
	"fmt"
	"github.com/dhax/go-base/api/admin"
	"github.com/dhax/go-base/api/app"
	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/auth/pwdless"
	"github.com/dhax/go-base/database"
	"github.com/dhax/go-base/email"
	"github.com/dhax/go-base/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

// New configures application resources and routes.
func New(enableCORS bool) (*chi.Mux, error) {
	logger := logging.NewLogger()
	db, err := database.DBConn()
	if err != nil {
		logger.Named("module").Named("database").Warn(err.Error())
		//return nil, err //TODO:
	}

	mailer, err := email.NewMailer()
	if err != nil {
		logger.Named("module").Named("email").Warn(err.Error())
		return nil, err
	}

	authStore := database.NewAuthStore(db)
	authResource, err := pwdless.NewResource(authStore, mailer)
	if err != nil {
		logger.Named("module").Named("auth").Warn(err.Error())
		return nil, err
	}

	adminAPI, err := admin.NewAPI(db)
	if err != nil {
		logger.Named("module").Named("admin").Warn(err.Error())
		return nil, err
	}

	appAPI, err := app.NewAPI(db)
	if err != nil {
		logger.Named("module").Named("app").Error(err.Error())
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Use(logging.NewStructuredLogger(logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// use CORS middleware if client is not served by this api, e.g. from other domain or CDN
	if enableCORS {
		r.Use(corsConfig().Handler)
	}

	//r.Get("/dash-live2/streams/1tv-dvr", func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Add("Accept-Ranges", "bytes")
	//	w.Header().Add("Age", "0")
	//	w.Header().Add("Cache-Control", "no-cache")
	//	w.Header().Add("Etag", "\"647d9bff-1016\"")
	//	w.Header().Add("Expires", "Mon, 05 Jun 2023 08:25:34 GMT")
	//	w.Header().Add("Last-Modified", "Mon, 05 Jun 2023 08:25:35 GMT")
	//	w.Header().Add("Server", "nginx")
	//	w.Header().Add("Vary", "Origin")
	//	w.Header().Add("Connection", "keep-alive")
	//	w.Header().Add("Content-Type", "application/octet-stream")
	//	w.Header().Add("Date", "Mon, 05 Jun 2023 08:25:40 GMT")
	//	w.Header().Add("Keep-Alive", "timeout=60")
	//
	//	datFile, err := ioutil.ReadFile("./files/1tv-dvr/1tvdash.mpd")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	w.Write(datFile)
	//})
	r.Route("/dash-live2/streams/1tv-dvr/{chanName}", func(r chi.Router) {
		//r.Use(streamCtx)
		r.Get("/", streamCtx)
	})
	//  "https://cdn.bitmovin.com/content/assets/art-of-motion-dash-hls-progressive/mpds/f08e80da-bf1d-4e3d-8899-f0f6155f6efa.mpd"

	r.Route("/content/assets/art-of-motion-dash-hls-progressive/mpds/{chanName}", func(r chi.Router) {
		//r.Use(streamCtx)
		r.Get("/", streamCtxSaverExample)
	})
	r.Route("/content/assets/art-of-motion-dash-hls-progressive/video/1080_4800000/dash/{chanName}", func(r chi.Router) {
		//r.Use(streamCtx)
		r.Get("/", streamCtxSaverVideo)
	})

	r.Mount("/auth", authResource.Router())
	r.Group(func(r chi.Router) { // Это добавит 2 маршрута к роутеру
		r.Use(authResource.TokenAuth.Verifier())
		r.Use(jwt.Authenticator)
		r.Mount("/admin", adminAPI.Router())
		r.Mount("/api", appAPI.Router())
	})
	//
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	//client := "./public"
	//r.Get("/*", SPAHandler(client))

	return r, nil
}

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}

// SPAHandler serves the public Single Page Application.
func SPAHandler(publicDir string) http.HandlerFunc {
	handler := http.FileServer(http.Dir(publicDir))
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage := path.Join(publicDir, "index.html")
		serviceWorker := path.Join(publicDir, "service-worker.js")

		requestedAsset := path.Join(publicDir, r.URL.Path)
		if strings.Contains(requestedAsset, "service-worker.js") {
			requestedAsset = serviceWorker
		}
		if _, err := os.Stat(requestedAsset); err != nil {
			http.ServeFile(w, r, indexPage)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

const getHlsSession = `_eump_fw_callback_43108(
 { "s" : "f58pZgYINfU8kTAIo0k5whAVZFMT%2FxqGuVknvanSDjPnCbZOwI8z0Gh210yY7mU9P%2FeSr5EdPDxjP5LYRIy8s%2BAysuFXs2B6%2BupcMFX70hWU2Eg%2Fx2EOleBlRQ4w4qBs72sNcAkIErGzh7%2BcBUZVuVUcEcjyRqrajh1fJdegw1wZrN%2BvJmd01DhWPxbmIUf%2BcUpHJMpLry5YTaGd8NnOkJapTIJZQNoYFT43rmdmFLw%3D" });
`

var numFile int = 6967

func streamCtx(w http.ResponseWriter, r *http.Request) {
	fName := chi.URLParam(r, "chanName") //1tvdash.mpd
	//paramCallback := r.URL.Query().Get("e") //1685952628
	//fmt.Println(fName)
	//fmt.Println(paramCallback)
	var fullFname string
	if fName == "1tvdash.mpd" {
		fullFname = fmt.Sprintf("./files/1tv-dvr/%s", fName)
	} else {
		initName := "1tvdash1-12hFrag-hd-5-20230602T074403init.mp4"
		if fName == initName {
			fullFname = fmt.Sprintf("./files/1tv-dvr/%s", initName)
		} else {
			fullFname = fmt.Sprintf("./files/1tv-dvr/1tvdash1-12hFrag-hd-5-20230602T074403_00005%d.mp4", numFile)
			numFile++
		}

		//apiIncidentsPrefix := fmt.Sprintf("/dash-live2/streams/1tv-dvr/%s", fName)
		//urlStr := ComparableUrl("edge2.1internet.tv", apiIncidentsPrefix)
		//body, err := DoRequest(urlStr)
		//if err != nil {
		//	app.ErrRender(err)
		//	return
		//}

		//if fName == "1tvdash.mpd" {
		//	numFile++
		//	fName = fmt.Sprintf("%s_%d_%s", fName, numFile, paramCallback)
		//	body = []byte(OneRoad(string(body)))
		//}
		//f, err := os.Create(fmt.Sprintf("./httfiles/%s", fName))
		//if err != nil {
		//	app.ErrRender(err)
		//	return
		//}
		//defer f.Close()
		//f.Write(body)
		//w.Write(body)
		//return
	}

	datFile, err := ioutil.ReadFile(fullFname)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("num:%d|ln:%d|fName:%s\n", numFile, len(datFile), fName)
	w.Write(datFile)
}

var filesSaved = map[string]int{}

func streamCtxSaverExample(w http.ResponseWriter, r *http.Request) {
	fName := chi.URLParam(r, "chanName") // dash.mpd
	fmt.Printf("req:%s\n", fName)
	finalName := fName
	_, ok := filesSaved[fName]
	if !ok {
		filesSaved[fName] = 0
	} else {
		filesSaved[fName]++
		finalName = fmt.Sprintf("%s_%d", finalName, filesSaved[fName])
	}

	apiIncidentsPrefix := fmt.Sprintf("/content/assets/art-of-motion-dash-hls-progressive/mpds/%s", fName)
	urlStr := ComparableUrl("https://", "cdn.bitmovin.com", apiIncidentsPrefix)
	body, err := DoRequest(urlStr)
	if err != nil {
		app.ErrRender(err)
		return
	}

	body = []byte(TwoRoad(string(body)))
	f, err := os.Create(fmt.Sprintf("./httfiles/%s", finalName))
	if err != nil {
		app.ErrRender(err)
		return
	}
	defer f.Close()
	f.Write(body)
	w.Write(body)
}

func streamCtxSaverVideo(w http.ResponseWriter, r *http.Request) {
	fName := chi.URLParam(r, "chanName") // dash.mpd
	fmt.Printf("req:%s\n", fName)
	finalName := fName
	_, ok := filesSaved[fName]
	if !ok {
		filesSaved[fName] = 0
	} else {
		filesSaved[fName]++
		finalName = fmt.Sprintf("%s_%d", finalName, filesSaved[fName])
	}

	apiIncidentsPrefix := fmt.Sprintf("/content/assets/art-of-motion-dash-hls-progressive/video/1080_4800000/dash/%s", fName)
	urlStr := ComparableUrl("https://", "cdn.bitmovin.com", apiIncidentsPrefix)
	body, err := DoRequest(urlStr)
	if err != nil {
		app.ErrRender(err)
		return
	}

	f, err := os.Create(fmt.Sprintf("./httfiles/video/%s", finalName))
	if err != nil {
		app.ErrRender(err)
		return
	}
	defer f.Close()
	f.Write(body)
	w.Write(body)
}

func ComparableUrl(prefix, baseURL, apiPrefix string) (urlStr string) { // "http://example.com",  "/path" --> "http://example.com/path?param1=value1&param2=value2"
	params := url.Values{}
	if prefix == "" {
		prefix = "http://"
	}
	u, _ := url.ParseRequestURI(fmt.Sprint(prefix, baseURL))
	u.Path = apiPrefix // path.Join(apiPrefix, urlMedhod)
	u.RawQuery = params.Encode()
	return fmt.Sprintf("%v", u)
}

func DoRequest(urlStr string) (body []byte, err error) { // --> {"incidents":[]}
	// Запрашиваем данные с провайдера
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("client request do: %s", err)
	}
	if resp != nil {
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("resp.ReadAll err:%v|status:%s|code:%d|url:%s", err, resp.Status, resp.StatusCode, urlStr)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("can't request status responce not status: %s, code:%d|method:GET|body:%s|urlStr:%s", resp.Status, resp.StatusCode, body, urlStr)
		}
	} else {
		return nil, fmt.Errorf("responce statistics api is nil err:%v|status:%s|code:%d|url:%s", err, resp.Status, resp.StatusCode, urlStr)
	}
	return body, nil
}

func OneRoad(text string) string {
	var re = regexp.MustCompile(`(?s)<Representation id="2".*?<\/Representation>`)
	s := re.ReplaceAllString(text, ``)
	re = regexp.MustCompile(`(?s)<Representation id="3".*?<\/Representation>`)
	s = re.ReplaceAllString(s, ``)
	re = regexp.MustCompile(`(?s)<Representation id="4".*?<\/Representation>`)
	s = re.ReplaceAllString(s, ``)
	re = regexp.MustCompile(`(?s)<Representation id="5".*?<\/Representation>`)
	s = re.ReplaceAllString(s, ``)
	re = regexp.MustCompile(`(?s)<AdaptationSet mimeType="audio/mp4".*?</AdaptationSet>`)
	s = re.ReplaceAllString(s, ``)
	re = regexp.MustCompile(`(?s)<AdaptationSet mimeType="application/mp4".*?</AdaptationSet>`)
	s = re.ReplaceAllString(s, ``)
	return s
}

func TwoRoad(text string) string {
	var re = regexp.MustCompile(`(?s)<Representation id="180_250000".*?/>`)
	s := re.ReplaceAllString(text, ``)
	re = regexp.MustCompile(`(?s)<Representation id="270_400000".*?/>`)
	s = re.ReplaceAllString(s, ``)
	re = regexp.MustCompile(`(?s)<Representation id="360_800000".*?/>`)
	s = re.ReplaceAllString(s, ``)
	re = regexp.MustCompile(`(?s)<Representation id="540_1200000".*?/>`)
	s = re.ReplaceAllString(s, ``)
	re = regexp.MustCompile(`(?s)<Representation id="720_2400000".*?/>`)
	s = re.ReplaceAllString(s, ``)
	re = regexp.MustCompile(`(?s)<AdaptationSet lang="en" mimeType="audio/mp4".*?</AdaptationSet>`)
	s = re.ReplaceAllString(s, ``)
	return s
}
