package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	CONFIG_FNAME          = "config.json"
	ARTEFACTS_FOLDER_NAME = "artefacts"
	VERIFY_SSL            = false
)

func main() {
	config := GetConfig(CONFIG_FNAME)
	results := make(map[string]*[]time.Duration)

	fmt.Println("Hello there")

	for _, task := range config.Tasks {
		results[task.Name] = runTask(&task, config)
		fmt.Printf("Done: %v\n", task.Name)
	}

	for key, value := range results {
		fmt.Println(key)
		for _, d := range *value {
			fmt.Println(d)
		}
	}

	fmt.Println("All done. Goodbye")
}

func runTask(task *Task, config *Config) *[]time.Duration {
	config.Printf("Running task %v\n", task.Name)

	c := make(chan int)

	// Prepare the target server
	go runCommand(c, task.Command, config.Verbose)

	results := make([]time.Duration, config.AttemptsPerTask)

	for pid := range c {
		config.Println("Waiting for the server")
		time.Sleep(20 * time.Second) // TODO: add healthcheck
		config.Println("Server is ready (probably)")

		// A hook to run additional actions before each request e.g. clear the cache.
		if task.PreRequestCommand != "" {
			runCommandSync(task.PreRequestCommand, config.Verbose)
		}

		for i := 0; i < config.AttemptsPerTask; i++ {
			results[i] = doRequest(i, task, config)
		}

		// TODO: kill the task even if there is an error happened before
		// (doRequest may fail with log.Fatal and I probably want to change that)
		config.Printf("Done. [%v]\n\n", pid)
		syscall.Kill(-pid, syscall.SIGKILL)
		break
	}

	close(c)
	return &results
}

// runCommand runs the command without waiting for it to complete.
func runCommand(c chan int, command string, verbose int) {
	cmd := exec.Command(command)
	if verbose > 1 {
		cmd.Stdout = os.Stdout
	}

	// To make it possible to kill it with all the children later.
	// https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	if verbose > 1 {
		log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
	}

	c <- cmd.Process.Pid
}

// runCommandSync runs the command and waits for it to complete.
func runCommandSync(command string, verbose int) {
	cmd := exec.Command(command)
	if verbose > 1 {
		cmd.Stdout = os.Stdout
	}

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// doRequest does a request to API endpoint and measure the time spent.
func doRequest(i int, task *Task, config *Config) time.Duration {
	// TODO: init only once?
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !VERIFY_SSL},
	}
	client := &http.Client{Transport: tr}

	start := time.Now()
	resp, err := client.Get(task.RequestURL)
	end := time.Now()

	if err != nil {
		log.Fatal(err) // TODO: return Error instead for the caller to handle
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Invalid status code: %s", resp.StatusCode)
	}

	duration := end.Sub(start)

	config.Printf("[%v] %v %v\n", i, duration, resp.ContentLength)
	defer resp.Body.Close()

	var body []byte

	rawBody, _ := ioutil.ReadAll(resp.Body)
	contentType := resp.Header.Get("Content-Type")

	if contentType == "application/json" {
		body = JSONPrettyfier(rawBody)
	} else {
		body = rawBody
	}

	DumpBody(body, ARTEFACTS_FOLDER_NAME, task.Name)

	return duration
}
