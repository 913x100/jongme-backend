package fbbot

import (
	"fmt"
	"jongme/app/config"
	"jongme/app/network"
)

type FB struct {
	Client network.Client
}

// func (f *FB) GetPageAccessToken(pageID, userToken string) (res network.Response) {
// 	fmt.Println("PageID: ", pageID, "UserToken : ", userToken)
// 	uri := config.FacebookAPIEndpoint + pageID + "?fields=id,name,access_token&access_token=" + userToken
// 	res = f.Client.Get(uri)
// 	checkFBErrors(&res)
// 	fmt.Println(res)
// 	return
// }

func (f FB) OauthWithCode(code string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "oauth/access_token?client_id=" + config.AppID + "&client_secret=" + config.AppSecret + "&redirect_uri=" + config.WebURL + "/auth&code=" + string(code)
	res = f.Client.Get(uri)
	// fmt.Println(res)
	return
}

func (f FB) GetLongLiveToken(accessToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "oauth/access_token?grant_type=fb_exchange_token&client_id=" + config.AppID + "&client_secret=" + config.AppSecret + "&fb_exchange_token=" + accessToken
	res = f.Client.Get(uri)
	// fmt.Println(res)
	return
}

// func (f *FB) CheckPermission(userToken string) network.Response {
// 	uri := config.FacebookAPIEndpoint + "me/permissions?access_token=" + userToken
// 	res := f.Client.Get(uri)
// 	// fmt.Println(res)
// 	checkFBErrors(&res)
// 	return res
// }

func (f FB) GetUserInfo(userToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "me?fields=id,name,email,picture.type(large)&access_token=" + userToken
	res = f.Client.Get(uri)
	// fmt.Println(res)
	checkFBErrors(&res)
	return
}

// func (f *FB) GetPages(userToken string) (res network.Response) {
// 	uri := config.FacebookAPIEndpoint + "me/accounts?fields=id,name,picture.type(large)&access_token=" + userToken
// 	res = f.Client.Get(uri)
// 	// fmt.Println(res)
// 	checkFBErrors(&res)
// 	return
// }

func (f FB) GetPages(userToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "me/accounts?fields=id,name&access_token=" + userToken
	res = f.Client.Get(uri)
	return res
}

func (f FB) GetPageToken(pageID, userToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + pageID + "/?fields=id,name,access_token&access_token=" + userToken
	res = f.Client.Get(uri)
	return res
}

func (f FB) SubscribeWebhook(pageID, accessToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + pageID + "/subscribed_apps?access_token=" + accessToken + "&subscribed_fields=" + "messages,messaging_postbacks,message_deliveries,message_reads,message_deliveries,messaging_referrals,standby"
	res = f.Client.Post(nil, uri)
	// fmt.Println(res)
	checkFBErrors(&res)
	return
}

func (f FB) EnabledGetStarted(pageAccessToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "/me/messenger_profile?access_token=" + pageAccessToken
	body := map[string]interface{}{
		"get_started": map[string]string{"payload": `{ "stepid":"0", "payload":"Get Started"}`},
	}
	res = f.Client.Post(body, uri)
	checkFBErrors(&res)
	return
}

func (f FB) AddPersistentMenus(pageAccessToken string, menus ...*Menu) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "/me/messenger_profile?access_token=" + pageAccessToken
	body := make(map[string]interface{})
	body["persistent_menu"] = menus
	res = f.Client.Post(body, uri)
	checkFBErrors(&res)
	return
}

func (f FB) SendQuickRepliesMessage(recipientID, pageAccessToken string, m *QuickRepliesMessage) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "me/messages?access_token=" + pageAccessToken
	body := make(map[string]interface{})
	body["messaging_type"] = "RESPONSE"
	body["recipient"] = map[string]string{"id": recipientID}
	body["message"] = m
	res = f.Client.Post(body, uri)
	checkFBErrors(&res)
	return
}

func (f FB) SendTextMessage(recipientID, pageAccessToken string, m *TextMessage) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "me/messages?access_token=" + pageAccessToken
	body := make(map[string]interface{})

	body["messaging_type"] = "RESPONSE"
	body["recipient"] = map[string]string{"id": recipientID}
	body["message"] = map[string]string{"text": m.Text}
	res = f.Client.Post(body, uri)
	checkFBErrors(&res)
	return
}
func (f FB) SendGenericMessage(recipientID, pageAccessToken string, m *GenericMessage) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "me/messages?access_token=" + pageAccessToken
	payload := make(map[string]interface{})
	payload["template_type"] = "generic"
	payload["elements"] = m.Elements

	attachment := make(map[string]interface{})
	attachment["type"] = "template"
	attachment["payload"] = payload

	body := make(map[string]interface{})
	body["messaging_type"] = "RESPONSE"
	body["recipient"] = map[string]string{"id": recipientID}
	body["message"] = map[string]interface{}{"attachment": attachment}

	res = f.Client.Post(body, uri)
	checkFBErrors(&res)
	return
}

func (f FB) SendButtonMessage(recipientID, pageAccessToken string, m *ButtonMessage) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "me/messages?access_token=" + pageAccessToken

	body := make(map[string]interface{})

	payload := make(map[string]interface{})
	payload["template_type"] = "button"
	payload["text"] = m.Text
	payload["buttons"] = m.Buttons

	attachment := make(map[string]interface{})
	attachment["type"] = "template"
	attachment["payload"] = payload

	body["messaging_type"] = "RESPONSE"
	body["recipient"] = map[string]string{"id": recipientID}
	body["message"] = map[string]interface{}{"attachment": attachment}

	res = f.Client.Post(body, uri)
	checkFBErrors(&res)
	return
}

func checkFBErrors(res *network.Response) {
	if res.Err != nil {
		err, ok := res.Response["error"].(map[string]interface{})
		if ok {
			code, _ := err["code"].(float64)
			subcode, _ := err["error_subcode"].(float64)
			fmt.Println("FB ErrorCode: ", code, subcode)
			delete(res.Response, "error")
			res.Response = getCodeErrorInfo(int(code), int(subcode))
		}
	}
}

func getCodeErrorInfo(code, subcode int) map[string]interface{} {
	// 3 try api again // 2 permission error 1 // login agin
	res := map[string]interface{}{}

	if code == 0 {
		return res
	}
	switch code {
	case 1:
		res["code"] = 3
		fmt.Println("unknown error")
	case 190: // accesstoken expired
		fmt.Println("accesstoken error")
		res["code"] = 1
	case 10: // permission error
		fmt.Println("permission error")
		res["code"] = 2
	case 200:
		fmt.Println("permission error")
		res["code"] = 2
	case 100:
		fmt.Println("unknown error")
		res["code"] = 3
	}

	if subcode == 0 {
		return res
	} // thannls
	switch subcode {
	case 458:
		fmt.Println("APP uninstalled")
	case 460:
		fmt.Println("Session Expired")
	}
	return res
}
