package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	log "github.com/sirupsen/logrus"
)

func Directory(w http.ResponseWriter, r *http.Request, devMode bool) {
	if devMode {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	//read directory
	files, err := ioutil.ReadDir("./logs")
	if err != nil {
		http.Error(w, "Unable to open logs folder. There might not be any logs created yet", http.StatusInternalServerError)
		log.Error("Unable to open logs folder. There might not be any logs created yet")
		return
	}

	//put them all in a struct
	var response config.DirectoryResponse

	//this initialises the array to empty, as otherwise if there are no folders it is null
	response.FolderNames = []string{}

	for _, f := range files {
		response.FolderNames = append(response.FolderNames, f.Name())
	}

	//convert struct to json and return the response

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Error while encoding the response of directory HTTP request: " + err.Error())
		return
	}
}
