package asana

import (
	"log"
)

func (c *Client) info(format string, args ...interface{}) {
	if len(c.Verbose) > 0 {
		log.Printf(format, args...)
	}
}

func (c *Client) trace(format string, args ...interface{}) {
	if len(c.Verbose) > 1 {
		log.Printf(format, args...)
	}
}

func (c *Client) debug(format string, args ...interface{}) {
	if IsTrue(c.DefaultOptions.Debug) {
		log.Printf(format, args...)
	}
}
