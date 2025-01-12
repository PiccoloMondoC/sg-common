// sg-common/secrets/secrets.go
/*
   Secrets, also called secret keys, are crucial components in this application.
   During development, we will store and fetch secrets from environment variables.
   In production, we need to use a secret management service such as GCP Secret
   Manager.

   The secrets package makes fetching secrets an abstract operation, independent
   of whether the secretscome from environment variables or Secret Manager. This
   allows us to  easily switch from using environment variables during local
   development to using Secret Manager in the production environment.

   Implementation: Secrets is implemented as a shared library that can be used by
   any microservice that needs it. Implementation involves 3 steps: Develop a
   Secret Fetching Interface; Implement the Interface for Both Methods (fetching
   secrets from environment variables, and fetching secrets from Secret manager);
   Decide on the Method at Runtime.
*/
package secrets

import (
	"context"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

type SecretFetcher interface {
	GetSecret(key string) (string, error)
}

type EnvVarSecretFetcher struct{}

func (f *EnvVarSecretFetcher) GetSecret(key string) (string, error) {
	// Fetch secret from environment variable
	return os.Getenv(key), nil
}

type GcpSecretManagerFetcher struct {
	client    *secretmanager.Client
	projectID string
}

func NewGcpSecretManagerFetcher(projectID string, credentialsFile string) (*GcpSecretManagerFetcher, error) {
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret manager client: %v", err)
	}

	return &GcpSecretManagerFetcher{
		client:    client,
		projectID: projectID,
	}, nil
}

func (f *GcpSecretManagerFetcher) GetSecret(secretID string) (string, error) {
	ctx := context.Background()

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", f.projectID, secretID),
	}

	// Call the API.
	result, err := f.client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	// Return the secret payload as a string.
	return string(result.Payload.Data), nil
}

func GetFetcher() SecretFetcher {
	useSecretManager := os.Getenv("USE_SECRET_MANAGER")
	if useSecretManager == "true" {
		projectID := os.Getenv("GCP_PROJECT_ID")
		credentialsFile := os.Getenv("GCP_CREDENTIALS_FILE")
		fetcher, err := NewGcpSecretManagerFetcher(projectID, credentialsFile)
		if err != nil {
			log.Printf("Failed to create GcpSecretManagerFetcher: %v", err)
			log.Println("Falling back to EnvVarSecretFetcher")
		} else {
			return fetcher
		}
	}
	return &EnvVarSecretFetcher{}
}

/*
    To fetch the secrets in the application, we can use the code below:

	secretFetcher := secrets.GetFetcher()
	apiKey, err := secretFetcher.GetSecret("API_KEY") // eg. to a secret named "API_KEY"

*/
