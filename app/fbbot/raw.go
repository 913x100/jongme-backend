package fbbot

type RawCallbackMessage struct {
	RawObject  string `json:"object"`
	RawEntries []struct {
		RawID        string           `json:"id"`
		RawEventTime int64            `json:"time"`
		RawMessaging []RawMessageData `json:"messaging"`
	} `json:"entry"`
}

// rawMessageData contains data related to a message
type RawMessageData struct {
	RawSender    User  `json:"sender"`
	RawRecipient Page  `json:"recipient"`
	RawTimestamp int64 `json:"timestamp"`

	RawMessage *RawMessage `json:"message"`
	Postback   *Postback   `json:"postback"`
}

// rawMessage is a Facebook message
type RawMessage struct {
	RawMid        string        `json:"mid"`
	RawSeq        int           `json:"seq"`
	RawText       string        `json:"text"`
	RawQuickreply RawQuickreply `json:"quick_reply"`
}

type RawQuickreply struct {
	Payload string `json:"payload"`
}

type Page struct {
	ID string `json:"id"`
}

type User struct {
	ID string `json:"id"`
}

func (cbMsg *RawCallbackMessage) Unbox() []interface{} {
	var messages []interface{}
	for _, entry := range cbMsg.RawEntries {
		for _, rawMessageData := range entry.RawMessaging {
			if rawMessageData.RawMessage != nil {
				messages = append(messages, buildMessage(rawMessageData))
			} else if rawMessageData.Postback != nil {
				rawMessageData.Postback.Sender = rawMessageData.RawSender
				messages = append(messages, rawMessageData.Postback)
			}
		}
	}
	return messages
}

func buildMessage(m RawMessageData) *Message {
	var msg Message
	msg.ID = m.RawMessage.RawMid
	msg.Page = m.RawRecipient
	msg.Sender = m.RawSender
	msg.Text = m.RawMessage.RawText
	msg.Seq = m.RawMessage.RawSeq
	msg.Timestamp = m.RawTimestamp
	msg.Quickreply = Quickreply{Payload: m.RawMessage.RawQuickreply.Payload}

	return &msg
}
