package pixie

import (
	"log"
	"os/user"

	"github.com/spf13/viper"
)

func init() {
	exeUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	userDir := exeUser.HomeDir

	viper.SetConfigName("pixie-config")
	viper.AddConfigPath(userDir)
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to load config file: %s", err.Error())
	}
}

// BotName returns the command to use as the bot's name inside your server;
// the default name is 'pixie'
func BotName() string {
	name := viper.GetString("bot_name")
	if name == "" {
		name = "pixie"
	}
	return name
}

// DiscordToken returns the Discord authentication token that should be loaded
// from the config file; this allows the bot to communicate with the server
func DiscordToken() string {
	return "Bot " + viper.GetString("discord_token")
}
