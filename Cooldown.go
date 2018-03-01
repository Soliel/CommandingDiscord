package CommandingDiscord

import (
	"time"
)

type (
	cooldownMap map[string]map[command]time.Time //Stores a map of each users current cooldowns.
)

func (handler CommandHandler) isOnCooldown(user string, cmd command) bool {
	if handler.cooldowns[user][cmd].IsZero() {
		return false
	}
	if handler.cooldowns[user][cmd].Sub(time.Now()) >= 0 {
		return true
	}
	return false
}

func (handler CommandHandler) startCooldown(user string, cmd command) {
	handler.cooldowns[user][cmd] = time.Now().Add(time.Duration(cmd.cooldown) * time.Second)
}

func (handler CommandHandler) startCooldownTicker() {
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
