.PHONY: build

functionName := function-name
zipFileName := main.zip

build:
	sam build

invoke:
	sam build && sam local invoke --env-vars env.json --event events/signup.json CustomMessageSenderFunction

invoke-sms:
	sam build && sam local invoke --env-vars env.json --event events/sms_signup.json CustomMessageSenderFunction

invoke-forgot:
	sam build && sam local invoke --env-vars env.json --event events/forgot_password.json CustomMessageSenderFunction

invoke-admin:
	sam build && sam local invoke --env-vars env.json --event events/admin_create_user.json CustomMessageSenderFunction

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
