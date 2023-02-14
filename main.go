package main

import (
	"io/ioutil"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"os"
	"log"
)

type ResponseData struct {
	Data string `bson:"data" json:"data"`
}

func main() {
	http.HandleFunc("/", serveFile)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func serveFile(w http.ResponseWriter, r *http.Request) {
	prefix := "FACT"

	fileComplete := prefix + "_" + r.URL.Query().Get("fel") + ".pdf"

	//file, err := os.Open("../facturas/FACT_002D81A0_2635155854.pdf")
	file, err := os.Open("../facturas/" + fileComplete)

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
		Data:  base64.StdEncoding.EncodeToString(data),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileResponse)
}

func filterFile() {}
