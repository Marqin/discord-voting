package main

import (
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
)

type userVotes struct {
	Votes []string
}

func decode(data []byte) (*userVotes, error) {
	var uv *userVotes
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&uv)
	if err != nil {
		return nil, err
	}
	return uv, nil
}

func (uv *userVotes) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(uv)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func voteCommand(s *discordgo.Session, m *discordgo.Message) (string, error) {

	valid, err := validUser(s, m.Author)
	if err != nil {
		return "", err
	}

	if !valid {
		return "You are not privileged to vote!", nil
	}

	splittedContent := strings.Fields(m.Content)

	if len(splittedContent) != 2 {
		return "Wrong usage!", nil
	}

	castedVote := strings.Title(strings.ToLower(splittedContent[1]))

	validVote := false
	for _, v := range conf.Choices {
		if castedVote == strings.Title(strings.ToLower(v)) {
			validVote = true
		}
	}

	if !validVote {
		return "This is not a valid outfit!", nil
	}

	response := ""

	err = db.Update(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte("votes"))
		if err != nil {
			return err
		}

		val := b.Get([]byte(m.Author.ID))
		uv := new(userVotes)

		if val == nil {
			uv.Votes = []string{castedVote}
		} else {
			uv, err = decode(val)
			if err != nil {
				return err
			}

			if len(uv.Votes) >= conf.VotesPerUser {
				response = "You already used all your votes!"
				return nil
			}

			for _, v := range uv.Votes {
				if castedVote == strings.Title(strings.ToLower(v)) {
					response = "You have already voted on that outfit!"
					return nil
				}
			}

			uv.Votes = append(uv.Votes, castedVote)
		}

		val, err = uv.encode()
		if err != nil {
			return err
		}

		return b.Put([]byte(m.Author.ID), val)
	})

	if response == "" {
		response = "Vote send to database!"
	}

	return response, err
}
