package eosapi

import (
	"encoding/hex"
	"encoding/json"
	"reflect"

	"github.com/eosioca/eosapi/ecc"
)

type InfoResp struct {
	ServerVersion            string      `json:"server_version"`              // "2cc40a4e"
	HeadBlockNum             uint32      `json:"head_block_num"`              // 2465669,
	LastIrreversibleBlockNum uint32      `json:"last_irreversible_block_num"` // 2465655
	HeadBlockID              string      `json:"head_block_id"`               // "00259f856bfa142d1d60aff77e70f0c4f3eab30789e9539d2684f9f8758f1b88",
	HeadBlockTime            JSONTime    `json:"head_block_time"`             //  "2018-02-02T04:19:32"
	HeadBlockProducer        AccountName `json:"head_block_producer"`         // "inita"
	RecentSlots              string      `json:"recent_slots"`                //  "1111111111111111111111111111111111111111111111111111111111111111"
	ParticipationRate        string      `json:"participation_rate"`          // "1.00000000000000000" // this should be a `double`, or a decimal of some sort..

}

type BlockResp struct {
	Previous              string           `json:"previous"`                // : "0000007a9dde66f1666089891e316ac4cb0c47af427ae97f93f36a4f1159a194",
	Timestamp             JSONTime         `json:"timestamp"`               // : "2017-12-04T17:12:08",
	TransactionMerkleRoot string           `json:"transaction_merkle_root"` // : "0000000000000000000000000000000000000000000000000000000000000000",
	Producer              AccountName      `json:"producer"`                // : "initj",
	ProducerChanges       []ProducerChange `json:"producer_changes"`        // : [],
	ProducerSignature     string           `json:"producer_signature"`      // : "203dbf00b0968bfc47a8b749bbfdb91f8362b27c3e148a8a3c2e92f42ec55e9baa45d526412c8a2fc0dd35b484e4262e734bea49000c6f9c8dbac3d8861c1386c0",
	Cycles                []Cycle          `json:"cycles"`                  // : [],
	ID                    string           `json:"id"`                      // : "0000007b677719bdd76d729c3ac36bed5790d5548aadc26804489e5e179f4a5b",
	BlockNum              uint64           `json:"block_num"`               // : 123,
	RefBlockPrefix        uint64           `json:"ref_block_prefix"`        // : 2624744919

}

type ProducerChange struct {
}

type Cycle struct {
}

type AccountResp struct {
	AccountName AccountName  `json:"account"`
	Permissions []Permission `json:"permissions"`
}

type CurrencyBalanceResp struct {
	EOSBalance        Asset    `json:"eos_balance"`
	StakedBalance     Asset    `json:"staked_balance"`
	UnstakingBalance  Asset    `json:"unstaking_balance"`
	LastUnstakingTime JSONTime `json:"last_unstaking_time"`
}

type GetTableRowsResp struct {
	More bool            `json:"more"`
	Rows json.RawMessage `json:"rows"` // defer loading, as it depends on `JSON` being true/false.
}

func (resp *GetTableRowsResp) JSONToStructs(v interface{}) error {
	return json.Unmarshal(resp.Rows, v)
}

func (resp *GetTableRowsResp) BinaryToStructs(v interface{}) error {
	var rows []string

	err := json.Unmarshal(resp.Rows, &rows)
	if err != nil {
		return err
	}

	outSlice := reflect.ValueOf(v).Elem()
	structType := reflect.TypeOf(v).Elem().Elem()

	for _, row := range rows {
		bin, err := hex.DecodeString(row)
		if err != nil {
			return err
		}

		// access the type of the `Slice`, create a bunch of them..
		newStruct := reflect.New(structType)
		if err := UnmarshalBinary(bin, newStruct.Interface()); err != nil {
			return err
		}

		outSlice = reflect.Append(outSlice, reflect.Indirect(newStruct))
	}

	reflect.ValueOf(v).Elem().Set(outSlice)

	return nil
}

type Currency struct {
	Precision uint8
	Name      CurrencyName
}

type GetRequiredKeysResp struct {
	RequiredKeys []*ecc.PublicKey `json:"required_keys"`
}

type PushTransactionResp struct {
	TransactionID string `json:"transaction_id"`
	Processed     bool   `json:"processed"` // WARN: is an `fc::variant` in server..
}

type WalletSignTransactionResp struct {
	// Ignore the rest of the transaction, so the wallet server
	// doesn't forge some transactions on your behalf, and you send it
	// to the network..  ... although.. it's better if you can trust
	// your wallet !

	Signatures []string `json:"signatures"`
}

type MyStruct struct {
	Currency
	Balance uint64
}

// NetConnectionResp
type NetConnectionsResp struct {
	Peer          string           `json:"peer"`
	Connecting    bool             `json:"connecting"`
	Syncing       bool             `json:"syncing"`
	LastHandshake HandshakeMessage `json:"last_handshake"`
}

// Decode the `Key`. FIXME: this is unsatisfactory.. we should be able to handle
// broken keys.. perhaps keep the raw PubKey data and decode when we need it instead of.
type HandshakeMessage struct {
	// net_plugin/protocol.hpp handshake_message
	NetworkVersion           int16         `json:"network_version"`
	ChainID                  HexBytes      `json:"chain_id"`
	NodeID                   HexBytes      `json:"node_id"` // sha256
	Key                      ecc.PublicKey `json:"key"`     // can be empty, producer key, or peer key
	Time                     int           `json:"time"`    // time?!
	Token                    HexBytes      `json:"token"`   // digest of time to prove we own the private `key`
	Signature                ecc.Signature `json:"sig"`     // can be empty if no key, signature of the digest above
	P2PAddress               string        `json:"p2p_address"`
	LastIrreversibleBlockNum uint32        `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  HexBytes      `json:"last_irreversible_block_id"`
	HeadNum                  uint32        `json:"head_num"`
	HeadID                   HexBytes      `json:"head_id"`
	OS                       string        `json:"os"`
	Agent                    string        `json:"agent"`
	Generation               int16         `json:"generaiton"`
}

type NetStatusResp struct {
}

type NetConnectResp struct {
}

type NetDisconnectResp struct {
}
