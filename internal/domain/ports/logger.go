package ports

// Logger define el puerto para logging
// Este es un puerto de SALIDA (driving port)
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

