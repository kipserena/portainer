package upgrade

import "os"

// portainerImagePrefixEnvVar represents the name of the environment variable used to define the image prefix for portainer-updater
// useful if there's a need to test PR images
const portainerImagePrefixEnvVar = "UPGRADE_PORTAINER_IMAGE_PREFIX"

func getEEImageName(version string) string {
	prefix := os.Getenv(portainerImagePrefixEnvVar)
	if prefix == "" {
		prefix = "portainer/portainer-ee"
	}
	return prefix + ":" + version
}

// updaterImageEnvVar represents the name of the environment variable used to define the updater image
// useful if there's a need to test a different updater
const updaterImageEnvVar = "UPGRADE_UPDATER_IMAGE"

func getUpdaterImageName() string {
	image := os.Getenv(updaterImageEnvVar)
	if image == "" {
		image = "portainer/portainer-updater:latest"
	}
	return image
}

// skipPullImageEnvVar represents the name of the environment variable used to define if the image pull should be skipped
// useful if there's a need to test local images
const skipPullImageEnvVar = "UPGRADE_SKIP_PULL_PORTAINER_IMAGE"

func shouldSkipPullImage() bool {
	skipPullImageEnv := os.Getenv(skipPullImageEnvVar)
	return skipPullImageEnv != ""
}
