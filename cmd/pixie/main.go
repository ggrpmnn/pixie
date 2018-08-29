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
	botCmd := "!" + pixie.BotName() + " "
	if !strings.HasPrefix(m.Content, botCmd) {
		return
	}
	cmdStr := strings.TrimPrefix(m.Content, botCmd)
	log.Printf("Received request from '%s': %s", m.Author, cmdStr)

	reply := fmt.Sprintf("<@%s> %s! Got it!", m.Author.ID, cmdStr)
	_, err := s.ChannelMessageSend(m.ChannelID, reply)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}
