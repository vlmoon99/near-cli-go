package main

import (
	"test1/a"
)

// @contract:state
type Contract struct {
	Message string
}

// Init initializes the contract with a specific message.
// It can only be called once, when the contract has no state.
//
// @contract:init
func (c *Contract) Init(msg string) {
	c.Message = msg
}

// GetMessage returns the current message combined with a greeting.
//
// @contract:public
func (c *Contract) GetMessage() string {
	msg := a.Hello()
	return c.Message + " " + msg
}

// SetMessage updates the message.
// If the contract was not initialized via Init(), this will initialize it
// implicitly with the new message and zero-values for other fields.
//
// @contract:public
// @contract:mutating
func (c *Contract) SetMessage(newMessage string) {
	c.Message = newMessage
}
