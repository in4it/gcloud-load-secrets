package secrets

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Secret struct {
	Name    string
	ID      string
	Version string
	Payload string
}

type ReadSecrets struct {
	Client    *secretmanager.Client
	ProjectID string
}

func NewReadSecrets() (*ReadSecrets, error) {
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		panic(err)
	}
	return &ReadSecrets{
		Client:    c,
		ProjectID: credentials.ProjectID,
	}, nil
}

func (r *ReadSecrets) ListSecrets(secretsPrefix, secretsLabel string) ([]Secret, error) {
	ctx := context.Background()

	req := &secretmanagerpb.ListSecretsRequest{
		Parent: fmt.Sprintf("projects/%s", r.ProjectID),
	}
	it := r.Client.ListSecrets(ctx, req)
	secrets := []Secret{}
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return secrets, err
		}
		secretElements := strings.Split(resp.Name, "/")
		if len(secretElements) < 4 {
			return secrets, fmt.Errorf("secret name in unexpected format: %s", resp.Name)
		}
		secretName := strings.Join(secretElements[3:], "/")
		if strings.HasPrefix(secretName, secretsPrefix) && r.MatchLabel(secretsLabel, resp.Labels) {
			secrets = append(secrets, Secret{ID: resp.Name, Name: secretName})
		}
	}
	return secrets, nil
}

func (r *ReadSecrets) GetSecretsValue(secrets []Secret) ([]Secret, error) {
	ctx := context.Background()

	for k, secret := range secrets {
		req := &secretmanagerpb.AccessSecretVersionRequest{
			Name: fmt.Sprintf("%s/versions/latest", secret.ID),
		}
		resp, err := r.Client.AccessSecretVersion(ctx, req)
		if err != nil {
			return secrets, err
		}
		secrets[k].Payload = string(resp.Payload.Data)
	}
	return secrets, nil
}

func (r *ReadSecrets) GetKV(secrets []Secret) []string {
	ret := []string{}
	for _, secret := range secrets {
		ret = append(ret, fmt.Sprintf("%s=%s", secret.Name, secret.Payload))
	}
	return ret
}

func (r *ReadSecrets) MatchLabel(label string, labels map[string]string) bool {
	if label == "" {
		return true
	}
	split := strings.Split(label, "=")
	if len(split) != 2 {
		return false
	}
	for k, v := range labels {
		if k == split[0] && v == split[1] {
			return true
		}
	}
	return false
}
