package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	discord "github.com/bwmarrin/discordgo"
	"github.com/ggrpmnn/pixie"
)

func init() {
	// check that a Discord token has been supplied
	if pixie.DiscordToken() == "" {
		log.Fatalf("Discord access token not found; exiting")
	}
}

func main() {
	// create and open Discord session
	session, err := discord.New(pixie.DiscordToken())
	if err != nil {
		log.Fatalf("Failed to start Discord session: %s", err.Error())
	}
	err = session.Open()
	if err != nil {
		log.Fatalf("Error opening connection to Discord server: %s", err.Error())
	}

	// add handlers for Discord events
	session.AddHandler(messageCreate)

	// run until an interrupt is sent
	log.Printf("HEY! LISTEN! Pixie is now online; running until told otherwise...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sigChan
	log.Printf("Received shutdown command; bye!")

	// close the session
	err = session.Close()
	if err != nil {
		log.Fatalf("Error closing session: %s", err.Error())
	}

	return
}

// messageCreate handles events triggered by the creation of a message in Discord
func messageCreate(s *discord.Session, m *discord.MessageCreate) {
	// check whether the bot should respond to the message
	botPrefix := "!" + pixie.BotName() + " "
	if !strings.HasPrefix(m.Content, botPrefix) {
		return
	}
	cmdStr := strings.TrimPrefix(m.Content, botPrefix)
	log.Printf("Received command request from '%s': %s", m.Author, cmdStr)

	// attempt to run the supplied command
	resp, err := pixie.Run(cmdStr)
	if err != nil {
		log.Printf(err.Error())
		sendReply(s, m, err.BotMessage())
		return
	}

	// resp has already been formatted; send it
	sendReply(s, m, resp)
	return
}

// sendReply sends a reply message to the server
func sendReply(s *discord.Session, m *discord.MessageCreate, msg string) {
	log.Printf("Sending reply to channel")
	reply := fmt.Sprintf("<@%s> %s", m.Author.ID, msg)
	_, err := s.ChannelMessageSend(m.ChannelID, reply)
	if err != nil {
		log.Printf("Failed to send reply to channel: %s", err.Error())
	}
}
