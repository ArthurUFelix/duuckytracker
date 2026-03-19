package tracker

import (
	"log"
	"time"

	"github.com/arthurufelix/duuckylcu/internal/api"
	"github.com/arthurufelix/duuckylcu/internal/lcu"
)

type FriendTracker struct {
	lcuClient *lcu.Client
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

func NewFriendTracker(lcuClient *lcu.Client, apiClient *api.Client) *FriendTracker {
	return &FriendTracker{
		lcuClient: lcuClient,
		apiClient: apiClient,
		tracking:  make(map[string]*FriendGameState),
	}
}

func (ft *FriendTracker) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Println("Started tracking friend list for ARAM games...")

	for range ticker.C {
		ft.checkFriends()
	}
}

func (ft *FriendTracker) checkFriends() {
	friends, err := ft.lcuClient.GetFriendList()
	if err != nil {
		log.Printf("failed to get friend list: %v", err)
		return
	}

	for _, friend := range friends {
		ft.checkFriend(friend)
	}
}

func (ft *FriendTracker) checkFriend(friend lcu.Friend) {
	if friend.Lol == (lcu.GameStatus{}) {
		ft.handleFriendNotInGame(friend.PUUID)
		return
	}

	if friend.Lol.GameStatus == "inGame" && lcu.IsARAM(&friend.Lol) {
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
