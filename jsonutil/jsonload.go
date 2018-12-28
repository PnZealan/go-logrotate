package jsonutil

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// JSONStruct json struct
type JSONStruct struct {
}

// Load json from file
func (js JSONStruct) Load(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("config file not found")
	}
	if err := json.Unmarshal(data, v); err != nil {
		log.Println("configration err")
	}
}
