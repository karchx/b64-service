package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/karchx/b64-service/config"
	"github.com/rs/cors"
)

type ResponseData struct {
	Data string `bson:"data" json:"data"`
}

type Config struct {
  QueryKey string
  Prefix string
  Path string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveFile)
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	prefix := "FACT"

	fileComplete := prefix + "_" + r.URL.Query().Get("fel") + ".pdf"

	file, err := os.Open("../assets-generics/" + fileComplete)

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

  c := Config{}
  log.Printf("%s", c.filterFile(r))

	fileResponse := ResponseData{
		Data: base64.StdEncoding.EncodeToString(data),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileResponse)
}

func (c *Config) filterFile(r *http.Request) string {
  c.loadConfig()
	return c.Prefix + "_" + r.URL.Query().Get(c.QueryKey) + ".pdf"
}

func (c *Config) loadConfig() {
	cfg, err := config.ParserConfig()
	if err != nil {
		log.Fatal(err)
	}
  c.QueryKey = cfg.Settings.Querys
  c.Prefix = cfg.Settings.Prefix
  c.Path = cfg.Settings.Path
}
