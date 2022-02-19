package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/tidwall/gjson"
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

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
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

	switch command[0] {
	case "!commands":
		var commands = [2]string{"!apexname", "!maprotation"}
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("Commands", "Here are the following commands:\n"+fmt.Sprint(commands)))

	case "!apexname":
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

			character := getCharacterInfo(string(body))

			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(character["Name"], "RP Score: "+character["RankScore"]+", Current Rank: "+character["rankName"]+" "+character["rankDiv"]))
		}
	case "!maprotation":

		var gameModes = [5]string{"battle_royale", "arenas", "ranked", "arenasRanked", "control"}
		var gameMode string = "battle_royale"
		if len(command) == 2 {
			if stringInSlice(command[1], gameModes[:]) {
				gameMode = command[1]
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("Error", "Please select from the following gamme modes:\n"+fmt.Sprint(gameModes)))
				return
			}
		}

		api_key := os.Getenv("API_KEY")
		resp, err := http.Get("https://api.mozambiquehe.re/maprotation?version=2&auth=" + api_key)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		mapInfo := getMapRotation(string(body), gameMode)

		if gameMode == "ranked" {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("Map Rotation", "Current Map: "+mapInfo["currentMap"]+"\n Next Map: "+mapInfo["nextMap"]))
		} else {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("Map Rotation", "Current Map: "+mapInfo["currentMap"]+", Time Remaining: "+mapInfo["remainingTimer"]+"\n Next Map: "+mapInfo["nextMap"]))
		}
	}
}

func getMapRotation(mapRotation string, gameMode string) map[string]string {
	mapInfo := make(map[string]string)

	mapInfo["currentMap"] = gjson.Get(mapRotation, gameMode+".current.map").String()
	if gameMode != "ranked" {
		mapInfo["remainingTimer"] = gjson.Get(mapRotation, gameMode+".current.remainingTimer").String()
	}
	mapInfo["nextMap"] = gjson.Get(mapRotation, gameMode+".next.map").String()

	return mapInfo
}

func getCharacterInfo(characterInfo string) map[string]string {
	character := make(map[string]string)

	character["Name"] = gjson.Get(characterInfo, "global.name").String()
	character["RankScore"] = gjson.Get(characterInfo, "global.rank.rankScore").String()
	character["rankName"] = gjson.Get(characterInfo, "global.rank.rankName").String()
	character["rankDiv"] = gjson.Get(characterInfo, "global.rank.rankDiv").String()

	// TODO: We need to look into downloading this file so we can manipulate it and send it back as part of the embedded response.
	character["rankImage"] = gjson.Get(characterInfo, "global.rank.rankImg").String()

	return character
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
