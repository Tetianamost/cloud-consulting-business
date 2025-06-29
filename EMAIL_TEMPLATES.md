# Email Templates for Netlify Dashboard

## 📧 Contact Form Template

**Use for form:** `contact-form`
**Subject:** `New Contact Form Submission - Cloud Partner Pro`
**Email:** `info@cloudpartner.pro`

### HTML Template
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Contact Form Submission</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f8fafc;
            color: #334155;
            line-height: 1.6;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header {
            background: linear-gradient(135deg, #1e40af 0%, #3b82f6 100%);
            color: white;
            padding: 32px 24px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
            font-weight: 600;
        }
        .header p {
            margin: 8px 0 0 0;
            opacity: 0.9;
            font-size: 16px;
        }
        .content {
            padding: 32px 24px;
        }
        .section {
            margin-bottom: 32px;
        }
        .section-title {
            font-size: 18px;
            font-weight: 600;
            color: #1e40af;
            margin-bottom: 16px;
            padding-bottom: 8px;
            border-bottom: 2px solid #e2e8f0;
        }
        .info-grid {
            display: table;
            width: 100%;
            border-collapse: collapse;
        }
        .info-row {
            display: table-row;
        }
        .info-label {
            display: table-cell;
            font-weight: 600;
            color: #475569;
            padding: 12px 16px 12px 0;
            vertical-align: top;
            width: 140px;
            border-bottom: 1px solid #f1f5f9;
        }
        .info-value {
            display: table-cell;
            padding: 12px 0;
            color: #1e293b;
            border-bottom: 1px solid #f1f5f9;
        }
        .message-box {
            background-color: #f8fafc;
            border: 1px solid #e2e8f0;
            border-radius: 8px;
            padding: 20px;
            margin-top: 16px;
        }
        .services-list {
            background-color: #f0f9ff;
            border-left: 4px solid #3b82f6;
            padding: 16px 20px;
            margin-top: 16px;
        }
        .footer {
            background-color: #f8fafc;
            padding: 24px;
            text-align: center;
            border-top: 1px solid #e2e8f0;
        }
        .footer p {
            margin: 0;
            color: #64748b;
            font-size: 14px;
        }
        .priority-badge {
            display: inline-block;
            background-color: #dc2626;
            color: white;
            padding: 4px 12px;
            border-radius: 16px;
            font-size: 12px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        @media only screen and (max-width: 600px) {
            .container {
                width: 100% !important;
            }
            .header {
                padding: 24px 16px;
            }
            .content {
                padding: 24px 16px;
            }
            .info-label, .info-value {
                display: block;
                width: 100%;
            }
            .info-label {
                font-weight: 600;
                margin-bottom: 4px;
            }
            .info-value {
                margin-bottom: 16px;
                padding-left: 0;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📬 New Contact Inquiry</h1>
            <p>Cloud Partner Pro - Professional Cloud Consulting</p>
        </div>
        
        <div class="content">
            <div class="section">
                <div class="section-title">
                    Contact Information
                    <span class="priority-badge">New Lead</span>
                </div>
                <div class="info-grid">
                    <div class="info-row">
                        <div class="info-label">Full Name:</div>
                        <div class="info-value">{{ name }}</div>
                    </div>
                    <div class="info-row">
                        <div class="info-label">Email:</div>
                        <div class="info-value">{{ email }}</div>
                    </div>
                    <div class="info-row">
                        <div class="info-label">Company:</div>
                        <div class="info-value">{{ company }}</div>
                    </div>
                    <div class="info-row">
                        <div class="info-label">Phone:</div>
                        <div class="info-value">{{ phone }}</div>
                    </div>
                </div>
            </div>

            <div class="section">
                <div class="section-title">Services of Interest</div>
                <div class="services-list">
                    <strong>{{ services }}</strong>
                </div>
            </div>

            <div class="section">
                <div class="section-title">Message Details</div>
                <div class="message-box">
                    {{ message }}
                </div>
            </div>
        </div>
        
        <div class="footer">
            <p><strong>Cloud Partner Pro</strong> | Professional Cloud Consulting Services</p>
            <p>This inquiry was submitted through the contact form on cloudpartner.pro</p>
            <p>Please respond within 24 hours for optimal customer experience.</p>
        </div>
    </div>
</body>
</html>
```

### Text Template
```
NEW CONTACT FORM SUBMISSION
Cloud Partner Pro - Professional Cloud Consulting

========================================
CONTACT INFORMATION
========================================
Name: {{ name }}
Email: {{ email }}
Company: {{ company }}
Phone: {{ phone }}

========================================
SERVICES OF INTEREST
========================================
{{ services }}

========================================
MESSAGE
========================================
{{ message }}

========================================
This inquiry was submitted through the contact form on cloudpartner.pro
Please respond within 24 hours for optimal customer experience.

Cloud Partner Pro | Professional Cloud Consulting Services
```

---

## 💰 Quote Request Form Template

**Use for form:** `quote-request-form`
**Subject:** `New Quote Request - Cloud Partner Pro`
**Email:** `info@cloudpartner.pro`

### HTML Template
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Quote Request</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f8fafc;
            color: #334155;
            line-height: 1.6;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header {
            background: linear-gradient(135deg, #059669 0%, #10b981 100%);
            color: white;
            padding: 32px 24px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
            font-weight: 600;
        }
        .header p {
            margin: 8px 0 0 0;
            opacity: 0.9;
            font-size: 16px;
        }
        .content {
            padding: 32px 24px;
        }
        .section {
            margin-bottom: 32px;
        }
        .section-title {
            font-size: 18px;
            font-weight: 600;
            color: #059669;
            margin-bottom: 16px;
            padding-bottom: 8px;
            border-bottom: 2px solid #e2e8f0;
        }
        .info-grid {
            display: table;
            width: 100%;
            border-collapse: collapse;
        }
        .info-row {
            display: table-row;
        }
        .info-label {
            display: table-cell;
            font-weight: 600;
            color: #475569;
            padding: 12px 16px 12px 0;
            vertical-align: top;
            width: 160px;
            border-bottom: 1px solid #f1f5f9;
        }
        .info-value {
            display: table-cell;
            padding: 12px 0;
            color: #1e293b;
            border-bottom: 1px solid #f1f5f9;
        }
        .pricing-summary {
            background: linear-gradient(135deg, #f0fdf4 0%, #ecfdf5 100%);
            border: 2px solid #10b981;
            border-radius: 12px;
            padding: 24px;
            margin-top: 20px;
        }
        .pricing-title {
            font-size: 20px;
            font-weight: 700;
            color: #059669;
            margin-bottom: 16px;
            text-align: center;
        }
        .pricing-breakdown {
            margin-bottom: 16px;
        }
        .pricing-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 8px 0;
            border-bottom: 1px solid #d1fae5;
        }
        .pricing-row:last-child {
            border-bottom: none;
            font-size: 18px;
            font-weight: 700;
            color: #059669;
            margin-top: 12px;
            padding-top: 16px;
            border-top: 2px solid #10b981;
        }
        .requirements-box {
            background-color: #f8fafc;
            border: 1px solid #e2e8f0;
            border-radius: 8px;
            padding: 20px;
            margin-top: 16px;
        }
        .footer {
            background-color: #f8fafc;
            padding: 24px;
            text-align: center;
            border-top: 1px solid #e2e8f0;
        }
        .footer p {
            margin: 0;
            color: #64748b;
            font-size: 14px;
        }
        .priority-badge {
            display: inline-block;
            background-color: #dc2626;
            color: white;
            padding: 4px 12px;
            border-radius: 16px;
            font-size: 12px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .quote-badge {
            display: inline-block;
            background-color: #059669;
            color: white;
            padding: 4px 12px;
            border-radius: 16px;
            font-size: 12px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        @media only screen and (max-width: 600px) {
            .container {
                width: 100% !important;
            }
            .header {
                padding: 24px 16px;
            }
            .content {
                padding: 24px 16px;
            }
            .info-label, .info-value {
                display: block;
                width: 100%;
            }
            .info-label {
                font-weight: 600;
                margin-bottom: 4px;
            }
            .info-value {
                margin-bottom: 16px;
                padding-left: 0;
            }
            .pricing-row {
                flex-direction: column;
                align-items: flex-start;
                gap: 4px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>💰 New Quote Request</h1>
            <p>Cloud Partner Pro - Professional Cloud Consulting</p>
        </div>
        
        <div class="content">
            <div class="section">
                <div class="section-title">
                    Client Information
                    <span class="priority-badge">Quote Request</span>
                </div>
                <div class="info-grid">
                    <div class="info-row">
                        <div class="info-label">Full Name:</div>
                        <div class="info-value">{{ name }}</div>
                    </div>
                    <div class="info-row">
                        <div class="info-label">Email:</div>
                        <div class="info-value">{{ email }}</div>
                    </div>
                    <div class="info-row">
                        <div class="info-label">Company:</div>
                        <div class="info-value">{{ company }}</div>
                    </div>
                    <div class="info-row">
                        <div class="info-label">Phone:</div>
                        <div class="info-value">{{ phone }}</div>
                    </div>
                </div>
            </div>

            <div class="section">
                <div class="section-title">
                    Service Requirements
                    <span class="quote-badge">Details</span>
                </div>
                <div class="info-grid">
                    <div class="info-row">
                        <div class="info-label">Service Type:</div>
                        <div class="info-value">{{ serviceType }}</div>
                    </div>
                    <div class="info-row">
                        <div class="info-label">Complexity Level:</div>
                        <div class="info-value">{{ complexity }}</div>
                    </div>
                    <div class="info-row">
                        <div class="info-label">Count:</div>
                        <div class="info-value">{{ count }}</div>
                    </div>
                </div>
            </div>

            <div class="section">
                <div class="section-title">Pricing Estimate</div>
                <div class="pricing-summary">
                    <div class="pricing-title">💲 Cost Breakdown</div>
                    <div class="pricing-breakdown">
                        <div class="pricing-row">
                            <span>Base Service Fee:</span>
                            <span><strong>${{ basePrice }}</strong></span>
                        </div>
                        <div class="pricing-row">
                            <span>Variable Cost:</span>
                            <span><strong>${{ variablePrice }}</strong></span>
                        </div>
                        <div class="pricing-row">
                            <span>Complexity Multiplier:</span>
                            <span><strong>{{ complexityMultiplier }}x</strong></span>
                        </div>
                        <div class="pricing-row">
                            <span>📊 TOTAL ESTIMATE:</span>
                            <span><strong>${{ totalEstimate }}</strong></span>
                        </div>
                    </div>
                </div>
            </div>

            <div class="section">
                <div class="section-title">Additional Requirements</div>
                <div class="requirements-box">
                    {{ requirements }}
                </div>
            </div>
        </div>
        
        <div class="footer">
            <p><strong>Cloud Partner Pro</strong> | Professional Cloud Consulting Services</p>
            <p>This quote request was submitted through the pricing calculator on cloudpartner.pro</p>
            <p><strong>Action Required:</strong> Please prepare and send detailed quote within 24-48 hours.</p>
        </div>
    </div>
</body>
</html>
```

### Text Template
```
NEW QUOTE REQUEST
Cloud Partner Pro - Professional Cloud Consulting

========================================
CLIENT INFORMATION
========================================
Name: {{ name }}
Email: {{ email }}
Company: {{ company }}
Phone: {{ phone }}

========================================
SERVICE REQUIREMENTS
========================================
Service Type: {{ serviceType }}
Complexity Level: {{ complexity }}
Count: {{ count }}

========================================
PRICING ESTIMATE
========================================
Base Service Fee: ${{ basePrice }}
Variable Cost: ${{ variablePrice }}
Complexity Multiplier: {{ complexityMultiplier }}x

TOTAL ESTIMATE: ${{ totalEstimate }}

========================================
ADDITIONAL REQUIREMENTS
========================================
{{ requirements }}

========================================
This quote request was submitted through the pricing calculator on cloudpartner.pro
ACTION REQUIRED: Please prepare and send detailed quote within 24-48 hours.

Cloud Partner Pro | Professional Cloud Consulting Services
```

## 📋 Configuration Instructions

1. Deploy your site to Netlify
2. Go to Netlify Dashboard → Your Site → Forms
3. Click on each form name (contact-form, quote-request-form)
4. Go to Settings & usage → Form notifications
5. Add email notification with the templates above
6. Test with form submissions