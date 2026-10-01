# Pulumi Go GCP Artifact Registry

This Pulumi project provisions a Google Cloud Artifact Registry repository (Docker format) with optional IAM bindings for read/write access.

## Resources Created

- **Artifact Registry Repository** (Docker format, Standard mode)
- **IAM Bindings** (optional): Reader and Writer roles for specified members
- **API Enablement**: `artifactregistry.googleapis.com`

## Prerequisites

- Pulumi CLI installed
- Go 1.21+
- GCP credentials configured (`gcloud auth application-default login` or service account)
- A GCP project with billing enabled

## Configuration

Copy `Pulumi.dev.yaml.example` to `Pulumi.dev.yaml` and update values:

```yaml
config:
  gcp:project: your-project-id
  gcp:region: us-central1
  artifactRegistry:repositoryId: docker-repo
  artifactRegistry:description: Docker repository for storing container images
  artifactRegistry:readerMembers: "user:alice@example.com,group:devs@example.com"
  artifactRegistry:writerMembers: "serviceAccount:ci-cd@my-project.iam.gserviceaccount.com"
```

### Config Options

| Key | Description | Default |
|-----|-------------|---------|
| `gcp:project` | GCP Project ID | **Required** |
| `gcp:region` | GCP Region (Free Tier: us-central1, us-east1, us-west1) | `us-central1` |
| `artifactRegistry:repositoryId` | Base repository ID (random suffix appended) | `docker-repo` |
| `artifactRegistry:description` | Repository description | `Docker repository for storing container images` |
| `artifactRegistry:readerMembers` | Comma-separated members for `roles/artifactregistry.reader` | `""` |
| `artifactRegistry:writerMembers` | Comma-separated members for `roles/artifactregistry.writer` | `""` |

## Deployment

```bash
# Install dependencies
go mod tidy

# Preview changes
pulumi preview

# Deploy
pulumi up
```

## Outputs

| Output | Description |
|--------|-------------|
| `repositoryId` | The repository ID (with random suffix) |
| `repositoryName` | Full resource name of the repository |
| `dockerRegistryUrl` | Docker registry URL for pushing/pulling images |
| `registryLocation` | The repository location/region |

## Docker Usage

```bash
# Configure Docker auth
gcloud auth configure-docker <region>-docker.pkg.dev

# Tag and push
docker tag my-image:latest <region>-docker.pkg.dev/<project>/<repository-id>/my-image:latest
docker push <region>-docker.pkg.dev/<project>/<repository-id>/my-image:latest
```

## Free Tier Notes

Artifact Registry has a free tier including:
- 0.5 GB storage/month
- 1 GB network egress/month (within Americas)

See [GCP Artifact Registry Pricing](https://cloud.google.com/artifact-registry/pricing) for details.

## Cleanup

```bash
pulumi destroy
```