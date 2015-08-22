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
	"github.com/julienschmidt/httprouter"
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

func ChannelList(c config) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func MessageList(c config) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		ORDER BY
			m.timestamp ASC
		`,
			p.ByName("channel"),
			fmt.Sprintf("%s-%s-%s", p.ByName("year"), p.ByName("month"), p.ByName("day")),
		)
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
		json.NewEncoder(w).Encode(&messageListResponse{Messages: messages, Channel: p.ByName("channel")})
	}
}

type StaticFileHandler struct {};

func (s *StaticFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/index.html")
}

func main() {
	c := configFromFile()

	router := httprouter.New()
	router.Handler("GET", "/", http.FileServer(http.Dir("./assets/")))
	router.NotFound = &StaticFileHandler{}
	router.ServeFiles("/js/*filepath", http.Dir("./assets/js/"))
	router.ServeFiles("/css/*filepath", http.Dir("./assets/css/"))
	router.GET("/channels", ChannelList(c))
	router.GET("/log/:channel/:year/:month/:day", MessageList(c))

	http.ListenAndServe(":8080", router)
}
