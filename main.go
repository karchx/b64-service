package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

type ResponseData struct {
	Data string `bson:"data" json:"data"`
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

	fileResponse := ResponseData{
		Data: base64.StdEncoding.EncodeToString(data),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileResponse)
}

func filterFile() {}
