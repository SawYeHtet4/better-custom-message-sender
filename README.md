# Better Custom Message Sender

AWS Lambda function for AWS Cognito Custom Message Sender that handles both Email and SMS notifications.

## Overview

This Lambda function is designed to work with AWS Cognito's custom message sender feature. It intercepts Cognito's authentication flows (signup, password reset, MFA, etc.) and sends customized email and SMS messages to users.

## Features

- **Email Support**: Multiple transport options
  - SMTP (MailTrap for development/testing)
  - AWS SES for production
  - AWS SQS for asynchronous processing
- **SMS Support**: Integration with external SMS API
- **Security**: Uses AWS KMS for decrypting OTP codes from Cognito
- **Customizable Templates**: HTML email templates with responsive design
- **Event Support**: Handles multiple Cognito trigger sources:
  - Sign Up
  - Forgot Password
  - Resend Code
  - Update User Attribute
  - Verify User Attribute
  - Admin Create User
  - Account Takeover Notification

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
│   │   └── message_sender/     # Email and SMS sending
│   └── usecase/                # Business logic
├── assets/
│   └── templates/
│       └── email/              # HTML email templates
└── events/                     # Sample event payloads for testing
```

## Prerequisites

- Go 1.24 or later
- AWS CLI configured with appropriate credentials
- AWS SAM CLI (for local testing and deployment)
- An AWS account with:
  - Cognito User Pool configured
  - KMS key for encryption
  - (Optional) SES verified domain or SQS queue
  - (Optional) SMS API credentials

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
   cp .env_example .env
   # Edit .env with your configuration
   ```

4. **Environment Variables**:
   - `APP`: Application name
   - `VERSION`: Application version
   - `ENV`: Environment (development, staging, production)
   - `EMAIL_FROM`: Sender email address
   - `MAIL_TRANSPORT`: Email transport method (smtp, SES, SQS)
   - `MAIL_TRAP_HOST`: SMTP host for MailTrap
   - `MAIL_TRAP_PORT`: SMTP port
   - `MAIL_TRAP_USER`: SMTP username
   - `MAIL_TRAP_PASSWORD`: SMTP password
   - `BLACK_LIST_EMAILS`: Comma-separated list of blocked emails
   - `SMS_HOST`: SMS API endpoint
   - `SMS_AUTH_KEY`: SMS API authentication key
   - `KEY_ID`: AWS KMS key ID for decryption
   - `SQS_QUEUE_URL`: (Optional) SQS queue URL for email processing

## Building

```bash
# Build with SAM
make build

# Or directly with Go
cd cmd/function/message_sender
GOOS=linux GOARCH=amd64 go build -o main
```

## Testing Locally

1. **Create env.json** with your environment variables:
   ```json
   {
     "CustomMessageSenderFunction": {
       "APP": "Better Custom Message Sender",
       "EMAIL_FROM": "noreply@example.com",
       "MAIL_TRANSPORT": "smtp",
       ...
     }
   }
   ```

2. **Invoke locally with SAM**:
   ```bash
   sam build
   sam local invoke --env-vars env.json --event events/singup.json CustomMessageSenderFunction
   ```

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
2. Select "Custom message sender" trigger
3. Choose this Lambda function
4. Configure KMS key for encryption

## Email Templates

HTML email templates are located in `assets/templates/email/`. The templates use Go's `html/template` package.

Available templates:
- `index.html`: Base template with header and footer
- `verify.html`: Email verification template
- `forgot_password.html`: Password reset template
- `admin_create_user.html`: Admin user creation template
- `welcome.html`: Welcome message template

## Customization

### Branding

Edit the email templates in `assets/templates/email/` to match your branding:
- Update colors in the CSS variables
- Change logo URLs
- Modify footer links and content

### SMS Messages

SMS message content can be modified in `pkg/usecase/message_sender.go` in the `sendSMS` function.

### Subject Lines

Email subjects are built in the `buildSubject` function in `pkg/usecase/message_sender.go`.

## Troubleshooting

### Common Issues

1. **KMS Decryption Errors**: Ensure the Lambda execution role has permissions to use the KMS key
2. **Email Not Sending**: Check SES domain verification and sending limits
3. **Template Errors**: Verify template syntax and ensure all referenced templates exist

### Logs

The function uses structured JSON logging. Check CloudWatch Logs for detailed execution logs.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

[Your License Here]

## Support

For issues and questions, please open an issue in the repository.
