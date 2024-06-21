package upgrade

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/filesystem"

	"github.com/cbroglie/mustache"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (service *service) upgradeDocker(environment *portainer.Endpoint, licenseKey, version string, envType string) error {
	ctx := context.TODO()

	image := getEEImageName(version)
	updaterImage := getUpdaterImageName()
	skipPullImage := shouldSkipPullImage()

	if err := service.checkImageForDocker(ctx, environment, image, skipPullImage); err != nil {
		return err
	}

	tempStack, err := service.generateUpdaterStackFromTemplate(generateOptions{
		image:         image,
		skipPullImage: skipPullImage,
		updaterImage:  updaterImage,
		license:       licenseKey,
		envType:       envType,
		version:       version,
	})

	if err != nil {
		return err
	}

	service.stackDeployer.DeployComposeStack(tempStack, environment, []portainer.Registry{}, true, !skipPullImage)

	// err = service.dockerComposeStackManager.Run(ctx, tempStack, environment, "updater", portainer.ComposeRunOptions{
	// 	Remove:   true,
	// 	Detached: true,
	// })

	// if err != nil {
	// 	return errors.Wrap(err, "failed to deploy upgrade stack")
	// }

	return nil
}

func (service *service) checkImageForDocker(ctx context.Context, environment *portainer.Endpoint, imageName string, skipPullImage bool) error {
	cli, err := service.dockerClientFactory.CreateClient(environment, "", nil)
	if err != nil {
		return errors.Wrap(err, "failed to create docker client")
	}

	// check in existing images on host
	if skipPullImage {
		filters := filters.NewArgs()
		filters.Add("reference", imageName)
		images, err := cli.ImageList(ctx, image.ListOptions{
			Filters: filters,
		})
		if err != nil {
			return errors.Wrap(err, "failed to list images")
		}

		if len(images) == 0 {
			return errors.Errorf("image %s not found locally", imageName)
		}

		return nil
	}

	// check if available on registry
	_, err = cli.DistributionInspect(ctx, imageName, "")
	if err != nil {
		return errors.Errorf("image %s not found on registry", imageName)
	}

	return nil

}

// always sync the template with fields names!
type generateOptions struct {
	image         string
	skipPullImage bool
	updaterImage  string
	license       string
	envType       string
	version       string
}

// mustacheUpgradeDockerTemplateFile represents the name of the template file for the docker upgrade
const mustacheUpgradeDockerTemplateFile = "upgrade-docker.yml.mustache"

func (service *service) generateUpdaterStackFromTemplate(generateOptions generateOptions) (*portainer.Stack, error) {

	templateFilePath := filesystem.JoinPaths(service.assetsPath, "mustache-templates", mustacheUpgradeDockerTemplateFile)

	composeFile, err := mustache.RenderFile(templateFilePath, generateOptions)
	log.Debug().
		Str("composeFile", composeFile).
		Msg("Compose file for upgrade")

	if err != nil {
		return nil, errors.Wrap(err, "failed to render upgrade template")
	}

	timeId := time.Now().Unix()
	fileName := fmt.Sprintf("upgrade-%d.yml", timeId)

	filePath, err := service.fileService.StoreStackFileFromBytes("upgrade", fileName, []byte(composeFile))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create upgrade compose file")
	}

	projectName := fmt.Sprintf(
		"portainer-upgrade-%d-%s",
		timeId,
		strings.ReplaceAll(generateOptions.version, ".", "-"),
	)

	return &portainer.Stack{
		Name:        projectName,
		ProjectPath: filePath,
		EntryPoint:  fileName,
	}, nil
}
