package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	fmt.Println("starting func main")

	fmt.Printf("PID: %d\n", os.Getpid())

	fmt.Println("starting lambda")
	lambda.StartWithOptions(
		Handler,
		lambda.WithEnableSIGTERM(func() {
			fmt.Println("Received SIGTERM and shutting down")
		}),
	)
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
