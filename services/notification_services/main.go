package main

import (
	"context"
	"fmt"
	"sync"

	database "github.com/pikapika/database"
	clock "github.com/pikapika/notification_services/timer"
	service "github.com/pikapika/notification_services/webhook"
	"github.com/project001/producer"
)

func main() {
	db, err := database.SetupDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}

	err = database.Initialize(db)
	if err != nil {
		fmt.Println("Failed to auto migrate:", err)
		return
	}
	fmt.Println("Database auto migration completed successfully")
	produce1 := producer.Newwriter()
	someevent := producer.Event{Message: []byte{34, 97, 98, 99, 100, 101, 34},
		Source:             "source1",
		DestinationAddress: "https://95b4c",
		Topics:             []string{"quickstart-events"},
		NotificationType:   "webhook"}

	checkerror := produce1.Produce(context.Background(), &someevent)
	fmt.Print(checkerror)
	if checkerror != nil {
		fmt.Println("Error producing message:", checkerror)
		return
	}

	fmt.Println("Message produced successfully")
	go service.Consuming()

	go clock.RetryClock()

	var wg sync.WaitGroup
	wg.Wait()
	//
	//fmt.Println("checking go routine")

	//for {
	//		time.Sleep(time.Minute)
	//		time.Sleep(time.Minute)

	//}

}
