package gchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Message represents a google-chat message.
type Message struct {
	Cards []CardV2 `json:"cardsV2"`
}

// CardV2 represents a google-chat card message field.
// https://developers.google.com/chat/api/reference/rest/v1/cards
type CardV2 struct {
	CardId string `json:"cardId"`
	Card Card `json:"card"`
}

type Card struct {
	// Header CardHeader `json:"header"`
	Sections []CardSection `json:"sections"`
}

type CardHeader struct {
	Title string `json:"title"`
	SubTitle string `json:"subtitle"`
	ImageUrl string `json:"imageUrl"`
	ImageType string `json:"imageType"`
	ImageAltText string `json:"imageAltText"`
}

type CardSection struct {
	Header string `json:"header,omitempty"`
	Collapsible bool `json:"collapsible"`
	Widgets []CardWidget `json:"widgets"`
}

type CardWidget struct {
	DecoratedText WidgetDecoratedText `json:"decoratedText"`
}

type WidgetDecoratedText struct {
	StartIcon *DecoratedTextIcon `json:"startIcon"`
	Text string `json:"text,omitempty"`
	WrapText bool `json:"wrapText"`
}

type DecoratedTextIcon struct {
	KnownIcon string `json:"knownIcon,omitempty"`
	IconUrl string `json:"iconUrl,omitempty"`
}

// Send sends the message to the webhook URL.
func Send(url string, m *Message) error {
	buf, err := json.Marshal(m)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}
