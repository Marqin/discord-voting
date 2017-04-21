package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	botURL    = "https://discordapp.com/oauth2/authorize?client_id=%v&scope=bot&permissions=0"
	cmdPrefix = "."
)

func startBot(conf *Config) error {
	bot, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		return err
	}

	bot.AddHandler(messageCreate)

	err = bot.Open()
	if err != nil {
		return err
	}

	err = bot.UpdateStatus(0, "Message me!")
	if err != nil {
		return err
	}

	u, err := bot.User("@me")
	if err != nil {
		return err
	}

	fmt.Println("Invite bot via: " + fmt.Sprintf(botURL, u.ID))

	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	me, err := s.User("@me")
	if err != nil {
		log.Println(err)
		return
	}

	if me.ID == m.Author.ID {
		return
	}

	splittedContent := strings.Fields(m.Content)

	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}

	if len(splittedContent) >= 1 {
		if strings.HasPrefix(splittedContent[0], cmdPrefix) || ch.IsPrivate { // ch.IsPrivate not needed here atm
			command := getCommand(strings.TrimPrefix(splittedContent[0], cmdPrefix), ch.IsPrivate)
			if command != nil {
				err := command(s, m.Message)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

}
