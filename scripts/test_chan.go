package main

import (
	"errors"
	"fmt"
	"time"
)

// Simulated validate function
func validate(msg string) error {
	fmt.Printf("Validating message: %s\n", msg)
	if msg == "error" {
		return errors.New("validation failed")
	}
	return nil
}

func main() {
	messageChan := make(chan string) // Unbuffered channel

	// Simulated message producer in a goroutine
	go func() {
		messages := []string{"hello", "world", "error", "last"}
		for _, msg := range messages {
			time.Sleep(1 * time.Second) // Simulate a delay between messages
			messageChan <- msg          // Send message to the channel
		}
		close(messageChan) // Close the channel when done
	}()

	// Main thread waits for messages
	for {
		msg, ok := <-messageChan
		if !ok {
			fmt.Println("Message channel closed. Exiting...")
			break // Exit the loop if the channel is closed
		}

		err := validate(msg)
		if err != nil {
			fmt.Println("Validation error:", err)
			break // Exit the loop if validate() returns an error
		}
	}
	fmt.Println("Main thread exiting.")
}
