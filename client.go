// Package msauth provides an authenticated http Client which has been authenticated using the Microsoft
// identity platform and OAuth 2.0 authorization flow. The resulting client can be used for the
// Microsoft Graph API, Power BI and other API's.
//
// See: https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow
package msauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dghubble/sling"
)

type Client struct {
	tenantID      string // See https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal#get-tenant-id
	applicationID string // See https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal#get-application-id-and-authentication-key
	clientSecret  string // See https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal#get-application-id-and-authentication-key

	token         Token // the current token to be used
	doer          sling.Doer
	defaultSetter func(client *sling.Sling) *sling.Sling

	// Microsoft Identity Platform URL in use. For available endpoints see https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	azureADAuthEndpoint string
	// Resource string. This usually is the URL for the service being requested access to
	resource string
}

type authError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCodes       []int  `json:"error_codes"`
	Timestamp        string `json:"timestamp"`
	TraceID          string `json:"trace_id"`
	CorrelationID    string `json:"correlation_id"`
	ErrorURI         string `json:"error_uri"`
}

// New creates a new AuthClient instance with the given parameters
// and grabs a token. Returns an error if the token cannot be initialized. The
// default Microsoft Identity Platfrom URL is used.
func New(tenantID, applicationID, clientSecret, resource string) (*Client, error) {
	return NewWithCustomEndpoint(tenantID, applicationID, clientSecret, AzureADAuthEndpointGlobal, resource)
}

// NewWithCustomEndpoint creates a new Microsoft Identity Platform client instance with the
// given parameters and tries to get a valid token. All available public endpoints
// for azureADAuthEndpoint and serviceRootEndpoint are available via msauth.azureADAuthEndpoint*
//
// For available endpoints from Microsoft, see documentation:
//   * Authentication Endpoints: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
func NewWithCustomEndpoint(tenantID, applicationID, clientSecret string, azureADAuthEndpoint, resource string) (*Client, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	g := Client{
		tenantID:            tenantID,
		applicationID:       applicationID,
		clientSecret:        clientSecret,
		azureADAuthEndpoint: azureADAuthEndpoint,
		doer:                httpClient,
		resource:            resource,
	}
	return &g, g.refreshToken()
}

// refreshToken refreshes the current Token. Grabs a new one and saves it within the GraphClient instance
func (g *Client) refreshToken() error {
	if g.tenantID == "" {
		return fmt.Errorf("tenant ID is empty")
	}

	path := g.azureADAuthEndpoint + fmt.Sprintf("/%s/oauth2/token", g.tenantID)

	payload := url.Values{}
	payload.Add("grant_type", "client_credentials")
	payload.Add("client_id", g.applicationID)
	payload.Add("client_secret", g.clientSecret)
	payload.Add("resource", g.resource)

	var newToken Token
	errResult := authError{}
	_, err := sling.New().
		Doer(g.doer).
		Post(path).
		BodyForm(payload).
		Receive(&newToken, &errResult)
	if err != nil || errResult.Error != "" {
		if err != nil {
			return err
		}
		return fmt.Errorf(fmt.Sprintf("%s: %s", errResult.Error, errResult.ErrorDescription))
	}

	g.token = newToken
	return nil
}

func (g *Client) SetDefaults(setter func(client *sling.Sling) *sling.Sling) {
	g.defaultSetter = setter
}

func (g *Client) GetClient() (*sling.Sling, error) {
	if g.token.WantsToBeRefreshed() { // Token not valid anymore?
		err := g.refreshToken()
		if err != nil {
			return nil, err
		}
	}

	client := sling.New().
		Doer(g.doer).
		Set("Authorization", g.token.GetAccessToken())

	if g.defaultSetter != nil {
		client = g.defaultSetter(client)
	}
	return client, nil
}

func (g *Client) Do(ctx context.Context, method, path string, headers http.Header, urlParams url.Values, bodyJSON, result, errResult interface{}) error {
	client, err := g.GetClient()
	if err != nil {
		return err
	}
	_, err = client.
		Path(path).
		Method(method).
		AddHeaders(headers).
		QueryValues(urlParams).
		BodyJSON(bodyJSON).
		ReceiveWithContext(ctx, result, errResult)
	return err
}
