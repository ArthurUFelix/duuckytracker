package tracker

import (
	"encoding/json"
	"log"
	"slices"

	"github.com/arthurufelix/duuckylcu/internal/api"
	"github.com/arthurufelix/duuckylcu/internal/lcu"
)

type FriendTracker struct {
	wsClient  *lcu.WebSocketClient
	apiClient *api.Client
	tracking  map[string]*FriendGameState
}

type FriendGameState struct {
	PUUID      string
	Name       string
	InGame     bool
	GameMode   string
	ChampionId string
}

func NewFriendTracker(wsClient *lcu.WebSocketClient, apiClient *api.Client) *FriendTracker {
	return &FriendTracker{
		wsClient:  wsClient,
		apiClient: apiClient,
		tracking:  make(map[string]*FriendGameState),
	}
}

func (ft *FriendTracker) Start() {
	log.Println("Started tracking friend list for ARAM games...")

	err := ft.wsClient.Listen(func(event lcu.LCUEvent) {
		ft.checkFriend(event)
	})
	if err != nil {
		log.Fatalf("websocket error: %v", err)
	}
}

func (ft *FriendTracker) checkFriend(event lcu.LCUEvent) {
	var friend lcu.Friend
	if err := json.Unmarshal(event.Data, &friend); err != nil {
		log.Printf("failed to unmarshal event: %v", err)
		return
	}

	tracking := []string{"Arthur#ツツツ", "Il L7 Il#L77", "Naba#naba1", "dream#heart"}
	summonerName := friend.GameName + "#" + friend.GameTag
	if !slices.Contains(tracking, summonerName) {
		return
	}

	if friend.Lol == (lcu.GameStatus{}) {
		ft.handleFriendNotInGame(friend.PUUID)
		return
	}

	if friend.Lol.GameStatus == "inGame" && IsARAM(&friend.Lol) {
		ft.handleFriendInGame(friend, &friend.Lol)
	} else {
		ft.handleFriendNotInGame(friend.PUUID)
	}
}

func (ft *FriendTracker) handleFriendInGame(friend lcu.Friend, gameStatus *lcu.GameStatus) {
	state, exists := ft.tracking[friend.PUUID]

	if !exists || !state.InGame {
		friendName := friend.Name
		if friendName == "" {
			friendName = friend.GameName + "#" + friend.GameTag
		}

		log.Printf("%s started an ARAM game!", friendName)

		ft.tracking[friend.PUUID] = &FriendGameState{
			PUUID:      friend.PUUID,
			Name:       friendName,
			InGame:     true,
			ChampionId: gameStatus.ChampionId,
			GameMode:   gameStatus.GameQueueType,
		}
	}
}

func (ft *FriendTracker) handleFriendNotInGame(puuid string) {
	state, exists := ft.tracking[puuid]

	if exists && state.InGame {
		log.Printf("%s finished their ARAM game", state.Name)

		if ft.apiClient.IsAuthenticated() {
			err := ft.apiClient.ReportGameEnded(state.Name, state.ChampionId)
			if err != nil {
				log.Printf("failed to report game end to API: %v", err)
			} else {
				log.Printf("sent game end notification to API")
			}
		}

		state.InGame = false
	}
}

func IsARAM(gameStatus *lcu.GameStatus) bool {
	return gameStatus.GameMode == "ARAM" || gameStatus.MapID == "12"
}
