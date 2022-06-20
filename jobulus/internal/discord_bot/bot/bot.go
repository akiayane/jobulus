package bot

import (
	"fmt"
	"jobulus/internal/discord_bot/config"
	"sync"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var BotId string
var goBot *discordgo.Session

type Bot struct {
	session *discordgo.Session
	wg      sync.WaitGroup
}

func Start() (Bot, error) {

	//creating new bot session
	goBot, err := discordgo.New("Bot " + config.Token)

	var wg sync.WaitGroup

	//Handling error
	if err != nil {
		fmt.Println(err.Error())
		return Bot{}, err
	}
	// Making our bot a user using User function .
	u, err := goBot.User("@me")
	//Handlinf error
	if err != nil {
		fmt.Println(err.Error())
		return Bot{}, err
	}
	// Storing our id from u to BotId .
	BotId = u.ID

	// Adding handler function to handle our messages using AddHandler from discordgo package. We will declare messageHandler function later.
	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	//Error handling
	if err != nil {
		fmt.Println(err.Error())
		return Bot{}, err
	}
	return Bot{goBot, wg}, nil
}

func (b *Bot) SendMessage(message string) {
	for _, guild := range b.session.State.Guilds {

		b.wg.Add(len(b.session.State.Guilds))

		go func() {
			// Get channels for this guild
			channels, _ := b.session.GuildChannels(guild.ID)

			for _, c := range channels {
				// Check if channel is a guild text channel and not a voice or DM channel
				if c.Type != discordgo.ChannelTypeGuildText {
					continue
				}

				// Send text message
				b.session.ChannelMessageSend(c.ID, message)
			}

			b.wg.Done()
		}()

		b.wg.Wait()

	}
}

//Definition of messageHandler function it takes two arguments first one is discordgo.Session which is s , second one is discordgo.MessageCreate which is m.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check.
	if m.Author.ID == BotId {
		return
	}
	//If we message ping to our bot in our discord it will return us pong .
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}
