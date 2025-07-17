package templates

const BillEmailTemplate = `<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hóa đơn #{{.Bill.ID}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f8f9fa;
        }
        
        .container {
            max-width: 800px;
            margin: 20px auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px 30px;
            text-align: center;
            position: relative;
        }
        
        .header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="white" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="white" opacity="0.1"/><circle cx="50" cy="10" r="1" fill="white" opacity="0.1"/><circle cx="10" cy="60" r="1" fill="white" opacity="0.1"/><circle cx="90" cy="40" r="1" fill="white" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
        }
        
        .header-content {
            position: relative;
            z-index: 1;
        }
        
        .header h1 {
            font-size: 2.8em;
            margin-bottom: 10px;
            font-weight: 300;
            text-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        
        .header p {
            font-size: 1.1em;
            opacity: 0.9;
        }
        
        .content {
            padding: 40px;
        }
        
        .bill-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 40px;
            flex-wrap: wrap;
            gap: 30px;
        }
        
        .company-info, .customer-info {
            flex: 1;
            min-width: 280px;
        }
        
        .info-section h3 {
            color: #667eea;
            margin-bottom: 15px;
            font-size: 1.3em;
            border-bottom: 2px solid #667eea;
            padding-bottom: 8px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .info-item {
            margin-bottom: 10px;
            color: #555;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .info-item strong {
            color: #333;
        }
        
        .bill-details {
            background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
            padding: 30px;
            border-radius: 12px;
            margin-bottom: 40px;
            border-left: 5px solid #667eea;
            position: relative;
        }
        
        .bill-details::before {
            content: '';
            position: absolute;
            top: 0;
            right: 0;
            width: 60px;
            height: 60px;
            background: #667eea;
            opacity: 0.1;
            border-radius: 0 12px 0 60px;
        }
        
        .bill-details-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 30px;
        }
        
        .detail-item {
            text-align: center;
        }
        
        .detail-item label {
            display: block;
            font-weight: 600;
            color: #667eea;
            margin-bottom: 8px;
            text-transform: uppercase;
            font-size: 0.9em;
            letter-spacing: 0.5px;
        }
        
        .detail-item .value {
            font-size: 1.2em;
            color: #333;
            font-weight: 600;
        }
        
        .amount-section {
            background: linear-gradient(135deg, #e8f5e8 0%, #f0f8f0 100%);
            padding: 30px;
            border-radius: 12px;
            margin-bottom: 30px;
            border-left: 5px solid #28a745;
            text-align: center;
        }
        
        .amount-section h3 {
            color: #28a745;
            margin-bottom: 15px;
            font-size: 1.4em;
        }
        
        .amount {
            font-size: 2.5em;
            font-weight: 700;
            color: #28a745;
            margin-bottom: 10px;
        }
        
        .amount-currency {
            font-size: 1.2em;
            color: #666;
            font-weight: 400;
        }
        
        .description-section {
            background: #fff3cd;
            padding: 25px;
            border-radius: 12px;
            border-left: 5px solid #ffc107;
            margin-bottom: 30px;
        }
        
        .description-section h4 {
            color: #856404;
            margin-bottom: 15px;
            font-size: 1.2em;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .description-section p {
            color: #856404;
            line-height: 1.6;
        }
        
        .status-badge {
            display: inline-block;
            padding: 8px 16px;
            border-radius: 20px;
            font-size: 0.9em;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        .status-paid {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        
        .status-pending {
            background: #fff3cd;
            color: #856404;
            border: 1px solid #ffeaa7;
        }
        
        .status-overdue {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        
        .payment-info {
            background: #e8f4f8;
            padding: 25px;
            border-radius: 12px;
            border-left: 5px solid #17a2b8;
            margin-bottom: 30px;
        }
        
        .payment-info h4 {
            color: #17a2b8;
            margin-bottom: 15px;
            font-size: 1.2em;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .payment-info p {
            color: #0c5460;
            margin-bottom: 8px;
        }
        
        .footer {
            background: #2c3e50;
            color: white;
            text-align: center;
            padding: 30px;
            position: relative;
        }
        
        .footer::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: linear-gradient(45deg, transparent 30%, rgba(255,255,255,0.1) 50%, transparent 70%);
        }
        
        .footer-content {
            position: relative;
            z-index: 1;
        }
        
        .footer h4 {
            margin-bottom: 15px;
            font-size: 1.3em;
        }
        
        .footer p {
            margin-bottom: 8px;
            opacity: 0.9;
        }
        
        .contact-info {
            margin-top: 20px;
            padding-top: 20px;
            border-top: 1px solid rgba(255,255,255,0.2);
        }
        
        @media (max-width: 768px) {
            .container {
                margin: 10px;
                border-radius: 8px;
            }
            
            .header {
                padding: 30px 20px;
            }
            
            .header h1 {
                font-size: 2.2em;
            }
            
            .content {
                padding: 30px 20px;
            }
            
            .bill-header {
                flex-direction: column;
                gap: 20px;
            }
            
            .bill-details-grid {
                grid-template-columns: 1fr;
                gap: 20px;
            }
            
            .amount {
                font-size: 2em;
            }
        }
        
        @media (max-width: 480px) {
            .header h1 {
                font-size: 1.8em;
            }
            
            .amount {
                font-size: 1.8em;
            }
            
            .bill-details,
            .amount-section,
            .description-section,
            .payment-info {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="header-content">
                <h1>💳 HÓA ĐƠN</h1>
                <p>Cảm ơn bạn đã sử dụng dịch vụ của chúng tôi</p>
            </div>
        </div>
        
        <div class="content">
            <div class="bill-header">
                <div class="company-info info-section">
                    <h3>🏢 Thông tin trường</h3>
                    <div class="info-item">
                        <strong>{{.Company.Name}}</strong>
                    </div>
                    <div class="info-item">
                        📍 {{.Company.Address}}
                    </div>
                    <div class="info-item">
                        📞 {{.Company.Phone}}
                    </div>
                    <div class="info-item">
                        ✉️ {{.Company.Email}}
                    </div>
                    <div class="info-item">
                        🌐 {{.Company.Website}}
                    </div>
                </div>
                
                <div class="customer-info info-section">
                    <h3>👤 Thông tin sinh viên</h3>
                    <div class="info-item">
                        <strong>{{.Bill.User.FullName}}</strong>
                    </div>
                    <div class="info-item">
                        ✉️ {{.Bill.User.Email}}
                    </div>
                    {{if .Bill.User.Phone}}
                    <div class="info-item">
                        📞 {{.Bill.User.Phone}}
                    </div>
                    {{end}}
                    {{if .Bill.User.Address}}
                    <div class="info-item">
                        📍 {{.Bill.User.Address}}
                    </div>
                    {{end}}
                </div>
            </div>
            
            <div class="bill-details">
                <div class="bill-details-grid">
                    <div class="detail-item">
                        <label>Mã hóa đơn</label>
                        <div class="value">#{{.Bill.ID}}</div>
                    </div>
                    <div class="detail-item">
                        <label>Ngày tạo</label>
                        <div class="value">{{.Bill.CreatedAt.Format "02/01/2006"}}</div>
                    </div>
                    <div class="detail-item">
                        <label>Ngày đến hạn</label>
                        <div class="value">{{.Bill.DueDate.Format "02/01/2006"}}</div>
                    </div>
                    <div class="detail-item">
                        <label>Trạng thái</label>
                        <div class="value">
                            {{if eq .Bill.Status "PAID"}}
                                <span class="status-badge status-paid">✅ Đã thanh toán</span>
                            {{else if eq .Bill.Status "PENDING"}}
                                <span class="status-badge status-pending">⏳ Chờ thanh toán</span>                            
                            {{else if eq .Bill.Status "OVERDUE"}}
                                <span class="status-badge status-overdue">⚠️ Quá hạn</span>
                            {{else}}
                                <span class="status-badge status-pending">{{.Bill.Status}}</span>
                            {{end}}
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="amount-section">
                <h3>💰 Tổng số tiền</h3>
                <div class="amount">{{printf "%.0f" .Bill.Amount}}</div>
                <div class="amount-currency">VNĐ</div>
            </div>
            
            {{if .Bill.Description}}
            <div class="description-section">
                <h4>📝 Mô tả dịch vụ</h4>
                <p>{{.Bill.Description}}</p>
            </div>
            {{end}}
            
            <div class="payment-info">
                <h4>💳 Thông tin thanh toán</h4>
                <p><strong>Phương thức:</strong> {{.Bill.Payment.PaymentMethod}}</p>
                <p><strong>Trạng thái thanh toán:</strong> 
                    {{if eq .Bill.Payment.Status "PAID"}}
                        <span style="color: #28a745;">✅ Hoàn thành</span>
                    {{else if eq .Bill.Payment.Status "PENDING"}}
                        <span style="color: #ffc107;">⏳ Đang xử lý</span>
                    {{else if eq .Bill.Payment.Status "FAILED"}}
                        <span style="color: #dc3545;">❌ Thất bại</span>
                    {{else}}
                        <span>{{.Bill.Payment.Status}}</span>
                    {{end}}
                </p>
            </div>
        </div>
        
        <div class="footer">
            <div class="footer-content">
                <h4>🎉 Cảm ơn bạn đã tin tưởng!</h4>
                <p>Chúng tôi rất trân trọng sự tin tưởng và đồng hành của bạn.</p>
                <p>Nếu có bất kỳ thắc mắc nào, vui lòng liên hệ với chúng tôi.</p>
                
                <div class="contact-info">
                    <p>📧 Email: {{.Company.Email}}</p>
                    <p>📞 Hotline: {{.Company.Phone}}</p>
                    <p>🌐 Website: {{.Company.Website}}</p>
                </div>
            </div>
        </div>
    </div>
</body>
</html>
`
