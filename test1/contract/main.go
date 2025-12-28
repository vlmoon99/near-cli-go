package main

import "github.com/vlmoon99/near-sdk-go/env"

// @contract:state
type Contract struct {
	Message string
}

// @contract:init
func (c *Contract) Init(msg string) {
	env.LogString("Init Method")
	c.Message = msg
}

// @contract:public
func (c *Contract) GetMessage() string {
	return c.Message
}

// @contract:public
// @contract:mutating
func (c *Contract) SetMessage(newMessage string) {
	c.Message = newMessage
}
