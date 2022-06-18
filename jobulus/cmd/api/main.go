package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jobulus/internal/data"
	"jobulus/internal/jsonlog"
	"os"
	"sync"
	"time"

	"jobulus/internal/discord_bot/bot"
	bot_setting "jobulus/internal/discord_bot/config"

	_ "github.com/go-sql-driver/mysql"
)

const version = "1.0.0"

// @title gin-api-template Swagger API
// @version 1.0
// @description Swagger API for gin backend template.

type config struct {
	Port int    `json:"port"` //Network Port
	Env  string `json:"env"`  //Current operating environment
	Db   struct {
		Dsn          string `json:"dsn"` //Database connection
		MaxOpenConns int    `json:"maxOpenConns"`
		MaxIdleConns int    `json:"maxIdleConns"`
		MaxIdleTime  string `json:"maxIdleTime"`
	} `json:"db"`
	Limiter struct {
		Rps     float64 `json:"rps"`      //Allowed requests per second
		Burst   int     `json:"burst"`    //Num of  maximum requests in single burst
		Enabled bool    `json:"disabled"` //Is Rate Limiter is on
	} `json:"limiter"`
	Token string `json:"token"` //discord bot token

	// cors struct {
	// 	trustedOrigins []string
	// }
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	wg     sync.WaitGroup
	offers chan data.Offer
}

func main() {
	//var cfg config
	conf, err := os.Open("./cmd/configs/config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer conf.Close()

	byteValue, _ := ioutil.ReadAll(conf)

	var configs config
	err = json.Unmarshal(byteValue, &configs)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(configs)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	offers := make(chan data.Offer)

	app := &application{
		config: configs,
		logger: logger,
		models: data.NewModes(db),
		offers: offers,
	}

	//separate thread for discord bot
	go serveDiscordBot(&app.config, app)

	err = app.Serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}

func serveDiscordBot(cfg *config, app *application) {
	err := bot_setting.SetToken(cfg.Token, "!")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jobulus := bot.Start()

	for {
		offer := <-app.offers
		jobulus.SendMessage(offer.Print())
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
