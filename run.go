package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/legolord208/stdutil"
	"github.com/legolord208/timeouts"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var ownID string
var TIMEOUT timeouts.Timeout

const DELAY_PASSIVE = 3
const DELAY_AGRESSIVE = 5

func doRun(args []string) {
	if len(args) < 1 {
		stdutil.PrintErr("No token supplied in arguments", nil)
		return
	}
	token := args[0]

	loadRules()
	TIMEOUT = timeouts.NewTimeout()
	fmt.Println("Starting...")

	session, err := discordgo.New(token)
	if err != nil {
		stdutil.PrintErr("Couldn't initialize bot", err)
		return
	}

	me, err := session.User("@me")
	if err != nil {
		stdutil.PrintErr("Couldn't fetch @me", err)
		return
	}

	ownID = me.ID
	session.AddHandler(messageCreate)

	err = session.Open()
	if err != nil {
		stdutil.PrintErr("Couldn't start bot", err)
		return
	}

	fmt.Println("Started!")

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("\nClosing!")
	session.Close()
}

func messageCreate(session *discordgo.Session, e *discordgo.MessageCreate) {
	message(session, e.Message)
}
func messageUpdate(session *discordgo.Session, e *discordgo.MessageUpdate) {
	message(session, e.Message)
}

func message(session *discordgo.Session, e *discordgo.Message) {
	if e.Author == nil {
		return
	}
	if e.Author.ID == ownID {
		return
	}

	if TIMEOUT.InTimeout(e.Author.ID) {
		TIMEOUT.SetTimeout(e.Author.ID, time.Duration(DELAY_AGRESSIVE)*time.Second)
		return
	}
	TIMEOUT.SetTimeout(e.Author.ID, time.Duration(DELAY_PASSIVE)*time.Second)

	content := strings.ToLower(strings.TrimSpace(e.Content))
rules:
	for _, rule := range rules {
		if rule.Exact && content != rule.Msg {
			continue
		} else if !rule.Exact && !strings.Contains(content, rule.Msg) {
			continue
		}

		// Channel checking
		for _, c := range rule.NotInChannel {
			if c == e.ChannelID {
				continue rules
			}
		}
		if len(rule.InChannel) > 0 {
			found := false
			for _, c := range rule.InChannel {
				if c == e.ChannelID {
					found = true
				}
			}

			if !found {
				continue
			}
		}

		// User checking
		for _, from := range rule.NotFrom {
			if from == e.Author.ID {
				continue rules
			}
		}
		if len(rule.From) > 0 {
			found := false
			for _, from := range rule.From {
				if from == e.Author.ID {
					found = true
				}
			}

			if !found {
				continue
			}
		}

		_, err := session.ChannelMessageSend(e.ChannelID, rule.Reply)
		if err != nil {
			stdutil.PrintErr("Could not send message", err)
		}
		break
	}
}
