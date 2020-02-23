package fbbot

type Message struct {
	ID         string
	Page       Page
	Sender     User
	Text       string
	Seq        int
	Timestamp  int64
	Quickreply Quickreply
}

type Quickreply struct {
	Payload string
}

// Postback
type Postback struct {
	Sender  User
	Title   string `json:"title"`
	Payload string `json:"payload"`
}
