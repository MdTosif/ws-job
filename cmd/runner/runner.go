package runner

import (
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"slices"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Runner struct {
	jobRunning []*Job
}

type Job struct {
	ID         string
	Cmd        *exec.Cmd
	Subscriber []*websocket.Conn
	// mutex to handle exit
	Exited bool

	mu sync.Mutex
}

func New() *Runner {
	return &Runner{}
}

func (j *Job) kill() {
	j.Cmd.Process.Signal(os.Kill)
	j.SetExited()
}

// SetExited safely sets the Exited flag to true using a mutex
func (j *Job) SetExited() {
	j.mu.Lock()         // Acquire the lock
	defer j.mu.Unlock() // Ensure the lock is released
	j.Exited = true     // Modify the shared resource
}

// IsExited safely checks the Exited flag using a mutex
func (j *Job) IsExited() bool {
	j.mu.Lock()         // Acquire the lock
	defer j.mu.Unlock() // Ensure the lock is released
	return j.Exited     // Read the shared resource
}

func (r *Runner) Run(command string) (io.ReadCloser, io.ReadCloser, error) {
	if command == "stop-all-running-jobs" {
		r.Stop()
		return nil, nil, nil
	}
	cmd := exec.Command("bash", "-c", command)

	currentJob := &Job{
		ID:         uuid.NewString(),
		Cmd:        cmd,
		Exited: false,
		Subscriber: []*websocket.Conn{},
	}

	r.jobRunning = append(r.jobRunning, currentJob)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}


	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}


	// Wait in a goroutine so it doesn’t block
	go func() {
		cmd.Wait()
		log.Println("cmd exited after waiting")
		currentJob.SetExited()
		// remove it from cmdRunning
		for i, c := range r.jobRunning {
			if c.Cmd == cmd {
				r.jobRunning = slices.Delete(r.jobRunning, i, i+1)
				break
			}
		}
	}()

	return stdoutPipe, stderrPipe, nil
}

func (r *Runner) Stop() {
	for _, job := range r.jobRunning {
		job.kill()
	}
}
