package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/bwmarrin/discordgo"
)

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

    bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
        if m.Author.Bot {
            return
        }

        if m.Content == "!ping" {
            s.ChannelMessageSend(m.ChannelID, "Pong!")
        }
    })

    err = bot.Open()
    if err != nil {
        fmt.Println("Erro ao conectar o bot:", err)
        return
    }

    fmt.Println("Bot está online. Pressione Ctrl+C para sair.")
    
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop

    bot.Close()
}
