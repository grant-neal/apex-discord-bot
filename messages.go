package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

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
