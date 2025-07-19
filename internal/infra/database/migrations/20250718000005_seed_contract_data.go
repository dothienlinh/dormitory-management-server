package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSeedContractData, downSeedContractData)
}

func upSeedContractData(ctx context.Context, tx *sql.Tx) error {
	// Insert sample contract terms for existing contracts
	seedContractTerms := `
	INSERT INTO contract_terms (contract_id, title, content, "order") 
	SELECT 
		c.id,
		'Thời hạn hợp đồng' as title,
		'Hợp đồng có hiệu lực từ ngày ' || to_char(c.start_date, 'DD/MM/YYYY') || ' đến ngày ' || to_char(c.end_date, 'DD/MM/YYYY') || '.' as content,
		1 as "order"
	FROM contracts c
	WHERE NOT EXISTS (
		SELECT 1 FROM contract_terms ct WHERE ct.contract_id = c.id
	);

	INSERT INTO contract_terms (contract_id, title, content, "order") 
	SELECT 
		c.id,
		'Phí ký túc xá' as title,
		'Phí hàng tháng: ' || COALESCE(c.price, 0) || ' VNĐ/tháng. Thanh toán trước ngày 15 hàng tháng.' as content,
		2 as "order"
	FROM contracts c
	WHERE NOT EXISTS (
		SELECT 1 FROM contract_terms ct WHERE ct.contract_id = c.id AND ct.title = 'Phí ký túc xá'
	);

	INSERT INTO contract_terms (contract_id, title, content, "order") 
	SELECT 
		c.id,
		'Quy định ở KTX' as title,
		'Sinh viên phải tuân thủ nội quy KTX, giữ gìn vệ sinh chung và bảo quản tài sản công.' as content,
		3 as "order"
	FROM contracts c
	WHERE NOT EXISTS (
		SELECT 1 FROM contract_terms ct WHERE ct.contract_id = c.id AND ct.title = 'Quy định ở KTX'
	);

	INSERT INTO contract_terms (contract_id, title, content, "order") 
	SELECT 
		c.id,
		'Chấm dứt hợp đồng' as title,
		'Sinh viên muốn chấm dứt hợp đồng trước thời hạn phải thông báo trước ít nhất 30 ngày.' as content,
		4 as "order"
	FROM contracts c
	WHERE NOT EXISTS (
		SELECT 1 FROM contract_terms ct WHERE ct.contract_id = c.id AND ct.title = 'Chấm dứt hợp đồng'
	);
	`

	if _, err := tx.ExecContext(ctx, seedContractTerms); err != nil {
		return err
	}

	// Insert sample payment histories for existing contracts
	seedPaymentHistories := `
	INSERT INTO payment_histories (contract_id, period, amount, status, due_date, payment_date, method) 
	SELECT 
		c.id,
		'Tháng ' || EXTRACT(MONTH FROM c.start_date) || '/' || EXTRACT(YEAR FROM c.start_date) as period,
		COALESCE(c.price, 850000) as amount,
		'paid' as status,
		c.start_date + INTERVAL '15 days' as due_date,
		c.start_date + INTERVAL '10 days' as payment_date,
		'bank_transfer' as method
	FROM contracts c
	WHERE c.status = 'active'
	AND NOT EXISTS (
		SELECT 1 FROM payment_histories ph WHERE ph.contract_id = c.id
	);

	INSERT INTO payment_histories (contract_id, period, amount, status, due_date) 
	SELECT 
		c.id,
		'Tháng ' || EXTRACT(MONTH FROM c.start_date + INTERVAL '1 month') || '/' || EXTRACT(YEAR FROM c.start_date + INTERVAL '1 month') as period,
		COALESCE(c.price, 850000) as amount,
		'pending' as status,
		c.start_date + INTERVAL '1 month 15 days' as due_date
	FROM contracts c
	WHERE c.status = 'active'
	AND c.start_date + INTERVAL '1 month' < NOW()
	AND NOT EXISTS (
		SELECT 1 FROM payment_histories ph 
		WHERE ph.contract_id = c.id 
		AND ph.period = 'Tháng ' || EXTRACT(MONTH FROM c.start_date + INTERVAL '1 month') || '/' || EXTRACT(YEAR FROM c.start_date + INTERVAL '1 month')
	);
	`

	if _, err := tx.ExecContext(ctx, seedPaymentHistories); err != nil {
		return err
	}

	return nil
}

func downSeedContractData(ctx context.Context, tx *sql.Tx) error {
	// Remove seeded data
	query := `
	DELETE FROM contract_terms 
	WHERE title IN ('Thời hạn hợp đồng', 'Phí ký túc xá', 'Quy định ở KTX', 'Chấm dứt hợp đồng');
	
	DELETE FROM payment_histories 
	WHERE method = 'bank_transfer' OR status = 'pending';
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
