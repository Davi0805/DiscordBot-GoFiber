package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"
    "io/ioutil"

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

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    if err != nil {
        return "", fmt.Errorf("erro ao criar a requisição: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("erro ao enviar a requisição: %w", err)
    }
    defer resp.Body.Close()

    // Verifica se a resposta foi 200 OK
    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return "", fmt.Errorf("erro na resposta da API. Status: %s, Corpo: %s", resp.Status, string(body))
    }

    // Processa o corpo da resposta
    var res OpenAIResponse
    if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
        return "", fmt.Errorf("erro ao decodificar a resposta JSON: %w", err)
    }

    if len(res.Choices) > 0 {
        return res.Choices[0].Message.Content, nil
    }

    return "Erro: nenhuma resposta encontrada", nil
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
                fmt.Println("Erro ao obter resposta da IA:", err) // Log do erro
                resposta = "Erro ao se conectar à IA. Tente novamente mais tarde."
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
