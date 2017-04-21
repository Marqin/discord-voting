package main

import (
	"github.com/bwmarrin/discordgo"
)

type msgCommand func(s *discordgo.Session, m *discordgo.Message) (string, error)
type command func(s *discordgo.Session, m *discordgo.Message) error

func getCommand(s string, private bool) command {

	if private {
		switch s {
		case "vote":
			return simpleRespondCommand(voteCommand)
		case "status":
			return simpleRespondCommand(statusCommand)
		case "help":
			return simpleRespondCommand(helpCommand)
		case "info":
			return simpleRespondCommand(info)
		default:
			return simpleRespondCommand(dummyCommand)
		}
	} else {
		switch s {
		case "turnout":
			return simpleRespondCommand(turnoutCommand)
		default:
			return nil
		}
	}

}

func simpleRespondCommand(cmd msgCommand) command {

	return func(s *discordgo.Session, m *discordgo.Message) error {
		msg, err := cmd(s, m)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "There was an error :-(")
			return err
		}
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return err
	}

}

func helpCommand(s *discordgo.Session, m *discordgo.Message) (string, error) {
	helpStr := "Usage:```help - this help\ninfo - give link to rules and outfit list\nstatus - show status of your votes\nvote HERE_PUT_OUTFIT_NAME - vote for selected outfit\n```"
	return helpStr, nil
}

func dummyCommand(s *discordgo.Session, m *discordgo.Message) (string, error) {
	return "Are you lost? Use command `help`.", nil
}

func info(s *discordgo.Session, m *discordgo.Message) (string, error) {
	return "To find rules and outfit list, go to: " + conf.RulesURL, nil
}
