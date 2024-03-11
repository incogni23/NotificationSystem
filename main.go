package main

import (
	"consumer"
	"fmt"

	"github.com/project001/producer"
	"golang.org/x/net/context"
)

func main() {
	var topicname string = "hellloo"
	produce1 := producer.Newwriter(topicname)
	someevent := producer.Events{Msg: "hello yoyooyo"}
	checkerror := produce1.Produce(context.Background(), someevent)

	fmt.Print(checkerror)
	brokeraddress := []string{"localhost:9092"}

	consume1 := consumer.NewReader(brokeraddress, topicname)
	consume1.Consume()

}
