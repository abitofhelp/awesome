package main

import (
	"fmt"
	messagesv1 "go.buf.build/protocolbuffers/go/abitofhelp/abcdapis/messages/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func main() {

	mypet := &messagesv1.Pet{
		BirthdayUtc: timestamppb.New(time.Now()),
		Name:        "Lassie",
	}

	fmt.Printf("\nPet: name '%s', birthday: '%s'\n",
		mypet.Name, mypet.BirthdayUtc.AsTime().Format(time.RFC3339))

}
