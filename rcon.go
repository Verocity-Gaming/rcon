package rcon

import (
	"fmt"
	"net"
)

type Conn struct {
	net.Conn
	key []byte // XOR key.
}

const msglen = 8196

// New returns a new HLL RCON client to set/get server parameters.
func New(addr string, password string) (*Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		Conn: conn,
	}

	// Retrieve the XOR key used to encrypt communications between client/server.
	c.key = make([]byte, msglen)

	n, err := c.Read(c.key)
	if err != nil {
		return nil, err
	}

	c.key = c.key[:n]

	return c, c.login(password)
}

func (c *Conn) login(password string) error {
	err := c.send("login " + password)
	if err != nil {
		return err
	}

	auth, err := c.read()
	if err != nil {
		return err
	}

	if auth != "SUCCESS" {
		return fmt.Errorf("rcon authentication failed for %s and password %x", c.RemoteAddr(), c.key)
	}

	return nil
}

func (c *Conn) read() (string, error) {
	b := make([]byte, msglen)

	n, err := c.Read(b)
	if err != nil {
		return "", err
	}

	return string(c.xor(b)[:n]), nil
}

func (c *Conn) send(s string) error {
	_, err := c.Write(c.xor([]byte(s)))
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) xor(b []byte) []byte {
	d := make([]byte, len(b))

	for i := range b {
		d[i] = b[i] ^ c.key[i%len(c.key)]
	}

	return d
}
