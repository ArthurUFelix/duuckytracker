package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type RabbitMQMessage struct {
	Pattern string      `json:"pattern"`
	Data    MessageData `json:"data"`
}

type MessageData struct {
	ID   string       `json:"id"`
	Data SummonerData `json:"data"`
}

type SummonerData struct {
	Summoner     string `json:"summoner"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	LPDelta      int    `json:"lpDelta"`
	LeaguePoints int    `json:"leaguePoints"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
	ChampionName string `json:"championName"`
}

type discordWebhook struct {
	Embeds []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title     string          `json:"title"`
	Color     int             `json:"color"`
	Fields    []embedField    `json:"fields"`
	Timestamp string          `json:"timestamp"`
	Thumbnail *embedThumbnail `json:"thumbnail"`
}

type embedThumbnail struct {
	Url string `json:"url"`
}

type embedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func (d *DiscordNotifier) SendNotification(msg RabbitMQMessage) error {
	summoner := msg.Data.Data
	color := 15158332 // Red for loss
	result := "💀 Derrota"

	if summoner.LPDelta > 0 {
		color = 3066993 // Green for win
		result = "🏆 Vitória"
	}

	totalGames := summoner.Wins + summoner.Losses
	winRate := 0.0
	if totalGames > 0 {
		winRate = (float64(summoner.Wins) / float64(totalGames)) * 100
	}

	lpDeltaText := fmt.Sprintf("📈 %+d PDL", summoner.LPDelta)
	if summoner.LPDelta < 0 {
		lpDeltaText = fmt.Sprintf("📉 %d PDL", summoner.LPDelta)
	}

	embed := discordEmbed{
		Title: fmt.Sprintf("Partida finalizada - %s", summoner.Summoner),
		Color: color,
		Fields: []embedField{
			{
				Name:   "Pontos",
				Value:  lpDeltaText,
				Inline: true,
			},
			{
				Name:   "Resultado",
				Value:  result,
				Inline: true,
			},
			{
				Name:   "KDA",
				Value:  fmt.Sprintf("%d/%d/%d", summoner.Kills, summoner.Deaths, summoner.Assists),
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  fmt.Sprintf("PDL: %d - Vitórias: %d - Derrotas: %d - WR: %.1f%%", summoner.LeaguePoints, summoner.Wins, summoner.Losses, winRate),
				Inline: true,
			},
		},
		Thumbnail: &embedThumbnail{
			Url: fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/img/champion/loading/%s_0.jpg", summoner.ChampionName),
		},
	}

	webhook := discordWebhook{
		Embeds: []discordEmbed{embed},
	}

	payload, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	resp, err := d.client.Post(
		d.webhookURL,
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("failed to send webhook rqeuest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("discord API returned status %d", resp.StatusCode)
	}

	return nil
}
