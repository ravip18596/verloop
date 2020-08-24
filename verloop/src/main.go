package main

import (
	"Constants"
	"Controllers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)


func main()  {
	log.SetOutput(os.Stdout)
	if Constants.Debug{
		log.SetLevel(log.DebugLevel)
	}else {
		log.SetLevel(log.InfoLevel)
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/add", Controllers.AddWords)
	router.HandleFunc("/stories/{story_id}", Controllers.FetchStory)
	router.HandleFunc("/stories", Controllers.FetchAllStories)

	//router.Use(recoverHandler)
	err := http.ListenAndServe(":8050", router)
	if err != nil {
		log.WithFields(log.Fields{
			"addr": ":8050",
			"err":  err,
		}).Error("Unable to create HTTP Service ")
	}else{
		log.Info("Service started at port 8050")
	}
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//Middleware should be able to handle panic

				log.WithFields(log.Fields{
					"URL": r.RequestURI,
				}).Errorf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}