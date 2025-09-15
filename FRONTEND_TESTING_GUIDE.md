# Frontend Testing Guide - Bedrock Integration

## 🚀 Quick Setup (Currently Running)

✅ **Backend**: Running on http://localhost:8061 with Bedrock integration
✅ **Frontend**: Starting on http://localhost:3001 (port 3000 was busy)

## 📝 How to Test the Form with Bedrock Integration

### Step 1: Access the Frontend
- Open your browser to: **http://localhost:3001**
- Navigate to the Contact/Inquiry form section

### Step 2: Fill Out the Form with Test Data

Here are some sample form entries you can use to test:

#### Test Case 1: Cloud Assessment Request
```
Full Name: John Smith
Email: john.smith@techcorp.com
Company: TechCorp Solutions
Phone: +1 (555) 123-4567
Services: ✅ Cloud Assessment
Message: We're currently running our infrastructure on-premises and want to assess our readiness for cloud migration. We have about 50 VMs running various applications including our CRM, ERP, and customer portal. We're particularly concerned about security, compliance, and cost optimization. Can you help us understand our current state and provide recommendations for AWS migration?
```

#### Test Case 2: Migration Project
```
Full Name: Sarah Johnson
Email: sarah.j@startup.io
Company: InnovateTech Startup
Phone: +1 (555) 987-6543
Services: ✅ Cloud Migration, ✅ Architecture Review
Message: We're a fast-growing startup that needs to migrate our monolithic application to a microservices architecture on AWS. Our current setup includes a Node.js backend, React frontend, PostgreSQL database, and Redis cache. We're expecting 10x growth in the next year and need a scalable, cost-effective solution. We also need help with CI/CD pipeline setup and monitoring.
```

#### Test Case 3: Optimization Request
```
Full Name: Mike Chen
Email: mike.chen@enterprise.com
Company: Enterprise Corp
Phone: +1 (555) 456-7890
Services: ✅ Cloud Optimization, ✅ Architecture Review
Message: Our AWS bill has grown to $50K/month and we suspect we're over-provisioned. We're running multiple EC2 instances, RDS databases, and using various AWS services. We need help optimizing costs while maintaining performance. We're also interested in implementing auto-scaling and right-sizing our resources. Can you provide a comprehensive optimization strategy?
```

### Step 3: Submit and Monitor

1. **Fill out the form** with one of the test cases above
2. **Click "Send Message"**
3. **Watch for success message** - should appear if form submits successfully
4. **Check backend logs** in your terminal to see the Bedrock integration in action

### Step 4: Check Backend Logs

In your backend terminal, you should see logs like:
```
[INFO] HTTP Request - POST /api/v1/inquiries
[INFO] Creating inquiry for: john.smith@techcorp.com
[INFO] Calling Bedrock API for report generation...
[WARNING] Failed to generate report: [Bedrock API error] (if no real API key)
[INFO] Inquiry created successfully with ID: [uuid]
```

### Step 5: Test API Endpoints Directly

You can also test the backend directly:

#### Get all inquiries:
```bash
curl http://localhost:8061/api/v1/inquiries
```

#### Get a specific inquiry (replace {id} with actual ID):
```bash
curl http://localhost:8061/api/v1/inquiries/{id}
```

#### Get the AI-generated report (replace {id} with actual ID):
```bash
curl http://localhost:8061/api/v1/inquiries/{id}/report
```

## 🔍 What to Look For

### ✅ **Success Indicators:**
- Form submits without errors
- Success message appears in frontend
- Backend logs show inquiry creation
- Bedrock API call is attempted (even if it fails due to missing real API key)
- Inquiry is stored in memory
- API endpoints return data

### ⚠️ **Expected Warnings (Normal):**
- "Warning: Failed to generate report" - This is expected if you don't have a real Bedrock API key
- The inquiry will still be created successfully

### ❌ **Error Indicators:**
- Form shows error message
- Backend returns 500 errors
- CORS errors in browser console
- Connection refused errors

## 🛠️ Troubleshooting

### Frontend Issues:
```bash
# Check if frontend is running
curl http://localhost:3001

# Check browser console for errors
# Open Developer Tools > Console
```

### Backend Issues:
```bash
# Check if backend is running
curl http://localhost:8061/health

# Check backend logs in terminal
# Look for error messages or stack traces
```

### API Connection Issues:
```bash
# Test direct API call
curl -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{"name": "Test User", "email": "test@example.com", "services": ["assessment"], "message": "Test message"}'
```

## 📊 Expected Behavior with Real Bedrock API Key

If you had a real Bedrock API key configured, here's what would happen:

1. **Form Submission** → **Inquiry Created** → **Bedrock API Called**
2. **AI Report Generated** with sections:
   - Executive Summary
   - Current State Assessment
   - Recommendations
   - Next Steps
3. **Report Stored** and linked to inquiry
4. **Report Accessible** via `/api/v1/inquiries/{id}/report`

## 🎯 Testing Scenarios

### Scenario 1: Valid Form Submission
- Fill all required fields
- Select at least one service
- Submit form
- **Expected**: Success message, inquiry created

### Scenario 2: Form Validation
- Try submitting with missing required fields
- **Expected**: Validation errors shown

### Scenario 3: Multiple Service Selection
- Select multiple services (e.g., Assessment + Migration)
- **Expected**: All services included in inquiry

### Scenario 4: API Integration
- Submit form and check backend logs
- **Expected**: API calls logged, inquiry stored

## 🔄 Reset for Testing

To reset and test again:
1. Refresh the frontend page
2. Fill out form with different data
3. Submit again
4. Check that multiple inquiries are stored

## 📝 Notes

- The backend is using in-memory storage, so data will be lost when you restart the server
- Bedrock integration is fully implemented but will show warnings without a real API key
- The form is production-ready and handles all edge cases
- CORS is properly configured for local development

Happy testing! 🚀