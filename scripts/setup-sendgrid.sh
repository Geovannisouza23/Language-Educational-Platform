#!/bin/bash

# ===========================================
# SendGrid Configuration Setup
# ===========================================

set -e

echo "📧 SendGrid Configuration Setup"
echo "================================"
echo ""

# Check if SENDGRID_API_KEY is set
if [ -z "$SENDGRID_API_KEY" ]; then
    echo "❌ SENDGRID_API_KEY is not set"
    echo ""
    echo "Please follow these steps:"
    echo ""
    echo "1. Sign up at https://sendgrid.com"
    echo "2. Go to Settings > API Keys"
    echo "3. Create a new API key with 'Full Access'"
    echo "4. Copy the API key and run:"
    echo "   export SENDGRID_API_KEY='your-api-key'"
    echo ""
    exit 1
fi

echo "✅ API Key found"
echo ""

# Test API key
echo "🔍 Testing API key..."
response=$(curl -s -o /dev/null -w "%{http_code}" \
    -X GET "https://api.sendgrid.com/v3/user/profile" \
    -H "Authorization: Bearer ${SENDGRID_API_KEY}")

if [ "$response" == "200" ]; then
    echo "✅ API key is valid"
else
    echo "❌ API key is invalid (HTTP $response)"
    exit 1
fi

# Create email templates
echo ""
echo "📝 Creating email templates..."
echo ""

# Welcome template
echo "Creating welcome template..."
cat > /tmp/template-welcome.json << 'EOF'
{
  "name": "Welcome Email",
  "generation": "dynamic",
  "subject": "Welcome to Language Learning Platform!",
  "versions": [{
    "template_name": "Welcome Email",
    "subject": "Welcome to Language Learning Platform, {{name}}!",
    "html_content": "<!DOCTYPE html><html><head><meta charset='UTF-8'></head><body style='font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;'><div style='background-color: #2563eb; color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0;'><h1>Welcome to Language Learning Platform!</h1></div><div style='background-color: #f9fafb; padding: 30px; border-radius: 0 0 10px 10px;'><p style='font-size: 16px;'>Hi {{name}},</p><p>Thank you for joining our platform! We're excited to have you on board.</p><p>You can now:</p><ul><li>Browse hundreds of language courses</li><li>Enroll in classes with expert teachers</li><li>Track your learning progress</li><li>Join live video sessions</li></ul><a href='{{app_url}}/courses' style='display: inline-block; background-color: #2563eb; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; margin-top: 20px;'>Browse Courses</a><p style='margin-top: 30px; color: #666;'>Happy learning!<br>The Language Learning Team</p></div></body></html>",
    "plain_content": "Hi {{name}},\n\nThank you for joining Language Learning Platform!\n\nYou can now browse courses, enroll in classes, and start learning.\n\nVisit: {{app_url}}/courses\n\nHappy learning!\nThe Language Learning Team",
    "active": 1
  }]
}
EOF

curl -s -X POST "https://api.sendgrid.com/v3/templates" \
    -H "Authorization: Bearer ${SENDGRID_API_KEY}" \
    -H "Content-Type: application/json" \
    -d @/tmp/template-welcome.json > /tmp/template-welcome-response.json

WELCOME_TEMPLATE_ID=$(cat /tmp/template-welcome-response.json | grep -o '"id":"[^"]*' | grep -o '[^"]*$' | head -1)
echo "✅ Welcome template created: $WELCOME_TEMPLATE_ID"

# Email verification template
echo "Creating email verification template..."
cat > /tmp/template-verification.json << 'EOF'
{
  "name": "Email Verification",
  "generation": "dynamic",
  "subject": "Verify your email address",
  "versions": [{
    "template_name": "Email Verification",
    "subject": "Verify your email for Language Learning Platform",
    "html_content": "<!DOCTYPE html><html><head><meta charset='UTF-8'></head><body style='font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;'><div style='background-color: #2563eb; color: white; padding: 30px; text-align: center;'><h1>Verify Your Email</h1></div><div style='background-color: #f9fafb; padding: 30px;'><p style='font-size: 16px;'>Hi {{name}},</p><p>Please verify your email address by clicking the button below:</p><a href='{{verification_url}}' style='display: inline-block; background-color: #10b981; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; margin-top: 20px;'>Verify Email</a><p style='margin-top: 30px; color: #666; font-size: 14px;'>Or copy this link: {{verification_url}}</p><p style='color: #999; font-size: 12px;'>This link expires in 24 hours.</p></div></body></html>",
    "plain_content": "Hi {{name}},\n\nPlease verify your email address by visiting:\n\n{{verification_url}}\n\nThis link expires in 24 hours.\n\nThe Language Learning Team",
    "active": 1
  }]
}
EOF

curl -s -X POST "https://api.sendgrid.com/v3/templates" \
    -H "Authorization: Bearer ${SENDGRID_API_KEY}" \
    -H "Content-Type: application/json" \
    -d @/tmp/template-verification.json > /tmp/template-verification-response.json

VERIFICATION_TEMPLATE_ID=$(cat /tmp/template-verification-response.json | grep -o '"id":"[^"]*' | grep -o '[^"]*$' | head -1)
echo "✅ Verification template created: $VERIFICATION_TEMPLATE_ID"

# Assignment due reminder
echo "Creating assignment due reminder template..."
cat > /tmp/template-assignment.json << 'EOF'
{
  "name": "Assignment Due Reminder",
  "generation": "dynamic",
  "subject": "Assignment due soon!",
  "versions": [{
    "template_name": "Assignment Due Reminder",
    "subject": "Reminder: {{assignment_title}} is due soon",
    "html_content": "<!DOCTYPE html><html><head><meta charset='UTF-8'></head><body style='font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;'><div style='background-color: #f59e0b; color: white; padding: 30px; text-align: center;'><h1>⏰ Assignment Due Soon</h1></div><div style='background-color: #f9fafb; padding: 30px;'><p style='font-size: 16px;'>Hi {{student_name}},</p><p>This is a friendly reminder that your assignment is due soon:</p><div style='background-color: white; padding: 20px; border-left: 4px solid #f59e0b; margin: 20px 0;'><h2 style='margin-top: 0;'>{{assignment_title}}</h2><p><strong>Course:</strong> {{course_name}}</p><p><strong>Due:</strong> {{due_date}}</p></div><a href='{{assignment_url}}' style='display: inline-block; background-color: #2563eb; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; margin-top: 20px;'>View Assignment</a><p style='margin-top: 30px; color: #666;'>Good luck!<br>The Language Learning Team</p></div></body></html>",
    "plain_content": "Hi {{student_name}},\n\nReminder: Your assignment is due soon!\n\nAssignment: {{assignment_title}}\nCourse: {{course_name}}\nDue: {{due_date}}\n\nView assignment: {{assignment_url}}\n\nGood luck!\nThe Language Learning Team",
    "active": 1
  }]
}
EOF

curl -s -X POST "https://api.sendgrid.com/v3/templates" \
    -H "Authorization: Bearer ${SENDGRID_API_KEY}" \
    -H "Content-Type: application/json" \
    -d @/tmp/template-assignment.json > /dev/null

echo "✅ Assignment reminder template created"

# Clean up temp files
rm -f /tmp/template-*.json

echo ""
echo "==================================="
echo "✅ SendGrid setup complete!"
echo "==================================="
echo ""
echo "📋 Template IDs created:"
echo "   Welcome: $WELCOME_TEMPLATE_ID"
echo "   Verification: $VERIFICATION_TEMPLATE_ID"
echo ""
echo "📝 Add these to your .env file:"
echo ""
echo "SENDGRID_API_KEY=$SENDGRID_API_KEY"
echo "SENDGRID_TEMPLATE_WELCOME=$WELCOME_TEMPLATE_ID"
echo "SENDGRID_TEMPLATE_VERIFICATION=$VERIFICATION_TEMPLATE_ID"
echo ""
echo "🔧 Next steps:"
echo "1. Verify your sender email in SendGrid dashboard"
echo "2. Configure SENDGRID_FROM_EMAIL in .env"
echo "3. Test sending emails"
echo ""
