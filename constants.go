package msauth

const (
	// Azure AD authentication endpoint "Global". Used to acquire a token for the ms graph API connection.
	//
	// Microsoft Documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	AzureADAuthEndpointGlobal string = "https://login.microsoftonline.com"

	// Azure AD authentication endpoint "Germany". Used to acquire a token for the ms graph API connection.
	//
	// Microsoft Documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	AzureADAuthEndpointGermany string = "https://login.microsoftonline.de"

	// Azure AD authentication endpoint "US Government". Used to acquire a token for the ms graph API connection.
	//
	// Microsoft Documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	AzureADAuthEndpointUSGov string = "https://login.microsoftonline.us"

	// Azure AD authentication endpoint "China by 21 Vianet". Used to acquire a token for the ms graph API connection.
	//
	// Microsoft Documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	AzureADAuthEndpointChina string = "https://login.partner.microsoftonline.cn"
)

// APIVersion represents the APIVersion of msauth used by this implementation
const APIVersion string = "v1.0"
