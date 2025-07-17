package tasks

import (
	"bytes"
	"context"
	"dormitory_management/internal/common"
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/templates"
	"dormitory_management/pkg/logger"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type EmailTask struct {
	config *config.Config
	logger logger.Logger
	repos  repository.Repositories
}

func NewEmailTask(logger logger.Logger, config *config.Config, repos repository.Repositories) *EmailTask {
	return &EmailTask{
		config: config,
		logger: logger,
		repos:  repos,
	}
}

func (et *EmailTask) SendOTP(ctx context.Context, t *asynq.Task) error {
	var payload entity.OtpCode
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		et.logger.Error("could not unmarshal OTP payload", zap.Error(err))
		return err
	}

	to := []string{payload.Identifier}
	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s", payload.OtpCode)

	msg := []byte(
		"To: " + payload.Identifier + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			body,
	)

	if err := common.SendMail(to, msg, et.config.Email); err != nil {
		et.logger.Error("failed to send OTP email", zap.Error(err))
		return err
	}

	et.logger.Info("OTP email sent", zap.String("email", payload.Identifier))
	return nil
}

func (et *EmailTask) SendVerifyAccount(ctx context.Context, t *asynq.Task) error {
	var payload entity.SendMailVerifyAccount
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		et.logger.Error("could not unmarshal payload:", zap.Error(err))
		return err
	}

	escapeToken := url.QueryEscape(payload.Token)
	escapeEmail := url.QueryEscape(payload.Email)
	verifyLink := fmt.Sprintf("%s/auth/verify-account/%s/%s", et.config.Client.ClientDomain, escapeToken, escapeEmail)

	to := []string{payload.Email}
	subject := "Verify Account"
	body := fmt.Sprintf(`
		<html>
			<body>
				<p>Hello,</p>
				<p>Click the link below to verify your account:</p>
				<p><a href="%s">Verify Account</a></p>
			</body>
		</html>
	`, verifyLink)

	msg := []byte(
		"To: " + payload.Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
			body,
	)

	if err := common.SendMail(to, msg, et.config.Email); err != nil {
		et.logger.Error("failed to send verification email", zap.Error(err))
		return err
	}

	return nil
}

func (et *EmailTask) SendEmailBill(ctx context.Context, t *asynq.Task) error {
	var payload entity.Bill
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		et.logger.Error("could not unmarshal bill payload", zap.Error(err))
		return err
	}

	bill := entity.Bill{ID: payload.ID}
	if err := et.repos.Bill().GetBill(ctx, &bill); err != nil {
		et.logger.Error("failed to get bill by ID", zap.Error(err))
		return err
	}

	user := &entity.User{ID: payload.UserID}

	if err := et.repos.User().GetByID(ctx, user); err != nil {
		et.logger.Error("failed to get user by ID", zap.Error(err))
		return err
	}

	to := []string{user.Email}
	subject := "Hóa đơn thanh toán"
	templateEmail, err := template.New("bill").Parse(templates.BillEmailTemplate)
	if err != nil {
		et.logger.Error("failed to parse bill template", zap.Error(err))
		return err
	}
	var body bytes.Buffer
	header := "MIME-Version: 1.0;\r\n"
	header += "Content-Type: text/html; charset=\"UTF-8\";\r\n"
	header += "Content-Disposition: inline;\r\n"
	header += "Content-Transfer-Encoding: 8bit;\r\n"
	header += "X-Mailer: Go\r\n"

	body.Write([]byte(fmt.Sprintf("Subject: %s\r\n%s\r\n", subject, header)))
	emailData := entity.BillEmailData{
		Bill: bill,
		Company: entity.CompanyInfo{
			Name:    "Trường cao đẳng Công nghệ Bách khoa Hà Nội",
			Address: "Số 18-20 Nhân Mỹ - Mỹ Đình 1 - Quận Nam Từ Liêm - TP. Hà Nội",
			Phone:   "0961224529",
			Email:   "truyenthong@bachkhoahanoi.edu.vn",
			Website: "https://bachkhoahanoi.edu.vn",
		},
	}
	et.logger.Info("bill email data", zap.Any("emailData", emailData))
	templateEmail.Execute(&body, emailData)

	msg := body.Bytes()

	if err := common.SendMail(to, msg, et.config.Email); err != nil {
		et.logger.Error("failed to send bill email", zap.Error(err))
		return err
	}

	et.logger.Info("bill email sent", zap.String("email", user.Email))
	return nil
}
