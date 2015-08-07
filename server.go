package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/naoina/denco"
	"github.com/naoina/toml"
)

type config struct {
	SlackToken string
	Datasource string
}

func configFromFile() config {
	f, err := os.Open("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var cfg config
	if err = toml.Unmarshal(buf, &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

type channelListResponse struct {
	Channels []channel `json:"channels"`
}

type channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ChannelList(c config) func(w http.ResponseWriter, r *http.Request, params denco.Params) {
	return func(w http.ResponseWriter, r *http.Request, params denco.Params) {
		db, err := sql.Open("mysql", c.Datasource)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query("SELECT id, name FROM channels")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		channels := make([]channel, 0)
		for rows.Next() {
			var id string
			var name string

			if err := rows.Scan(&id, &name); err != nil {
				log.Fatal(err)
			}
			channels = append(channels, channel{ID: id, Name: name})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&channelListResponse{Channels: channels})
	}
}

type user struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type message struct {
	User      user   `json:"user"`
	Text      string `json:"text"`
	Timestamp string `json:"ts"`
}

type messageListResponse struct {
	Channel  string    `json:"channel"`
	Messages []message `json:"messages"`
}

func MessageList(c config) func(w http.ResponseWriter, r *http.Request, params denco.Params) {
	return func(w http.ResponseWriter, r *http.Request, params denco.Params) {
		db, err := sql.Open("mysql", c.Datasource)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query(`
		SELECT
		    u.name AS user_name,
			u.icon_url AS user_icon,
			m.text, m.timestamp
		FROM
		    messages m
			INNER JOIN channels c ON (m.channel_id = c.id)
			INNER JOIN users u ON (m.user_id = u.id)
	    WHERE
		    c.name = ? AND
			DATE(m.timestamp) = ?
		`, params.Get("channel"), fmt.Sprintf("%s-%s-%s", params.Get("year"), params.Get("month"), params.Get("day")))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		messages := make([]message, 0)
		for rows.Next() {
			var user_name string
			var user_icon string
			var text string
			var ts string

			if err := rows.Scan(&user_name, &user_icon, &text, &ts); err != nil {
				log.Fatal(err)
			}
			msg := message{User: user{Name: user_name, Icon: user_icon}, Text: text, Timestamp: ts}
			messages = append(messages, msg)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&messageListResponse{Messages: messages, Channel: params.Get("channel")})
	}
}

func main() {
	c := configFromFile()
	mux := denco.NewMux()

	handler, err := mux.Build([]denco.Handler{
		mux.GET("/channels", ChannelList(c)),
		mux.GET("/log/:channel/:year/:month/:day", MessageList(c)),
	})

	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":8080", handler)
}