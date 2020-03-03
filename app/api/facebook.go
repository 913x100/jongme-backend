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
	GetBookingsAccordingFilter(query []bson.M) ([]*model.Booking, error)
	UpdateBooking(booking *model.Booking) (*model.Booking, error)
	DeleteBookingByID(id primitive.ObjectID) error
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
	BookingID string `json:"booking_id"`
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

func (f *FbAPI) SendSuccesMessage(ctx *fasthttp.RequestCtx) {

	if !ctx.IsPost() {
		return
	}

	booking := model.Booking{}

	if err := json.Unmarshal(ctx.PostBody(), &booking); err != nil {
		return
	}

	booking_id := primitive.ObjectID.Hex(booking.ID)

	msg := fmt.Sprintf("การจองสำเร็จ\n\n%s-%s-%s\nat %s\nYes or no", booking.Year, booking.Month, booking.Day, booking.Time)
	message := fbbot.NewQuickRepliesMessage(msg)

	message.AddQuickRepliesItems(
		fbbot.NewQuickRepliesText("Yes",
			fmt.Sprintf(`{"step_id":%d, "page_id":"%s", "user_id":"%s", "service_id":"%s", "booking_id":"%s"}`, 3, booking.PageID, booking.UserID, booking.ServiceID, booking_id)),
		fbbot.NewQuickRepliesText("No",
			fmt.Sprintf(`{"step_id":%d, "page_id":"%s", "user_id":"%s", "service_id":"%s", "booking_id":"%s"}`, 4, booking.PageID, booking.UserID, booking.ServiceID, booking_id)),
	)
	page, _ := f.DB.GetPageByID(booking.PageID)

	f.FB.SendQuickRepliesMessage(booking.UserID, page.AccessToken, message)

	return
}

func (f *FbAPI) Process(r fbbot.User, payload Payload) {
	fmt.Println(payload)
	page, _ := f.DB.GetPageByID(payload.PageID)

	switch payload.StepID {
	case 1:
		message := f.step_service(payload.PageID, r.ID)
		f.Send(r, page.AccessToken, message)
		break
	case 2:
		fmt.Println(2)
		message := f.step_booking(payload.PageID, payload.ServiceID, r.ID)
		f.Send(r, page.AccessToken, message)
		break
	case 3:
		fmt.Println(3)
		message := f.step_confirm_booking(payload.BookingID)
		f.Send(r, page.AccessToken, message)
		break
	case 4:
		fmt.Println(3)
		message := f.step_cancel_booking1(payload.BookingID)
		f.Send(r, page.AccessToken, message)
		break
	case -1:
		fmt.Println(-1)
		message := f.step_cancel_booking2(payload.PageID, r.ID)
		f.Send(r, page.AccessToken, message)
		break
	case -2:
		fmt.Println(-2)
		message := f.step_cancel_success(payload.PageID, r.ID, payload.BookingID)
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

func (f *FbAPI) step_service(pageID, userID string) interface{} {

	page, _ := f.DB.GetPageByID(pageID)

	// fmt.Println("Page")
	// fmt.Println(page)

	if !page.IsActive {
		message := fbbot.NewTextMessage("ขออภัย ขณะนี้เพจปิดให้บริการชั่วคราว")
		return message
	}

	// t := time.Now().Weekday()
	// if (t == 0 && page.Sun == false) ||
	// 	(t == 1 && page.Mon == false) ||
	// 	(t == 2 && page.Tue == false) ||
	// 	(t == 3 && page.Wed == false) ||
	// 	(t == 4 && page.Thu == false) ||
	// 	(t == 5 && page.Fri == false) ||
	// 	(t == 6 && page.Sat == false) {
	// 	message := fbbot.NewTextMessage("ขออภัย ขณะนี้ไม่อยู่ในช่วงเวลาทำการ")
	// 	return message
	// }

	// a := time.Now().Format("15:04:05")

	// if (a < page.StartTime || a > page.EndTime) ||
	// 	(a >= page.BreakStart && a < page.BreakEnd) {
	// 	message := fbbot.NewTextMessage("ขออภัย ขณะนี้ไม่อยู่ในช่วงเวลาทำการ")
	// 	return message
	// }
	filter := []bson.M{}

	filter = append(filter, bson.M{"page_id": bson.M{"$eq": pageID}})

	services, _ := f.DB.GetServicesAccordingFilter(filter)
	message := fbbot.NewQuickRepliesMessage("คุณต้องการจองบริการใด?")

	for _, service := range services {
		message.AddQuickRepliesItems(
			fbbot.NewQuickRepliesText(service.Name,
				fmt.Sprintf(`{"step_id":%d, "page_id":"%s", "user_id":"%s", "service_id":"%s"}`, 2, pageID, userID, primitive.ObjectID.Hex(service.ID))),
		)
	}
	return message
}

func (f *FbAPI) step_booking(pageID, serviceID, userID string) interface{} {
	message := fbbot.NewButtonMessage("กรุณาเลือกวันและเวลา")

	message.AddWebURLButton("Jongme", fmt.Sprintf("%s/booking/%s/%s/%s", config.WebURL, pageID, serviceID, userID))
	return message
}

func (f *FbAPI) step_cancel_booking2(pageID, userID string) interface{} {

	filter := []bson.M{}

	filter = append(filter, bson.M{"user_id": bson.M{"$eq": userID}})
	filter = append(filter, bson.M{"status": bson.M{"$eq": 0}})

	bookings, _ := f.DB.GetBookingsAccordingFilter(filter)
	message := fbbot.NewQuickRepliesMessage("กรุณาเลือกบริการที่ต้องการยกเลิก")

	for _, booking := range bookings {
		// fmt.Printf("%+v\n", booking)
		message.AddQuickRepliesItems(
			fbbot.NewQuickRepliesText(
				// booking.name,
				fmt.Sprintf(`%s %s`, booking.Name, booking.Time),
				fmt.Sprintf(`{"step_id":%d, "page_id":"%s", "user_id":"%s", "booking_id":"%s"}`, -2, pageID, userID, primitive.ObjectID.Hex(booking.ID))),
		)
	}

	return message
}

func (f *FbAPI) step_cancel_success(pageID, userID, bookingID string) interface{} {
	fmt.Println("choose")
	id, _ := primitive.ObjectIDFromHex(bookingID)
	err := f.DB.DeleteBookingByID(id)

	var message *fbbot.TextMessage
	if err != nil {
		message = fbbot.NewTextMessage("พบข้อผิดพลาด กรุณาลองอีกครั้ง")
	} else {
		message = fbbot.NewTextMessage("ล้างรายการสำเร็จ")
	}

	return message
}

func (f *FbAPI) step_confirm_booking(bookingID string) interface{} {

	id, _ := primitive.ObjectIDFromHex(bookingID)
	filter := []bson.M{}

	filter = append(filter, bson.M{"_id": bson.M{"$eq": id}})
	bookings, _ := f.DB.GetBookingsAccordingFilter(filter)

	booking := bookings[0]

	booking.Status = 1

	_, _ = f.DB.UpdateBooking(booking)

	message := fbbot.NewTextMessage("ยืนยันสำเร็จ")

	return message

}

func (f *FbAPI) step_cancel_booking1(bookingID string) interface{} {

	id, _ := primitive.ObjectIDFromHex(bookingID)
	filter := []bson.M{}

	filter = append(filter, bson.M{"_id": bson.M{"$eq": id}})
	bookings, _ := f.DB.GetBookingsAccordingFilter(filter)

	booking := bookings[0]

	_ = f.DB.DeleteBookingByID(booking.ID)

	message := fbbot.NewTextMessage("ล้างรายการสำเร็จ")

	return message

}
