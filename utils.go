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
func DumpBody(body []byte, dirname, name string) string {
	if err := os.Mkdir(dirname, os.ModeDir|0766); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	hashString := fmt.Sprintf("%x", md5.Sum(body))
	fname := path.Join(dirname, fmt.Sprintf("%v-%v.txt", name, hashString))
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}

	f.Write(body)
	f.Close()

	return hashString
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

// ClearDirectory removes all the files in the directory specified
func ClearDirectory(dirname string) error {
	d, err := os.Open(dirname)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(path.Join(dirname, name))
		if err != nil {
			return err
		}
	}

	return nil
}
