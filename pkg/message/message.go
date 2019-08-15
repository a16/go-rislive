package rislive

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type RisLiveMessage struct {
	Type       string                  `json:"type"`
	BgpMsgType string                  `json:"-"`
	Data       RisLiveMessageInterface `json:"data,omitempty"`
}

type RisLiveMessageInterface interface {
	Dummy()
}

func (m *RisLiveMessage) UnmarshalJSON(buf []byte) error {
	type Alias RisLiveMessage
	a := struct {
		Data json.RawMessage `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(buf, &a); err != nil {
		return err
	}
	switch m.Type {
	case "ris_subscribe":
		// To Be Implemented
	case "ris_unsubscribe":
		// To Be Implemented
	case "request_rrc_list":
		m.Data = nil
	case "ping":
		m.Data = nil
	case "ris_message":
		l := map[string]json.RawMessage{}
		if err := json.Unmarshal(a.Data, &l); err != nil {
			return err
		}
		m.BgpMsgType, _ = strconv.Unquote(string(l["type"]))
		switch m.BgpMsgType {
		case "OPEN":
			var o RisMessageOpen
			if err := json.Unmarshal(a.Data, &o); err != nil {
				return err
			}
			m.Data = &o
		case "UPDATE":
			var u RisMessageUpdate
			if err := json.Unmarshal(a.Data, &u); err != nil {
				return err
			}
			m.Data = &u
		case "KEEPALIVE":
			var k RisMessageKeepalive
			if err := json.Unmarshal(a.Data, &k); err != nil {
				return err
			}
			m.Data = &k
		case "NOTIFICATION":
			var n RisMessageNotification
			if err := json.Unmarshal(a.Data, &n); err != nil {
				return err
			}
			m.Data = &n
		case "RIS_PEER_STATE":
			var r RisMessageRisPeerState
			if err := json.Unmarshal(a.Data, &r); err != nil {
				return err
			}
			m.Data = &r
		default:
			return fmt.Errorf("unknown type in ris_message: %s", m.BgpMsgType)
		}
	case "ris_error":
		var re RisError
		if err := json.Unmarshal(a.Data, &re); err != nil {
			return err
		}
		m.Data = &re
	case "ris_rrc_list":
		var rrl RisRrcList
		if err := json.Unmarshal(a.Data, &rrl); err != nil {
			return err
		}
		m.Data = rrl
	case "pong":
		m.Data = nil
	default:
		return fmt.Errorf("unknown type: %s", m.Type)
	}
	//DebugOutput(m)
	return nil
}

type Filter struct {
	Host          string            `json:"host,omitempty"`
	Type          string            `json:"type,omitempty"`
	Require       string            `json:"require,omitempty"`
	Peer          string            `json:"peer,omitempty"`
	Path          string            `json:"path,omitempty"`
	Prefix        string            `json:"prefix,omitempty"`
	MoreSpecific  bool              `json:"moreSpecific,omitempty"`
	LessSpecific  bool              `json:"lessSpecific,omitempty"`
	SocketOptions *RisSocketOptions `json:"socketOptions,omitempty"`
}

type RisSocketOptions struct {
	IncludeRaw bool `json:"includeRaw,omitempty"`
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) Dummy() {
}

func (f *Filter) SetHost(rrc string) {
	f.Host = rrc
}

func (f *Filter) SetType(msgType string) {
	f.Type = msgType
}

func (f *Filter) SetRequire(key string) {
	f.Require = key
}

func (f *Filter) SetPeer(ip string) {
	f.Peer = ip
}

func (f *Filter) SetPrefix(prefix string, moreSpecific, lessSpecific bool) {
	f.Prefix = prefix
	f.MoreSpecific = moreSpecific
	f.LessSpecific = lessSpecific
}

func (f *Filter) SetSocketOptions(includeRaw bool) {
	f.SocketOptions = &RisSocketOptions{
		IncludeRaw: includeRaw,
	}
}

func NewRisSubscribe(filter *Filter) *RisLiveMessage {
	return &RisLiveMessage{
		Type: "ris_subscribe",
		Data: filter,
	}
}

func NewRisUnsubscribe(filter *Filter) *RisLiveMessage {
	return &RisLiveMessage{
		Type: "ris_unsubscribe",
		Data: filter,
	}
}

func NewRisRequestRrcList() *RisLiveMessage {
	return &RisLiveMessage{
		Type: "request_rrc_list",
	}
}

func NewRisPing() *RisLiveMessage {
	return &RisLiveMessage{
		Type: "ping",
	}
}

type RisMessageInterface interface {
	GetTimestamp() time.Time
	GetType() string
}

type RisMessageCommon struct {
	Type      string  `json:"type"`
	Timestamp float64 `json:"timestamp"`
	Peer      string  `json:"peer"`
	PeerASN   string  `json:"peer_asn"`
	ID        string  `json:"id"`
	Host      string  `json:"host"`
	Raw       string  `json:"raw,omitempty"`
}

func (m RisMessageCommon) Dummy() {}

type RisMessageOpen struct {
	*RisMessageCommon
	Direction    string                         `json:"direction"`
	RouterID     string                         `json:"router_id"`
	Version      int                            `json:"version"`
	Capabilities map[string]CapabilityInterface `json:"capabilities"`
	HoldTime     int                            `json:"hold_time"`
}

// TODO: Capabilities are currently not implemented
type CapabilityInterface interface {
}

type RisMessageUpdate struct {
	*RisMessageCommon
	Path          []uint32       `json:"path,omitempty"`
	Communities   [][]uint32     `json:"community,omitempty"`
	Origin        string         `json:"origin,omitempty"`
	MED           uint32         `json:"med,omitempty"`
	Announcements []Announcement `json:"announcements,omitempty"`
	Withdrawals   []string       `json:"withdralwals,omitempty"`
}

type RisMessageNotification struct {
	*RisMessageCommon
	Notification struct {
		Code    uint8  `json:"code"`
		Subcode uint8  `json:"subcode"`
		Data    string `json:"data"`
	} `json:"notification"`
}

type RisMessageKeepalive struct {
	*RisMessageCommon
}

type RisMessageRisPeerState struct {
	*RisMessageCommon
	State string `json:"state"`
}

type Announcement struct {
	NextHop  string   `json:"next_hop"`
	Prefixes []string `json:"prefixes"`
}

type RisError struct {
	CommandType string `json:"command_type"`
	Message     string `json:"message"`
}

func (m RisError) Dummy() {
}

type RisRrcList []string

func (m RisRrcList) Dummy() {
}

type RisPong struct {
}

func (m RisPong) Dummy() {
}

func DebugOutput(m *RisLiveMessage) {
	j, _ := json.Marshal(m)
	fmt.Println(string(j))
}
