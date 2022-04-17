package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func HandleRequest(ctx context.Context, raw_event json.RawMessage) (string, error) {

	// unlikely this'd be the approach for a real function, that'd most likely be aligned to the trigger + event, but this'll work for any event as a proof-of-concept
	var unknown_event map[string]interface{}
	err := json.Unmarshal([]byte(raw_event), &unknown_event)
	if err != nil {
		error := fmt.Errorf("exception whilst unmarshalling the raw event: %w", err)
		return "", errors.New(error.Error())
	} else {
		log.Println("Unmarshalled event: ", unknown_event)
	}

	timestamp := time.Now().Unix()

	deadline_epoch, _ := ctx.Deadline()
	lc, _ := lambdacontext.FromContext(ctx)

	log.Println("AWS Request ID:", lc.AwsRequestID)
	log.Println("Trace ID:", os.Getenv("_X_AMZN_TRACE_ID"))
	log.Println("Context Deadline:", deadline_epoch)

	log.Println("Function name:", lambdacontext.FunctionName)
	log.Println("Function ARN:", lc.InvokedFunctionArn)
	log.Println("Handler:", os.Getenv("_HANDLER"))
	log.Println("Memory limit in MB:", lambdacontext.MemoryLimitInMB)

	log.Println("Log Group:", lambdacontext.LogGroupName)
	log.Print("Log Stream:", lambdacontext.LogStreamName)

	log.Println("Region:", os.Getenv("AWS_REGION"))
	log.Println("AWS Access Key ID:", os.Getenv("AWS_ACCESS_KEY_ID"))

	cw_logs := fmt.Sprintf("aws logs get-log-events --log-group-name '%s' --log-stream-name '%s' --start-time %d", lambdacontext.LogGroupName, lambdacontext.LogStreamName, timestamp)
	return cw_logs, nil
}

func main() {
	lambda.Start(HandleRequest)
}
