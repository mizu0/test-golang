package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/aws/aws-sdk-go/service/elbv2"
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
	lb := elbv2.New(sess)
	cd := codedeploy.New(sess)

	in := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{
			aws.String("ecs-bluegreen-test-lb"),
		},
	}
	out, err := lb.DescribeLoadBalancers(in)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// blue/green でのtarget groupの切り替わり見たいので、sleep入れてみる
	time.Sleep(time.Second * 30)

	// test listener の port
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/hello-world", *out.LoadBalancers[0].DNSName))
	//resp, err := http.Get("http://ecs-bluegreen-test-lb-1428833294.ap-northeast-1.elb.amazonaws.com:8080/hello-world")
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

	output, err := cd.PutLifecycleEventHookExecutionStatus(&codedeploy.PutLifecycleEventHookExecutionStatusInput{
		DeploymentId: deploymentID,
		LifecycleEventHookExecutionId:lifecycleEventHookExecutionID,
		Status: aws.String(validationTestResult),
	})
	fmt.Println(output)

	return nil
}

func main() {
	lambda.Start(handler)
}