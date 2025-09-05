# 🚀 Deploy Backend NOW - Step by Step

Open your **Windows Command Prompt** or **PowerShell** where you have gcloud installed and follow these steps:

## Step 1: Navigate to Backend Directory
```cmd
cd D:\DASYIN\CRM\das-crm-go-backend
```

## Step 1.5: (Optional) Set up Cloud SQL PostgreSQL
If you need a managed PostgreSQL database:
```cmd
gcloud sql instances create das-crm-db --database-version=POSTGRES_15 --tier=db-f1-micro --region=us-central1
gcloud sql databases create crm --instance=das-crm-db
gcloud sql users create crmuser --instance=das-crm-db --password=YOUR_SECURE_PASSWORD
```
**Note:** This creates a free-tier PostgreSQL instance. Skip if using external PostgreSQL.

## Step 2: Set Your Project
```cmd
gcloud config set project das-crm-frontend
```

## Step 3: Enable Required APIs
```cmd
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
```

## Step 4: Build and Submit to Cloud Build
```cmd
gcloud builds submit --tag gcr.io/das-crm-frontend/das-crm-backend
```
**This will:**
- Build your Docker image in the cloud
- Push it to Google Container Registry
- Take approximately 2-3 minutes

## Step 5: Deploy to Cloud Run
```cmd
gcloud run deploy das-crm-backend --image gcr.io/das-crm-frontend/das-crm-backend --platform managed --region us-central1 --allow-unauthenticated --port 8080 --memory 1Gi --cpu 1 --min-instances 0 --max-instances 10
```
**When prompted:**
- Confirm the deployment by pressing Enter or typing 'y'

## Step 6: Set Environment Variables
```cmd
gcloud run services update das-crm-backend --region us-central1 --update-env-vars DATABASE_URL="postgresql://user:password@host:5432/database?sslmode=require",JWT_SECRET="your-secure-jwt-secret-here",ENVIRONMENT="production"
```
**Replace:**
- `user:password@host:5432/database` with your actual PostgreSQL connection details
- `your-secure-jwt-secret-here` with a secure JWT secret (min 32 characters)

### PostgreSQL Connection Examples:
- **External PostgreSQL**: `postgresql://user:pass@host:5432/dbname?sslmode=require`
- **Google Cloud SQL**: `postgresql://user:pass@/dbname?host=/cloudsql/PROJECT:REGION:INSTANCE`
- **Local development**: `postgresql://user:pass@localhost:5432/dbname?sslmode=disable`

## Step 7: Get Your Service URL
```cmd
gcloud run services describe das-crm-backend --region us-central1 --format "value(status.url)"
```

---

## 🎯 Expected Output

After successful deployment, you should see:
```
Service [das-crm-backend] revision [das-crm-backend-00001-xxx] has been deployed and is serving 100 percent of traffic.
Service URL: https://das-crm-backend-xxxxx-uc.a.run.app
```

## ✅ Test Your Deployment

1. **Health Check:**
   Open your browser and visit:
   ```
   https://das-crm-backend-xxxxx-uc.a.run.app/health
   ```
   You should see a success response.

2. **Update Frontend:**
   Update your frontend environment to use the new backend URL.

---

## 🚨 If You Encounter Issues

### "Permission denied" error
```cmd
gcloud auth login
```

### "APIs not enabled" error
Run Step 3 again to enable the APIs.

### Build fails
Check that Docker files are correct:
- Dockerfile exists
- .dockerignore exists
- No syntax errors

### Service won't start
Check logs:
```cmd
gcloud logs read --project das-crm-frontend --limit 50
```

---

## 🎉 Success!

Your backend is now:
- ✅ Running on Google Cloud Run
- ✅ Auto-scaling from 0 to 10 instances
- ✅ Accessible via HTTPS
- ✅ Ready for production use

**Next:** Update your frontend to use the new backend URL!