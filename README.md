# msauth
A Go client library for easy Microsoft authentication

This package is a heavily modified Microsoft Identity Platform client borrowing ideas and code from https://github.com/open-networks/go-msgraph package. In order to enable sharing authentication with a `msgraph` package, a `mspowerbi` package and any other Microsoft package which uses the same authentication flow, this has been built into its own package.

Other projects of interest:
1. https://github.com/AzureAD/microsoft-authentication-library-for-go - a Microsoft supported Go authentication package. Still in preview.
2. https://github.com/open-networks/go-msgraph - a Microsoft Graph client which has the auth client built in.
