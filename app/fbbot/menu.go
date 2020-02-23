package fbbot

type Menu struct {
	Locale                string      `json:"locale"`
	ComposerInputDisabled bool        `json:"composer_input_disabled"`
	CallToActions         []*MenuItem `json:"call_to_actions"`
}

func NewMenu() *Menu {
	return &Menu{
		Locale: "default",
	}
}

func (m *Menu) AddMenuItems(items ...*MenuItem) {
	m.CallToActions = append(m.CallToActions, items...)
}

type MenuItem struct {
	Title              string      `json:"title"`
	Type               string      `json:"type"`
	URL                string      `json:"url,omitempty"`
	WebviewHeightRatio string      `json:"webview_height_ratio,omitempty"`
	Payload            string      `json:"payload,omitempty"`
	CallToActions      []*MenuItem `json:"call_to_actions,omitempty"`
}

func NewPostbackMenuItem(title, payload string) *MenuItem {
	return &MenuItem{
		Title:   title,
		Type:    "postback",
		Payload: payload,
	}
}
