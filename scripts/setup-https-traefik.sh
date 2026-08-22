#!/bin/bash
# SSL Certificate Setup Script for Traefik
# This script initializes Let's Encrypt certificates using Traefik

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║   SuProxy HTTPS Setup - Traefik + Let's Encrypt           ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check if .env.production exists
if [ ! -f ".env.production" ]; then
    echo "❌ Error: .env.production file not found"
    echo "Please copy .env.production.example to .env.production and configure it"
    exit 1
fi

# Load environment variables
source .env.production

# Check required variables
if [ -z "$DOMAIN" ] || [ "$DOMAIN" == "yourdomain.com" ]; then
    echo "❌ Error: DOMAIN not configured in .env.production"
    echo "Please set DOMAIN=yourdomain.com in .env.production"
    exit 1
fi

if [ -z "$ACME_EMAIL" ] || [ "$ACME_EMAIL" == "admin@yourdomain.com" ]; then
    echo "❌ Error: ACME_EMAIL not configured in .env.production"
    echo "Please set ACME_EMAIL=your-email@domain.com in .env.production"
    exit 1
fi

echo "📋 Configuration:"
echo "   Domain: $DOMAIN"
echo "   Email: $ACME_EMAIL"
echo ""

# Create acme.json with correct permissions
echo "🔐 Setting up Let's Encrypt storage..."
docker volume create traefik_letsencrypt
echo "✅ Volume created"
echo ""

# Start Traefik with staging first (testing)
echo "⚠️  IMPORTANT: Testing with Let's Encrypt staging first"
echo "   This prevents rate limiting during testing"
echo ""
read -p "Press Enter to start with staging certificates..."

# Uncomment staging caServer in traefik.yml
sed -i 's/# caServer: https:\/\/acme-staging/caServer: https:\/\/acme-staging/' traefik/traefik.yml

echo "🚀 Starting services with staging certificates..."
docker-compose -f docker-compose.production-https.yml up -d

echo ""
echo "⏳ Waiting 30 seconds for certificate generation..."
sleep 30

# Check if certificate was generated
echo "🔍 Checking certificate status..."
docker-compose -f docker-compose.production-https.yml logs traefik | grep -i "certificate"

echo ""
echo "✅ Staging setup complete!"
echo ""
echo "📝 Next steps:"
echo "   1. Verify your site works at https://api.$DOMAIN"
echo "   2. Check for certificate warnings (expected with staging)"
echo "   3. If everything works, run: ./scripts/enable-production-certs.sh"
echo ""
echo "🔒 Production certificates will be enabled after verification"
