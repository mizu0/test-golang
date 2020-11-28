package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
)



func handler(event codedeploy.PutLifecycleEventHookExecutionStatusInput) error {
	deploymentID := event.DeploymentId
	lifecycleEventHookExecutionID := event.LifecycleEventHookExecutionId
	validationTestResult := "Failed"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
	}))
	svc := codedeploy.New(sess)

	resp, err := http.Get("http://ecs-bluegreen-test-lb-1428833294.ap-northeast-1.elb.amazonaws.com:8080/hoge")
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 && string(body) == "Hello World!" {
		validationTestResult = "Succeeded"
		fmt.Println("Succeeded")
	} else {
		fmt.Println("Failed")
	}

	input := &codedeploy.PutLifecycleEventHookExecutionStatusInput{
		DeploymentId: deploymentID,
		LifecycleEventHookExecutionId:lifecycleEventHookExecutionID,
		Status: aws.String(validationTestResult),
	}
	output, err := svc.PutLifecycleEventHookExecutionStatus(input)
	fmt.Println(output)

	return nil
}

func main() {
	lambda.Start(handler)
}