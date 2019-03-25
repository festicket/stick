package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	AttemptsPerTask int `json:"attempts_per_task"`
	Verbose         int
	Tasks           []Task
}

type Task struct {
	Name              string
	Command           string
	PreRequestCommand string `json:"pre_request_command"`
	RequestURL        string `json:"url"`
}

func (t *Task) GetCommand() string {
	return t.GetRelativeName(t.Command)
}

func (t *Task) GetPreRequestCommand() string {
	return t.GetRelativeName(t.PreRequestCommand)
}

func (t *Task) GetRelativeName(name string) string {
	base, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path.Join(base, name)
}

func GetConfig(name string) *Config {
	f, err := os.Open(name)

	if err != nil {
		panic(err)
	}

	rawData, err := ioutil.ReadAll(f)

	fmt.Println(rawData)

	if err != nil {
		panic(err)
	}

	var config Config

	if err := json.Unmarshal(rawData, &config); err != nil {
		panic(err)
	}

	return &config
}
