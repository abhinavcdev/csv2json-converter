# Deployment Guide

This guide covers multiple options to deploy your CSV2JSON converter to the internet.

## Quick Start (Recommended: Railway)

### 1. Railway - Easiest Option
1. Push your code to GitHub
2. Go to [railway.app](https://railway.app)
3. Click "Deploy from GitHub repo"
4. Select your repository
5. Railway will auto-detect the Dockerfile and deploy
6. Your app will be available at `https://your-app-name.railway.app`

### 2. Render - Free Option
1. Push code to GitHub
2. Go to [render.com](https://render.com)
3. Create new "Web Service"
4. Connect your GitHub repo
5. Use these settings:
   - **Build Command:** `docker build -t csv2json .`
   - **Start Command:** `./csv2json -server`
   - **Port:** 8080
6. Deploy (free tier available)

### 3. Fly.io - Global Edge Deployment
```bash
# Install flyctl
brew install flyctl

# Login and deploy
flyctl auth login
flyctl launch --no-deploy
flyctl deploy
```

### 4. DigitalOcean App Platform
1. Go to [cloud.digitalocean.com](https://cloud.digitalocean.com)
2. Create new App
3. Connect GitHub repo
4. Select "Dockerfile" as build method
5. Set port to 8080
6. Deploy ($5/month)

### 5. Google Cloud Run
```bash
# Build and push to Google Container Registry
gcloud builds submit --tag gcr.io/PROJECT-ID/csv2json
gcloud run deploy --image gcr.io/PROJECT-ID/csv2json --platform managed
```

## Environment Variables
Most platforms will automatically detect port 8080. If needed, set:
- `PORT=8080`

## Domain Setup
After deployment, you can:
1. Use the provided subdomain (e.g., `your-app.railway.app`)
2. Add a custom domain in your platform's dashboard
3. SSL certificates are automatically provided

## Cost Comparison
- **Render:** Free tier available
- **Railway:** Free tier, then $5/month
- **Fly.io:** Generous free tier
- **DigitalOcean:** $5/month minimum
- **Google Cloud Run:** Pay per request (very cheap for low traffic)

## Files Created
- `render.yaml` - Render configuration
- `fly.toml` - Fly.io configuration  
- `railway.json` - Railway configuration
- `Dockerfile` - Already exists, works with all platforms

Choose Railway for the easiest deployment experience!
