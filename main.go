package main

import (
	"context"
	"fmt"

	"consumer"

	database "github.com/pikapika/database"
	"github.com/pikapika/models"
	service "github.com/pikapika/notification_services"
	clock "github.com/pikapika/timer"

	"github.com/project001/producer"
)

func consumertoDB(p *consumer.Event) *models.Event {
	d := models.Event{
		Message:            p.Message,
		Source:             p.Source,
		DestinationAddress: p.DestinationAddress,
		Topics:             p.Topics,
		NotificationType:   p.NotificationType,
		Status:             p.Status,
		Attempts:           p.Attempts,
		NextRetry:          p.NextRetry,
		IdempotencyKey:     p.IdempotencyKey}
	return &d
}

func main() {
	_, err := database.SetupDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}

	err = database.Initialize()
	if err != nil {
		fmt.Println("Failed to auto migrate:", err)
		return
	}
	fmt.Println("Database auto migration completed successfully")

	dummyEvent := producer.Event{
		Message:            []byte{34, 97, 98, 99, 100, 101, 29, 83, 85, 34},
		Source:             "source1",
		DestinationAddress: "https://webhook-test.com/16362ae996d3c6a3212a3ccb712459cb",
		Topics:             []string{"quickstart-events", "topic2"},
		IdempotencyKey:     "",
		NotificationType:   "abc",
	}

	produce1 := producer.Newwriter()
	var tp string = "quickstart-events"
	checkerror := produce1.Produce(context.Background(), &dummyEvent)

	if checkerror != nil {
		fmt.Println("Error producing message:", checkerror)
		return
	}
	fmt.Println("Message produced successfully")
	brokeraddress := []string{"localhost:9092"}
	consume1 := consumer.NewReader(brokeraddress, tp)
	failedEvent := consume1.Consume()
	//service.Automate() for loop

	err1 := service.Deliver(failedEvent, failedEvent.NotificationType)
	if err1 != nil {
		fmt.Print("failed to deliver the event", err1)
		dbEvent := consumertoDB(&failedEvent)
		dbEvent.Status = "Not Completed Yet"
		errr := database.CreateEvent(dbEvent) //

		if errr != nil {
			fmt.Printf("Failed to insert dummy data:%s", errr)
			return
		}

		fmt.Println("Dummy data inserted successfully")

	}
	clock.RetryClock()

	fmt.Println("checking go routine")

	//for {
	//time.Sleep(time.Minute)
	//}

}
