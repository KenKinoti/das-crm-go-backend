# DAS CRM Backend - Google Cloud Run Deployment Guide

## 🚀 Quick Deploy

### Prerequisites
- Google Cloud CLI installed and authenticated
- Docker Desktop running
- Project `das-crm-frontend` exists in Google Cloud

### Option 1: Automated Script
```bash
# Windows
./deploy-to-cloudrun.bat

# Linux/Mac
./deploy-to-cloudrun.sh
```

### Option 2: Manual Commands

1. **Set up project**
```bash
gcloud config set project das-crm-frontend
```

2. **Enable APIs**
```bash
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
```

3. **Build and push image**
```bash
gcloud builds submit --tag gcr.io/das-crm-frontend/das-crm-backend
```

4. **Deploy to Cloud Run**
```bash
gcloud run deploy das-crm-backend \
  --image gcr.io/das-crm-frontend/das-crm-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --port 8080 \
  --memory 1Gi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 10 \
  --timeout 300 \
  --concurrency 80
```

5. **Get service URL**
```bash
gcloud run services describe das-crm-backend --region=us-central1 --format="value(status.url)"
```

## 📋 Environment Variables

Set these environment variables in Cloud Run:

### Required
- `DATABASE_URL`: Your database connection string
- `JWT_SECRET`: JWT signing secret
- `ENVIRONMENT`: Set to `production`

### Optional
- `PORT`: Will default to 8080 (Cloud Run requirement)
- `GIN_MODE`: Will be set to `release` automatically

## 🗄️ Database Setup

### Option 1: Google Cloud SQL
```bash
# Create PostgreSQL instance
gcloud sql instances create das-crm-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=us-central1

# Create database
gcloud sql databases create crm --instance=das-crm-db

# Create user
gcloud sql users create crmuser --instance=das-crm-db --password=YOUR_PASSWORD
```

### Option 2: External Database
- Use your existing PostgreSQL/MySQL database
- Update `DATABASE_URL` accordingly

## ⚙️ Configuration

### Cloud Run Settings
- **Memory**: 1Gi (can be adjusted based on load)
- **CPU**: 1 (can be scaled up if needed)
- **Min instances**: 0 (cost-effective)
- **Max instances**: 10 (adjust based on expected traffic)
- **Timeout**: 300 seconds
- **Concurrency**: 80 requests per instance

### Security
- Service allows unauthenticated requests (adjust based on your needs)
- CORS is configured in the application
- JWT authentication handled at application level

## 🔧 Updating Environment Variables

```bash
gcloud run services update das-crm-backend \
  --region us-central1 \
  --set-env-vars DATABASE_URL=your_database_url,JWT_SECRET=your_jwt_secret
```

## 📊 Monitoring

- **Logs**: `gcloud logs tail "projects/das-crm-frontend/logs/run.googleapis.com"`
- **Metrics**: Available in Google Cloud Console > Cloud Run
- **Health checks**: Automatic with Cloud Run

## 🚨 Troubleshooting

### Common Issues
1. **Build fails**: Check Dockerfile syntax and dependencies
2. **Service won't start**: Verify PORT environment variable (should be 8080)
3. **Database connection fails**: Check DATABASE_URL format and network access
4. **CORS issues**: Ensure frontend URL is properly configured

### Debug Commands
```bash
# View service details
gcloud run services describe das-crm-backend --region us-central1

# View logs
gcloud logs read "projects/das-crm-frontend/logs/run.googleapis.com" --limit 50

# Test deployment locally
docker build -t das-crm-backend .
docker run -p 8080:8080 -e PORT=8080 das-crm-backend
```

## 🔄 CI/CD Integration

For automated deployments, consider:
- GitHub Actions with Cloud Run deploy
- Cloud Build triggers
- Firebase hosting with Cloud Run backend

## 📈 Scaling

Cloud Run automatically scales based on traffic:
- **Cold starts**: ~1-2 seconds
- **Auto-scaling**: 0 to max instances
- **Pay per use**: Only pay for actual requests

## 🔐 Security Best Practices

1. **Environment Variables**: Store secrets in Google Secret Manager
2. **IAM**: Use least privilege principles
3. **HTTPS**: Automatically provided by Cloud Run
4. **VPC**: Consider VPC connector for database access

---

## 🎉 Expected Result

After successful deployment:
- Backend API available at: `https://das-crm-backend-[hash]-uc.a.run.app`
- Health check endpoint: `/api/v1/health`
- Auto-scaling enabled
- HTTPS termination
- Logs available in Cloud Console