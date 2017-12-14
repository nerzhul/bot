package internal

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func runTwitterClient() {
	config := oauth1.NewConfig(gconfig.Twitter.ConsumerKey, gconfig.Twitter.ConsumerSecret)
	token := oauth1.NewToken(gconfig.Twitter.Token, gconfig.Twitter.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	tweets, resp, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 20,
	})

	if err != nil {
		log.Errorf("Failed to get home timeline: %s", err.Error())
		return
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get home timeline: %d %s", resp.StatusCode, resp.Status)
		return
	}

	for i, tweet := range tweets {
		log.Infof("tweet %d: %v", i, tweet)
	}
}
