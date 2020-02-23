package fbbot

type TextMessage struct {
	Text string
}

func NewTextMessage(text string) *TextMessage {
	var t TextMessage
	t.Text = text
	return &t
}

type QuickRepliesMessage struct {
	Text  string              `json:"text"`
	Items []*QuickRepliesItem `json:"quick_replies"`
}

func NewQuickRepliesMessage(text string) *QuickRepliesMessage {
	return &QuickRepliesMessage{
		Text: text,
	}
}

func (m *QuickRepliesMessage) AddQuickRepliesItems(items ...*QuickRepliesItem) {
	m.Items = append(m.Items, items...)
}

type QuickRepliesItem struct {
	ContentType string `json:"content_type"`        // 'text' or 'location'
	Title       string `json:"title,omitempty"`     // empty when ContentType='location'
	Payload     string `json:"payload,omitempty"`   // empty when ContentType='location'
	ImageURL    string `json:"image_url,omitempty"` // optional, empty when ContentType='location'
}

func NewQuickRepliesText(title string, payload string) *QuickRepliesItem {
	return &QuickRepliesItem{
		ContentType: "text",
		Title:       title,
		Payload:     payload,
	}
}

type GenericMessage struct {
	Elements []*GenericMessageElement
}

func NewGenericMessage() *GenericMessage {
	return &GenericMessage{}
}

func (m *GenericMessage) AddGenericElement(elements ...*GenericMessageElement) {
	m.Elements = append(m.Elements, elements...)
}

type GenericMessageElement struct {
	Title         string               `json:"title"`
	ImageUrl      string               `json:"image_url"`
	SubTitle      string               `json:"subtitle"`
	DefaultAction GenericMessageAction `json:"default_action"`
	Buttons       []*Button            `json:"buttons"`
}

type GenericMessageAction struct {
	Type    string `json:"type"`
	Url     string `json:"url"`
	WebView string `json:"webview_height_ratio"`
}

func NewGenericMessageElement(title, imageUrl, subTitle string) *GenericMessageElement {
	return &GenericMessageElement{
		Title:    title,
		ImageUrl: imageUrl,
		SubTitle: subTitle,
	}
}

func (m *GenericMessageElement) AddGenericMessageButton(buttons ...*Button) {
	m.Buttons = append(m.Buttons, buttons...)
}

type ButtonMessage struct {
	Text    string
	Buttons []Button
}

func NewButtonMessage(text string) *ButtonMessage {
	return &ButtonMessage{
		Text: text,
	}
}

func (m *ButtonMessage) AddWebURLButton(title, URL string) {
	b := NewWebURLButton(title, URL)
	m.Buttons = append(m.Buttons, b)
}

type Button struct {
	Type    string `json:"type"` // web_url or postback
	Title   string `json:"title,omitempty"`
	URL     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

func NewWebURLButton(title, URL string) Button {
	return Button{
		Type:  "web_url",
		Title: title,
		URL:   URL,
	}
}

func NewPostbackButton(title, payload string) Button {
	return Button{
		Type:    "postback",
		Title:   title,
		Payload: payload,
	}
}
