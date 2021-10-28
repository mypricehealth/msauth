package msauth

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dghubble/sling"
)

// get graph client config from environment
var (
	// Optional: Azure AD Authentication Endpoint, defaults to msgraph.AzureADAuthEndpointGlobal: https://login.microsoftonline.com
	adAuthEndpoint string
	// Microsoft Azure AD tenant ID
	adTenantID string
	// Microsoft Azure AD Application ID
	adApplicationID string
	// Microsoft Azure AD Client Secret
	adClientSecret string
	// resource which is being requested access to
	resource string
)

func getEnvOrPanic(key string) string {
	var val = os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Expected %s to be set, but is empty", key))
	}
	return val
}

func TestMain(m *testing.M) {
	adTenantID = getEnvOrPanic("AD_TENANT_ID")
	adApplicationID = getEnvOrPanic("AD_APPLICATION_ID")
	adClientSecret = getEnvOrPanic("AD_CLIENT_SECRET")
	resource = "https://graph.microsoft.com"

	if adAuthEndpoint = os.Getenv("AD_AUTH_ENDPOINT"); adAuthEndpoint == "" {
		adAuthEndpoint = AzureADAuthEndpointGlobal
	}

	os.Exit(m.Run())
}

func TestNewGraphClient(t *testing.T) {
	if adAuthEndpoint != AzureADAuthEndpointGlobal {
		t.Skip("Skipping TestNewGraphClient because the endpoint is not the default - global - endpoint")
	}
	type args struct {
		tenantID      string
		applicationID string
		clientSecret  string
		resource      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "GraphClient from Environment-variables",
			args:    args{adTenantID, adApplicationID, adClientSecret, resource},
			wantErr: false,
		}, {
			name:    "GraphClient fail - wrong tenant ID",
			args:    args{"wrong tenant id", adApplicationID, adClientSecret, resource},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong application ID",
			args:    args{adTenantID, "wrong application id", adClientSecret, resource},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong client secret",
			args:    args{adTenantID, adApplicationID, "wrong client secret", resource},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong resource",
			args:    args{adTenantID, adApplicationID, adClientSecret, "invalid URL"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.tenantID, tt.args.applicationID, tt.args.clientSecret, tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGraphClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewGraphClientWithCustomEndpoint(t *testing.T) {
	type args struct {
		tenantID            string
		applicationID       string
		clientSecret        string
		azureADAuthEndpoint string
		resource            string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "GraphClient from Environment-variables",
			args:    args{adTenantID, adApplicationID, adClientSecret, adAuthEndpoint, resource},
			wantErr: false,
		}, {
			name:    "GraphClient fail - wrong tenant ID",
			args:    args{"wrong tenant id", adApplicationID, adClientSecret, adAuthEndpoint, resource},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong application ID",
			args:    args{adTenantID, "wrong application id", adClientSecret, adAuthEndpoint, resource},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong client secret",
			args:    args{adTenantID, adApplicationID, "wrong client secret", adAuthEndpoint, resource},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong Azure AD Authentication Endpoint",
			args:    args{adTenantID, adApplicationID, adClientSecret, "completely invalid URL", resource},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong Azure AD Authentication Endpoint",
			args:    args{adTenantID, adApplicationID, adClientSecret, "https://iguess-this-does-not.exist.in.this.dksfowe3834myksweroqiwer.world", resource},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong Service Root Endpoint",
			args:    args{adTenantID, adApplicationID, adClientSecret, adAuthEndpoint, "invalid URL"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWithCustomEndpoint(tt.args.tenantID, tt.args.applicationID, tt.args.clientSecret, tt.args.azureADAuthEndpoint, tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGraphClientWithCustomEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMakeAPICall(t *testing.T) {
	var user = struct {
		Context           string `json:"@odata.context"`
		DisplayName       string `json:"displayName"`
		GivenName         string `json:"givenName"`
		JobTitle          string `json:"jobTitle"`
		Mail              string `json:"mail"`
		MobilePhone       string `json:"mobilePhone"`
		OfficeLocation    string `json:"officeLocation"`
		PreferredLanguage string `json:"preferredLanguage"`
		Surname           string `json:"surname"`
		UserPrincipalName string `json:"userPrincipalName"`
		ID                string `json:"id"`
	}{}

	graphClient, err := New(adTenantID, adApplicationID, adClientSecret, resource)
	if err != nil {
		t.Fatal(err)
	}

	graphClient.SetDefaults(func(client *sling.Sling) *sling.Sling {
		return client.Base(fmt.Sprintf("%s/%s/", resource, "v1.0"))
	})
	graphClient.token.ExpiresOn = time.Now().Add(-time.Second) // expire token so it'll get new token

	// This will always fail with an application authentication flow as it is only allowed in delegated flow
	if err := graphClient.Do(context.Background(), "GET", "me", nil, nil, nil, &user, nil); err == nil {
		t.Error(err)
	}
}
