package entities

// Connection representa el estado de la conexión de WhatsApp
type Connection struct {
	IsConnected bool
	IsLoggedIn  bool
	QRCode      string
	Session     *Session
}

// NewConnection crea una nueva conexión
func NewConnection() *Connection {
	return &Connection{
		IsConnected: false,
		IsLoggedIn:  false,
	}
}

// SetQRCode establece el código QR para autenticación
func (c *Connection) SetQRCode(qr string) {
	c.QRCode = qr
}

// MarkAsConnected marca la conexión como establecida
func (c *Connection) MarkAsConnected(session *Session) {
	c.IsConnected = true
	c.IsLoggedIn = true
	c.Session = session
	if session != nil {
		session.Connect()
	}
}

// MarkAsDisconnected marca la conexión como desconectada
func (c *Connection) MarkAsDisconnected() {
	c.IsConnected = false
	if c.Session != nil {
		c.Session.Disconnect()
	}
}
