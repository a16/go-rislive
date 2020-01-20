package rislive

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var examples = []struct {
	Description           string
	Type                  string
	BgpMsgType            string
	Timestamp             float64
	Peer                  string
	PeerASN               string
	ID                    string
	Raw                   string
	Host                  string
	State                 string
	ReceivedMsg           string
	ExpectedRawBGPMessage string
}{
	{
		Description: "Unmarshal ris_message(OPEN)",
		Type:        "ris_message",
		Timestamp:   1562841440.23,
		Peer:        "2001:7f8:4::1ad2:1",
		PeerASN:     "6866",
		ID:          "2001:7f8:4::1ad2:1-1562841440.23-403701",
		Raw:         "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF004F01041AD200B4C30E986532020601040002000102028000020202000206410400001AD202084006007800020100020E050C000100010002000100020002",
		Host:        "rrc01",
		BgpMsgType:  "OPEN",
		State:       "",
		ReceivedMsg: `{
			"type": "ris_message",
			"data": {
		 			"timestamp": 1562841440.23,
					"peer": "2001:7f8:4::1ad2:1",
					"peer_asn": "6866",
					"id": "2001:7f8:4::1ad2:1-1562841440.23-403701",
					"raw": "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF004F01041AD200B4C30E986532020601040002000102028000020202000206410400001AD202084006007800020100020E050C000100010002000100020002",
					"host": "rrc01",
					"type": "OPEN",
					"direction": "received",
					"version": 4,
					"asn": 6866,
					"hold_time": 180,
					"router_id": "195.14.152.101",
					"capabilities": {
							"1": {
								"name": "multiprotocol",
								"families": ["ipv6/unicast"]
							},
							"2": {
								"name": "route-refresh",
								"variant": "RFC"
							},
							"5": {
									"name": "unknown",
									"iana": "unknown",
									"value": 5,
									"raw": "000100010002000100020002"
							},
							"64": {
									"name": "graceful restart",
									"time": 120,
									"address family flags": {
										"ipv6/unicast": []
									},
									"restart flags": []
							},
							"65": {
									"name": "asn4",
									"asn4": 6866
							},
							"128": {
									"name": "route-refresh",
									"variant": "RFC"
							}
					}
			}
		}`,
		ExpectedRawBGPMessage: "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF004F01041AD200B4C30E986532020601040002000102028000020202000206410400001AD202084006007800020100020E050C000100010002000100020002",
	},
	{
		Description: "Unmarshal ris_message(UPDATE)",
		Type:        "ris_message",
		Timestamp:   1562822233.68,
		Peer:        "195.208.208.147",
		PeerASN:     "28917",
		ID:          "195.208.208.147-1562822233.68-150306082",
		Raw:         "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF006A020004148D8820002F400101004002160205000070F500000CB9000005130004155D000402ED400304C3D0D093C0080870F50FA070F50FA318B1177418B1177718B1177018A879C518B1260D18A879C718B1260A18B1260F",
		Host:        "rrc13",
		BgpMsgType:  "UPDATE",
		State:       "",
		ReceivedMsg: `{
			"type": "ris_message",
			"data": {
				"timestamp": 1562822233.68,
				"peer": "195.208.208.147",
				"peer_asn": "28917",
				"id": "195.208.208.147-1562822233.68-150306082",
				"raw": "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF006A020004148D8820002F400101004002160205000070F500000CB9000005130004155D000402ED400304C3D0D093C0080870F50FA070F50FA318B1177418B1177718B1177018A879C518B1260D18A879C718B1260A18B1260F",
				"host": "rrc13",
				"type": "UPDATE",
				"path": [28917, 3257, 1299, 267613, 262893],
				"community": [[28917, 4000], [28917, 4003]],
				"origin": "igp",
				"announcements": [{"next_hop": "195.208.208.147", "prefixes": ["177.23.116.0/24", "177.23.119.0/24", "177.23.112.0/24", "168.121.197.0/24", "177.38.13.0/24", "168.121.199.0/24", "177.38.10.0/24", "177.38.15.0/24"]}],
				"withdrawals": ["141.136.32.0/20"]
			}
		}`,
		ExpectedRawBGPMessage: "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF006A020004148D8820002F400101004002160205000070F500000CB9000005130004155D000402ED400304C3D0D093C0080870F50FA070F50FA318B1177418B1177718B1177018A879C518B1260D18A879C718B1260A18B1260F",
	},
	{
		Description: "Unmarshal ris_message(NOTIFICATION)",
		Type:        "ris_message",
		Timestamp:   1562822895.4,
		Peer:        "2606:6d00:eb0::254",
		PeerASN:     "1403",
		ID:          "2606:6d00:eb0::254-1562822895.4-519878",
		Raw:         "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0015030605",
		Host:        "rrc00",
		BgpMsgType:  "NOTIFICATION",
		State:       "",
		ReceivedMsg: `{
			"type": "ris_message",
			"data": {
				"timestamp": 1562822895.4,
				"peer": "2606:6d00:eb0::254",
				"peer_asn": "1403",
				"id": "2606:6d00:eb0::254-1562822895.4-519878",
				"raw": "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0015030605",
				"host": "rrc00",
				"type": "NOTIFICATION",
				"notification": {
					"code": 6,
					"subcode": 5,
					"data": "0605"
				}
			}
		}`,
		ExpectedRawBGPMessage: "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0015030605",
	},
	{
		Description: "Unmarshal ris_message(KEEPALIVE)",
		Type:        "ris_message",
		Timestamp:   1562822767.1,
		Peer:        "195.66.224.31",
		PeerASN:     "32787",
		ID:          "195.66.224.31-1562822767.1-1248612",
		Raw:         "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF001304",
		Host:        "rrc01",
		BgpMsgType:  "KEEPALIVE",
		State:       "",
		ReceivedMsg: `{
			"type": "ris_message",
			"data": {
				"timestamp": 1562822767.1,
				"peer": "195.66.224.31",
				"peer_asn": "32787",
				"id": "195.66.224.31-1562822767.1-1248612",
				"raw": "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF001304",
				"host": "rrc01",
				"type": "KEEPALIVE"
			}
		}`,
		ExpectedRawBGPMessage: "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF001304",
	},
	{
		Description: "Unmarshal ris_message(RIS_PEER_STATE)",
		Type:        "ris_message",
		Timestamp:   1562823052.55,
		Peer:        "2001:43f8:6d0::55",
		PeerASN:     "327991",
		ID:          "2001:43f8:6d0::55-1562823052.55-1007659",
		Raw:         "",
		Host:        "rrc19",
		BgpMsgType:  "RIS_PEER_STATE",
		State:       "connected",
		ReceivedMsg: `{
			"type": "ris_message",
			"data": {
				"timestamp": 1562823052.55,
				"peer": "2001:43f8:6d0::55",
				"peer_asn": "327991",
				"id": "2001:43f8:6d0::55-1562823052.55-1007659",
				"host": "rrc19",
				"type": "RIS_PEER_STATE",
				"state": "connected"
			}
		}`,
	},
	{
		Description: "Unmarshal ris_rrc_list",
		Type:        "ris_rrc_list",
		Timestamp:   0,
		Peer:        "",
		PeerASN:     "",
		ID:          "",
		Raw:         "",
		Host:        "",
		BgpMsgType:  "",
		State:       "",
		ReceivedMsg: `{
			"type": "ris_rrc_list",
			"data": [
				"rrc00",
				"rrc01"
			]
		}`,
	},
	{
		Description: "Unmarshal ris_error",
		Type:        "ris_error",
		Timestamp:   0,
		Peer:        "",
		PeerASN:     "",
		ID:          "",
		Raw:         "",
		Host:        "",
		BgpMsgType:  "",
		State:       "",
		ReceivedMsg: `{
			"type": "ris_error",
			"data": {
				"message": "Unknown command type",
				"command_type":"wrong"
			}
		}`,
	},
	{
		Description: "Unmarshal pong",
		Type:        "pong",
		Timestamp:   0,
		Peer:        "",
		PeerASN:     "",
		ID:          "",
		Raw:         "",
		Host:        "",
		BgpMsgType:  "",
		State:       "",
		ReceivedMsg: `{
			"type":"pong"
		}`,
	},
	{
		Description: "Unmarshal ris_error",
		Type:        "ris_error",
		Timestamp:   0,
		Peer:        "",
		PeerASN:     "",
		ID:          "",
		Raw:         "",
		Host:        "",
		BgpMsgType:  "",
		State:       "",
		ReceivedMsg: `{
			"type": "ris_error",
			"data": {
				"message": "Closing connection after being behind by more than 262144000 bytes over 30 seconds",
				"bufferSize":313026734
			}
		}`,
	},
}

func TestUnmarshalJSON(t *testing.T) {
	for _, ex := range examples {
		t.Run(ex.Description, func(t *testing.T) {
			assert := assert.New(t)
			var r RisLiveMessage
			err := json.Unmarshal([]byte(ex.ReceivedMsg), &r)
			assert.NoError(err)
			assert.Equal(ex.Type, r.Type)
			assert.Equal(ex.Timestamp, r.Timestamp)
			assert.Equal(ex.Peer, r.Peer)
			assert.Equal(ex.PeerASN, r.PeerASN)
			assert.Equal(ex.ID, r.ID)
			assert.Equal(ex.Host, r.Host)
			assert.Equal(ex.BgpMsgType, r.BgpMsgType)
		})
	}
}
