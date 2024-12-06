package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	fmt.Println("starting func main")

	fmt.Printf("PID: %d\n", os.Getpid())

	fmt.Println("starting graceful shutdown function")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("waiting for signal")
		sig := <-sigChan
		fmt.Printf("Received signal: %s, shutting down\n", sig)

		time.Sleep(2 * time.Second)

		fmt.Println("shutdown complete")

		// Reset the signal handlers and re-send the signal to the current process
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)
		if err := syscall.Kill(os.Getpid(), sig.(syscall.Signal)); err != nil {
			fmt.Printf("failed to send signal: %v\n", err)
		}
	}()

	fmt.Println("starting lambda")
	lambda.Start(Handler)
}

func Handler() (events.APIGatewayProxyResponse, error) {
	fmt.Println("request received")

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"message": "Hello World"}`,
	}, nil
}
