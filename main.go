package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Mapa para armazenar canais de notificação para cada servidor (guild)
var channelMap = make(map[string]string)
var mu sync.Mutex // Mutex para evitar concorrência

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Println("Por favor, defina a variável de ambiente DISCORD_BOT_TOKEN")
		return
	}

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Erro ao iniciar o bot:", err)
		return
	}

	// Ativar intents necessárias
	bot.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuildPresences

	// Comando para definir o canal
	bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}

		// Comando !ping
		if m.Content == "!ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}

		// Comando !setchannel para definir canal dinamicamente
		if m.Content == "!setchannel" {
			mu.Lock()
			channelMap[m.GuildID] = m.ChannelID
			mu.Unlock()
			s.ChannelMessageSend(m.ChannelID, "✅ Canal de notificações definido!")
		}
	})

	// Detectar quando um usuário inicia um jogo
	bot.AddHandler(func(s *discordgo.Session, p *discordgo.PresenceUpdate) {
		if len(p.Activities) > 0 {
			for _, activity := range p.Activities {
				if activity.Type == discordgo.ActivityTypeGame {
					mu.Lock()
					channelID, exists := channelMap[p.GuildID]
					mu.Unlock()
					if exists {
						msg := fmt.Sprintf("🚀 %s começou a jogar **%s**! Será que ele é homoafetivo?", p.User.Username, activity.Name)
						s.ChannelMessageSend(channelID, msg)
					}
				}
			}
		}
	})

	// Conectar o bot
	err = bot.Open()
	if err != nil {
		fmt.Println("Erro ao conectar o bot:", err)
		return
	}

	fmt.Println("Bot está online. Pressione Ctrl+C para sair.")

	// Aguardar sinal de interrupção
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	bot.Close()
}
