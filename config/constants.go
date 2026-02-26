package config

const (
	SMSMessageType = "SMS"

	ServiceHeaderKey = "X-SERVICE-Id"
)

const (
	// 200 - Success
	StatusCode200Success            = "R0000211"
	StatusCode200SuccessMessage     = "ok"
	StatusCode200SuccessDescription = "data is ok"

	StatusCode200ServiceToServiceSuccess            = "R0000212"
	StatusCode200ServiceToServiceSuccessMessage     = "ok"
	StatusCode200ServiceToServiceSuccessDescription = "service to service is ok"

	StatusCode200ExternalSuccess            = "R0000213"
	StatusCode200ExternalSuccessMessage     = "ok"
	StatusCode200ExternalSuccessDescription = "external integration is ok"

	// 201 - Created
	StatusCode201Created            = "R0000221"
	StatusCode201CreatedMessage     = "created"
	StatusCode201CreatedDescription = "data is created"

	StatusCode201ServiceToServiceCreated            = "R0000222"
	StatusCode201ServiceToServiceCreatedMessage     = "created"
	StatusCode201ServiceToServiceCreatedDescription = "service to service is created"

	StatusCode201ExternalCreated            = "R0000223"
	StatusCode201ExternalCreatedMessage     = "created"
	StatusCode201ExternalCreatedDescription = "external integration is created"

	// 400 - Bad Request
	StatusCode400AlreadyExists            = "R0000411"
	StatusCode400AlreadyExistsMessage     = "Bad Request"
	StatusCode400AlreadyExistsDescription = "it already exists in the system."

	StatusCode400InvalidCharacters            = "R0000412"
	StatusCode400InvalidCharactersMessage     = "Bad Request"
	StatusCode400InvalidCharactersDescription = "does not have valid characters"

	StatusCode400InvalidFormat            = "R0000413"
	StatusCode400InvalidFormatMessage     = "Bad Request"
	StatusCode400InvalidFormatDescription = "request is not in a valid given format"

	StatusCode400RequestIdExceeds            = "R0000414"
	StatusCode400RequestIdExceedsMessage     = "Bad Request"
	StatusCode400RequestIdExceedsDescription = "length of RequestId exceeds"

	// 401 - Unauthorized
	StatusCode401NotProvisioned            = "R0000421"
	StatusCode401NotProvisionedMessage     = "Unauthorize"
	StatusCode401NotProvisionedDescription = "not provisioned to access the API"

	StatusCode401PolicyScopeMismatch            = "R0000422"
	StatusCode401PolicyScopeMismatchMessage     = "Unauthorize"
	StatusCode401PolicyScopeMismatchDescription = "Policy and scope is not match"

	StatusCode401AccountRestricted            = "R0000423"
	StatusCode401AccountRestrictedMessage     = "Unauthorize"
	StatusCode401AccountRestrictedDescription = "Account is restricted"

	// 403 - Forbidden
	StatusCode403ResourceForbidden            = "R0000441"
	StatusCode403ResourceForbiddenMessage     = "Forbidden"
	StatusCode403ResourceForbiddenDescription = "requested resource is forbidden"

	StatusCode403InvalidClient            = "R0000442"
	StatusCode403InvalidClientMessage     = "Forbidden"
	StatusCode403InvalidClientDescription = "invalid client"

	StatusCode403InvalidRelationship            = "R0000443"
	StatusCode403InvalidRelationshipMessage     = "Forbidden"
	StatusCode403InvalidRelationshipDescription = "relationship resource is invalid"

	// 404 - Not Found
	StatusCode404ResourceNotFound            = "R0000451"
	StatusCode404ResourceNotFoundMessage     = "Not Found"
	StatusCode404ResourceNotFoundDescription = "resource not found"

	StatusCode404RequestIdNotFound            = "R0000452"
	StatusCode404RequestIdNotFoundMessage     = "Not Found"
	StatusCode404RequestIdNotFoundDescription = "RequestId not found"

	StatusCode404RelationalResourceNotFound            = "R0000453"
	StatusCode404RelationalResourceNotFoundMessage     = "Not Found"
	StatusCode404RelationalResourceNotFoundDescription = "Relational resource not found"

	// 405 - Method Not Allowed
	StatusCode405MethodNotAllowed            = "R0000461"
	StatusCode405MethodNotAllowedMessage     = "Method not allow"
	StatusCode405MethodNotAllowedDescription = "Method not allow"

	StatusCode405InvalidMethod            = "R0000462"
	StatusCode405InvalidMethodMessage     = "Method not allow"
	StatusCode405InvalidMethodDescription = "invalid given method"

	// 422 - Unprocessable Content
	StatusCode422SingleInvalidData            = "R0000471"
	StatusCode422SingleInvalidDataMessage     = "Unprocessable Content"
	StatusCode422SingleInvalidDataDescription = "single invalid request data"

	StatusCode422MultipleInvalidData            = "R0000472"
	StatusCode422MultipleInvalidDataMessage     = "Unprocessable Content"
	StatusCode422MultipleInvalidDataDescription = "multiple invalid request data"

	// 500 - Server Error
	StatusCode500ServerError            = "R0000511"
	StatusCode500ServerErrorMessage     = "Server Error"
	StatusCode500ServerErrorDescription = "Server side error"

	StatusCode500InternalIntegrationError            = "R0000512"
	StatusCode500InternalIntegrationErrorMessage     = "Server Error"
	StatusCode500InternalIntegrationErrorDescription = "internal side Integration error"

	StatusCode500ExternalIntegrationError            = "R0000513"
	StatusCode500ExternalIntegrationErrorMessage     = "Server Error"
	StatusCode500ExternalIntegrationErrorDescription = "external side integration error"

	// 503 - Service Unavailable
	StatusCode503ServiceUnavailable            = "R0000521"
	StatusCode503ServiceUnavailableMessage     = "Service Unavailable"
	StatusCode503ServiceUnavailableDescription = "Service Unavailable"

	StatusCode503OtherUnavailable            = "R0000522"
	StatusCode503OtherUnavailableMessage     = "Service Unavailable"
	StatusCode503OtherUnavailableDescription = "others unavailable"
)
