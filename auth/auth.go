package auth

import "github.com/lestrrat-go/jwx/jwk"

// OpenIDConfig well known public OpenID config.
type OpenIDConfig struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	JWKSURI                          string   `json:"jwks_uri"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	RequireRequestURIRegistration    string   `json:"require_request_uri_registration"`
	UserinfoEndpoint                 string   `json:"userinfo_endpoint"`
	RegistrationEndpoint             string   `json:"registration_endpoint"`
	ScopesSupported                  []string `json:"scopes_supported"`
	ClaimsSupported                  []string `json:"claims_supported"`
}

// JWKS contains a list of public keys
type JWKS struct {
	Keys []jwk.Key `json:"keys"`
}
