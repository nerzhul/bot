package internal

import (
	"github.com/nlopes/slack"
	"github.com/satori/go.uuid"
)

var slackAPI *slack.Client
var slackRTM *slack.RTM

func runSlackClient() {
	log.Infof("Starting slack client.")
	slackAPI = slack.New(gconfig.Slack.APIKey)
	slackRTM = slackAPI.NewRTM()

	go slackRTM.ManageConnection()

	for msg := range slackRTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
		case *slack.ConnectedEvent:
			break
		case *slack.MessageEvent:
			// Ignore non command
			if len(ev.Text) < 2 || ev.Text[0] != '!' {
				break
			}

			event := commandEvent{
				ev.Text[1:],
				ev.Channel,
				ev.User,
			}

			log.Infof("User %s sent command on channel %s: %s", event.User, event.Channel, event.Command)

			if !verifyPublisher() {
				log.Error("Failed to verify publisher, no command sent to broker")
				break
			}

			if !verifyConsumer() {
				log.Error("Failed to verify consumer, no command sent to broker")
				break
			}

			consumerCfg := gconfig.RabbitMQ.GetConsumer("commands")
			if consumerCfg == nil {
				log.Fatalf("RabbitMQ consumer configuration 'commands' not found, aborting.")
			}

			rabbitmqPublisher.Publish(
				&event,
				"command",
				uuid.NewV4().String(),
				consumerCfg.RoutingKey,
				300000,
			)
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
