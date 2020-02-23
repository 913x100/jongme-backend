package api

import (
	"encoding/json"
	"fmt"
	"jongme/app/config"
	"jongme/app/fbbot"
	"jongme/app/model"
	"jongme/app/network"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FbDatabase interface {
	GetPageByID(id string) (*model.Page, error)
	GetServicesAccordingFilter(query []bson.M) ([]*model.Service, error)
}

type FbBot interface {
	// Process(messages []interface{})
	SendTextMessage(recipientID, pageAccessToken string, m *fbbot.TextMessage) network.Response
	SendQuickRepliesMessage(recipientID, pageAccessToken string, m *fbbot.QuickRepliesMessage) network.Response
	SendGenericMessage(recipientID, pageAccessToken string, m *fbbot.GenericMessage) network.Response
	SendButtonMessage(recipientID, pageAccessToken string, m *fbbot.ButtonMessage) network.Response
}

type FbAPI struct {
	DB FbDatabase
	FB FbBot
}

type Payload struct {
	StepID    int    `json:"step_id"`
	PageID    string `json:"page_id"`
	UserID    string `json:"user_id"`
	ServiceID string `json:"service_id"`
}

func (f *FbAPI) Webhook(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")
	mode := string(ctx.FormValue("hub.mode"))
	verifyToken := string(ctx.FormValue("hub.verify_token"))
	if mode == "subscribe" && verifyToken == config.ValidationToken {
		fmt.Println("Validating webhook")
		ctx.Write(ctx.FormValue("hub.challenge"))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusForbidden)
}

func (f *FbAPI) RecieveWebhook(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var rawMsg fbbot.RawCallbackMessage
	if err := json.Unmarshal(ctx.PostBody(), &rawMsg); err != nil {
		fmt.Println(err)
		return
	}
	messages := rawMsg.Unbox()

	var payload Payload
	for _, m := range messages {
		switch m := m.(type) {
		case *fbbot.Message:
			fmt.Println("Message")
			fmt.Printf("%+v\n", m)
			_ = json.Unmarshal(([]byte)(m.Quickreply.Payload), &payload)
			fmt.Println(m.Quickreply.Payload)
			f.Process(m.Sender, payload)
			break
		case *fbbot.Postback:
			fmt.Println("Postback")
			fmt.Printf("%+v\n", m)
			_ = json.Unmarshal(([]byte)(m.Payload), &payload)
			f.Process(m.Sender, payload)
			break
		}
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	return
}

func (f *FbAPI) Process(r fbbot.User, payload Payload) {
	fmt.Println(payload)
	page, _ := f.DB.GetPageByID(payload.PageID)

	switch payload.StepID {
	case 1:
		message := f.step1(payload.PageID, r.ID)
		f.Send(r, page.AccessToken, message)
		break
	case 2:
		fmt.Println(2)
		message := f.step2(payload.PageID)
		f.Send(r, page.AccessToken, message)
		break
	}

}

func (f *FbAPI) Send(r fbbot.User, pageAccessToken string, message interface{}) {
	switch m := message.(type) {
	case *fbbot.TextMessage:
		f.FB.SendTextMessage(r.ID, pageAccessToken, m)
		break
	case *fbbot.QuickRepliesMessage:
		fmt.Println("Quick")
		f.FB.SendQuickRepliesMessage(r.ID, pageAccessToken, m)
		break
	case *fbbot.GenericMessage:
		fmt.Println("Web")
		f.FB.SendGenericMessage(r.ID, pageAccessToken, m)
		break
	case *fbbot.ButtonMessage:
		fmt.Println("Web")
		f.FB.SendButtonMessage(r.ID, pageAccessToken, m)
		break
	}

}

func (f *FbAPI) step1(pageID, userID string) interface{} {
	filter := []bson.M{}

	filter = append(filter, bson.M{"page_id": bson.M{"$eq": pageID}})

	services, _ := f.DB.GetServicesAccordingFilter(filter)

	message := fbbot.NewQuickRepliesMessage("Which services do you want?")
	for _, service := range services {
		message.AddQuickRepliesItems(
			fbbot.NewQuickRepliesText(service.Name,
				fmt.Sprintf(`{"step_id":%d, "page_id":"%s", "user_id":"%s", "service_id":"%s"}`, 2, pageID, userID, primitive.ObjectID.Hex(service.ID))),
		)
	}
	return message
}

func (f *FbAPI) step2(pageID string) interface{} {
	message := fbbot.NewButtonMessage("Please select date and time")
	message.AddWebURLButton("Jongme", "https://bit.ly/3bIWNyP")
	return message
}
