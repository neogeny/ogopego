package heartbeat

type Heart interface {
	Beat(mark, total int)
}

type NullHeart struct{}

func (NullHeart) Beat(_ int, _ int) {}
