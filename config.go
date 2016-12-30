package main

import (
	r "github.com/GoRethink/gorethink"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

var (
	db   *r.Session
	e    *echo.Echo
	mqtt MQTT.Client
)

func config() {

	e = echo.New()

	e.Logger.SetLevel(log.Lvl(viper.GetInt("LOG_LEVEL")))

	if err := godotenv.Load(); err != nil {
		e.Logger.Warn("Error loading .env file")
	}

	viper.AutomaticEnv()
	viper.SetDefault("PORT", "5000")
	viper.SetDefault("LOG_LEVEL", 1)
	viper.SetDefault("RETHINKDB_URL", nil)
	viper.SetDefault("RETHINKDB_DATABASE", nil)
	viper.SetDefault("RETHINKDB_USERNAME", nil)
	viper.SetDefault("RETHINKDB_PASSWORD", nil)
	viper.SetDefault("MQTT_URL", nil)
	viper.SetDefault("MQTT_CLIENT_ID", nil)

	if url := viper.GetString("RETHINKDB_URL"); url != "" {
		var err error
		if db, err = r.Connect(r.ConnectOpts{
			Address:  url,
			Database: viper.GetString("RETHINKDB_DATABASE"),
			Username: viper.GetString("RETHINKDB_USERNAME"),
			Password: viper.GetString("RETHINKDB_PASSWORD"),
		}); err != nil {
			e.Logger.Fatal(err)
		}
	}

	if viper.Get("MQTT_URL") != nil {
		var debugMessageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
			e.Logger.Debugf("TOPIC: %s\n", msg.Topic())
			e.Logger.Debugf("MSG: %s\n", msg.Payload())
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
			e.Logger.Fatal(token.Error())
		}
	}

	e.Logger.Fatal(e.Start(":" + viper.GetString("PORT")))

}
