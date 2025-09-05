#!/bin/bash

echo "===================================="
echo "  DAS CRM Backend - Cloud Run Deploy"
echo "===================================="
echo

echo "Step 1: Set project variables"
PROJECT_ID="das-crm-frontend"
SERVICE_NAME="das-crm-backend"
REGION="us-central1"
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

echo "Project ID: $PROJECT_ID"
echo "Service Name: $SERVICE_NAME"
echo "Region: $REGION"
echo "Image Name: $IMAGE_NAME"
echo

echo "Step 2: Configure gcloud project"
gcloud config set project $PROJECT_ID

echo
echo "Step 3: Enable required APIs"
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com

echo
echo "Step 4: Build and push image to Container Registry"
gcloud builds submit --tag $IMAGE_NAME

echo
echo "Step 5: Deploy to Cloud Run"
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080 \
  --memory 1Gi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 10 \
  --timeout 300 \
  --concurrency 80

echo
echo "Step 6: Get service URL"
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)")
echo "Service URL: $SERVICE_URL"

echo
echo "===================================="
echo "  Deployment Complete!"
echo "===================================="
echo "Your backend is now live at: $SERVICE_URL"