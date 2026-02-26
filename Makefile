.PHONY: build

functionName := function-name
zipFileName := main.zip

build:
	sam build

invoke-sms-signup:
	sam build && sam local invoke --env-vars env.json --event events/sms_signup.json CustomMessageSenderFunction

invoke-sms-forgot:
	sam build && sam local invoke --env-vars env.json --event events/sms_forgot_password.json CustomMessageSenderFunction

invoke-sms-admin:
	sam build && sam local invoke --env-vars env.json --event events/sms_admin_create.json CustomMessageSenderFunction

invoke-sms-resend:
	sam build && sam local invoke --env-vars env.json --event events/sms_resend_code.json CustomMessageSenderFunction

test:
	go test -v -cover \
	 	./pkg/domain/e/... \
		./pkg/domain/models/... \
		./pkg/usecase/...

makeZip:
	./scripts/make_zip.sh

updateCode:
	aws lambda --profile ${profile} update-function-code \
    --function-name  $(functionName) \
    --zip-file fileb://$(zipFileName) && rm -rf $(zipFileName)
