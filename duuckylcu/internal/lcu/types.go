package lcu

import "encoding/json"

type LockfileData struct {
	ProcessName string
	PID         string
	Port        string
	Password    string
	Protocol    string
}

type LCUEvent struct {
	EventType string          `json:"event_type"`
	URI       string          `json:"uri"`
	Data      json.RawMessage `json:"data"`
}

type Friend struct {
	PUUID        string     `json:"puuid"`
	SummonerID   int        `json:"summonerId"`
	GameName     string     `json:"gameName"`
	GameTag      string     `json:"gameTag"`
	Name         string     `json:"name"`
	Availability string     `json:"availability"`
	Lol          GameStatus `json:"lol"`
}

type GameStatus struct {
	GameStatus    string `json:"gameStatus"`
	GameQueueType string `json:"gameQueueType"`
	GameMode      string `json:"gameMode"`
	MapID         string `json:"mapId"`
	ChampionId    string `json:"championId"`
}
