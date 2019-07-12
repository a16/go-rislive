package main

import (
	"bufio"
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"

	rislive "bitbucket.org/a16/go-rislive/pkg/message"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func ConvertTime(f float64) time.Time {
	t1, t2 := math.Modf(f)
	t2 = math.Ceil(t2 * math.Pow10(6))
	return time.Unix(int64(t1), int64(t2)).UTC()
}

func main() {
	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	})
	log.Out = os.Stdout

	u := url.URL{Scheme: "https", Host: "ris-live.ripe.net", Path: "/v1/stream/"}
	q := u.Query()
	q.Set("format", "json")
	q.Set("client", "go-rislive")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		buf := scanner.Bytes()
		var msg rislive.RisLiveMessage
		err = json.Unmarshal(buf, &msg)
		if err != nil {
			log.Fatalf("%v: %v\n", err, scanner.Text())
		}
		switch msg.Type {
		case "ris_message":
			switch msg.BgpMsgType {
			case "UPDATE":
				update := msg.Data.(*rislive.RisMessageUpdate)
				fields := logrus.Fields{
					"Type":     update.Type,
					"RcvdTime": ConvertTime(update.Timestamp),
					"Peer":     update.Peer,
					"PeerASN":  update.PeerASN,
				}
				for _, a := range update.Announcements {
					for _, prefix := range a.Prefixes {
						fields["NextHop"] = a.NextHop
						fields["Prefix"] = prefix
						fields["AsPath"] = update.Path
						fields["Origin"] = update.Origin
						fields["AnnouncementOrWithdrawal"] = "Announcement"
						log.WithFields(fields).Info()
					}
				}
				for _, w := range update.Withdrawals {
					for _, pfx := range w {
						fields["Prefix"] = pfx
						fields["AnnouncementOrWithdrawal"] = "Withdrawal"
						log.WithFields(fields).Info()
					}
				}
			}
		}
	}
}
