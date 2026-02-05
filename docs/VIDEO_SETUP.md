# Video Conferencing Setup Guide

## Overview

The platform supports three video conferencing providers:
- **Zoom** (Recommended for production)
- **Agora** (Best for scalability)
- **Jitsi** (Free, open-source)

---

## Option 1: Zoom Setup

### Step 1: Create Zoom Account

1. Go to https://marketplace.zoom.us
2. Sign in or create account
3. Click "Develop" → "Build App"

### Step 2: Create Server-to-Server OAuth App

1. Choose "Server-to-Server OAuth"
2. Fill in app information:
   - App Name: "Language Learning Platform"
   - Company Name: Your company
   - Developer Contact: your-email@domain.com

### Step 3: Get Credentials

After creating the app, you'll get:
- Account ID
- Client ID
- Client Secret

Add to `config/env/.env.production`:
```bash
ZOOM_ACCOUNT_ID=your-account-id
ZOOM_CLIENT_ID=your-client-id
ZOOM_CLIENT_SECRET=your-client-secret
VIDEO_PROVIDER=zoom
```

### Step 4: Configure Scopes

Add these scopes to your Zoom app:
- `meeting:write:admin` - Create meetings
- `meeting:read:admin` - Read meeting info
- `recording:write:admin` - Manage recordings
- `user:read:admin` - Read user info

### Step 5: Activate App

1. Click "Activation" tab
2. Click "Activate your app"

### Step 6: Test Integration

```bash
# Test creating a meeting
curl -X POST https://api.yourdomain.com/api/v1/sessions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": "course-uuid",
    "title": "Spanish Conversation Practice",
    "description": "Live conversation class",
    "scheduled_at": "2026-02-10T15:00:00Z",
    "duration": 60
  }'
```

---

## Option 2: Agora Setup

### Step 1: Create Agora Account

1. Go to https://console.agora.io
2. Sign up for free account
3. Verify email

### Step 2: Create Project

1. Click "Projects" in sidebar
2. Click "Create New Project"
3. Choose "Secured mode: APP ID + Token"
4. Name: "Language Learning Platform"

### Step 3: Get Credentials

You'll receive:
- App ID
- App Certificate
- Customer ID
- Customer Secret

Add to `config/env/.env.production`:
```bash
AGORA_APP_ID=your-app-id
AGORA_APP_CERTIFICATE=your-app-certificate
AGORA_CUSTOMER_ID=your-customer-id
AGORA_CUSTOMER_SECRET=your-customer-secret
VIDEO_PROVIDER=agora
```

### Step 4: Enable Services

In Agora Console:
1. Go to "Products & Usage"
2. Enable:
   - Real-Time Communication
   - Cloud Recording
   - Interactive Live Streaming

### Step 5: Configure Webhook

1. Go to "Notifications"
2. Add webhook URL: `https://api.yourdomain.com/webhooks/agora`
3. Select events:
   - Channel Created
   - Channel Destroyed
   - Recording Started
   - Recording Stopped

### Step 6: Implement Token Generation

Agora requires server-side token generation. Example implementation:

```go
// services/video-service/internal/agora/token.go
package agora

import (
    "time"
    rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

func GenerateRTCToken(channelName string, uid uint32, role rtctokenbuilder.Role) (string, error) {
    appID := os.Getenv("AGORA_APP_ID")
    appCertificate := os.Getenv("AGORA_APP_CERTIFICATE")
    
    expireTime := uint32(time.Now().Unix()) + 3600 // 1 hour
    
    token, err := rtctokenbuilder.BuildTokenWithUid(
        appID,
        appCertificate,
        channelName,
        uid,
        role,
        expireTime,
        expireTime,
    )
    
    return token, err
}
```

---

## Option 3: Jitsi Setup

### Step 1: Choose Deployment

**Option A: Use public Jitsi (meet.jit.si)**
```bash
JITSI_DOMAIN=meet.jit.si
VIDEO_PROVIDER=jitsi
```

**Option B: Self-hosted (Recommended for production)**

### Step 2: Install Jitsi Meet (Self-hosted)

```bash
# On Ubuntu 22.04
sudo apt update
sudo apt install -y apt-transport-https

# Add Jitsi repository
curl https://download.jitsi.org/jitsi-key.gpg.key | sudo sh -c 'gpg --dearmor > /usr/share/keyrings/jitsi-keyring.gpg'
echo 'deb [signed-by=/usr/share/keyrings/jitsi-keyring.gpg] https://download.jitsi.org stable/' | sudo tee /etc/apt/sources.list.d/jitsi-stable.list

# Install
sudo apt update
sudo apt install -y jitsi-meet

# Configure SSL
sudo /usr/share/jitsi-meet/scripts/install-letsencrypt-cert.sh
```

### Step 3: Configure JWT Authentication

```bash
# Install JWT plugin
sudo apt install -y jitsi-meet-tokens

# Edit config
sudo nano /etc/prosody/conf.avail/meet.jitsi.yourdomain.com.cfg.lua
```

Add:
```lua
VirtualHost "meet.jitsi.yourdomain.com"
    authentication = "token"
    app_id = "language_platform"
    app_secret = "your-generated-secret"
```

### Step 4: Configure Environment

```bash
JITSI_DOMAIN=meet.jitsi.yourdomain.com
JITSI_APP_ID=language_platform
JITSI_APP_SECRET=your-generated-secret
VIDEO_PROVIDER=jitsi
```

### Step 5: Generate JWT Tokens

```go
// services/video-service/internal/jitsi/jwt.go
package jitsi

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

func GenerateJitsiToken(roomName, userName string) (string, error) {
    claims := jwt.MapClaims{
        "aud": "jitsi",
        "iss": os.Getenv("JITSI_APP_ID"),
        "sub": os.Getenv("JITSI_DOMAIN"),
        "room": roomName,
        "context": map[string]interface{}{
            "user": map[string]string{
                "name": userName,
            },
        },
        "exp": time.Now().Add(2 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JITSI_APP_SECRET")))
}
```

---

## Comparison

| Feature | Zoom | Agora | Jitsi |
|---------|------|-------|-------|
| **Ease of Setup** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| **Cost** | $$ | $ | Free |
| **Scalability** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Recording** | ✅ Built-in | ✅ Cloud | ✅ Requires setup |
| **Mobile Support** | ✅ Excellent | ✅ Excellent | ✅ Good |
| **Video Quality** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Reliability** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Customization** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

### Recommendations

- **Small/Medium Platforms (<1000 users)**: Zoom
- **Large Platforms (>1000 users)**: Agora
- **Budget-Conscious/Open-Source**: Jitsi

---

## Implementation in Services

### Video Service Endpoints

```go
// POST /api/v1/sessions - Create video session
// Returns meeting URL and credentials

type CreateSessionResponse struct {
    SessionID   string    `json:"session_id"`
    MeetingURL  string    `json:"meeting_url"`
    MeetingID   string    `json:"meeting_id"`
    Password    string    `json:"password,omitempty"`
    Token       string    `json:"token,omitempty"` // For Agora/Jitsi
    ScheduledAt time.Time `json:"scheduled_at"`
}
```

### Frontend Integration

**Zoom:**
```typescript
// Use Zoom Web SDK
import ZoomVideo from '@zoom/videosdk'

const client = ZoomVideo.createClient()
await client.join(sessionName, token, userName)
```

**Agora:**
```typescript
// Use Agora Web SDK
import AgoraRTC from 'agora-rtc-sdk-ng'

const client = AgoraRTC.createClient({ mode: 'rtc', codec: 'vp8' })
await client.join(appId, channel, token, uid)
```

**Jitsi:**
```typescript
// Use Jitsi Meet API
const api = new JitsiMeetExternalAPI('meet.jitsi.yourdomain.com', {
    roomName: 'session-123',
    jwt: token,
    parentNode: document.getElementById('jitsi-container')
})
```

---

## Testing

### Test Script

```bash
#!/bin/bash

echo "Testing video conferencing setup..."

# Test session creation
response=$(curl -s -X POST https://api.yourdomain.com/api/v1/sessions \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": "test-course-id",
    "title": "Test Session",
    "scheduled_at": "2026-02-10T15:00:00Z",
    "duration": 30
  }')

echo "Response: $response"

# Extract meeting URL
meeting_url=$(echo $response | jq -r '.meeting_url')
echo "Meeting URL: $meeting_url"

# Test joining (manual verification needed)
echo "Please test joining the meeting at: $meeting_url"
```

---

## Monitoring

### Zoom Webhooks

Configure webhooks at https://marketplace.zoom.us/develop/apps/YOUR_APP_ID/webhooks

Events to monitor:
- `meeting.started`
- `meeting.ended`
- `meeting.participant_joined`
- `meeting.participant_left`
- `recording.completed`

### Agora Callbacks

Configure in Agora Console → Notifications:
- Channel events
- Recording events
- User events

### Jitsi Logs

```bash
# View Jitsi logs
sudo tail -f /var/log/jitsi/jicofo.log
sudo tail -f /var/log/jitsi/jvb.log
```

---

## Troubleshooting

### Zoom: "Invalid API Key"
- Verify credentials in .env
- Check app is activated
- Ensure OAuth scopes are correct

### Agora: "Token expired"
- Check token generation logic
- Verify system time is synchronized
- Increase token expiration time

### Jitsi: "Room creation failed"
- Check JWT configuration
- Verify Prosody is running: `sudo systemctl status prosody`
- Check firewall rules for ports 443, 4443, 10000

---

## Security Best Practices

1. **Never expose API keys in frontend**
2. **Use short-lived tokens** (1-2 hours)
3. **Validate meeting participants** server-side
4. **Enable waiting rooms** (Zoom)
5. **Implement rate limiting** on session creation
6. **Log all video session activity**
7. **Encrypt recordings** at rest

---

## Cost Estimation

### Zoom
- Basic: $0 (40 min limit)
- Pro: $149.90/year/host
- Business: $199.90/year/host

### Agora
- Pay-as-you-go
- ~$0.99 per 1000 minutes
- Free 10,000 minutes/month

### Jitsi
- Self-hosted: Server costs only
- ~$50-200/month for VPS

---

## Next Steps

1. Choose your video provider
2. Sign up and get credentials
3. Add credentials to `config/env/.env.production`
4. Test video session creation
5. Integrate frontend SDK
6. Configure webhooks
7. Set up monitoring
8. Train support team

For support: devops@yourdomain.com
