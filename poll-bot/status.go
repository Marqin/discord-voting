package main

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
)

func statusCommand(s *discordgo.Session, m *discordgo.Message) (string, error) {

	valid, err := validUser(s, m.Author)
	if err != nil {
		return "", err
	}

	if !valid {
		return "You are not privileged to vote!", nil
	}

	var value *[]byte // default value "nil"

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("votes"))

		if b != nil {
			value = new([]byte)
			*value = b.Get([]byte(m.Author.ID))
		}
		return nil
	})

	if value == nil || *value == nil {
		str := fmt.Sprintf("You have not voted!\nVotes left: %d", conf.VotesPerUser)
		return str, nil
	}

	uv := new(userVotes)

	uv, err = decode(*value)
	if err != nil {
		return "", err
	}

	str := fmt.Sprintf("You have already voted for: %v\nVotes left: %d", uv.Votes, conf.VotesPerUser-len(uv.Votes))
	return str, nil
}
