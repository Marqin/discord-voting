package main

import "github.com/bwmarrin/discordgo"

func validUser(s *discordgo.Session, user *discordgo.User) (bool, error) {
	g, err := s.Guild(conf.Guild)
	if err != nil {
		return false, err
	}

	for _, member := range g.Members {
		if member.User.ID == user.ID {
			for _, s := range member.Roles {
				if conf.Role == s {
					return true, nil
				}
			}
		}
	}

	return false, nil

}
