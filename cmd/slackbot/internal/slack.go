package internal

import (
	"github.com/nlopes/slack"
	"github.com/satori/go.uuid"
)

func runSlackClient() {
	api := slack.New(gconfig.Slack.ApiKey)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
		case *slack.ConnectedEvent:
			break
		case *slack.MessageEvent:
			// Ignore non command
			if ev.Text[0] != '!' {
				break
			}

			event := commandEvent{
				ev.Text[1:],
				ev.Channel,
				ev.User,
			}

			log.Infof("User %s sent command on channel %s: %s", event.User, event.Channel, event.Command)

			if !verifyPublisher() {
				log.Error("Failed to publish command event")
				break
			}

			rabbitmqPublisher.Publish(&event, "command", uuid.NewV4().String(), gconfig.RabbitMQ.EventRoutingKey)
			break
		case *slack.PresenceChangeEvent:
		case *slack.LatencyReport:
		case *slack.RTMError:
		case *slack.InvalidAuthEvent:
			break
		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
