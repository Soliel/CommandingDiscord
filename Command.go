package CommandingDiscord

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type (
	Context struct {
		Msg     *discordgo.MessageCreate
		Session *discordgo.Session
		Guild   *discordgo.Guild
		Channel *discordgo.Channel
		Args    []string
	}

	command struct {
		cmdFunc  func(Context)
		cooldown int
	}

	commandMap map[string]command

	CommandHandler struct {
		commands  commandMap
		truncate  bool
		truncLen  int
		cooldowns cooldownMap
		prefix    string
	}
)

//TODO: Export settings to an external file and pull from it to find truncate and truncate length.

func NewCommandHandler() (error, *CommandHandler) {
	commandHandler := &CommandHandler{make(commandMap), false, 0, make(cooldownMap), "~"}

	if commandHandler == nil {
		return errors.New("unable to create command handler"), nil
	} else {
		return nil, commandHandler
	}
}

func (handler CommandHandler) ToString() string {
	handlerstring := fmt.Sprintf("Prefix: %s\nTruncate: %t\nTruncate Length: %d\ncommandMap: %v\ncooldownMap: %v",
		handler.prefix,
		handler.truncate,
		handler.truncLen,
		handler.commands,
		handler.cooldowns)

	return handlerstring
}

func (handler CommandHandler) IsTrunc() bool {
	return handler.truncate
}

func (handler CommandHandler) GetTruncLength() int {
	return handler.truncLen
}

func (handler CommandHandler) SetTruncLength(newlength int) error {
	handler.truncLen = newlength
	if handler.truncLen != newlength {
		return errors.New("unable to set new trunclength")
	} else {
		return nil
	}
}

func (handler CommandHandler) GetCommands() commandMap {
	return handler.commands
}

func (handler CommandHandler) Get(name string) (*command, bool) {
	cmd, found := handler.commands[name]
	return &cmd, found
}

func (handler CommandHandler) Register(name string, cmd func(Context), cooldown int) error {
	_, exists := handler.commands[name]
	if exists {
		return errors.New("cannot register multiple commands to a single name. The cmd: " + name + " will not be registered.")
	}

	handler.commands[name] = command{cmd, cooldown}
	if handler.truncate {
		if len(name) < handler.truncLen {
			handler.commands[name[:handler.truncLen]] = command{cmd, cooldown}
		}
	}
	return nil
}

func (cmd command) HasCooldown() bool {
	if cmd.cooldown > 0 {
		return true
	}
	return false
}
