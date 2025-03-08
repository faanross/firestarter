package main

import (
	"firestarter/internal/factory"
	"fmt"
	"time"
)

var serverPorts = []string{"7777", "8888", "9999"}

func main() {

	listenerFactory := factory.NewListenerFactory()

	for _, port := range serverPorts {
		time.Sleep(2 * time.Second)
		l, err := listenerFactory.CreateListener(port)
		if err != nil {
			fmt.Printf("Error creating service: %v\n", err)
			continue
		}
		time.Sleep(2 * time.Second)
		go func(l *factory.Listener) {

			err := l.Start()
			if err != nil {
				fmt.Printf("Error starting listener %s: %v\n", l.ID, err)
			}
		}(l)
	}

	select {}

}
