package CommandingDiscord

import (
	"time"
)

type (
	cooldownMap map[string]map[string]time.Time //Stores a map of each users current cooldowns.
)

func (handler CommandHandler) IsOnCooldown(user string, cmd string) bool {
	if handler.cooldowns[user][cmd].IsZero() {
		return false
	}
	if handler.cooldowns[user][cmd].Sub(time.Now()) >= 0 {
		return true
	}
	return false
}

func (handler CommandHandler) StartCooldown(user string, cmd string) {
	handler.cooldowns[user][cmd] = time.Now().Add(time.Duration(handler.commands[cmd].cooldown) * time.Second)
}

func (handler CommandHandler) StartCooldownTicker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for tick := range ticker.C {
		for key, value := range handler.cooldowns {
			for key2, value2 := range value {
				if tick.Sub(value2) >= 0 {
					delete(handler.cooldowns[key], key2)
					if len(handler.cooldowns[key]) <= 0 {
						delete(handler.cooldowns, key)
					}
				}
			}
		}
	}
}
