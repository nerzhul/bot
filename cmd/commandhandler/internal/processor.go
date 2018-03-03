package internal

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func runProcessor() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	shouldStop := false
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Infof("Signal received, notify to stop.")
		shouldStop = true
		done <- true
	}()

	log.Infof("Starting processor.")
	for !shouldStop {
		runStep()
	}
	log.Infof("Stopping processor.")

	<-done
}

func runStep() {
	time.Sleep(time.Second * 2)
	if !asyncClient.VerifyPublisher() {
		return
	}

	if !asyncClient.VerifyConsumer() {
		return
	}
}
