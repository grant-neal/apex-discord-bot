package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

type Global struct {
	name                string
	uid                 uint32
	avatar              string
	platform            string
	level               uint32
	toNextLevelPercent  uint32
	internalUpdateCount uint32
	bans                Bans
	rank                Rank
	battlepass          Battlepass
	badges              string
}

type Bans struct {
	isActive         bool
	remainingSeconds uint32
	last_banReason   string
}

type Rank struct {
	RankScore         uint32
	RankName          string
	RankDiv           int32
	ladderPosPlatform int32
	rankImg           string
	rankedSeason      string
}

type Arena struct {
	rankScore         uint32
	rankName          string
	rankDiv           uint32
	ladderPosPlatform int32
	rankImg           string
	rankedSeason      string
}

type Battlepass struct {
	level   string
	history History
}

type History struct {
	season1  int32
	season2  int32
	season3  int32
	season4  int32
	season5  int32
	season6  int32
	season7  int32
	season8  int32
	season9  int32
	season10 int32
	season11 int32
	season12 int32
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

	fmt.Println(command[0])

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

			var global Global
			json.Unmarshal([]byte(string(body)), &global)

			fmt.Println("Name: %s", "Rank: %s", global.name, global.rank.RankScore)

			// s.ChannelMessageSend(m.ChannelID, string(body))

			// fmt.Println(string(body))
		}
	}
}
