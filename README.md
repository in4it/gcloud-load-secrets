# Google Cloud load secrets
Execute binary with secrets loaded as environment variables

## Usage
Run bash command and inject environment variables that start with "myapp"
```
export GOOGLE_APPLICATION_CREDENTIALS=credentials.json
./gcloud-load-secrets-darwin-amd64 -prefix myapp -cmd '/bin/bash -c ls -ahl' -debug true
```
