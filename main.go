package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/bwmarrin/discordgo"
)

type OpenAIResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}

func getOpenAIResponse(prompt string) (string, error) {
    apiKey := os.Getenv("OPENAI_API_KEY")
    url := "https://api.openai.com/v1/chat/completions"

    requestBody, _ := json.Marshal(map[string]interface{}{
        "model": "gpt-3.5-turbo",
        "messages": []map[string]string{
            {"role": "system", "content": "Você é um assistente útil."},
            {"role": "user", "content": prompt},
        },
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var res OpenAIResponse
    json.NewDecoder(resp.Body).Decode(&res)

    if len(res.Choices) > 0 {
        return res.Choices[0].Message.Content, nil
    }
    return "Erro ao obter resposta", nil
}

func main() {
    bot, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
    if err != nil {
        fmt.Println("Erro ao iniciar o bot:", err)
        return
    }

    bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
        if m.Author.Bot {
            return
        }

        if strings.HasPrefix(m.Content, "!chat ") {
            pergunta := strings.TrimPrefix(m.Content, "!chat ")
            if pergunta == "" {
                s.ChannelMessageSend(m.ChannelID, "Por favor, insira uma pergunta depois de `!chat`.")
                return
            }

            resposta, err := getOpenAIResponse(pergunta)
            if err != nil {
                resposta = "Erro ao se conectar à IA"
            }
            s.ChannelMessageSend(m.ChannelID, resposta)
        }
    })

    err = bot.Open()
    if err != nil {
        fmt.Println("Erro ao conectar o bot:", err)
        return
    }

    fmt.Println("Bot está online. Digite !chat [pergunta] para falar com a IA.")
    select {} // Mantém o bot rodando
}
