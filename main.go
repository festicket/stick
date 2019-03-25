package main

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"
)

func runner(c chan int, command string, verbose int) {
	cmd := exec.Command(command)

	if verbose > 0 {
		cmd.Stdout = os.Stdout
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	if verbose > 0 {
		log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
	}

	c <- cmd.Process.Pid
}

func runnerSync(command string, verbose int) {
	cmd := exec.Command(command)

	if verbose > 0 {
		cmd.Stdout = os.Stdout
	}

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func request(i int, task *Task) time.Duration {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	start := time.Now()
	resp, err := client.Get(task.RequestURL)
	end := time.Now()

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Invalid status code: %s", resp.StatusCode)
	}

	duration := end.Sub(start)

	log.Printf("[%v] %v %v", i, duration, resp.ContentLength)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	dirName := task.GetRelativeName("artefacts")

	dump(body, dirName, task.Name)

	return duration
}

func dump(body []byte, dirname, name string) {
	if err := os.Mkdir(dirname, os.ModeDir|0766); err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	fname := path.Join(dirname, fmt.Sprintf("%v-%x.json", name, md5.Sum(body)))
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		panic(err)
	}

	f.Write(body)
	f.Close()
}

func runTask(task *Task, config *Config) *[]time.Duration {
	fmt.Printf("Runnin task %v", task.Name)
	c := make(chan int)
	go runner(c, task.GetCommand(), config.Verbose)

	results := make([]time.Duration, config.AttemptsPerTask)

	for pid := range c {
		if config.Verbose > 0 {
			fmt.Println("Waiting for the server")
		}
		time.Sleep(20 * time.Second) // TODO: add healthcheck

		if config.Verbose > 0 {
			fmt.Println("Server is ready (probably)")
		}

		if task.PreRequestCommand != "" {
			runnerSync(task.GetPreRequestCommand(), config.Verbose)
		}

		for i := 0; i < config.AttemptsPerTask; i++ {
			results[i] = request(i, task)
		}

		fmt.Printf("Done. [%v]\n\n", pid)
		syscall.Kill(-pid, syscall.SIGKILL)
		break
	}

	close(c)

	return &results
}

func main() {
	base, _ := os.Getwd()
	configName := path.Join(base, "config.json")
	config := GetConfig(configName)

	results := make(map[string]*[]time.Duration)

	for _, task := range config.Tasks {
		results[task.Name] = runTask(&task, config)
	}

	for key, value := range results {
		fmt.Println(key)

		for _, d := range *value {
			fmt.Println(d)
		}
	}

	fmt.Println("All done. Goodbye")
}
