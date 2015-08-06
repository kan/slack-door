package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/naoina/toml"
	"github.com/nlopes/slack"
)

func updateChannels(s *slack.Slack, db *sql.DB) error {
	channels, err := s.GetChannels(false)
	if err != nil {
		return err
	}

	for _, c := range channels {
		fmt.Printf("%s : %s\n", c.Name, c.ID)
		_, err := db.Exec(
			`INSERT INTO channels(id, name, created_at, updated_at)
			 VALUES(?, ?, NOW(), NOW())
			 ON DUPLICATE KEY UPDATE name = VALUES(name), updated_at = NOW()
			`,
			c.ID, c.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateUsers(s *slack.Slack, db *sql.DB) error {
	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		fmt.Printf("%s : %s\n", u.Name, u.ID)
		_, err := db.Exec(
			`INSERT INTO users(id, name, icon_url, created_at, updated_at)
			 VALUES(?, ?, ?, NOW(), NOW())
			 ON DUPLICATE KEY UPDATE
			     name = VALUES(name), 
			     icon_url = VALUES(icon_url), 
				 updated_at = NOW()
			`,
			u.ID, u.Name, u.Profile.Image72)
		if err != nil {
			return nil
		}
	}

	return nil
}

func addChannelHistory(s *slack.Slack, db *sql.DB, channel, oldest string) (string, error) {
	params := slack.NewHistoryParameters()
	params.Oldest = oldest
	params.Count = 1000
	history, err := s.GetChannelHistory(channel, params)
	if err != nil {
		return "", err
	}
	for _, msg := range history.Messages {
		if msg.SubType != "" {
			continue
		}
		fmt.Printf("%s: %s\n", msg.User, msg.Text)
		ts := timestampToTime(msg.Timestamp)
		_, err := db.Exec(
			`INSERT INTO messages(channel_id, user_id, text, timestamp, created_at, updated_at)
			 VALUES(?, ?, ?, ?, NOW(), NOW())
			`,
			channel, msg.User, msg.Text, ts.Format("2006-01-02 15:04:05"))
		if err != nil {
			return "", err
		}
	}

	if len(history.Messages) > 0 {
		return history.Messages[0].Timestamp, nil
	} else {
		return "", nil
	}
}

func timestampToTime(timestamp string) time.Time {
	ts, _ := strconv.ParseFloat(timestamp, 64)
	return time.Unix(int64(ts), int64((ts-math.Floor(ts))*1000000000))
}

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

func main() {
	c := configFromFile()
	s := slack.New(c.SlackToken)
	db, err := sql.Open("mysql", c.Datasource)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = updateChannels(s, db)
	if err != nil {
		log.Fatal(err)
	}

	err = updateUsers(s, db)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, name, latest FROM channels")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name string
		var latest sql.NullString

		if err := rows.Scan(&id, &name, &latest); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("import history: #%s\n", name)
		if !latest.Valid {
			latest.String = ""
		}
		l, err := addChannelHistory(s, db, id, latest.String)
		if err != nil {
			log.Fatal(err)
		}
		if l != "" {
			_, err = db.Exec(`UPDATE channels SET latest = ? WHERE id = ?`, l, id)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
