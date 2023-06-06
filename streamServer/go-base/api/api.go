// Package api configures an http server for administration and application resources.
package api

import (
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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
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

	r.Get("/get_hls_session", func(w http.ResponseWriter, r *http.Request) {
		paramCallback := r.URL.Query().Get("callback")
		respData := strings.Replace(getHlsSession, "_eump_fw_callback_43108", paramCallback, -1)
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "text/html")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Keep-Alive", "timeout=15")
		w.Header().Add("Server", "QRATOR")
		w.Header().Add("Set-Cookie", "s=f58pZgYINfU8kTAIo0k5whAVZFMT/xqGuVknvanSDjPnCbZOwI8z0Gh210yY7mU9P/eSr5EdPDxjP5LYRIy8s+AysuFXs2B6+upcMFX70hWU2Eg/x2EOleBlRQ4w4qBs72sNcAkIErGzh7+cBUZVuVUcEcjyRqrajh1fJdegw1wZrN+vJmd01DhWPxbmIUf+cUpHJMpLry5YTaGd8NnOkJapTIJZQNoYFT43rmdmFLw=;Max-Age=14400")
		w.Header().Add("Transfer-Encoding", "chunked")
		// ?callback=_eump_fw_callback_43108&rnd=1685693402042379
		w.Write([]byte(respData))
	})
	r.Get("/api/playlist/1tvch-dvr-hls-12h_as_array.json", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "text/json; charset=UTF-8")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Expires", "Thu, 01 Jan 1970 00:00:01 GMT")
		w.Header().Add("Keep-Alive", "timeout=15")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Server", "QRATOR")
		w.Header().Add("X-Powered-By", "Express")
		datFile, err := ioutil.ReadFile("./files/playlist/1tvch-dvr-hls-12h_as_array.json")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(datFile)
	})
	r.Get("/api/playlist/1tvch-v1_as_array.json", func(w http.ResponseWriter, r *http.Request) {
		//
		w.Header().Add("Cache-Control", "public, max-age=0, no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "text/json; charset=UTF-8")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Expires", "Thu, 01 Jan 1970 00:00:01 GMT")
		w.Header().Add("Keep-Alive", "timeout=15")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Server", "QRATOR")
		w.Header().Add("X-Powered-By", "Express")
		datFile, err := ioutil.ReadFile("./files/playlist/1tvch-v1_as_array.json")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(datFile)
	})
	r.Get("/api/schedule.json", func(w http.ResponseWriter, r *http.Request) {
		//
		w.Header().Add("Cache-Control", "max-age=60")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "text/json; charset=UTF-8")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Expires", "Thu, 01 Jan 1970 00:00:01 GMT")
		w.Header().Add("Keep-Alive", "timeout=15")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Server", "QRATOR")
		w.Header().Add("X-Powered-By", "Express")
		datFile, err := ioutil.ReadFile("./files/schedule.json")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(datFile)
	})
	r.Get("/api/restrictions.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept-Ranges", "bytes")

		w.Header().Add("Cache-Control", "max-age=2")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Expires", "Thu, 01 Jan 1970 00:00:01 GMT")
		w.Header().Add("Keep-Alive", "timeout=15")
		w.Header().Add("Etag", "\"ac8029-49-5fd2118e59340\"")
		w.Header().Add("Last-Modified", "Fri, 02 Jun 2023 08:11:17 GMT")
		w.Header().Add("Server", "QRATOR")
		datFile, err := ioutil.ReadFile("./files/restrictions.json")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(datFile)
	})
	r.Get("/api/com-live.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept-Ranges", "bytes")

		w.Header().Add("Age", "48")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Add("Keep-Alive", "timeout=15")
		w.Header().Add("Server", "QRATOR")
		w.Header().Add("Status", "200 OK")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("Via", "nginx")
		w.Header().Add("X-Cache", "HIT")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-Download-Options", "noopen")
		w.Header().Add("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Add("X-Request-Id", "ab6c9313-3af0-473f-9991-a6deb87b9a04")
		w.Header().Add("X-Xss-Protection", "1; mode=block")

		datFile, err := ioutil.ReadFile("./files/com-live.json")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(datFile)
	})
	r.Get("/api/com-inject.json", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Cache-Control", "max-age=2")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Keep-Alive", "timeout=15")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Etag", "W/\"2b48090-aa9-5fd2118a88a40\"")
		w.Header().Add("Expires", "Fri, 02 Jun 2023 08:11:36 GMT")
		w.Header().Add("Last-Modified", "Fri, 02 Jun 2023 08:11:13 GMT")
		w.Header().Add("Server", "QRATOR")
		w.Header().Add("Transfer-Encoding", "chunked")

		datFile, err := ioutil.ReadFile("./files/com-inject.json")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(datFile)
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
