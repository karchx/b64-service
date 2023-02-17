package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/karchx/b64-service/config"
	"github.com/rs/cors"
)

type ResponseData struct {
	Data string `bson:"data" json:"data"`
}

type Config struct {
	Name     string
	QueryKey string
	Prefix   string
	Path     string
}

type Services struct {
  config []Config
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveFile)
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	/*prefix := "FACT"

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
	}*/

  service := Services{}
  service.loadConfig()
  var pathF string

  for _, service := range service.config {
    fmt.Println(service.Name)
    if service.Name == r.URL.String() {
	    pathF = filterFile(r, service)
    }

	  fmt.Printf("%s", pathF)
  }
 

	fileResponse := ResponseData{
    Data: "working...",
		//Data: base64.StdEncoding.EncodeToString(data),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileResponse)
}

func filterFile(r *http.Request, config Config) string {
  return config.Prefix + "_" + r.URL.Query().Get(config.QueryKey) + ".pdf"
}

func (s *Services) loadConfig() {
	cfg, err := config.ParserConfig()
  var configServices []Config
	if err != nil {
		log.Fatal(err)
	}

  for key, service := range cfg.Services {
    config := Config {
      Name: key,
      Prefix: service.Prefix,
      Path: service.Path,
      QueryKey: service.Querys,
    }
    configServices = append(configServices, config)
  }

  s.config = configServices
}
