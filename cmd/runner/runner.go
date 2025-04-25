package runner

import (
	"bytes"
	"os/exec"
)

// Runner struct that will execute commands
type Runner struct{}

// New creates a new Runner instance
func New() *Runner {
	return &Runner{}
}

// Run executes a shell command and returns a channel that will provide output as it comes in
func (r *Runner) Run(command string) (<-chan string, chan error) {
	// Create a channel for output
	outputChan := make(chan string)
	// Create a channel for errors
	errorChan := make(chan error)

	// Start a goroutine to execute the command
	go func() {
		defer close(outputChan)
		defer close(errorChan)

		cmd := exec.Command("bash", "-c", command)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		// Start the command and stream its output
		err := cmd.Start()
		if err != nil {
			errorChan <- err
			return
		}

		// Create a goroutine that will read the output and send it to the channel
		go func() {
			// Continuously read the output line by line and send it through the channel
			for {
				// Read output as it comes in (line by line)
				data := make([]byte, 10)
				_, err := out.Read(data)
				if err != nil {
					// Break the loop if there's no more output
					break
				}
				// Send the line to the output channel, but only if it's not closed
				select {
				case outputChan <- string(data): // Safe to send data to the channel
				default: // Channel is closed, no sending
					return
				}
			}
		}()

		// Wait for the command to finish
		err = cmd.Wait()
		if err != nil {
			errorChan <- err
		}
	}()

	// Return both channels
	return outputChan, errorChan
}
