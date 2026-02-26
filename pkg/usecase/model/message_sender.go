package model

import "github.com/aws/aws-lambda-go/events"

type CognitoEventUserPoolsCustomMessage struct {
	events.CognitoEventUserPoolsHeader
	Request  CognitoEventUserPoolsCustomMessageRequest         `json:"request"`
	Response events.CognitoEventUserPoolsCustomMessageResponse `json:"response"`
}

// CognitoEventUserPoolsCustomMessageRequest contains the request portion of a CustomMessage event
type CognitoEventUserPoolsCustomMessageRequest struct {
	UserAttributes    map[string]interface{} `json:"userAttributes"`
	CodeParameter     string                 `json:"code"`
	UsernameParameter string                 `json:"usernameParameter"`
	ClientMetadata    map[string]string      `json:"clientMetadata"`
}
