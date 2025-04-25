package runner

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
)

type Runner struct{
	cmdRunning []*exec.Cmd
}

type Job struct {
	ID         string
	Cmd        string
	Subscriber []*websocket.Conn
}

func New() *Runner {
	return &Runner{}
}

func (r *Runner) Run(command string) (<-chan string, <-chan string, error) {
	cmd := exec.Command("bash", "-c", command)
	r.cmdRunning = append(r.cmdRunning, cmd)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	// Function to stream output line by line
	stream := func(pipe io.ReadCloser, ch chan<- string) {
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}

	go stream(stdoutPipe, stdoutChan)
	go stream(stderrPipe, stderrChan)

	// Wait in a goroutine so it doesn’t block
	go func() {
		cmd.Wait()
		// remove it from cmdRunning
		for i, c := range r.cmdRunning {
			if c == cmd {
				r.cmdRunning = append(r.cmdRunning[:i], r.cmdRunning[i+1:]...)
				break
			}
		}
	}()

	return stdoutChan, stderrChan, nil
}

func (r *Runner) Stop() {
	for _, cmd := range r.cmdRunning {
		cmd.Process.Signal(os.Kill)
	}
}	
