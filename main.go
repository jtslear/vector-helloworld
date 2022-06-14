package main

import (
	"context"
	"log"
	"os"

	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
)

func main() {

	log.Printf("Welcome!\n")

	v, err := vector.New(
		vector.WithTarget(os.Getenv("BOT_TARGET")),
		vector.WithToken(os.Getenv("BOT_TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Script configured to talk to Robot.\n")
	}

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	log.Printf("Starting up behavior control...\n")
	go func() {
		err = v.BehaviorControl(ctx, start, stop)

		if err != nil {
			log.Fatalf("Failed to gain control... %v\n", err)
		} else {
			log.Printf("Behavior control begin...\n")
		}
	}()

	anim := &vectorpb.AnimationTrigger{Name: "GreetAfterLongTime"}

	for {
		select {
		case <-start:
			animList, _ := v.Conn.ListAnimationTriggers(
				ctx, &vectorpb.ListAnimationTriggersRequest{},
			)
			log.Printf("List of triggers:\n %v\n", animList)

			log.Printf("drive off charger...\n")
			_, _ = v.Conn.DriveOffCharger(
				ctx,
				&vectorpb.DriveOffChargerRequest{},
			)

			log.Printf("play hello animation...\n")
			_, _ = v.Conn.PlayAnimationTrigger(
				ctx,
				&vectorpb.PlayAnimationTriggerRequest{
					AnimationTrigger: anim,
					Loops:            1,
				},
			)

			//log.Printf("look for me...\n")
			//robot.behavior.find_faces()

			log.Printf("say start message...\n")
			_, _ = v.Conn.SayText(
				ctx,
				&vectorpb.SayTextRequest{
					Text:           "Hey Skarbek!  It's time to deploy!  Enjoy it!",
					UseVectorVoice: true,
					DurationScalar: 1.0,
				},
			)

			log.Printf("drive on charger...\n")
			_, _ = v.Conn.DriveOnCharger(
				ctx,
				&vectorpb.DriveOnChargerRequest{},
			)
			stop <- true
			return
		}
	}

}
