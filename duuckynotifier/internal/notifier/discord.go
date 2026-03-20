package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	ID   string          `json:"id"`
	Data json.RawMessage `json:"data"`
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

type MatchRankingData struct {
	MatchesRanking []MatchData `json:"matchesRanking"`
}

type MatchData struct {
	Summoner     string `json:"summoner"`
	MatchesCount int    `json:"matchesCount"`
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

func (d *DiscordNotifier) SendNotification(msg RabbitMQMessage, queueName string) error {
	var webhook discordWebhook
	var err error

	switch queueName {
	case "summoners":
		webhook, err = buildSummonerWebhook(msg)
	case "matches":
		webhook, err = buildMatchWebhook(msg)
	default:
		return fmt.Errorf("invalid queuename: %s", queueName)
	}

	if err != nil {
		return fmt.Errorf("failed to build webhook: %w", err)
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

func buildSummonerWebhook(msg RabbitMQMessage) (discordWebhook, error) {
	var summoner SummonerData
	if err := json.Unmarshal(msg.Data.Data, &summoner); err != nil {
		return discordWebhook{}, fmt.Errorf("failed to parse summoner data: %w", err)
	}

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

	return discordWebhook{
		Embeds: []discordEmbed{embed},
	}, nil
}

func buildMatchWebhook(msg RabbitMQMessage) (discordWebhook, error) {
	var rankings MatchRankingData
	if err := json.Unmarshal(msg.Data.Data, &rankings); err != nil {
		return discordWebhook{}, fmt.Errorf("failed to parse match data: %w", err)
	}

	color := 3447003 // Blue

	var rankingText strings.Builder
	medals := []string{"🥇", "🥈", "🥉"}
	totalGames := 0

	for i, entry := range rankings.MatchesRanking {
		// Use medals for top 3, numbers for the rest
		position := fmt.Sprintf("%d.", i+1)
		if i < 3 {
			position = medals[i]
		}

		rankingText.WriteString(fmt.Sprintf("%s **%s** - %d partida", position, entry.Summoner, entry.MatchesCount))

		if entry.MatchesCount != 1 {
			rankingText.WriteString("s")
		}

		rankingText.WriteString("\n")
		totalGames += entry.MatchesCount
	}

	embed := discordEmbed{
		Title: "🏆 Ranking ARAMs diários",
		Color: color,
		Fields: []embedField{
			{
				Name:   "",
				Value:  rankingText.String(),
				Inline: false,
			},
		},
	}

	return discordWebhook{
		Embeds: []discordEmbed{embed},
	}, nil
}
