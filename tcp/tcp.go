package tcp
import (
  "net"
  "log"
)

var message string = "Hello, World!"

func SetMessage(newmessage string) {
  log.Println("Setting new message:", newmessage)
  message = newmessage
}

type Server struct {
  address string
}

type Client struct {
  conn net.Conn
}

func (c *Client) Listen() {
  log.Println("Accepted new client connection:", c.conn.RemoteAddr())
  c.conn.Write([]byte(message))
  c.conn.Close()
}

func (s *Server) Listen() {
  listener, err := net.Listen("tcp", s.address)
  if err != nil {
    log.Fatal("Error starting TCP server.")
  }
  defer listener.Close()

  for {
    conn, _ := listener.Accept()
    client := &Client{conn: conn}
    go client.Listen()
  }
}

func New(address string) *Server {
  log.Println("Creating server with address", address)
  server := &Server{address: address}
  return server
}
