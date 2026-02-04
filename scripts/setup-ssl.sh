#!/bin/bash

# ===========================================
# SSL Certificate Setup with Let's Encrypt
# ===========================================

set -e

echo "🔒 SSL Certificate Setup"
echo "========================"
echo ""

# Check if certbot is installed
if ! command -v certbot &> /dev/null; then
    echo "📦 Installing certbot..."
    if command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo apt-get install -y certbot python3-certbot-nginx
    elif command -v yum &> /dev/null; then
        sudo yum install -y certbot python3-certbot-nginx
    else
        echo "❌ Please install certbot manually"
        exit 1
    fi
fi

# Load environment
if [ -f .env.production ]; then
    source .env.production
else
    echo "❌ .env.production not found"
    exit 1
fi

# Validate domain
if [ -z "$PRIMARY_DOMAIN" ]; then
    echo "❌ PRIMARY_DOMAIN not set in .env.production"
    exit 1
fi

echo "🌐 Domain: $PRIMARY_DOMAIN"
echo "📧 Email: ${LETSENCRYPT_EMAIL:-admin@$PRIMARY_DOMAIN}"
echo ""

# Get certificates
echo "📜 Obtaining SSL certificates..."
echo ""

domains=(
    "$PRIMARY_DOMAIN"
    "api.$PRIMARY_DOMAIN"
    "cdn.$PRIMARY_DOMAIN"
    "admin.$PRIMARY_DOMAIN"
)

# Build domain args
domain_args=""
for domain in "${domains[@]}"; do
    domain_args="$domain_args -d $domain"
done

# Run certbot
if [ "$LETSENCRYPT_STAGING" == "true" ]; then
    staging_arg="--staging"
    echo "⚠️  Using Let's Encrypt staging environment"
else
    staging_arg=""
fi

sudo certbot certonly \
    --nginx \
    $staging_arg \
    --non-interactive \
    --agree-tos \
    --email "${LETSENCRYPT_EMAIL:-admin@$PRIMARY_DOMAIN}" \
    $domain_args

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ SSL certificates obtained successfully!"
    echo ""
    echo "📂 Certificate location:"
    echo "   /etc/letsencrypt/live/$PRIMARY_DOMAIN/"
    echo ""
    echo "📝 Update your NGINX configuration:"
    echo ""
    cat << EOF
server {
    listen 443 ssl http2;
    server_name $PRIMARY_DOMAIN;

    ssl_certificate /etc/letsencrypt/live/$PRIMARY_DOMAIN/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$PRIMARY_DOMAIN/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Your location blocks here
}

server {
    listen 80;
    server_name $PRIMARY_DOMAIN;
    return 301 https://\$server_name\$request_uri;
}
EOF
    echo ""
    echo "🔄 Auto-renewal setup:"
    echo "   Certbot auto-renewal is configured via systemd timer"
    echo "   Test renewal: sudo certbot renew --dry-run"
    echo ""
else
    echo "❌ Failed to obtain certificates"
    exit 1
fi

# Setup auto-renewal hook
echo "🔧 Setting up renewal hooks..."
sudo mkdir -p /etc/letsencrypt/renewal-hooks/deploy

cat << 'EOF' | sudo tee /etc/letsencrypt/renewal-hooks/deploy/reload-nginx.sh > /dev/null
#!/bin/bash
# Reload NGINX after certificate renewal
systemctl reload nginx
EOF

sudo chmod +x /etc/letsencrypt/renewal-hooks/deploy/reload-nginx.sh

echo "✅ Renewal hooks configured"
echo ""
echo "==================================="
echo "✅ SSL setup complete!"
echo "==================================="
echo ""
