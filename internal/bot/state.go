package bot

type State string

const (
	StateNone         State = "none"
	StateAwaitingName State = "awaitingSourceName"
	StateAwaitingURL  State = "awaitingSourceURL"

	StateAwaitingChannelName   State = "awaitingChannelName"
	StateAwaitingUnlinkChannel State = "awaitingUnlinkChannel"	
	StateAwaitingPostTime      State = "awaitingPostTime"
	StateAwaitingPostCount     State = "awaitingPostCount"
)

type UserState struct {
	Current   State
	TempName  string
	TempValue string
}

var userStates = make(map[int64]*UserState)

func GetUserState(userID int64) *UserState {
	if _, exists := userStates[userID]; !exists {
		userStates[userID] = &UserState{Current: StateNone}
	}
	return userStates[userID]
}

func ResetUserState(userID int64) {
	userStates[userID] = &UserState{Current: StateNone}
}
