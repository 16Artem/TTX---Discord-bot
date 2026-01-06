// bot/bot.go
package bot

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session       *discordgo.Session
	Token         string
	stopChan      chan struct{}
	isRunning     bool
	isRunningLock sync.RWMutex
	reconnectLock sync.Mutex
}

func NewBot(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Session:   session,
		Token:     token,
		stopChan:  make(chan struct{}),
		isRunning: false,
	}, nil
}

func (b *Bot) Start() error {
	b.isRunningLock.Lock()
	b.isRunning = true
	b.isRunningLock.Unlock()

	defer func() {
		b.isRunningLock.Lock()
		b.isRunning = false
		b.isRunningLock.Unlock()
	}()

	b.Session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	b.Session.AddHandler(b.ready)
	b.Session.AddHandler(b.interactionCreate)
	b.Session.AddHandler(b.disconnectHandler)

	for {
		select {
		case <-b.stopChan:
			log.Println("Получен сигнал остановки бота")
			return b.Session.Close()
		default:
			err := b.connectAndListen()
			if err != nil {
				log.Printf("Ошибка соединения: %v", err)
			}

			if !b.isRunning {
				return nil
			}

			log.Println("Попытка переподключения через 5 секунд...")
			time.Sleep(5 * time.Second)
		}
	}
}

func (b *Bot) connectAndListen() error {
	b.reconnectLock.Lock()
	defer b.reconnectLock.Unlock()

	err := b.Session.Open()
	if err != nil {
		return err
	}

	log.Println("Discord соединение установлено")

	b.registerCommands()

	select {
	case <-b.stopChan:
		return b.Session.Close()
	case <-time.After(1 * time.Hour):
	}

	return nil
}

func (b *Bot) Stop() error {
	b.isRunningLock.Lock()
	b.isRunning = false
	b.isRunningLock.Unlock()

	close(b.stopChan)
	return b.Session.Close()
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "TTX-Bot | /help")
	log.Printf("Бот вошел как: %s", event.User.Username)
}

func (b *Bot) disconnectHandler(s *discordgo.Session, event *discordgo.Disconnect) {
	log.Printf("Discord соединение разорвано: %v", event)
}

func (b *Bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	if data.Name == "help" {
		b.handleHelpCommand(s, i)
	}
}

func (b *Bot) GetStats() (servers int, commands int, users int) {
	if b.Session == nil || b.Session.State == nil {
		return 0, 0, 0
	}

	guilds := b.Session.State.Guilds
	servers = len(guilds)

	users = 0
	for _, guild := range guilds {
		users += guild.MemberCount
	}

	commands = 1

	return servers, commands, users
}
