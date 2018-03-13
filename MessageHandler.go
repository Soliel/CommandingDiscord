package CommandingDiscord

import (
	"strings"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

//This is an optional message handler and must be wrapped into a handler with scope allowing for BotID and handler to be passed in without explicitly stating it in the onMessageRecieved Function
/*
Ideally the implementation will be something along the lines of

var (
	BotID = "blah blah"
	handler = CommandHandler
)

func onMessageRecieved(s *discordgo.Session, m *discordgo.MessageCreate) {
	handleMessages(s, m, BotID, handler)
}

 */
func handleMessages(s *discordgo.Session, m *discordgo.MessageCreate, BotID string, handler CommandHandler) {

	if m.Author.ID == BotID {
		return
	}

	if len(m.Content) < len(handler.prefix) {
		return
	}

	if m.Content[:len(handler.prefix)] != handler.prefix {
		return
	}

	content := m.Content[len(handler.prefix):]
	if len(content) < 1 {
		return
	}

	content = strings.ToLower(content)

	var args []string

	//If someones name has spaces this allows them to fix it.
	if strings.Contains(content, "\"") {
		tempArgs := strings.Split(content, "\"")
		for s := range tempArgs {
			tempArgs[s] = strings.TrimSpace(tempArgs[s])
			if tempArgs[s] != "" {
				args = append(args, tempArgs[s])
			}
		}

		args[0] = strings.TrimPrefix(args[0], " ")
		args[0] = strings.TrimSuffix(args[0], " ")
		if strings.Contains(args[0], " ") {
			firstArgs := strings.Fields(args[0])
			args = append(firstArgs[:], args[1:]...)
		}
	} else {
		args = strings.Fields(content)
	}
	name := args[0]

	command, found := handler.get(name)
	if !found {
		return
	}

	if command.hasCooldown() {
		if handler.isOnCooldown(m.Author.ID, *command) {
			return
		}
		handler.startCooldown(m.Author.ID, *command)
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel, ", err)
		return
	}

	//set up my context to pass to whatever function is called.
	ctx := new(Context)
	ctx.Args = args[1:]
	ctx.Session = s
	ctx.Msg = m
	ctx.Channel = channel

	guild, err := s.State.Guild(channel.GuildID)
	if err == nil {
		ctx.Guild = guild
	}


	//pass command pointer and run the function
	c := command.cmdFunc
	go c(*ctx)
}