package upgrade

import (
	"fmt"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	dockerclient "github.com/portainer/portainer/api/docker/client"
	kubecli "github.com/portainer/portainer/api/kubernetes/cli"
	plf "github.com/portainer/portainer/api/platform"
	"github.com/portainer/portainer/api/stacks/deployments"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Upgrade(platform plf.ContainerPlatform, environment *portainer.Endpoint, licenseKey string) error
}

type service struct {
	kubernetesClientFactory   *kubecli.ClientFactory
	dockerClientFactory       *dockerclient.ClientFactory
	dockerComposeStackManager portainer.ComposeStackManager
	fileService               portainer.FileService

	isUpdating bool

	assetsPath    string
	stackDeployer deployments.StackDeployer
}

func NewService(
	assetsPath string,
	kubernetesClientFactory *kubecli.ClientFactory,
	dockerClientFactory *dockerclient.ClientFactory,
	dockerComposeStackManager portainer.ComposeStackManager,
	dataStore dataservices.DataStore,
	fileService portainer.FileService,
	stackDeployer deployments.StackDeployer,
) (Service, error) {

	return &service{
		assetsPath:                assetsPath,
		kubernetesClientFactory:   kubernetesClientFactory,
		dockerClientFactory:       dockerClientFactory,
		dockerComposeStackManager: dockerComposeStackManager,
		fileService:               fileService,
		stackDeployer:             stackDeployer,
	}, nil
}

func (service *service) Upgrade(platform plf.ContainerPlatform, environment *portainer.Endpoint, licenseKey string) error {
	service.isUpdating = true
	log.Debug().
		Str("platform", string(platform)).
		Msg("Starting upgrade process")

	switch platform {
	case plf.PlatformDockerStandalone:
		return service.upgradeDocker(environment, licenseKey, portainer.APIVersion, "standalone")
	case plf.PlatformDockerSwarm:
		return service.upgradeDocker(environment, licenseKey, portainer.APIVersion, "swarm")
	case plf.PlatformKubernetes:
		return service.upgradeKubernetes(environment, licenseKey, portainer.APIVersion)
	}

	service.isUpdating = false
	return fmt.Errorf("unsupported platform %s", platform)
}
