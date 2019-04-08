package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

// DumpBody saves the raw data returned from the server to a file.
func DumpBody(body []byte, dirname, name string) {
	if err := os.Mkdir(dirname, os.ModeDir|0766); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	fname := path.Join(dirname, fmt.Sprintf("%v-%x.txt", name, md5.Sum(body)))
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}

	f.Write(body)
	f.Close()
}

// JSONPrettyfier converts ugly formatted JSON into something better.
func JSONPrettyfier(body []byte) []byte {
	var parsed map[string]interface{}
	err := json.Unmarshal(body, &parsed)
	if err != nil {
		log.Fatalf("Invalid JSON %v", err)
	}

	prettyBody, _ := json.MarshalIndent(parsed, "", "  ")

	return prettyBody
}
