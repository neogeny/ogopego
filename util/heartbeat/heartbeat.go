package heartbeat

type Heartbeat interface {
	Tick(mark, total int)
}

type NullHeartbeat struct{}

func (NullHeartbeat) Tick(_ int, _ int) {}
