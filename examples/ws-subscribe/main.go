package main

import (
	"encoding/json"
	"log"
	"math"
	"net/url"
	"time"

	rislive "github.com/a16/go-rislive/pkg/message"
	"github.com/gorilla/websocket"
)

func main() {
	queueLen := 0
	maxWorkers := 1

	values := url.Values{}
	values.Add("client", "go-rislive-gorilla")
	u := url.URL{
		Scheme:   "wss",
		Host:     "ris-live.ripe.net:443",
		Path:     "/v1/ws/",
		RawQuery: values.Encode(),
	}
	log.Printf("connecting to %s\n", u.String())

	log.Printf("Connecting to RIS Live server")
	var err error
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("Could not connect RIS Live server: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("Connected to RIS Live server")

	queue := make(chan *rislive.RisLiveMessage, queueLen)
	defer close(queue)
	doneCh := make(chan struct{})

	for i := 0; i < maxWorkers; i++ {
		go risliveWorker(queue, doneCh)
	}

	filter := rislive.NewFilter()
	filter.SetSocketOptions(true)
	subReq := rislive.NewRisSubscribe(filter)

	if err := conn.WriteJSON(subReq); err != nil {
		log.Printf("Failed to write ris_subscribe message: %v", err)
		return
	}

	go func() {
		defer func() {
			doneCh <- struct{}{}
		}()
		for {
			var msg rislive.RisLiveMessage
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Panicf("In ReadMessage: %v, %v", err, string(p))
				continue
			}
			if err := json.Unmarshal(p, &msg); err != nil {
				log.Panicf("In Unmarshal: %v, %v", err, string(p))
			}
			queue <- &msg
		}
	}()

	for {
		select {
		case <-doneCh:
			return
		}
	}
}

func FloatToTime(f float64) time.Time {
	t1, t2 := math.Modf(f)
	t2 = math.Ceil(t2 * math.Pow10(6))
	return time.Unix(int64(t1), int64(t2)).UTC()
}

func risliveWorker(queue chan *rislive.RisLiveMessage, doneCh chan struct{}) {
	defer func() {
		doneCh <- struct{}{}
	}()

	t := time.NewTicker(1 * time.Minute)
	defer t.Stop()
	counter := 0
	for msg := range queue {
		switch msg.Type {
		case "ris_error":
			risErr := msg.Data.(*rislive.RisError)
			log.Printf("ris_error: %v, %v", risErr.CommandType, risErr.Message)
		case "ris_message":
			counter += 1
			switch msg.BgpMsgType {
			case "OPEN":
				risMsgOpen := msg.Data.(*rislive.RisMessageOpen)
				log.Printf("ris_message(OPEN): %v, %v", risMsgOpen.Timestamp, risMsgOpen.Raw)
			case "UPDATE":
				risMsgUpdate := msg.Data.(*rislive.RisMessageUpdate)
				log.Println(msg.Data)
				log.Printf("ris_message(UPDATE): %v, %v", risMsgUpdate.Timestamp, risMsgUpdate.Raw)
			case "KEEPALIVE":
				risMsgKeepalive := msg.Data.(*rislive.RisMessageKeepalive)
				log.Printf("ris_message(KEEPALIVE): %v, %v", risMsgKeepalive.Timestamp, risMsgKeepalive.Raw)
			case "NOTIFICATION":
				risMsgNotification := msg.Data.(*rislive.RisMessageNotification)
				log.Printf("ris_message(NOTIFICATION): %v, %v", risMsgNotification.Timestamp, risMsgNotification.Raw)
			case "RIS_PEER_STATE":
				risMsgRisPeerState := msg.Data.(*rislive.RisMessageRisPeerState)
				log.Printf("ris_message(PEER_STATE): %v", FloatToTime(risMsgRisPeerState.Timestamp))
			default:
				log.Printf("UNKNOWN: %#v", msg.Data)
				return
			}
		}
	}
}
