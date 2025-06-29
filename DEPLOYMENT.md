# DigitalOcean App Platform Deployment Guide

## Prerequisites

1. A DigitalOcean account
2. Your repository pushed to GitHub
3. Environment variables configured

## Environment Configurations

This project includes configurations for multiple environments:

- **Development**: `.do/app.dev.yaml` - Uses `develop` branch, smallest instance
- **Test**: `.do/app.test.yaml` - Uses `test` branch, smallest instance
- **Staging**: `.do/app.staging.yaml` - Uses `staging` branch, medium instance
- **Production**: `.do/app.prod.yaml` - Uses `main` branch, larger instance with 2 replicas

## Required Environment Variables

You'll need to set these environment variables in each DigitalOcean App Platform environment:

### Shared Variables (same across all environments):

- `OPENAI_API_KEY` - Your OpenAI API key
- `OPENAI_MODEL` - OpenAI model to use (e.g., "gpt-4")
- `ALGOLIA_APP_ID` - Your Algolia application ID
- `ALGOLIA_ADMIN_API_KEY` - Your Algolia admin API key

### Environment-Specific Variables:

- `JWT_SECRET_DEV` / `JWT_SECRET_TEST` / `JWT_SECRET_STAGING` / `JWT_SECRET_PROD` - Secret keys for JWT token signing
- `DATABASE_URL_DEV` / `DATABASE_URL_TEST` / `DATABASE_URL_STAGING` / `DATABASE_URL_PROD` - PostgreSQL connection strings
- `ROLAND_DATABASE_URL_DEV` / `ROLAND_DATABASE_URL_TEST` / `ROLAND_DATABASE_URL_STAGING` / `ROLAND_DATABASE_URL_PROD` - Roland database connection strings

## Deployment Steps

### For each environment:

1. **Create the app in DigitalOcean:**

   - Go to DigitalOcean App Platform
   - Click "Create App"
   - Connect your GitHub repository
   - Choose the appropriate branch for the environment

2. **Configure the app:**

   - Use the corresponding `.do/app.{env}.yaml` file
   - DigitalOcean will automatically configure the deployment based on the file

3. **Set environment variables:**

   - In the App Platform dashboard, go to your app settings
   - Add all the required environment variables for that environment
   - Make sure to use the environment-specific variable names

4. **Deploy:**
   - DigitalOcean will automatically build and deploy your app
   - The health check endpoint `/health` will be used to verify the app is running

## Environment Differences

| Environment | Branch  | Instance Size | Instances | Rate Limit | Purpose                   |
| ----------- | ------- | ------------- | --------- | ---------- | ------------------------- |
| Development | develop | basic-xxs     | 1         | 1 req/s    | Local development testing |
| Test        | test    | basic-xxs     | 1         | 1 req/s    | Automated testing         |
| Staging     | staging | basic-xs      | 1         | 2 req/s    | Pre-production testing    |
| Production  | main    | basic-s       | 2         | 5 req/s    | Live production           |

## Accessing Your APIs

Once deployed, your gRPC APIs will be available at:

- **Development**: `https://study-guides-api-dev.ondigitalocean.app`
- **Test**: `https://study-guides-api-test.ondigitalocean.app`
- **Staging**: `https://study-guides-api-staging.ondigitalocean.app`
- **Production**: `https://study-guides-api-prod.ondigitalocean.app`

Each environment supports both HTTP/1.1, HTTP/2, and gRPC protocols.

## Monitoring

- Check the App Platform dashboard for logs and metrics for each environment
- The health check endpoint will help monitor app status
- Set up alerts for any deployment issues
- Monitor rate limiting and performance metrics
