package consul

//import "encoding/json"



// APIClient create a api client to the panel.
type KVPut struct {
	Verb                string                `json:"Verb"`
	Key                 string                `json:"Key"`
	Value               string                `json:"Value"`
	Flags               int                   `json:"Flags"`
	Index               int                   `json:"Index"`
	Session             string                `json:"Session"`
}

type KVRes struct {
    LockIndex           int                   `json:"LockIndex"`
    Key                 string                `json:"Key"`
    Flags               int                   `json:"Flags"`
    Value               string                `json:"Value"`
    CreateIndex         int                   `json:"CreateIndex"`
    ModifyIndex         int                   `json:"ModifyIndex"`
}

type DatePut struct {
	KV                   KVPut                `json:"KV"`
}

type DateRes struct {
	KV                   KVRes                `json:"KV"`
}

type Errors struct {
	OpIndex               int                 `json:"OpIndex"`
	What                  string              `json:"What"`
}

type Response struct {
	Results              []DateRes            `json:"Results"`
	Errors               []Errors             `json:"Errors"`
}

      