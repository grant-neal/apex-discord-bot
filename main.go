package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tidwall/gjson"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := strings.Fields(m.Content)

	if command[0] == "!apexname" {

		// Lots more error handling can be done to stop abuse. But be sensible for now yeah?
		if len(command) == 1 {
			s.ChannelMessageSend(m.ChannelID, "No name included")
		} else {
			api_key := os.Getenv("API_KEY")
			resp, err := http.Get("https://api.mozambiquehe.re/bridge?version=5&platform=PC&player=" + command[1] + "&auth=" + api_key)
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Fatal(err)
			}

			character := make(map[string]string)
			
			characterName := gjson.Get(string(body), "global.name")
			character["Name"] = characterName.String()

			characterRankScore := gjson.Get(string(body), "global.rank.rankScore")
			character["RankScore"] = characterRankScore.String()

			characterRankName := gjson.Get(string(body), "global.rank.rankName")
			character["rankName"] = characterRankName.String()

			characterRankDiv := gjson.Get(string(body), "global.rank.rankDiv")
			character["rankDiv"] = characterRankDiv.String()

			s.ChannelMessageSend(m.ChannelID, "Name: " + character["Name"] + ", RP Score: " + character["RankScore"] + ", Current Rank: " + character["rankName"] + " " + character["rankDiv"])
		}
	}
}
