package internal

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot"
	"time"
)

func runTwitterClient() {
	lastTweetID := int64(0)
	for {
		config := oauth1.NewConfig(gconfig.Twitter.ConsumerKey, gconfig.Twitter.ConsumerSecret)
		token := oauth1.NewToken(gconfig.Twitter.Token, gconfig.Twitter.TokenSecret)
		httpClient := config.Client(oauth1.NoContext, token)

		// Twitter client
		client := twitter.NewClient(httpClient)

		htp := &twitter.HomeTimelineParams{
			Count: 20,
		}

		if lastTweetID != 0 {
			htp.SinceID = lastTweetID
		}

		tweets, resp, err := client.Timelines.HomeTimeline(htp)

		if err != nil {
			log.Errorf("Failed to get home timeline: %s", err.Error())
			time.Sleep(time.Second * 60)
			continue
		}

		if resp.StatusCode != 200 {
			log.Errorf("Failed to get home timeline: %d %s", resp.StatusCode, resp.Status)
			time.Sleep(time.Second * 60)
			continue
		}

		for _, tweet := range tweets {
			tm := &bot.TweetMessage{
				Message:        tweet.Text,
				Username:       tweet.User.Name,
				UserScreenName: tweet.User.ScreenName,
				Date:           tweet.CreatedAt,
			}

			if !verifyPublisher() {
				log.Errorf("Failed to verify publisher, ignoring current messages.")
				break
			}

			rabbitmqPublisher.Publish(
				tm,
				"tweet",
				&bot.EventOptions{
					CorrelationID: uuid.NewV4().String(),
					ExpirationMs:  3600000,
				},
			)

			// Try to update the last ID only after publishing
			if tweet.ID > lastTweetID {
				lastTweetID = tweet.ID
			}
		}

		time.Sleep(time.Second * 60)
	}
}
