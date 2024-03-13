package main

import (
	"fmt"

	service "github.com/project001/Services/notification_service"
	"github.com/project001/producer"
	"golang.org/x/net/context"
)

func main() {

	produce1 := producer.Newwriter()
	//var tp string = "quickstart-events"

	someevent := producer.Event{Message: []byte{34, 97, 98, 99, 100, 101, 34},
		Source:             "source1",
		DestinationAddress: "https://webhook.site/26408297-d7e0-4400-8031-669b26e66d13",
		Topics:             []string{"quickstart-events"},
		NotificationType:   "webhook"}

	checkerror := produce1.Produce(context.Background(), someevent)

	//brokeraddress := []string{"localhost:9092"}

	//consume1 := consumer.NewReader(brokeraddress, tp)

	//consume1.Consume()

	fmt.Print(checkerror)
	if checkerror != nil {
		fmt.Println("Error producing message:", checkerror)
		return
	}

	fmt.Println("Message produced successfully")
	service.Automate()

}
