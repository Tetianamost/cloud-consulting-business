# Quick Start: Deploy to AWS App Runner with Custom Domain

## 🚀 One-Command Deployment

Follow these steps to deploy your Cloud Consulting Platform to AWS App Runner with your custom domain:

### Prerequisites ✅

1. **AWS CLI configured** with appropriate permissions
2. **Domain ready** (e.g., `cloudpartner.pro`)
3. **GitHub repository** pushed to `kiro-dev` branch

### Step 1: Deploy Everything 🎯

```bash
# 1. Deploy RDS database (takes 5-10 minutes)
./deploy-apprunner-rds.sh

# 2. Deploy App Runner service (takes 3-5 minutes)
./deploy-now.sh

# 3. Update environment variables
./update-env-vars.sh

# 4. Set up custom domain (requires DNS configuration)
./setup-custom-domain.sh

# 5. Check deployment status
./check-deployment-status.sh
```

### Step 2: Configure DNS 🌐

After running `setup-custom-domain.sh`, you'll get DNS records to add to your domain:

**For cloudpartner.pro:**
1. Go to your domain registrar (GoDaddy, Namecheap, etc.)
2. Add the CNAME records provided by the script
3. Wait for DNS propagation (up to 48 hours)

### Step 3: Verify Deployment ✅

```bash
# Check if everything is working
./check-deployment-status.sh

# Test your endpoints
curl https://[your-apprunner-url]/health
curl https://cloudpartner.pro/health  # After DNS propagation
```

## 🎯 What Gets Deployed

### Infrastructure:
- **RDS PostgreSQL** (db.t3.micro) - ~$13/month
- **App Runner Service** (0.25 vCPU, 0.5GB RAM) - ~$5/month + requests
- **SSL Certificate** (free via AWS Certificate Manager)
- **Custom Domain** (cloudpartner.pro)

### Features Enabled:
- ✅ AI-powered chat system (AWS Bedrock)
- ✅ Email notifications (AWS SES)
- ✅ Admin dashboard
- ✅ Contact forms
- ✅ Report generation
- ✅ Performance monitoring
- ✅ Auto-scaling
- ✅ HTTPS with custom domain

## 🔧 Configuration

### Environment Variables Set:
- Database connection to RDS
- AWS Bedrock AI integration
- AWS SES email service
- CORS for your domain
- Production security settings
- Performance optimizations

### Security Features:
- JWT authentication
- HTTPS encryption
- Database encryption at rest
- Secure environment variables
- CORS protection

## 📊 Monitoring

### Check Status:
```bash
# Overall status
./check-deployment-status.sh

# View logs
aws logs tail /aws/apprunner/cloud-consulting-prod --region us-east-1 --follow

# Service metrics
aws apprunner describe-service --service-arn [your-service-arn] --region us-east-1
```

### Key URLs:
- **App Runner URL**: `https://[random-id].us-east-1.awsapprunner.com`
- **Custom Domain**: `https://cloudpartner.pro`
- **Health Check**: `/health`
- **API Health**: `/api/health`
- **Admin Dashboard**: `/admin`

## 🚨 Troubleshooting

### Common Issues:

1. **"Service not found"**
   ```bash
   # Check if service exists
   aws apprunner list-services --region us-east-1
   ```

2. **"Database connection failed"**
   ```bash
   # Check RDS status
   ./check_rds_status.sh
   ```

3. **"Custom domain not working"**
   - Check DNS propagation: `dig cloudpartner.pro`
   - Verify CNAME records are added
   - Wait up to 48 hours for full propagation

4. **"Build failed"**
   ```bash
   # Check build logs
   aws logs tail /aws/apprunner/cloud-consulting-prod --region us-east-1
   ```

### Get Help:
```bash
# Service status
aws apprunner describe-service --service-arn [arn] --region us-east-1

# Recent operations
aws apprunner list-operations --service-arn [arn] --region us-east-1

# Custom domain status
aws apprunner describe-custom-domains --service-arn [arn] --region us-east-1
```

## 💰 Cost Estimate

### Monthly Costs:
- **RDS db.t3.micro**: ~$13
- **App Runner**: ~$5 base + $0.40/million requests
- **Data Transfer**: ~$1-5 depending on traffic
- **SSL Certificate**: Free
- **Route 53** (if used): $0.50/hosted zone

**Total**: ~$20-25/month for low-medium traffic

## 🔄 Updates and Maintenance

### Auto-Deployment:
- Push to `kiro-dev` branch triggers automatic deployment
- No downtime deployments
- Rollback capability

### Manual Updates:
```bash
# Update environment variables
./update-env-vars.sh

# Force new deployment
aws apprunner start-deployment --service-arn [arn] --region us-east-1
```

## 🎉 Success!

Once deployed, your Cloud Consulting Platform will be live at:
- **https://cloudpartner.pro** (your custom domain)
- **https://www.cloudpartner.pro** (www subdomain)

Features available:
- 🤖 AI consultant chat
- 📧 Contact forms with email notifications
- 📊 Admin dashboard
- 📱 Mobile-responsive design
- 🔒 Secure HTTPS
- ⚡ Auto-scaling
- 📈 Performance monitoring

Your professional cloud consulting platform is now live and ready for clients! 🚀