package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	logger "github.com/gothew/l-og"
	"github.com/karchx/b64-service/config"
	"github.com/rs/cors"
)

type ResponseData struct {
	Data string `bson:"data" json:"data"`
}

type Services struct {
	service map[string]config.SettingsConfig
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveFile)
	handler := cors.Default().Handler(mux)
	logger.Info("Server run")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func serveFile(w http.ResponseWriter, r *http.Request) {

	service := Services{}
	service.loadConfig()

	pathFile, err := service.filterFile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := os.Open(pathFile)
	if err != nil {
		log.Printf("%s: ", err)
		http.Error(w, "Can't open file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("%s: ", err)
		http.Error(w, "Can't open file", http.StatusInternalServerError)
		return
	}

	fileResponse := ResponseData{
		Data: base64.StdEncoding.EncodeToString(data),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileResponse)
}

func (s *Services) filterFile(r *http.Request) (string, error) {
	if r.URL.Path != "/favicon.ico" {

		filterService := strings.Replace(r.URL.Path, "/", "", -1)

		service, ok := s.service[filterService]

		if ok {
			return service.Path + "/" + service.Prefix + "_" + r.URL.Query().Get(service.Querys) + ".pdf", nil
		} else {
			return "", errors.New("Service not config")
		}
	}

	return "", errors.New("Service not config")
}

func (s *Services) loadConfig() {
	cfg, err := config.ParserConfig()
	if err != nil {
		log.Fatal(err)
	}

	s.service = cfg.Services
}
