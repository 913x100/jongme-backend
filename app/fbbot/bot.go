package fbbot

import (
	"fmt"
	"jongme/app/network"
)

type Bot struct {
	Client network.Client
}

func New(client network.Client) Bot {
	return Bot{
		Client: client,
	}
}

func (b *Bot) Process(messages []interface{}) {
	for _, m := range messages {
		switch m := m.(type) {
		case *Message:
			fmt.Println("Message", m)
			// msg := NewTextMessage("Test")
			// b.Send(m.Sender, pageAccessToken, msg)
		case *Postback:
			fmt.Println("Postback", m)
		}
	}
}

// func (b *Bot) Send(r User, pageAccessToken string, m interface{}) error {
// 	switch m := m.(type) {
// 	case *TextMessage:
// 		// _ = b.SendTextMessage(r, pageAccessToken, m)
// 	// case *QuickRepliesMessage:
// 	// 	return b.sendQuickRepliesMessage(r, m)
// 	default:
// 		return errors.New("unknown message type")
// 	}
// 	return nil
// }
