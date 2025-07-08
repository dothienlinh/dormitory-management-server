package tasks

type TaskType string

const (
	TaskSendCodeEmail           TaskType = "task:send_code_email"
	TaskSendEmailVerifyAccount  TaskType = "task:send_email_verify_account"
	TaskCreateLinkPaymentVietQR TaskType = "task:create_link_payment_vietqr"
	TaskSendEmailForgotPassword TaskType = "task:send_email_forgot_password"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)
