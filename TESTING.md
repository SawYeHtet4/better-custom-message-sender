# Testing Guide

## Overview

This guide explains how to test the Custom Message Sender Lambda function locally and in AWS.

## Important Note About Testing

⚠️ **The Lambda function requires AWS KMS to decrypt the OTP codes from Cognito.** This means:
- Local testing requires valid AWS credentials with KMS permissions
- The `code` field in test events must be an actual encrypted code from Cognito
- You cannot use dummy/mock codes for end-to-end testing

## Prerequisites

1. **AWS Credentials**: Configure AWS CLI with credentials
   ```bash
   aws configure
   ```

2. **Environment Variables**: Create `env.json` in the project root
   ```json
   {
     "CustomMessageSenderFunction": {
       "APP": "Yoma Fleet Better",
       "VERSION": "1.0.0",
       "ENV": "development",
       "EMAIL_FROM": "noreply@example.com",
       "MAIL_TRANSPORT": "smtp",
       "MAIL_TRAP_HOST": "smtp.mailtrap.io",
       "MAIL_TRAP_PORT": "2525",
       "MAIL_TRAP_USER": "your_mailtrap_user",
       "MAIL_TRAP_PASSWORD": "your_mailtrap_password",
       "BLACK_LIST_EMAILS": "",
       "SMS_HOST": "https://your-sms-api.com",
       "SMS_AUTH_KEY": "your_sms_api_key",
       "KEY_ID": "arn:aws:kms:us-east-1:123456789012:key/your-key-id",
       "SQS_QUEUE_URL": ""
     }
   }
   ```

## Test Event Files

The following test event files are available in the `events/` directory:

| File | Trigger | Description |
|------|---------|-------------|
| `signup.json` | CustomEmailSender_SignUp | Email verification on signup |
| `forgot_password.json` | CustomEmailSender_ForgotPassword | Password reset email |
| `admin_create_user.json` | CustomEmailSender_AdminCreateUser | Admin user creation email |
| `sms_signup.json` | CustomSMSSender_SignUp | SMS verification on signup |

### Updating Test Events

To test with real Cognito events:

1. **Set up Cognito User Pool** with custom message sender
2. **Trigger an actual event** (signup, password reset, etc.)
3. **Capture the Lambda input** from CloudWatch Logs
4. **Copy the encrypted code** to your test event file

Replace `YOUR_ENCRYPTED_CODE_HERE_FROM_COGNITO` with the actual encrypted code from Cognito.

## Local Testing with SAM

### 1. Build the Lambda

```bash
sam build
```

### 2. Test Email Signup Flow

```bash
sam local invoke \
  --env-vars env.json \
  --event events/signup.json \
  CustomMessageSenderFunction
```

### 3. Test Password Reset Flow

```bash
sam local invoke \
  --env-vars env.json \
  --event events/forgot_password.json \
  CustomMessageSenderFunction
```

### 4. Test Admin Create User Flow

```bash
sam local invoke \
  --env-vars env.json \
  --event events/admin_create_user.json \
  CustomMessageSenderFunction
```

### 5. Test SMS Flow

```bash
sam local invoke \
  --env-vars env.json \
  --event events/sms_signup.json \
  CustomMessageSenderFunction
```

## Testing Strategy

### Option 1: Unit Testing (Without AWS Dependencies)

For testing business logic without AWS dependencies, you can create unit tests:

```bash
# Run unit tests
make test
```

Create test files in `pkg/usecase/` and mock the interfaces.

### Option 2: Integration Testing (With MailTrap)

1. **Set MAIL_TRANSPORT** to `smtp` in env.json
2. **Configure MailTrap** credentials
3. **Run local invoke** with test events
4. **Check MailTrap inbox** for received emails

### Option 3: End-to-End Testing (Deployed Lambda)

1. **Deploy to AWS**:
   ```bash
   sam deploy --guided
   ```

2. **Configure Cognito**:
   - Go to your Cognito User Pool
   - Navigate to "Triggers" → "Custom message sender"
   - Select your deployed Lambda function
   - Set KMS key for encryption

3. **Test via Cognito**:
   - Trigger signup in your application
   - Check CloudWatch Logs for Lambda execution
   - Verify email/SMS delivery

## Verification Checklist

- [ ] Lambda builds successfully (`sam build`)
- [ ] Environment variables configured in `env.json`
- [ ] AWS credentials configured
- [ ] KMS key permissions granted to Lambda execution role
- [ ] Email transport configured (SMTP/SES/SQS)
- [ ] SMS API configured (if using SMS)
- [ ] Test events updated with real encrypted codes
- [ ] Lambda invokes successfully locally
- [ ] Emails received in MailTrap/inbox
- [ ] SMS messages delivered (if configured)

## Common Issues

### 1. KMS Decryption Failure

**Error**: `AccessDeniedException` or decryption fails

**Solution**:
- Ensure Lambda execution role has `kms:Decrypt` permission
- Verify KEY_ID matches the Cognito User Pool KMS key
- Check that the encrypted code is from a real Cognito event

### 2. Email Not Sending

**Error**: SMTP connection failed or SES access denied

**Solution**:
- For SMTP: Verify MailTrap credentials
- For SES: Verify sender email address
- For SES: Check SES sending limits and sandbox status
- Check EMAIL_FROM is a verified email in SES

### 3. Template Rendering Error

**Error**: Template not found or parsing error

**Solution**:
- Ensure `assets/templates/email/` directory exists
- Verify all HTML template files are present
- Check template syntax for Go template errors

### 4. SMS API Failure

**Error**: SMS sending failed

**Solution**:
- Verify SMS_HOST is accessible
- Check SMS_AUTH_KEY is valid
- Ensure phone number format is correct
- Check SMS API response in logs

## Debugging

### Enable Verbose Logging

The Lambda uses structured JSON logging. Check CloudWatch Logs for:
- Event details
- Decrypted code (in development only)
- Trigger source
- Email/SMS send status

### Local Debugging

```bash
# Build with debug info
sam build --debug

# Invoke with verbose output
sam local invoke \
  --env-vars env.json \
  --event events/signup.json \
  --debug \
  CustomMessageSenderFunction
```

## Testing Checklist by Transport

### SMTP/MailTrap Testing
- [ ] MAIL_TRANSPORT set to `smtp`
- [ ] MailTrap credentials configured
- [ ] Test email signup event
- [ ] Verify email in MailTrap inbox
- [ ] Check HTML rendering

### AWS SES Testing
- [ ] MAIL_TRANSPORT set to `SES`
- [ ] SES domain/email verified
- [ ] Lambda has SES send permissions
- [ ] EMAIL_FROM is verified in SES
- [ ] Test email event
- [ ] Check SES sending statistics

### AWS SQS Testing
- [ ] MAIL_TRANSPORT set to `SQS`
- [ ] SQS queue created
- [ ] Lambda has SQS send permissions
- [ ] SQS_QUEUE_URL configured
- [ ] Test email event
- [ ] Verify message in SQS queue
- [ ] Verify downstream processor handles message

## Next Steps

After successful local testing:

1. Deploy to development environment
2. Test with real Cognito events
3. Monitor CloudWatch Logs
4. Validate email/SMS delivery
5. Deploy to production

## Support

For issues or questions:
- Check CloudWatch Logs for detailed error messages
- Verify all configuration values
- Ensure AWS permissions are correctly set
- Review the code in `pkg/usecase/message_sender.go` for business logic
