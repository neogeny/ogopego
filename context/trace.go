package context

type Tracer interface {
}

type NullTracer struct{}

type ConsoleTracer struct{}
