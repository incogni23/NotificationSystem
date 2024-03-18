package main

import (
	"context"
	"fmt"

	"consumer"

	"github.com/pikapika/database"
	service "github.com/pikapika/notification_services"
	clock "github.com/pikapika/timer"

	"github.com/project001/producer"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDB() (*gorm.DB, error) {
	dsn := "host=localhost user=pikapika password=Ankita@2307 dbname=db1 port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := setupDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}

	database.SetDB(db)

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

		errr := db.Model(&producer.Event{}).Create(&failedEvent).Error

		if errr != nil {
			fmt.Println("Failed to insert dummy data:", errr)
			return
		}
		fmt.Println("Dummy data inserted successfully")

		clock.RetryClock()
	}

	//for {
	//time.Sleep(time.Minute)
	//}

}
