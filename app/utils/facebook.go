package utils

import (
	"fmt"
	"jongme/app/config"
	"jongme/app/network"
)

type FB struct {
	Client network.Client
}

func (f *FB) GetPageAccessToken(pageID, userToken string) (res network.Response) {
	fmt.Println("PageID: ", pageID, "UserToken : ", userToken)
	uri := config.FacebookAPIEndpoint + pageID + "?fields=id,name,access_token&access_token=" + userToken
	res = f.Client.Get(uri)
	checkFBErrors(&res)
	fmt.Println(res)
	return
}

func (f *FB) OauthWithCode(code string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "oauth/access_token?client_id=" + config.AppID + "&client_secret=" + config.AppSecret + "&redirect_uri=" + config.WebURL + "/auth&code=" + string(code)
	res = f.Client.Get(uri)
	// fmt.Println(res)
	return
}

func (f *FB) GetLongLiveToken(accessToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "oauth/access_token?grant_type=fb_exchange_token&client_id=" + config.AppID + "&client_secret=" + config.AppSecret + "&fb_exchange_token=" + accessToken
	res = f.Client.Get(uri)
	// fmt.Println(res)
	return
}

func (f *FB) CheckPermission(userToken string) network.Response {
	uri := config.FacebookAPIEndpoint + "me/permissions?access_token=" + userToken
	res := f.Client.Get(uri)
	// fmt.Println(res)
	checkFBErrors(&res)
	return res
}

func (f *FB) GetUserInfo(userToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "me?fields=id,name,email,picture.type(large)&access_token=" + userToken
	res = f.Client.Get(uri)
	// fmt.Println(res)
	checkFBErrors(&res)
	return
}

// func (f *FB) GetPages(userToken string) (res network.Response) {
// 	uri := config.FacebookAPIEndpoint + "me/accounts?fields=id,name,picture.type(large)&access_token=" + userToken
// 	res = f.Client.Get(uri)
// 	fmt.Println(res)
// 	checkFBErrors(&res)
// 	return
// }

// func (f *FB) GetPages2(userID, userToken string) (res network.Response) {
// 	uri := config.FacebookAPIEndpoint + "me/accounts?fields=id,name,access_token,picture.type(large)&access_token=" + userToken
// 	res = f.Client.Get(uri)
// 	fmt.Println(res)
// 	checkFBErrors(&res)
// 	return
// }

func (f *FB) GetPageWithToken(userToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "me/accounts?fields=id,name,access_token,picture.type(large)&access_token=" + userToken
	res = f.Client.Get(uri)
	checkFBErrors(&res)

	return
}

func (f *FB) SubscribeWebhook(pageID, pageAccessToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + pageID + "/subscribed_apps?access_token=" + pageAccessToken + "&subscribed_fields=" + "feed,messages,messaging_postbacks,message_deliveries,message_reads,message_deliveries,messaging_referrals"
	res = f.Client.Post(nil, uri)
	// fmt.Println(res)
	checkFBErrors(&res)
	return
}

func (f *FB) CheckWebhookSubscription(pageID, pageAccessToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + pageID + "/subscribed_apps?access_token=" + pageAccessToken
	res = f.Client.Get(uri)
	// fmt.Println(res)
	return
}

func (f *FB) RemoveWebhookSuscription(pageID, appAccessToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + pageID + "/subscribed_apps?access_token" + appAccessToken
	// res = f.Client.DELETE(uri) // implemnt DELETE
	fmt.Println(uri)
	return
}

func (f *FB) GetStartedPayload(pageAccessToken string) (res network.Response) {
	uri := config.FacebookAPIEndpoint + "/me/messenger_profile?access_token=" + pageAccessToken
	body := map[string]interface{}{
		"get_started": map[string]interface{}{"payload": `{ "stepid":"", "flowid":"", "payload":"getstarted"}`},
	}
	res = f.Client.Post(body, uri)
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
