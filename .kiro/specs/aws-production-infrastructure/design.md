# AWS Production Infrastructure Design

## Overview

This document provides a comprehensive design for deploying the Cloud Consulting Platform on AWS. The architecture follows AWS Well-Architected Framework principles, ensuring security, reliability, performance efficiency, cost optimization, and operational excellence.

## Architecture

### High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        Internet Gateway                          │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                    Application Load Balancer                    │
│                     (SSL Termination)                          │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                      Amazon EKS Cluster                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Frontend      │  │    Backend      │  │     Redis       │ │
│  │   (React)       │  │     (Go)        │  │   (ElastiCache) │ │
│  │   Pods          │  │    Pods         │  │                 │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                     Amazon RDS PostgreSQL                      │
│                      (Multi-AZ)                                │
└─────────────────────────────────────────────────────────────────┘

External Services:
├── AWS Bedrock (AI)
├── AWS SES (Email)
├── AWS Secrets Manager
├── CloudWatch (Monitoring)
└── Route 53 (DNS)
```

### Network Architecture

#### VPC Configuration
- **VPC CIDR**: 10.0.0.0/16
- **Availability Zones**: 3 AZs for high availability
- **Subnets**:
  - Public Subnets: 10.0.1.0/24, 10.0.2.0/24, 10.0.3.0/24 (ALB, NAT Gateways)
  - Private Subnets: 10.0.11.0/24, 10.0.12.0/24, 10.0.13.0/24 (EKS Nodes)
  - Database Subnets: 10.0.21.0/24, 10.0.22.0/24, 10.0.23.0/24 (RDS, ElastiCache)

#### Security Groups
```yaml
# ALB Security Group
ALB-SG:
  Inbound:
    - Port 80 (HTTP) from 0.0.0.0/0
    - Port 443 (HTTPS) from 0.0.0.0/0
  Outbound:
    - All traffic to EKS-SG

# EKS Nodes Security Group
EKS-SG:
  Inbound:
    - Port 30000-32767 from ALB-SG
    - Port 443 from ALB-SG
    - All traffic from EKS-SG (self-reference)
  Outbound:
    - All traffic to 0.0.0.0/0

# RDS Security Group
RDS-SG:
  Inbound:
    - Port 5432 from EKS-SG
  Outbound:
    - None

# ElastiCache Security Group
REDIS-SG:
  Inbound:
    - Port 6379 from EKS-SG
  Outbound:
    - None
```

## Components and Interfaces

### 1. Amazon EKS Cluster

#### Cluster Configuration
```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cloud-consulting-prod
  region: us-east-1
  version: "1.28"

vpc:
  cidr: "10.0.0.0/16"
  nat:
    gateway: HighlyAvailable

managedNodeGroups:
  - name: main-nodes
    instanceType: t3.medium
    minSize: 2
    maxSize: 10
    desiredCapacity: 3
    volumeSize: 50
    ssh:
      enableSsm: true
    iam:
      withAddonPolicies:
        imageBuilder: true
        autoScaler: true
        externalDNS: true
        certManager: true
        appMesh: true
        ebs: true
        fsx: true
        efs: true
        awsLoadBalancerController: true
        xRay: true
        cloudWatch: true

addons:
  - name: vpc-cni
  - name: coredns
  - name: kube-proxy
  - name: aws-ebs-csi-driver

cloudWatch:
  clusterLogging:
    enableTypes: ["*"]
```

#### Application Deployments
Based on existing Kubernetes manifests in `k8s/` directory:

**Backend Deployment:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  namespace: cloud-consulting
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    spec:
      containers:
      - name: backend
        image: your-registry/cloud-consulting-backend:latest
        ports:
        - containerPort: 8061
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: database-url
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: redis-url
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8061
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8061
          initialDelaySeconds: 5
          periodSeconds: 5
```

**Frontend Deployment:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
  namespace: cloud-consulting
spec:
  replicas: 2
  selector:
    matchLabels:
      app: frontend
  template:
    spec:
      containers:
      - name: frontend
        image: your-registry/cloud-consulting-frontend:latest
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "200m"
```

### 2. Amazon RDS PostgreSQL

#### Database Configuration
```yaml
Engine: postgres
EngineVersion: "13.13"
DBInstanceClass: db.t3.medium
AllocatedStorage: 100
StorageType: gp2
StorageEncrypted: true
MultiAZ: true
BackupRetentionPeriod: 7
PreferredBackupWindow: "03:00-04:00"
PreferredMaintenanceWindow: "sun:04:00-sun:05:00"
DeletionProtection: true
EnablePerformanceInsights: true
MonitoringInterval: 60
```

#### Database Schema
The system will use the existing migration scripts:
- `backend/scripts/init.sql` - Core tables (inquiries, reports, activities)
- `backend/scripts/chat_migration.sql` - Chat system tables
- `backend/scripts/email_events_migration.sql` - Email tracking tables

#### Connection Configuration
```go
// Database connection pool settings
MaxOpenConnections: 25
MaxIdleConnections: 5
ConnMaxLifetime: 30 * time.Minute
```

### 3. Amazon ElastiCache Redis

#### Redis Configuration
```yaml
Engine: redis
EngineVersion: "6.2"
NodeType: cache.t3.micro
NumCacheNodes: 1
Port: 6379
ParameterGroupName: default.redis6.x
SecurityGroupIds: [redis-sg]
SubnetGroupName: redis-subnet-group
AtRestEncryptionEnabled: true
TransitEncryptionEnabled: true
```

#### Usage Patterns
- Session storage for JWT tokens
- Chat message caching
- AI response caching
- Performance optimization caching

### 4. AWS Application Load Balancer

#### ALB Configuration
```yaml
Type: application
Scheme: internet-facing
IpAddressType: ipv4
SecurityGroups: [alb-sg]
Subnets: [public-subnet-1, public-subnet-2, public-subnet-3]

Listeners:
  - Port: 80
    Protocol: HTTP
    DefaultActions:
      - Type: redirect
        RedirectConfig:
          Protocol: HTTPS
          Port: 443
          StatusCode: HTTP_301
  
  - Port: 443
    Protocol: HTTPS
    SslPolicy: ELBSecurityPolicy-TLS-1-2-2017-01
    Certificates: [arn:aws:acm:us-east-1:account:certificate/cert-id]
    DefaultActions:
      - Type: forward
        TargetGroupArn: frontend-tg

Rules:
  - Priority: 100
    Conditions:
      - Field: path-pattern
        Values: ["/api/*"]
    Actions:
      - Type: forward
        TargetGroupArn: backend-tg
```

### 5. AWS Certificate Manager

#### SSL Certificate Configuration
```yaml
DomainName: cloudpartner.pro
SubjectAlternativeNames:
  - "*.cloudpartner.pro"
  - "api.cloudpartner.pro"
ValidationMethod: DNS
```

### 6. Amazon Route 53

#### DNS Configuration
```yaml
HostedZone: cloudpartner.pro

Records:
  - Name: cloudpartner.pro
    Type: A
    AliasTarget: ALB DNS name
  
  - Name: api.cloudpartner.pro
    Type: A
    AliasTarget: ALB DNS name
  
  - Name: www.cloudpartner.pro
    Type: CNAME
    Value: cloudpartner.pro

HealthChecks:
  - Name: backend-health
    Type: HTTPS
    ResourcePath: /health
    FailureThreshold: 3
```

## Data Models

### Environment Variables and Secrets

#### AWS Secrets Manager Configuration
```json
{
  "database-url": "postgresql://username:password@rds-endpoint:5432/consulting",
  "redis-url": "redis://elasticache-endpoint:6379",
  "jwt-secret": "production-jwt-secret-key",
  "aws-access-key-id": "AKIA...",
  "aws-secret-access-key": "...",
  "bedrock-api-key": "...",
  "ses-sender-email": "info@cloudpartner.pro"
}
```

#### Kubernetes ConfigMap
```yaml
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
```

### IAM Roles and Policies

#### EKS Node Group Role
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "bedrock:InvokeModel",
        "bedrock:InvokeModelWithResponseStream"
      ],
      "Resource": "arn:aws:bedrock:us-east-1::foundation-model/amazon.nova-lite-v1:0"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ses:SendEmail",
        "ses:SendRawEmail",
        "ses:GetSendQuota",
        "ses:GetSendStatistics"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": "arn:aws:secretsmanager:us-east-1:account:secret:cloud-consulting/*"
    }
  ]
}
```

## Error Handling

### Database Connection Failures
```go
// Implement connection retry logic with exponential backoff
func (c *DatabaseConnection) connectWithRetry(maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        if err := c.db.Ping(); err == nil {
            return nil
        }
        
        backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
        time.Sleep(backoff)
    }
    return fmt.Errorf("failed to connect after %d retries", maxRetries)
}
```

### AWS Service Failures
```go
// Implement circuit breaker pattern for external services
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFailure time.Time
    state       string // "closed", "open", "half-open"
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = "half-open"
        } else {
            return fmt.Errorf("circuit breaker is open")
        }
    }
    
    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }
    
    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

### Graceful Degradation
```go
// Email service with fallback
func (e *EmailService) SendWithFallback(ctx context.Context, email *Email) error {
    // Try primary SES region
    if err := e.sendViaSES(ctx, email, "us-east-1"); err == nil {
        return nil
    }
    
    // Fallback to secondary region
    if err := e.sendViaSES(ctx, email, "us-west-2"); err == nil {
        return nil
    }
    
    // Log failure and queue for retry
    e.logger.Error("Email delivery failed, queuing for retry", "email", email.ID)
    return e.queueForRetry(email)
}
```

## Testing Strategy

### Infrastructure Testing
```bash
# Test EKS cluster connectivity
kubectl cluster-info

# Test database connectivity
kubectl exec -it backend-pod -- /bin/sh -c "pg_isready -h $DATABASE_HOST -p 5432"

# Test Redis connectivity
kubectl exec -it backend-pod -- /bin/sh -c "redis-cli -h $REDIS_HOST ping"

# Test external services
kubectl exec -it backend-pod -- /bin/sh -c "curl -f http://localhost:8061/health"
```

### Load Testing
```yaml
# K6 load test configuration
apiVersion: v1
kind: ConfigMap
metadata:
  name: load-test-script
data:
  script.js: |
    import http from 'k6/http';
    import { check } from 'k6';
    
    export let options = {
      stages: [
        { duration: '2m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '2m', target: 200 },
        { duration: '5m', target: 200 },
        { duration: '2m', target: 0 },
      ],
    };
    
    export default function() {
      let response = http.get('https://api.cloudpartner.pro/health');
      check(response, {
        'status is 200': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
      });
    }
```

### Disaster Recovery Testing
```bash
# Test RDS failover
aws rds reboot-db-instance --db-instance-identifier cloud-consulting-prod --force-failover

# Test EKS node failure
kubectl drain node-name --ignore-daemonsets --delete-emptydir-data

# Test application recovery
kubectl rollout restart deployment/backend -n cloud-consulting
kubectl rollout restart deployment/frontend -n cloud-consulting
```

## Monitoring and Observability

### CloudWatch Dashboards
```json
{
  "widgets": [
    {
      "type": "metric",
      "properties": {
        "metrics": [
          ["AWS/ApplicationELB", "RequestCount", "LoadBalancer", "cloud-consulting-alb"],
          ["AWS/ApplicationELB", "TargetResponseTime", "LoadBalancer", "cloud-consulting-alb"],
          ["AWS/RDS", "CPUUtilization", "DBInstanceIdentifier", "cloud-consulting-db"],
          ["AWS/ElastiCache", "CPUUtilization", "CacheClusterId", "cloud-consulting-redis"],
          ["AWS/EKS", "cluster_failed_request_count", "cluster_name", "cloud-consulting-prod"]
        ],
        "period": 300,
        "stat": "Average",
        "region": "us-east-1",
        "title": "Application Performance"
      }
    }
  ]
}
```

### CloudWatch Alarms
```yaml
# High CPU Alarm
CPUAlarm:
  Type: AWS::CloudWatch::Alarm
  Properties:
    AlarmName: HighCPUUtilization
    MetricName: CPUUtilization
    Namespace: AWS/RDS
    Statistic: Average
    Period: 300
    EvaluationPeriods: 2
    Threshold: 80
    ComparisonOperator: GreaterThanThreshold
    AlarmActions:
      - !Ref SNSTopic

# Application Error Rate Alarm
ErrorRateAlarm:
  Type: AWS::CloudWatch::Alarm
  Properties:
    AlarmName: HighErrorRate
    MetricName: HTTPCode_Target_5XX_Count
    Namespace: AWS/ApplicationELB
    Statistic: Sum
    Period: 300
    EvaluationPeriods: 2
    Threshold: 10
    ComparisonOperator: GreaterThanThreshold
```

### Log Aggregation
```yaml
# Fluent Bit configuration for log forwarding
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluent-bit-config
data:
  fluent-bit.conf: |
    [SERVICE]
        Flush         1
        Log_Level     info
        Daemon        off
        Parsers_File  parsers.conf
        HTTP_Server   On
        HTTP_Listen   0.0.0.0
        HTTP_Port     2020

    [INPUT]
        Name              tail
        Path              /var/log/containers/*.log
        Parser            docker
        Tag               kube.*
        Refresh_Interval  5
        Mem_Buf_Limit     50MB
        Skip_Long_Lines   On

    [OUTPUT]
        Name  cloudwatch_logs
        Match *
        region us-east-1
        log_group_name /aws/eks/cloud-consulting/cluster
        log_stream_prefix from-fluent-bit-
        auto_create_group true
```

## Security Considerations

### Network Security
- All database and cache instances in private subnets
- Security groups with minimal required access
- VPC Flow Logs enabled for network monitoring
- WAF rules for application protection

### Data Encryption
- RDS encryption at rest using AWS KMS
- ElastiCache encryption in transit and at rest
- EBS volume encryption for EKS nodes
- Secrets Manager for sensitive configuration

### Access Control
- IAM roles with least privilege principle
- Service accounts for Kubernetes workloads
- Regular access reviews and rotation of credentials
- Multi-factor authentication for AWS console access

### Compliance
- CloudTrail logging for audit trails
- Config rules for compliance monitoring
- Regular security assessments and penetration testing
- Data retention policies aligned with business requirements

This design provides a robust, scalable, and secure foundation for deploying the Cloud Consulting Platform on AWS, leveraging managed services to reduce operational overhead while maintaining high availability and performance.