package session

type SessionData struct {
	Task string
	Slots map[string]interface{}
	StepPath []int
}

type Session interface {
	GetSlot(slot string) (interface{})
	SetSlot(slot string, value interface{})
	SetTask(task string)
	GetTask() (string)
	SetStepPath(path []int)
	GetStepPath() ([]int)
}

type SessionProvider interface {
	GetSession(sessionId string) (Session)
}

func makeSession() (SessionData) {
	return SessionData{
		Slots: make(map[string]interface{}),
	}
}