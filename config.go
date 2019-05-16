package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	AttemptsPerTask int `json:"attempts_per_task"`
	Verbose         int
	Tasks           []Task
	KeepArtefacts   bool `json:"keep_artefacts"`
}

func (c *Config) Println(msg string) {
	if c.Verbose > 0 {
		fmt.Println(msg)
	}
}

func (c *Config) Printf(s string, a ...interface{}) {
	if c.Verbose > 0 {
		fmt.Printf(s, a)
	}
}

type Task struct {
	Name              string
	Command           string
	PreRequestCommand string `json:"pre_request_command"`
	RequestURL        string `json:"url"`
}

func GetConfig(name string) *Config {
	f, err := os.Open(name)

	if err != nil {
		log.Fatal(err)
	}

	rawData, err := ioutil.ReadAll(f)

	if err != nil {
		log.Fatal(err)
	}

	var config Config

	if err := json.Unmarshal(rawData, &config); err != nil {
		log.Fatal(err)
	}

	return &config
}

type TaskResult struct {
	Duration   time.Duration
	HashString string
}
