#!/bin/bash
# SSL Certificate Setup Script for Nginx
# This script initializes Let's Encrypt certificates using Certbot

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║   SuProxy HTTPS Setup - Nginx + Let's Encrypt             ║"
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
echo "   Domain: api.$DOMAIN"
echo "   Email: $ACME_EMAIL"
echo ""

# Update nginx configuration with actual domain
echo "📝 Updating Nginx configuration..."
sed -i "s/yourdomain.com/$DOMAIN/g" nginx/suproxy.conf
echo "✅ Configuration updated"
echo ""

# Create directories
echo "📁 Creating directories..."
mkdir -p certbot/conf certbot/www
echo "✅ Directories created"
echo ""

# Start backend services first (without nginx)
echo "🚀 Starting backend services..."
docker-compose -f docker-compose.production-nginx.yml up -d postgres api prometheus grafana
echo "⏳ Waiting for services to be ready..."
sleep 20
echo "✅ Backend services started"
echo ""

# Obtain certificate using standalone mode
echo "🔐 Obtaining SSL certificate from Let's Encrypt..."
echo "   This may take a minute..."
docker run --rm \
  -v "$(pwd)/certbot/conf:/etc/letsencrypt" \
  -v "$(pwd)/certbot/www:/var/www/certbot" \
  -p 80:80 \
  certbot/certbot certonly \
  --standalone \
  --preferred-challenges http \
  --email "$ACME_EMAIL" \
  --agree-tos \
  --no-eff-email \
  -d "api.$DOMAIN"

if [ $? -eq 0 ]; then
    echo "✅ Certificate obtained successfully!"
else
    echo "❌ Failed to obtain certificate"
    echo "   Please check:"
    echo "   1. Domain DNS points to this server"
    echo "   2. Port 80 is accessible from internet"
    echo "   3. No firewall blocking connections"
    exit 1
fi
echo ""

# Start nginx and certbot
echo "🌐 Starting Nginx reverse proxy..."
docker-compose -f docker-compose.production-nginx.yml up -d nginx certbot
echo "✅ Nginx started"
echo ""

# Test nginx configuration
echo "🔍 Testing Nginx configuration..."
docker exec suproxy-nginx nginx -t
if [ $? -eq 0 ]; then
    echo "✅ Nginx configuration is valid"
else
    echo "❌ Nginx configuration has errors"
    exit 1
fi
echo ""

echo "╔════════════════════════════════════════════════════════════╗"
echo "║   HTTPS Setup Complete! 🎉                                 ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "✅ Your API is now accessible at: https://api.$DOMAIN"
echo ""
echo "🔒 SSL Certificate Info:"
echo "   Certificate will auto-renew every 12 hours (if needed)"
echo "   Certificate valid for: 90 days"
echo ""
echo "📝 Next Steps:"
echo "   1. Test your API: curl https://api.$DOMAIN/health"
echo "   2. Update frontend to use: https://api.$DOMAIN"
echo "   3. Monitor logs: docker-compose -f docker-compose.production-nginx.yml logs -f"
echo ""
