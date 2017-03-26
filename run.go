package main;

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/legolord208/stdutil"
	"os"
	"os/signal"
	"strings"
)

var ownID string;

func doRun(args []string){
	if(len(args) < 1){
		stdutil.PrintErr("No token supplied in arguments", nil);
		return;
	}
	token := args[0];

	loadRules();
	fmt.Println("Starting...");

	session, err := discordgo.New(token);
	if(err != nil){
		stdutil.PrintErr("Couldn't initialize bot", err);
		return;
	}

	me, err := session.User("@me");
	if(err != nil){
		stdutil.PrintErr("Couldn't fetch @me", err);
		return;
	}

	ownID = me.ID;
	session.AddHandler(messageCreate);

	err = session.Open();
	if(err != nil){
		stdutil.PrintErr("Couldn't start bot", err);
		return;
	}

	fmt.Println("Started!");

	c := make(chan os.Signal, 1);
	signal.Notify(c, os.Interrupt);

	for _ = range c{
		fmt.Println("Closing!");
		session.Close();
		return;
	}
}

func messageCreate(session *discordgo.Session, e *discordgo.MessageCreate){
	message(session, e.Message);
}
func messageUpdate(session *discordgo.Session, e *discordgo.MessageUpdate){
	message(session, e.Message);
}

func message(session *discordgo.Session, e *discordgo.Message){
	if(e.Author == nil){ return; }
	if(e.Author.ID == ownID){ return; }

	for _, rule := range rules{
		content := strings.ToLower(strings.TrimSpace(e.Content));
		if(rule.Exact && content != rule.Msg){
			continue;
		} else if(!rule.Exact && !strings.Contains(content, rule.Msg)){
			continue;
		}

		if(len(rule.NotFrom) > 0){
			for _, from := range rule.NotFrom{
				if(from == e.Author.ID){
					continue;
				}
			}
		}
		if(len(rule.From) > 0){
			found := false;
			for _, from := range rule.From{
				if(from == e.Author.ID){
					found = true;
				}
			}

			if(found){
				continue;
			}
		}

		_, err := session.ChannelMessageSend(e.ChannelID, rule.Reply);
		if(err != nil){
			stdutil.PrintErr("Could not send message", err);
		}
		break;
	}
}
