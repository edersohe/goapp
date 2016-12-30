package main

import (
	r "github.com/GoRethink/gorethink"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	glog "github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

var (
	db   *r.Session
	e    *echo.Echo
	mqtt MQTT.Client
	log  echo.Logger
)

func config() {

	e = echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	log = e.Logger

	viper.AutomaticEnv()
	viper.SetDefault("PORT", "5000")
	viper.SetDefault("LOG_LEVEL", 5)
	viper.SetDefault("RETHINKDB_URL", nil)
	viper.SetDefault("RETHINKDB_DATABASE", nil)
	viper.SetDefault("RETHINKDB_USERNAME", nil)
	viper.SetDefault("RETHINKDB_PASSWORD", nil)
	viper.SetDefault("MQTT_URL", nil)
	viper.SetDefault("MQTT_CLIENT_ID", nil)

	log.SetLevel(glog.Lvl(viper.GetInt("LOG_LEVEL")))

	if err := godotenv.Load(); err != nil {
		log.Warn("Error loading .env file")
	}

	if url := viper.GetString("RETHINKDB_URL"); url != "" {
		var err error
		if db, err = r.Connect(r.ConnectOpts{
			Address:  url,
			Database: viper.GetString("RETHINKDB_DATABASE"),
			Username: viper.GetString("RETHINKDB_USERNAME"),
			Password: viper.GetString("RETHINKDB_PASSWORD"),
		}); err != nil {
			log.Fatal(err)
		}
	}

	if viper.Get("MQTT_URL") != nil {
		var debugMessageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
			log.Debugf("TOPIC: %s\n", msg.Topic())
			log.Debugf("MSG: %s\n", msg.Payload())
		}

		opts := MQTT.NewClientOptions().AddBroker(viper.GetString("MQTT_URL"))
		if clientID := viper.GetString("MQTT_CLIENT_ID"); clientID != "" {
			opts.SetClientID(clientID)
		}
		if viper.GetInt("LOG_LEVEL") == 1 {
			opts.SetDefaultPublishHandler(debugMessageHandler)
		}

		mqtt = MQTT.NewClient(opts)
		if token := mqtt.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}

}
