package tasks

type TaskType string

const (
	TypeSendCodeEmail          TaskType = "email:send_code_email"
	TypeSendEmailVerifyAccount TaskType = "email:send_email_verify_account"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)
