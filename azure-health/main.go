package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	ct "github.com/daviddengcn/go-colortext"
)

//  reference: https://github.com/Azure/azure-sdk-for-go/blob/master/services/monitor/mgmt/2017-09-01/insights/activitylogs.go
func main() {
	// create an authorizer from using Azure CLI
	authorizer, err := auth.NewAuthorizerFromCLI()

	// Create an ActivityLogs instance.
	logsClient := insights.NewActivityLogsClient("<subscriptionID>")
	logsClient.Authorizer = authorizer

	// Create the OData filter for a time interval and the Azure.Health Provider.
	endTime := time.Now().UTC()
	startTime := endTime.Add(time.Duration(-24) * time.Hour)
	filter := fmt.Sprintf(
		"eventTimestamp ge '%s' and eventTimestamp le '%s' and category eq '%s'",
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		"ServiceHealth")

	// Get the Events from Azure.
	result, err := logsClient.List(context.Background(), filter, "")

	if err != nil {
		panic(err)
	}

	for result.NotDone() {
		for _, eventData := range result.Values() {
			if *eventData.Status.Value != "Resolved" && (eventData.Level == insights.Critical || eventData.Level == insights.Error) {
				ct.Foreground(ct.Red, false)
			} else if *eventData.Status.Value == "Resolved" {
				ct.Foreground(ct.Green, false)
			} else {
				ct.Foreground(ct.White, false)
			}

			fmt.Println(fmt.Sprintf(
				"%s - %s - %s",
				eventData.EventTimestamp.Local(),
				*eventData.ResourceProviderName.Value,
				*eventData.OperationName.Value))
			fmt.Println(fmt.Sprintf("Status:\t %s", *eventData.Status.Value))
			fmt.Println(fmt.Sprintf("Level:\t %s", eventData.Level))
			fmt.Println(fmt.Sprintf("CorrelationId:\t %s", *eventData.CorrelationID))
			if eventData.ResourceType.Value != nil {
				fmt.Println(fmt.Sprintf("Resource Type:\t %s", *eventData.ResourceType.Value))
			}
			fmt.Println(fmt.Sprintf("Description:\t %s", *eventData.Description))
		}

		// Get more events if available.
		result.NextWithContext(context.Background())
	}

	ct.Foreground(ct.Green, false)
	fmt.Println("No more events...")
	ct.Foreground(ct.White, false)
}
