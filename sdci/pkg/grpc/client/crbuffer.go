package client

import (
	"bytes"
)

type crbuffer struct {
	buf []byte
}

func (c *crbuffer) String() string {
	return string(c.buf)
}

func (c *crbuffer) Write(p []byte) (int, error) {
	for _, b := range p {
		switch b {
		case '\r':
			c.buf = c.buf[:bytes.LastIndexByte(c.buf, '\n')+1]
		default:
			c.buf = append(c.buf, b)
		}
	}
	return len(p), nil
}
