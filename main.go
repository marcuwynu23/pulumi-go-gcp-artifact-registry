package main

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		project := config.Get(ctx, "gcp:project")
		if project == "" {
			return fmt.Errorf("gcp:project is required")
		}

		region := config.Get(ctx, "gcp:region")
		if region == "" {
			region = "us-central1"
		}

		repositoryID := config.Get(ctx, "artifactRegistry:repositoryId")
		if repositoryID == "" {
			repositoryID = "docker-repo"
		}

		description := config.Get(ctx, "artifactRegistry:description")
		if description == "" {
			description = "Docker repository for storing container images"
		}

		readerMembersStr := config.Get(ctx, "artifactRegistry:readerMembers")
		writerMembersStr := config.Get(ctx, "artifactRegistry:writerMembers")

		// Enable Artifact Registry API
		_, err := projects.NewService(ctx, "artifactregistry-api", &projects.ServiceArgs{
			Service:          pulumi.String("artifactregistry.googleapis.com"),
			DisableOnDestroy: pulumi.Bool(false),
			Project:          pulumi.String(project),
		})
		if err != nil {
			return err
		}

		// Generate random suffix for unique repository ID
		randomSuffix, err := random.NewRandomId(ctx, "suffix", &random.RandomIdArgs{
			ByteLength: pulumi.Int(4),
		})
		if err != nil {
			return err
		}

		fullRepositoryID := pulumi.Sprintf("%s-%s", repositoryID, randomSuffix.Hex)

		// Create Artifact Registry repository
		repo, err := artifactregistry.NewRepository(ctx, "docker-repo", &artifactregistry.RepositoryArgs{
			Location:     pulumi.String(region),
			RepositoryId: fullRepositoryID,
			Description:  pulumi.String(description),
			Format:       pulumi.String("DOCKER"),
			Mode:         pulumi.String("STANDARD_REPOSITORY"),
			Project:      pulumi.String(project),
		}, pulumi.DependsOn([]pulumi.Resource{}))
		if err != nil {
			return err
		}

		// Grant read access to reader members
		if readerMembersStr != "" {
			readerMembers := strings.Split(readerMembersStr, ",")
			for i, member := range readerMembers {
				member = strings.TrimSpace(member)
				if member != "" {
					_, err := artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("reader-%d", i), &artifactregistry.RepositoryIamMemberArgs{
						Location:   repo.Location,
						Repository: repo.Name,
						Role:       pulumi.String("roles/artifactregistry.reader"),
						Member:     pulumi.String(member),
					}, pulumi.DependsOn([]pulumi.Resource{repo}))
					if err != nil {
						return err
					}
				}
			}
		}

		// Grant write access to writer members
		if writerMembersStr != "" {
			writerMembers := strings.Split(writerMembersStr, ",")
			for i, member := range writerMembers {
				member = strings.TrimSpace(member)
				if member != "" {
					_, err := artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("writer-%d", i), &artifactregistry.RepositoryIamMemberArgs{
						Location:   repo.Location,
						Repository: repo.Name,
						Role:       pulumi.String("roles/artifactregistry.writer"),
						Member:     pulumi.String(member),
					}, pulumi.DependsOn([]pulumi.Resource{repo}))
					if err != nil {
						return err
					}
				}
			}
		}

		// Export outputs
		ctx.Export("repositoryId", repo.RepositoryId)
		ctx.Export("repositoryName", repo.Name)
		ctx.Export("dockerRegistryUrl", pulumi.Sprintf("%s-docker.pkg.dev/%s/%s", region, project, repo.RepositoryId))
		ctx.Export("registryLocation", repo.Location)

		return nil
	})
}
