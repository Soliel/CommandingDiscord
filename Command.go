package CommandingDiscord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
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
		cdTick    *time.Ticker
	}
)

//TODO: Export settings to an external file and pull from it to find truncate and truncate length.

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{make(commandMap), false, 0, make(cooldownMap), time.NewTicker(time.Second)}
}

func (handler CommandHandler) isTrunc() bool {
	return handler.truncate
}

func (handler CommandHandler) getTruncLength() int {
	return handler.truncLen
}

func (handler CommandHandler) getCommands() commandMap {
	return handler.commands
}

func (handler CommandHandler) get(name string) (*command, bool) {
	cmd, found := handler.commands[name]
	return &cmd, found
}

func (handler CommandHandler) register(name string, cmd func(Context), cooldown int) {
	_, exists := handler.commands[name]
	if exists {
		fmt.Println("Cannot register multiple commands to a single name. The cmd: ", name, " will not be registered.")
		return
	}

	handler.commands[name] = command{cmd, cooldown}
	if handler.truncate {
		if len(name) < handler.truncLen {
			handler.commands[name[:handler.truncLen]] = command{cmd, cooldown}
		}
	}
}

func (cmd command) hasCooldown() bool {
	if cmd.cooldown > 0 {
		return true
	}
	return false
}
