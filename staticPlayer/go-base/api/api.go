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

	r.Mount("/auth", authResource.Router())

	r.Get("/eump/initializers/interactive.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=1200")
		w.Header().Add("Connection", "keep-alive")

		w.Header().Add("Content-Type", "application/javascript")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:35 GMT")
		w.Header().Add("Etag", "W/\"647721bc-278a\"")
		w.Header().Add("Expires", "Wed, 31 May 2023 11:06:25 GMT")
		w.Header().Add("Last-Modified", "Wed, 31 May 2023 10:30:20 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T07:56:20+00:00")
		w.Header().Add("X-Id", "m9p-up-gc23")

		w.Write([]byte(interactiveJs))
	})
	r.Get("/eump/embeds/interactive.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=1200")
		w.Header().Add("Connection", "keep-alive")

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:35 GMT")
		w.Header().Add("Etag", "W/\"6478637c-292\"")
		w.Header().Add("Expires", "Wed, 31 May 2023 11:06:25 GMT")
		w.Header().Add("Last-Modified", "Wed, 31 May 2023 10:30:20 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T07:56:19+00:00")
		w.Header().Add("X-Id", "m9p-up-gc23")
		w.Write([]byte(interactiveHtml))
	})
	r.Get("/eump/versions/v18.40.4_9.37.5_81/eump-1tv.all.min.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=1200")
		w.Header().Add("Connection", "keep-alive")

		w.Header().Add("Content-Type", "text/css")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:35 GMT")
		w.Header().Add("Etag", "W/\"647721b6-1e502\"")
		w.Header().Add("Expires", "Wed, 31 May 2023 11:14:06 GMT")
		w.Header().Add("Last-Modified", "Wed, 31 May 2023 10:30:14 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T07:58:14+00:00")
		w.Header().Add("X-Id", "m9p-up-gc23")
		w.Write([]byte(tvallmincss))
	})
	r.Get("/eump/configs/1tv_live.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=1200")
		w.Header().Add("Connection", "keep-alive")

		w.Header().Add("Content-Type", "application/javascript")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:35 GMT")
		w.Header().Add("Etag", "W/\"647721bc-2fff\"")
		w.Header().Add("Expires", "Wed, 31 May 2023 11:21:58 GMT")
		w.Header().Add("Last-Modified", "Wed, 31 May 2023 10:30:20 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T07:58:14+00:00")
		w.Header().Add("X-Id", "m9p-up-gc85")
		w.Write([]byte(tv_liveJs))
	})
	r.Get("/eump/versions/v18.40.4_9.37.5_81/eump-1tv.all.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=1200")
		w.Header().Add("Connection", "keep-alive")

		w.Header().Add("Content-Type", "application/javascript")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:35 GMT")
		w.Header().Add("Etag", "W/\"647721b6-1a37ca\"")
		w.Header().Add("Expires", "Wed, 31 May 2023 11:14:06 GMT")
		w.Header().Add("Last-Modified", "Wed, 31 May 2023 10:30:14 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T07:58:14+00:00")
		w.Header().Add("X-Id", "m9p-up-gc23")
		file, err := ioutil.ReadFile("./files/eump-1tv.all.min.js")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(file)
	})
	r.Get("/eump/versions/v18.40.4_9.37.5_81/eump-live1tv.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=1200")
		w.Header().Add("Connection", "keep-alive")

		w.Header().Add("Content-Type", "application/javascript")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:35 GMT")
		w.Header().Add("Etag", "W/\"647721b6-f111\"")
		w.Header().Add("Expires", "Wed, 31 May 2023 11:21:18 GMT")
		w.Header().Add("Last-Modified", "Wed, 31 May 2023 10:30:14 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T08:10:06+00:00")
		w.Header().Add("X-Id", "m9p-up-gc97")
		file, err := ioutil.ReadFile("./files/eump-live1tv.min.js")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(file)
	})
	r.Get("/eump/configs/1tv_live.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=1200")
		w.Header().Add("Connection", "keep-alive")

		w.Header().Add("Content-Type", "application/javascript")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:35 GMT")
		w.Header().Add("Etag", "W/\"647721bc-2fff\"")
		w.Header().Add("Expires", "Wed, 31 May 2023 11:21:58 GMT")
		w.Header().Add("Last-Modified", "Wed, 31 May 2023 10:30:20 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T08:00:41+00:00")
		w.Header().Add("X-Id", "m9p-up-gc85")
		file, err := ioutil.ReadFile("./files/configs/1tv_live.js")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(file)
	})
	r.Get("/eump/initializers/1tv_live.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "STALE")
		w.Header().Add("Cache-Control", "max-age=1200")
		w.Header().Add("Connection", "keep-alive")

		w.Header().Add("Content-Type", "application/javascript")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:35 GMT")
		w.Header().Add("Etag", "W/\"647721bc-1f2b\"")
		w.Header().Add("Expires", "Wed, 31 May 2023 11:07:59 GMT")
		w.Header().Add("Last-Modified", "Wed, 31 May 2023 10:30:20 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T07:46:03+00:00")
		w.Header().Add("X-Id", "m9p-up-gc16")
		file, err := ioutil.ReadFile("./files/initializers/1tv_live.js")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(file)
	})
	r.Get("/eump/embeds/1tv_live.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "STALE")
		w.Header().Add("Cache-Control", "max-age=1200")
		w.Header().Add("Connection", "keep-alive")

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:35 GMT")
		w.Header().Add("Etag", "W/\"647721bd-dda\"")
		w.Header().Add("Expires", "Wed, 31 May 2023 11:29:25 GMT")
		w.Header().Add("Last-Modified", "Wed, 31 May 2023 10:30:21 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T07:51:07+00:00")
		w.Header().Add("X-Id", "m9p-up-gc23")
		file, err := ioutil.ReadFile("./files/1tv_live.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(file)
	})
	r.Get("/player/config/banner.gif", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept-Ranges", "bytes")
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=7200")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "image/gif")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:36 GMT")
		w.Header().Add("Etag", "\"5800b749-2b\"")
		w.Header().Add("Expires", "Fri, 11 Nov 2022 14:45:56 GMT")
		w.Header().Add("Last-Modified", "Fri, 14 Oct 2016 10:45:29 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("X-Cached-Since", "2023-06-02T06:59:45+00:00")
		w.Header().Add("X-Id", "m9p-up-gc16")
		file, err := ioutil.ReadFile("./files/banner.gif")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(file)
	})
	r.Get("/player/fonts/Montserrat-VariableFont_wght.ttf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept-Ranges", "bytes")
		w.Header().Add("Access-Control-Allow-", "*")
		w.Header().Add("Headers", "")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Cache", "MISS")
		w.Header().Add("Cache-Control", "max-age=7200")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:36 GMT")
		w.Header().Add("Etag", "\"637c9e12-6039c\"")
		w.Header().Add("Expires", "Fri, 02 Jun 2023 10:11:36 GMT")
		w.Header().Add("Last-Modified", "Tue, 22 Nov 2022 10:01:54 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("X-Id", "m9p-up-gc47")
		file, err := ioutil.ReadFile("./files/Montserrat-VariableFont_wght.ttf")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(file)
	})
	r.Get("/player/teleport_media/stable/teleport.shaka.1tv.bundle.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept-Ranges", "bytes")
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=3600")

		w.Header().Add("Content-Type", "application/javascript")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Etag", "\"641156aa-1f9bd\"")
		w.Header().Add("Expires", "Fri, 02 Jun 2023 07:47:55 GMT")
		w.Header().Add("Last-Modified", "Wed, 15 Mar 2023 05:24:58 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("X-Cached-Since", "2023-06-02T07:47:59+00:00")
		w.Header().Add("X-Id", "m9p-up-gc19")
		file, err := ioutil.ReadFile("./files/teleport.shaka.1tv.bundle.js")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(file)
	})
	r.Get("/player/eump1tv-current/shaka.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=7200")
		w.Header().Add("Connection", "keep-alive")
		//
		w.Header().Add("Content-Type", "application/javascript")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Etag", "W/\"5df0d7d0-3689c\"")
		w.Header().Add("Expires", "Fri, 11 Nov 2022 14:46:41 GMT")
		w.Header().Add("Last-Modified", "Wed, 11 Dec 2019 11:49:36 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("Transfer-Encoding", "chunked")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("X-Cached-Since", "2023-06-02T07:47:59+00:00")
		w.Header().Add("X-Id", "m9p-up-gc19")
		file, err := ioutil.ReadFile("./files/shaka.min.js")
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Add("Content-Length", fmt.Sprint(len(file)))
		w.Write(file)
	})
	r.Get("/uploads/video/material/splash/2023/06/02/799411/big/799411_big_885825f81d.jpg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept-Ranges", "bytes")
		w.Header().Add("Cache", "HIT")
		w.Header().Add("Cache-Control", "max-age=604800")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Content-Type", "image/jpeg")
		w.Header().Add("Date", "Fri, 02 Jun 2023 08:11:37 GMT")
		w.Header().Add("Etag", "\"6479a0a9-48e0a\"")
		w.Header().Add("Expires", "Fri, 09 Jun 2023 08:00:07 GMT")
		w.Header().Add("Last-Modified", "Fri, 02 Jun 2023 07:56:25 GMT")
		w.Header().Add("Server", "nginx")
		w.Header().Add("X-Cached-Since", "2023-06-02T08:00:07+00:00")
		w.Header().Add("X-Id", "m9p-up-gc16")
		file, err := ioutil.ReadFile("./files/799411_big_885825f81d.jpg")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(file)
	})
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
