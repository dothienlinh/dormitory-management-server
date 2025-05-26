package tasks

type TaskType string

const (
	TypeSendCodeEmail TaskType = "email:send_code_email"
	TypeCallCodeSMS   TaskType = "stringee:call_code_sms"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)
