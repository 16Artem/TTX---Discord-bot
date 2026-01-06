// bot/commands.go
package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) registerCommands() {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Показать список команд",
		},
	}

	for _, cmd := range commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", cmd)
		if err != nil {
			log.Printf("Ошибка создания команды %s: %v", cmd.Name, err)
		}
	}
}

func (b *Bot) handleHelpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "TTX-Bot - Помощь",
		Description: "Доступные команды:",
		Color:       0x7289DA,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Основные команды",
				Value:  "`/help` - Показать это сообщение",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "TTX-Bot",
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
