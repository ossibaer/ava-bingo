package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	BotPrefix     string
	bingoSessions map[string]*Bingo
	bingoMessages map[string]string
)

func Start(token string) {
	goBot, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatal(err.Error())
	}

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Bot is running fine!")

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	message, isCommand := strings.CutPrefix(m.Content, BotPrefix)
	if !isCommand {
		return
	}

	var i int
	for i = 0; i < len(message); i++ {
		if message[i] == ' ' {
			break
		}
	}
	command := message[:i]
	args := strings.Split(message[i:], ",")

	switch command {
	case "start":
		if bingoSessions != nil {
			s.ChannelMessageSend(m.ChannelID, "There currently is an active Bingo-Session! You stop it with `!stop`.")
		}
		newBingoSession()
		s.ChannelMessageSend(m.ChannelID, "New Bingo-Session started! You can join with `!join`.")
	case "stop":
		stop()
		s.ChannelMessageSend(m.ChannelID, "Bingo-Session stopped! You can start a new one with `!start`.")
	case "join":
		if bingoSessions == nil {
			s.ChannelMessageSend(m.ChannelID, "There currently is no active Bingo-Session! You can start a new one with `!start`.")
		}
		joinBingoSession(m.Author, s)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint(m.Author.Mention(), " joined! Check your DMs for your Bingo-Card."))
	case "strike":
		if bingoSessions == nil {
			s.ChannelMessageSend(m.ChannelID, "There currently is no active Bingo-Session! You can start a new one with `!start`.")
		}
		strikeThrough(args, s, m.ChannelID)
	case "leave":
		if bingoSessions == nil {
			s.ChannelMessageSend(m.ChannelID, "There currently is no active Bingo-Session! You can start a new one with `!start`.")
		}
		leaveBingoSession(m.Author)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint(m.Author.Mention(), " left the Bingo-Session."))
	default:
		log.Println("Received unknown command", command, "from", m.Author.Username)
	}
}

func newBingoSession() {
	bingoSessions = make(map[string]*Bingo)
	bingoMessages = make(map[string]string)
}

func stop() {
	bingoSessions = nil
	bingoMessages = nil
}

func joinBingoSession(user *discordgo.User, ds *discordgo.Session) {
	b := NewBingo()
	bingoSessions[user.ID] = &b
	channel, err := ds.UserChannelCreate(user.ID)
	if err != nil {
		log.Println("Could not create channel to user ", user.ID)
	}
	message, err := ds.ChannelMessageSend(channel.ID, fmt.Sprint("```", b.String(), "```"))
	if err != nil {
		log.Println("Could not send Bingo Card to user ", user.ID)
	}
	bingoMessages[user.ID] = message.ID
}

func leaveBingoSession(user *discordgo.User) {
	delete(bingoSessions, user.ID)
}

func strikeThrough(items []string, ds *discordgo.Session, channelId string) {
	if len(items) < 1 {
		return
	}
	bingoCount := make(map[string]int)
	for _, s := range items {
		for k, v := range bingoSessions {
			bingoCount[k] += v.cross(strings.Trim(s, " "))
		}
	}
	for k, v := range bingoCount {
		userChannel, err := ds.UserChannelCreate(k)
		if err != nil {
			log.Println("Could not create channel to user ", k)
		}

		ds.ChannelMessageEdit(userChannel.ID, bingoMessages[k], fmt.Sprint("```", bingoSessions[k].String(), "```"))

		if v == 0 {
			continue
		}

		user, err := ds.User(k)
		if err != nil {
			log.Println("Could not get user ", k)
		}
		ds.ChannelMessageSend(
			channelId,
			fmt.Sprint(user.Mention(), " just got ", v, " Bingo(s)!"),
		)
	}
}

func main() {
	if _, err := os.Stat(".env"); !errors.Is(err, os.ErrNotExist) {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file. Ignore this, if there is not .env file.")
		}
	}

	BotPrefix = os.Getenv("AVA_BOT_PREFIX")
	botToken := os.Getenv("AVA_BOT_TOKEN")

	Start(botToken)

	<-make(chan struct{})
}
