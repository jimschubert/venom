package venom

// Logger allows the user to provide any logger fulfilling this interface
type Logger interface {
	// Printf is a common signature used by log.Logger, logrus.Logger, and others
	Printf(format string, v ...any)
}
