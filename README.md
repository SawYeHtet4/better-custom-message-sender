# Better Custom Message Sender

AWS Lambda function for AWS Cognito Custom Message Sender that handles SMS notifications.

**Author**: Saw Ye Htet
**Email**: sawyehtet@yomafleet.com

## Overview

This Lambda function is designed to work with AWS Cognito's custom SMS sender feature. It intercepts Cognito's authentication flows (signup, password reset, MFA, etc.) and sends customized SMS messages to users.

## Features

- **SMS Support**: Integration with external SMS API
- **Security**: Uses AWS KMS for decrypting OTP codes from Cognito
- **Custom Messages**: Different SMS templates for different trigger types
- **Event Support**: Handles multiple Cognito SMS trigger sources:
  - Sign Up
  - Forgot Password
  - Resend Code
  - Update User Attribute
  - Verify User Attribute
  - Admin Create User

## Architecture

The project follows clean architecture principles:

```
.
├── cmd/
│   └── function/
│       └── message_sender/      # Lambda handler entry point
├── config/                      # Configuration management
├── pkg/
│   ├── domain/                  # Domain contracts and error types
│   │   ├── contract/           # Interface definitions
│   │   └── e/                  # Domain error types
│   ├── interface/              # Infrastructure implementations
│   │   ├── encrypter/          # AWS KMS encryption/decryption
│   │   ├── logger/             # Structured logging
│   │   └── message_sender/     # SMS sending
│   └── usecase/                # Business logic
└── events/                     # Sample event payloads for testing
```

## Prerequisites

- Go 1.24 or later
- AWS CLI configured with appropriate credentials
- AWS SAM CLI (for local testing and deployment)
- An AWS account with:
  - Cognito User Pool configured
  - KMS key for encryption
  - SMS API credentials

## Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd better-custom-message-sender
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Configure environment variables**:
   ```bash
   cp env.json.example env.json
   # Edit env.json with your configuration
   ```

4. **Environment Variables**:
   - `APP`: Application name
   - `VERSION`: Application version
   - `ENV`: Environment (development, staging, production)
   - `SMS_HOST`: SMS API endpoint
   - `SMS_AUTH_KEY`: SMS API authentication key
   - `KEY_ID`: AWS KMS key ARN for decryption

## SMS Messages

The Lambda sends different messages based on trigger type:

### Sign Up
```
Welcome to Yoma Fleet Better! Your verification code is {code}. DO NOT share it with anyone.
```

### Forgot Password
```
Your Yoma Fleet Better password reset code is {code}. DO NOT share it with anyone.
```

### Resend Code
```
Your Yoma Fleet Better verification code is {code}. DO NOT share it with anyone.
```

### Admin Create User
```
Welcome to Yoma Fleet Better! Your temporary password is {code}. Please change it after login.
```

## Building

```bash
# Build with SAM
make build

# Or directly with Go
cd cmd/function/message_sender
GOOS=linux GOARCH=amd64 go build -o main
```

## Testing Locally

1. **Create env.json** with your environment variables (see env.json.example)

2. **Test different SMS triggers**:

```bash
# Test SMS signup
make invoke-sms-signup

# Test SMS forgot password
make invoke-sms-forgot

# Test SMS admin create user
make invoke-sms-admin

# Test SMS resend code
make invoke-sms-resend
```

**Important Note About Test Events**:
- The event files contain `"YOUR_ENCRYPTED_CODE_HERE_FROM_COGNITO"` as a placeholder
- This is **not** a real encrypted code and will cause a base64 decryption error
- To get real encrypted codes, you need to:
  1. Deploy the Lambda to AWS
  2. Configure it with a Cognito User Pool
  3. Trigger an actual event (signup, password reset, etc.)
  4. Capture the real encrypted code from CloudWatch Logs
  5. Replace the placeholder in your test event file
- Local testing without real codes will verify the build and event flow, but will fail at decryption

## Running Tests

```bash
make test
```

## Deployment

### Using SAM

1. **Build**:
   ```bash
   sam build
   ```

2. **Deploy**:
   ```bash
   sam deploy --guided
   ```

3. **Update Lambda Code** (for existing function):
   ```bash
   make makeZip
   make updateCode profile=your-aws-profile functionName=your-function-name
   ```

## Cognito Integration

1. In your Cognito User Pool settings, navigate to "Triggers"
2. Select "Custom SMS sender" trigger
3. Choose this Lambda function
4. Configure KMS key for encryption

## SMS API Integration

The Lambda expects an SMS API with the following specifications:

**Endpoint**: `POST /api/v2/send`

**Request Headers**:
```
Content-Type: application/json
Authorization: Bearer {SMS_AUTH_KEY}
```

**Request Body**:
```json
{
  "to": "+1234567890",
  "message": "Your OTP code...",
  "sender": "Yoma Fleet"
}
```

**Response**:
```json
{
  "status": true
}
```

Compatible with most SMS providers (Twilio, Vonage, MessageBird, etc.)

## Troubleshooting

### Common Issues

1. **KMS Decryption Errors**: Ensure the Lambda execution role has permissions to use the KMS key
2. **SMS Not Sending**: Check SMS_HOST and SMS_AUTH_KEY configuration
3. **Phone Number Format**: Must include country code (e.g., +1234567890)

### Logs

The function uses structured JSON logging. Check CloudWatch Logs for detailed execution logs.

## Event Files

Test event files are available in the `events/` directory:

- `sms_signup.json` - SMS verification on signup
- `sms_forgot_password.json` - Password reset SMS
- `sms_admin_create.json` - Admin user creation SMS
- `sms_resend_code.json` - Resend verification code

### Using Test Events

**Important**: All test event files contain placeholder values that need to be replaced with real data:

1. **Encrypted Code**: Replace `"YOUR_ENCRYPTED_CODE_HERE_FROM_COGNITO"` with a real encrypted code from Cognito
   - This can only be obtained from actual Cognito events after deployment
   - Without a real encrypted code, local testing will fail at KMS decryption

2. **Phone Number**: Replace `"+1234567890"` with your test phone number
   - Must include country code (e.g., `+959xxxxxxxxx` for Myanmar)

3. **User Pool ID**: Replace `"us-east-1_EXAMPLE123"` with your actual Cognito User Pool ID

4. **KMS Key**: Ensure `KEY_ID` in `env.json` matches your Cognito User Pool's KMS key ARN

**Expected Behavior**:
- ✅ With placeholder codes: Lambda will fail at decryption (expected)
- ✅ With real codes from Cognito: Lambda will decrypt and send SMS successfully

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

[Your License Here]

## Support

For issues and questions, please contact the author or open an issue in the repository.
