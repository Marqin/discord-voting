package main

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
)

func turnout() (int, error) {

	turnout := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("votes"))

		if b != nil {
			turnout = b.Stats().KeyN
		}

		return nil
	})

	return turnout, err
}

func turnoutCommand(s *discordgo.Session, m *discordgo.Message) (string, error) {

	turnout, err := turnout()
	if err != nil {
		return "", err
	}

	allowedToVote := 0

	g, err := s.Guild(conf.Guild)
	if err != nil {
		return "", err
	}

	for _, member := range g.Members {
		for _, s := range member.Roles {
			if conf.Role == s {
				allowedToVote++
			}
		}
	}

	message := fmt.Sprintf("%d/%d (%d%%) Crew members voted.", turnout, allowedToVote, (100*turnout)/allowedToVote)

	return message, err
}
