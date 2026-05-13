//go:build ignore

type Logger struct {
    name   string
    stream io.Writer
}

// Option is a function that modifies the Logger
type Option func(*Logger)

func WithStream(s io.Writer) Option {
    return func(l *Logger) {
        l.stream = s
    }
}

// NewLogger is the constructor
func NewLogger(name string, opts ...Option) *Logger {
    l := &Logger{
        name:   name,
        stream: os.Stderr, // Default value
    }
    for _, opt := range opts {
        opt(l)
    }
    return l
}

// Usage:
log := NewLogger("tatsu", WithStream(mySyncStream))
