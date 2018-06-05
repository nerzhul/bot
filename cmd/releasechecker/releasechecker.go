package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/pborman/getopt/v2"
	"gitlab.com/nerzhul/bot/cmd/releasechecker/internal"
)

var configFile = ""

func init() {
	getopt.FlagLong(&configFile, "config", 'c', "Configuration file")
}

func main() {
	getopt.Parse()
	internal.StartApp(configFile)
	client := github.NewClient(nil)
	tags, response, err := client.Repositories.ListTags(context.Background(), "minetest", "minetest", nil)
	if err != nil {
		println(err)
		return
	}

	for _, t := range tags {
		println(fmt.Sprintf("%v", *t.Name))
	}
	println(fmt.Sprintf("%v", response))
}
