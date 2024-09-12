package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/glue"
)

func main() {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create a Glue client
	client := glue.NewFromConfig(cfg)

	// List Glue triggers
	listTriggersInput := &glue.ListTriggersInput{}
	listTriggersOutput, err := client.ListTriggers(context.TODO(), listTriggersInput)
	if err != nil {
		log.Fatalf("failed to list triggers, %v", err)
	}

	// Get details for each trigger
	for _, triggerName := range listTriggersOutput.TriggerNames {
		getTriggerInput := &glue.GetTriggerInput{
			Name: aws.String(triggerName),
		}
		getTriggerOutput, err := client.GetTrigger(context.TODO(), getTriggerInput)
		if err != nil {
			log.Printf("failed to get trigger %s, %v", triggerName, err)
			continue
		}

		// Print the trigger schedule
		if getTriggerOutput.Trigger.Schedule != nil {
			fmt.Printf("Trigger: %s, Schedule: %s\n", triggerName, *getTriggerOutput.Trigger.Schedule)
		} else {
			fmt.Printf("Trigger: %s has no schedule\n", triggerName)
		}
	}
}
