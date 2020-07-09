package mapper

import (
	"net/url"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

func mapURL(s string) string {
	s = url.PathEscape(s)
	s = strings.ReplaceAll(s, url.PathEscape(","), ",")

	return "https://www.google.com/maps/dir//" + s
}

func nick(m *dg.Message) string {
	if m.Member != nil {
		if m.Member.Nick != "" {
			return m.Member.Nick
		}

		if m.Member.User != nil {
			return m.Member.User.Username
		}
	}

	if m.Author != nil {
		return m.Author.Username
	}

	return "NAME UNKNOWN"
}
