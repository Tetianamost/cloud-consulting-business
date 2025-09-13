# AWS Production Infrastructure Implementation Plan

## Task Overview

This implementation plan provides step-by-step tasks to deploy the Cloud Consulting Platform on AWS production infrastructure. Each task builds incrementally and includes specific commands, configurations, and validation steps.

## **Current AWS Setup Status:**

✅ **Already Configured:**

- AWS Account with credentials (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
- AWS Bedrock with Nova Lite model access (AWS_BEARER_TOKEN_BEDROCK)
- AWS SES with info@cloudpartner.pro sender email
- Region: us-east-1

🔧 **Need to Set Up for Production:**

- Amazon RDS PostgreSQL database
- Amazon ElastiCache Redis cluster
- Amazon EKS cluster for container orchestration
- Application Load Balancer with SSL
- Route 53 DNS (if you want to use cloudpartner.pro domain)
- Production monitoring and backups

## Implementation Tasks

- [ ] 1. Verify existing AWS configuration and enable additional services

  - Verify AWS CLI is configured with your existing credentials
  - Enable Amazon EKS service in us-east-1 region
  - Enable Amazon RDS service in us-east-1 region
  - Enable Amazon ElastiCache service in us-east-1 region
  - Verify AWS Bedrock and SES are already working (as per your .env)
  - Set up billing alerts for cost monitoring
  - _Requirements: 8.1, 8.2_

- [ ] 2. Configure domain and DNS infrastructure

  - Purchase or transfer domain to Route 53 (cloudpartner.pro)
  - Create hosted zone in Route 53
  - Configure DNS records for domain validation
  - Set up health checks for monitoring
  - _Requirements: 12.1, 12.2, 12.3_

- [ ] 3. Set up SSL certificates with AWS Certificate Manager

  - Request SSL certificate for cloudpartner.pro and \*.cloudpartner.pro
  - Validate domain ownership through DNS validation
  - Configure certificate for use with Application Load Balancer
  - Test certificate installation and renewal
  - _Requirements: 6.2, 12.4_

- [ ] 4. Create VPC and networking infrastructure

  - Create VPC with CIDR 10.0.0.0/16 across 3 availability zones
  - Set up public subnets (10.0.1.0/24, 10.0.2.0/24, 10.0.3.0/24)
  - Set up private subnets (10.0.11.0/24, 10.0.12.0/24, 10.0.13.0/24)
  - Set up database subnets (10.0.21.0/24, 10.0.22.0/24, 10.0.23.0/24)
  - Configure Internet Gateway and NAT Gateways
  - Set up route tables and security groups
  - _Requirements: 5.3, 8.4_

- [ ] 5. Deploy Amazon RDS PostgreSQL database

  - Create RDS subnet group using database subnets
  - Configure security group for database access (port 5432 from EKS only)
  - Deploy RDS PostgreSQL instance with Multi-AZ configuration
  - Configure automated backups with 7-day retention
  - Enable encryption at rest and performance insights
  - Test database connectivity and create initial database
  - _Requirements: 1.1, 1.2, 1.3, 1.5, 8.5_

- [ ] 6. Deploy Amazon ElastiCache Redis cluster

  - Create ElastiCache subnet group using database subnets
  - Configure security group for Redis access (port 6379 from EKS only)
  - Deploy Redis cluster with encryption in transit and at rest
  - Configure backup snapshots and monitoring
  - Test Redis connectivity and basic operations
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 8.5_

- [ ] 7. Configure AWS SES for production email delivery

  - Verify that info@cloudpartner.pro is already verified in AWS SES
  - Configure SPF, DKIM, and DMARC DNS records for cloudpartner.pro domain
  - Request production access to move SES out of sandbox mode (if not already done)
  - Set up SES configuration set for tracking delivery events
  - Configure SES webhooks for bounce and complaint handling
  - Test email delivery using your existing SES configuration
  - _Requirements: 3.1, 3.2, 3.3, 3.5, 3.6_

- [ ] 8. Verify and optimize AWS Bedrock configuration

  - Verify Amazon Nova Lite model (amazon.nova-lite-v1:0) access is working
  - Ensure IAM permissions are properly configured for production workloads
  - Test Bedrock API connectivity with your existing bearer token
  - Set up CloudWatch monitoring for Bedrock usage and costs
  - Configure appropriate rate limits and error handling for production
  - _Requirements: 4.1, 4.2, 4.3, 4.4_

- [ ] 9. Create and configure Amazon EKS cluster

  - Install eksctl and kubectl tools
  - Create EKS cluster configuration file with managed node groups
  - Deploy EKS cluster with 3 t3.medium nodes across multiple AZs
  - Configure cluster autoscaling and enable CloudWatch Container Insights
  - Install AWS Load Balancer Controller addon
  - Install and configure Fluent Bit for log forwarding to CloudWatch
  - _Requirements: 5.1, 5.2, 5.4, 5.6_

- [ ] 10. Set up AWS Secrets Manager for configuration

  - Create secrets for database connection string
  - Create secrets for Redis connection string
  - Create secrets for AWS credentials (SES, Bedrock)
  - Create secrets for JWT signing key
  - Configure IAM permissions for EKS to access secrets
  - Test secret retrieval from within EKS cluster
  - _Requirements: 8.3_

- [ ] 11. Deploy application to EKS cluster

  - Create Kubernetes namespace (cloud-consulting)
  - Deploy ConfigMap with environment-specific configuration
  - Deploy Secrets referencing AWS Secrets Manager
  - Deploy backend application using existing k8s/backend-deployment.yaml
  - Deploy frontend application using existing k8s/frontend-deployment.yaml
  - Configure horizontal pod autoscaling for both services
  - _Requirements: 5.5_

- [ ] 12. Configure Application Load Balancer and ingress

  - Create target groups for backend and frontend services
  - Deploy Application Load Balancer in public subnets
  - Configure SSL termination using ACM certificate
  - Set up routing rules (/ to frontend, /api/\* to backend)
  - Configure health checks using /health endpoint
  - Test load balancer functionality and SSL certificate
  - _Requirements: 6.1, 6.3, 6.4, 6.5_

- [ ] 13. Run database migrations and initialize data

  - Connect to RDS instance from EKS cluster
  - Run init.sql migration to create core tables
  - Run chat_migration.sql to create chat system tables
  - Run email_events_migration.sql to create email tracking tables
  - Verify all tables and indexes are created correctly
  - Test application database connectivity
  - _Requirements: 1.4_

- [ ] 14. Configure monitoring and alerting

  - Set up CloudWatch dashboards for application metrics
  - Create CloudWatch alarms for CPU, memory, and error rates
  - Configure SNS topics for alert notifications
  - Set up log aggregation and analysis in CloudWatch Logs
  - Configure AWS X-Ray for distributed tracing
  - Test monitoring and alerting functionality
  - _Requirements: 1.6, 2.5, 7.1, 7.2, 7.3, 7.4_

- [ ] 15. Implement backup and disaster recovery

  - Configure automated RDS backups with cross-region replication
  - Set up EBS snapshot policies for EKS node volumes
  - Create backup procedures for Kubernetes configurations
  - Document disaster recovery procedures and test restoration
  - Configure backup retention policies (30 days for production)
  - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5_

- [ ] 16. Set up CI/CD pipeline for automated deployment

  - Configure GitHub Actions workflow for building Docker images
  - Set up ECR repositories for storing container images
  - Create deployment pipeline for staging and production environments
  - Configure automated testing and security scanning
  - Set up blue-green deployment strategy for zero-downtime updates
  - _Requirements: 11.3_

- [ ] 17. Configure cost optimization and monitoring

  - Set up AWS Cost Explorer and billing alerts
  - Configure AWS Budgets for cost control
  - Implement cluster autoscaling and horizontal pod autoscaling
  - Set up spot instances for non-critical workloads
  - Configure resource requests and limits for optimal cost
  - _Requirements: 10.1, 10.2, 10.3, 10.5_

- [ ] 18. Perform security hardening and compliance

  - Enable VPC Flow Logs for network monitoring
  - Configure AWS WAF rules for application protection
  - Set up AWS Config for compliance monitoring
  - Enable CloudTrail for audit logging
  - Perform security assessment and penetration testing
  - _Requirements: 8.1, 8.2, 8.4_

- [ ] 19. Set up staging environment

  - Create separate VPC or use separate AWS account for staging
  - Deploy smaller instance sizes while maintaining same architecture
  - Configure environment-specific configurations and secrets
  - Set up automated deployment from development to staging
  - _Requirements: 11.1, 11.2, 11.4_

- [ ] 20. Conduct end-to-end testing and validation
  - Perform load testing using K6 or similar tools
  - Test disaster recovery procedures and failover scenarios
  - Validate all application features work correctly in production
  - Test email delivery, AI report generation, and chat functionality
  - Perform security testing and vulnerability assessment
  - Document operational procedures and troubleshooting guides
  - _Requirements: All requirements validation_

## Detailed Implementation Commands

### Task 1: Verify AWS Configuration

```bash
# Verify your existing AWS CLI configuration
aws sts get-caller-identity
aws configure list

# Test your existing SES configuration
aws ses get-send-quota --region us-east-1

# Test your existing Bedrock access
aws bedrock list-foundation-models --region us-east-1

# Enable additional services if needed
aws eks describe-cluster --name test-cluster --region us-east-1 2>/dev/null || echo "EKS not yet configured"
```

### Task 2: Domain and DNS Setup

```bash
# Create hosted zone (if domain not already in Route 53)
aws route53 create-hosted-zone \
    --name cloudpartner.pro \
    --caller-reference $(date +%s)

# Get hosted zone ID
HOSTED_ZONE_ID=$(aws route53 list-hosted-zones-by-name \
    --dns-name cloudpartner.pro \
    --query 'HostedZones[0].Id' \
    --output text)

echo "Hosted Zone ID: $HOSTED_ZONE_ID"
```

### Task 3: SSL Certificate Setup

```bash
# Request certificate
aws acm request-certificate \
    --domain-name cloudpartner.pro \
    --subject-alternative-names "*.cloudpartner.pro" "api.cloudpartner.pro" \
    --validation-method DNS \
    --region us-east-1

# Get certificate ARN
CERT_ARN=$(aws acm list-certificates \
    --query 'CertificateSummaryList[?DomainName==`cloudpartner.pro`].CertificateArn' \
    --output text)

echo "Certificate ARN: $CERT_ARN"
```

### Task 4: VPC Creation

```bash
# Create VPC
VPC_ID=$(aws ec2 create-vpc \
    --cidr-block 10.0.0.0/16 \
    --query 'Vpc.VpcId' \
    --output text)

# Enable DNS hostnames
aws ec2 modify-vpc-attribute \
    --vpc-id $VPC_ID \
    --enable-dns-hostnames

# Create Internet Gateway
IGW_ID=$(aws ec2 create-internet-gateway \
    --query 'InternetGateway.InternetGatewayId' \
    --output text)

# Attach Internet Gateway to VPC
aws ec2 attach-internet-gateway \
    --internet-gateway-id $IGW_ID \
    --vpc-id $VPC_ID

echo "VPC ID: $VPC_ID"
echo "Internet Gateway ID: $IGW_ID"
```

### Task 5: RDS Database Setup

```bash
# Create DB subnet group
aws rds create-db-subnet-group \
    --db-subnet-group-name cloud-consulting-db-subnet-group \
    --db-subnet-group-description "Subnet group for Cloud Consulting DB" \
    --subnet-ids subnet-xxx subnet-yyy subnet-zzz

# Create RDS instance
aws rds create-db-instance \
    --db-instance-identifier cloud-consulting-prod \
    --db-instance-class db.t3.medium \
    --engine postgres \
    --engine-version 13.13 \
    --master-username postgres \
    --master-user-password "$(openssl rand -base64 32)" \
    --allocated-storage 100 \
    --storage-type gp2 \
    --storage-encrypted \
    --multi-az \
    --db-subnet-group-name cloud-consulting-db-subnet-group \
    --vpc-security-group-ids sg-xxx \
    --backup-retention-period 7 \
    --preferred-backup-window "03:00-04:00" \
    --preferred-maintenance-window "sun:04:00-sun:05:00" \
    --enable-performance-insights \
    --deletion-protection
```

### Task 9: EKS Cluster Creation

```bash
# Install eksctl
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin

# Create cluster configuration file
cat > cluster-config.yaml << EOF
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cloud-consulting-prod
  region: us-east-1
  version: "1.28"

vpc:
  id: $VPC_ID
  subnets:
    private:
      us-east-1a: { id: subnet-private-1a }
      us-east-1b: { id: subnet-private-1b }
      us-east-1c: { id: subnet-private-1c }
    public:
      us-east-1a: { id: subnet-public-1a }
      us-east-1b: { id: subnet-public-1b }
      us-east-1c: { id: subnet-public-1c }

managedNodeGroups:
  - name: main-nodes
    instanceType: t3.medium
    minSize: 2
    maxSize: 10
    desiredCapacity: 3
    volumeSize: 50
    privateNetworking: true
    iam:
      withAddonPolicies:
        autoScaler: true
        awsLoadBalancerController: true
        cloudWatch: true

addons:
  - name: vpc-cni
  - name: coredns
  - name: kube-proxy
  - name: aws-ebs-csi-driver

cloudWatch:
  clusterLogging:
    enableTypes: ["*"]
EOF

# Create EKS cluster
eksctl create cluster -f cluster-config.yaml
```

### Task 11: Application Deployment

```bash
# Create namespace
kubectl create namespace cloud-consulting

# Create ConfigMap
kubectl apply -f - << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: cloud-consulting
data:
  GIN_MODE: "release"
  LOG_LEVEL: "info"
  PORT: "8061"
  BEDROCK_REGION: "us-east-1"
  BEDROCK_MODEL_ID: "amazon.nova-lite-v1:0"
  AWS_SES_REGION: "us-east-1"
  CORS_ALLOWED_ORIGINS: "https://cloudpartner.pro,https://www.cloudpartner.pro"
  CHAT_MODE: "polling"
  CHAT_POLLING_INTERVAL: "3000"
  ENABLE_EMAIL_EVENTS: "true"
EOF

# Deploy applications using existing manifests
kubectl apply -f k8s/ -n cloud-consulting
```

### Task 13: Database Migration

```bash
# Get RDS endpoint
RDS_ENDPOINT=$(aws rds describe-db-instances \
    --db-instance-identifier cloud-consulting-prod \
    --query 'DBInstances[0].Endpoint.Address' \
    --output text)

# Run migrations from a pod
kubectl run migration-pod \
    --image=postgres:13 \
    --rm -it --restart=Never \
    --namespace=cloud-consulting \
    --env="PGPASSWORD=$DB_PASSWORD" \
    -- psql -h $RDS_ENDPOINT -U postgres -d postgres -f /scripts/init.sql
```

## Validation Steps

### Infrastructure Validation

```bash
# Test EKS cluster
kubectl cluster-info
kubectl get nodes

# Test database connectivity
kubectl exec -it deployment/backend -n cloud-consulting -- \
    pg_isready -h $RDS_ENDPOINT -p 5432

# Test Redis connectivity
kubectl exec -it deployment/backend -n cloud-consulting -- \
    redis-cli -h $REDIS_ENDPOINT ping

# Test application health
curl -f https://api.cloudpartner.pro/health
curl -f https://cloudpartner.pro
```

### Application Feature Testing

```bash
# Test AI report generation
curl -X POST https://api.cloudpartner.pro/api/v1/inquiries \
    -H "Content-Type: application/json" \
    -d '{"name":"Test User","email":"test@example.com","services":["assessment"],"message":"Test inquiry"}'

# Test admin authentication
curl -X POST https://api.cloudpartner.pro/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"cloudadmin"}'

# Test chat functionality
curl -X POST https://api.cloudpartner.pro/api/v1/admin/simple-chat/messages \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"message":"Hello AI","sessionId":"test-session"}'
```

## Cost Estimation

### Monthly AWS Costs (Approximate)

- **EKS Cluster**: $73/month (cluster) + $45/month (3 t3.medium nodes)
- **RDS PostgreSQL**: $35/month (db.t3.medium Multi-AZ)
- **ElastiCache Redis**: $15/month (cache.t3.micro)
- **Application Load Balancer**: $23/month
- **Data Transfer**: $10-50/month (depending on usage)
- **CloudWatch/Monitoring**: $10-30/month
- **Route 53**: $0.50/month per hosted zone
- **Certificate Manager**: Free
- **Total Estimated**: $210-285/month

### Cost Optimization Recommendations

- Use Spot instances for development/staging environments
- Implement auto-scaling to reduce costs during low usage
- Use Reserved Instances for predictable workloads
- Monitor and optimize data transfer costs
- Set up billing alerts and budgets

This implementation plan provides a comprehensive roadmap for deploying your Cloud Consulting Platform on AWS with production-ready infrastructure, security, monitoring, and cost optimization.
