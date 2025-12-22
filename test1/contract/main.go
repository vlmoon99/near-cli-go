package main

import (
	"test1/a"
)

// @contract:state
type Contract struct {
	Message string
}

// @contract:init
func (c *Contract) Init(msg string) {
	c.Message = msg
}

// @contract:public
func (c *Contract) GetMessage() string {
	msg := a.Hello()
	return c.Message + " " + msg
}

// @contract:public
// @contract:mutating
func (c *Contract) SetMessage(newMessage string) {
	c.Message = newMessage
}
