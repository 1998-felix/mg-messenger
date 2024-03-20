package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	sdk "github.com/absmach/magistrala/pkg/sdk/go"
	"github.com/caarlos0/env/v10"
)

var (
	channelID = "1409a692-0d37-49f2-987d-a9a279002f27"          // Valid channel ID
	thingKeys = []string{"betty", "bettier", "bettiest"}        // thing secrets
	topics    = []string{"temperature", "humidity", "pressure"} // publisher topics
)

type config struct {
	HTTPAdapterURL  string          `env:"MG_HTTP_ADAPTER_URL"    envDefault:"http://localhost:8008"`
	ReaderURL       string          `env:"MG_READER_URL"          envDefault:"http://localhost:9011"`
	ThingsURL       string          `env:"MG_THINGS_URL"          envDefault:"http://localhost:9000"`
	UsersURL        string          `env:"MG_USERS_URL"           envDefault:"http://localhost:9002"`
	HostURL         string          `env:"MG_UI_HOST_URL"         envDefault:"http://localhost:9097"`
	BootstrapURL    string          `env:"MG_BOOTSTRAP_URL"       envDefault:"http://localhost:9013"`
	DomainsURL      string          `env:"MG_DOMAINS_URL"         envDefault:"http://localhost:8189"`
	InvitationsURL  string          `env:"MG_INVITATIONS_URL"     envDefault:"http://localhost:9020"`
	MsgContentType  sdk.ContentType `env:"MG_UI_CONTENT_TYPE"     envDefault:"application/senml+json"`
	TLSVerification bool            `env:"MG_UI_VERIFICATION_TLS" envDefault:"false"`
}

type myService struct {
	sdk sdk.SDK
}

type Message struct {
	BaseTime float64 `json:"bt"`
	BaseUnit string  `json:"bu"`
	Name     string  `json:"n"`
	Unit     string  `json:"u"`
	Value    float64 `json:"v"`
}

func New(sdk sdk.SDK) *myService {
	return &myService{sdk: sdk}
}

func (ms *myService) Publish(channelID, thingKey string, message Message) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return errors.New("Failed to marshal message")
	}

	messageArray := "[" + string(jsonMessage) + "]"

	if err := ms.sdk.SendMessage(channelID, messageArray, thingKey); err != nil {
		return err
	}

	return nil
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf(err.Error())
	}

	sdkConfig := sdk.Config{
		HTTPAdapterURL:  cfg.HTTPAdapterURL,
		ReaderURL:       cfg.ReaderURL,
		ThingsURL:       cfg.ThingsURL,
		UsersURL:        cfg.UsersURL,
		HostURL:         cfg.HostURL,
		MsgContentType:  cfg.MsgContentType,
		TLSVerification: cfg.TLSVerification,
		BootstrapURL:    cfg.BootstrapURL,
		DomainsURL:      cfg.DomainsURL,
		InvitationsURL:  cfg.InvitationsURL,
	}

	sdk := sdk.NewSDK(sdkConfig)
	svc := New(sdk)

	for {
		for _, thingKey := range thingKeys {
			for _, topic := range topics {
				message := generateMessage(topic)
				err := svc.Publish(channelID, thingKey, message)
				if err == nil {
					fmt.Printf("Published message: %+v\n", message)
				} else {
					fmt.Printf("Failed to publish message: %v\n", err)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func generateMessage(topic string) Message {
	randomValue := rand.Intn(100)
	baseUnit := "C" // default
	unit := "C"     // default

	switch topic {
	case "temperature":
		baseUnit = "C"
		unit = "C"
	case "humidity":
		baseUnit = "%"
		unit = "%"
	case "pressure":
		baseUnit = "hPa"
		unit = "hPa"
	}

	return Message{
		BaseTime: float64(time.Now().Local().UnixNano()),
		BaseUnit: baseUnit,
		Name:     topic,
		Unit:     unit,
		Value:    float64(randomValue),
	}
}
