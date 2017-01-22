package switchers

import "fmt"

// PlayerRepository persists player information
type PlayerRepository interface {
	GetOrCreatePlayer(string) (*Player, bool, error)
	SavePlayer(*Player) error
}

const (
	stateNew     = ""
	stateAskName = "askname"
	stateIdle    = "idle"

	commandResign = "/resign"

	commandSetName = "/setname"
	commandPause   = "/pause"
	commandResume  = "/resume"
)

// Player is a player
type Player struct {
	ID      string
	IsTrump bool
	State   string
	Name    string
	Paused  bool
	Score   int
}

// ExecuteCommand takes text command from a player, updates internal state and returns response
func (p *Player) ExecuteCommand(command string) string {
	if p.IsTrump {
		return p.executeTrumpCommand(command)
	}
	return p.executePlayerCommand(command)
}

func (p *Player) executeTrumpCommand(command string) string {
	switch command {
	case commandResign:
		p.IsTrump = false
		return "Отставка принята."
	}

	return fmt.Sprintf("Добро пожаловать, господин президент. Издайте какой-нибудь указ:\n\n%s — подать в отставку", commandResign)
}

func (p *Player) executePlayerCommand(command string) string {
	if command == commandPause {
		p.Paused = true
	}
	if p.Paused {
		if command != commandResume {
			return fmt.Sprintf("Участие в игре приостановлено. Ничего не сможешь делать, пока не напишешь %s.", commandResume)
		}
		p.Paused = false
		return "Участие в игре возобновлено. Продолжай как ни в чем не бывало."
	}

	switch p.State {
	case stateNew:
		p.State = stateAskName
		return "Привет! Чтобы стать участником Свитчеров, напиши в ответ свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду."

	case stateAskName:
		p.State = stateIdle
		p.Name = command
		return fmt.Sprintf("Приятно познакомиться, %s. Теперь жди инструкции. Они могут приходить в любой момент, так что держи телефон включенным! Чтобы приостановить участие в игре, напиши /pause.", p.Name)

	case stateIdle:
		if command == "/setname" {
			p.State = stateAskName
			return "Напиши свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду."
		}
	}

	return fmt.Sprintf("Жди инструкции или напиши какую-нибудь команду. Я понимаю:\n\n%s — изменить имя\n%s — приостановить участие в игре", commandSetName, commandPause)
}
