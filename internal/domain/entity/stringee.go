package entity

const TokenKeyStringee string = "token_key_stringee"

type TypeCallParty string

const (
	TypeCallPartyExternal TypeCallParty = "external"
)

type TypeAction string

const (
	TypeActionTalk TypeAction = "talk"
)

type CallParty struct {
	Type   TypeCallParty `json:"type"`
	Number string        `json:"number"`
	Alias  string        `json:"alias"`
}

type CallAction struct {
	Action TypeAction `json:"action"`
	Text   string     `json:"text"`
}

type OutboundCall struct {
	From      CallParty    `json:"from"`
	To        []CallParty  `json:"to"`
	AnswerURL string       `json:"answer_url"`
	Actions   []CallAction `json:"actions"`
}

type CallSMSCode struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifySMSCode struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required"`
}
