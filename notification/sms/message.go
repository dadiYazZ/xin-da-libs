package sms

import (
	"github.com/dadiYazZ/xin-da-libs/notification/contract"
)

type Message struct {
	contract.MessageInterface

	To          []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}
