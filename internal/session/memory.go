package session

type MemorySessionProvider struct {
	memorySession map[string]*MemorySession
}

type MemorySession struct {
	sessionId string
	data SessionData
}

func NewMemorySessionProvider() *MemorySessionProvider {
	return &MemorySessionProvider{
		memorySession: map[string]*MemorySession {},
	}
}

func (provider *MemorySessionProvider) GetSession(sessionId string) (*MemorySession) {
	if _, ok := provider.memorySession[sessionId]; !ok {
		provider.memorySession[sessionId] = &MemorySession{
			sessionId: sessionId,
			data: makeSession(),
		}
	}

	return provider.memorySession[sessionId]
}

func (session *MemorySession) GetSlot(slot string) (interface{}) {
	return session.data.Slots[slot]
}

func (session *MemorySession) SetSlot(slot string, value interface{}) {
	session.data.Slots[slot] = value
}

func (session *MemorySession) GetSlots() (map[string]interface{}) {
	return session.data.Slots
}

func (session *MemorySession) SetTask(task string) {
	session.data.Task = task
	session.data.StepPath = []int{0}
}

func (session *MemorySession) GetTask() (string) {
	return session.data.Task
}

func (session *MemorySession) SetStepPath(path []int) {
	session.data.StepPath = path
}

func (session *MemorySession) GetStepPath() ([]int) {
	return session.data.StepPath
}
