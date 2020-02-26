package config

import (
	"os"
)

var (
	AppEnvironment = func() string {
		if os.Getenv("APP_ENV") == "" {
			os.Exit(1)
		}
		return os.Getenv("APP_ENV")
	}()
	NumberOfProducers        = 6
	NumberOfConsumers        = 6
	MongoHost                = os.Getenv("MONGO_HOST")
	MongoPort                = os.Getenv("MONGO_PORT")
	MongoUserName            = os.Getenv("MONGO_USERNAME")
	MongoPassword            = os.Getenv("MONGO_PASSWORD")
	AppID                    = os.Getenv("APP_ID")
	AppSecret                = os.Getenv("APP_SECRET")
	WebURL                   = os.Getenv("WEB_URL")
	FacebookAPIEndpoint      = "https://graph.facebook.com/v5.0/"
	FacebookVideoAPIEndpoint = "https://graph-video.facebook.com/"
	FacebookAPIVersion       = "v4.0"
	BasicAuthPrefix          = []byte("Bearer ")
	TokenIssuer              = os.Getenv("ISSUER")
	TokenSecret              = []byte(os.Getenv("JWT_SECRET"))
	Environment              = os.Getenv("APP_ENV")
	ValidationToken          = os.Getenv("VALIDATION_TOKEN")
)
