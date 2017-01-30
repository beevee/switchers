package switchers

// Player states
const (
	PlayerStateNew        = ""
	PlayerStateAskName    = "askname"
	PlayerStateIdle       = "idle"
	PlayerStateInGame     = "ingame"
	PlayerStateModerating = "moderating"
)

// Player is a player
type Player struct {
	ID                  string
	Trump               bool
	State               string
	Name                string
	Paused              bool
	Score               int
	ModeratingTeamIndex int
}

// IsEligible detects if player can participate in a new round
func (p Player) IsEligible() bool {
	return !p.Trump && !p.Paused && p.State != PlayerStateAskName && p.State != PlayerStateNew && p.Name != ""
}
