// Copyright 2020 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package framework

import (
	"fmt"
	"strconv"

	"github.com/kiegroup/kogito-cloud-operator/pkg/apis/app/v1alpha1"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/kubernetes"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/meta"
	"github.com/kiegroup/kogito-cloud-operator/pkg/framework"
	"github.com/kiegroup/kogito-cloud-operator/test/config"
)

// WaitForService waits that the service has a certain number of replicas
func WaitForService(namespace string, serviceName string, replicas int, timeoutInMin int) error {
	return WaitForOnOpenshift(namespace, serviceName+" running", timeoutInMin,
		func() (bool, error) {
			deployment, err := GetDeployment(namespace, serviceName)
			if err != nil {
				return false, err
			}
			if deployment == nil {
				return false, nil
			}
			return deployment.Status.Replicas == int32(replicas) && deployment.Status.AvailableReplicas == int32(replicas), nil
		})
}

// InstallService the Kogito Service component
func InstallService(namespace string, serviceName string, replicas int, installerType InstallerType, image v1alpha1.Image, cliName string, crResource meta.ResourceObject, persistence bool) error {
	GetLogger(namespace).Infof("%s install %s with %d replicas and persistence %v", serviceName, installerType, replicas, persistence)
	switch installerType {
	case CLIInstallerType:
		return cliInstall(namespace, replicas, image, cliName, persistence)
	case CRInstallerType:
		return crInstall(crResource)
	default:
		panic(fmt.Errorf("Unknown installer type %s", installerType))
	}
}

// BuildImage creates an Image object from an image tag
func BuildImage(imageTag string) v1alpha1.Image {
	image := framework.ConvertImageTagToImage(imageTag)
	if len(image.Tag) == 0 || image.Tag == "latest" {
		image.Tag = config.GetServicesImageVersion()
	}

	return image
}

func crInstall(crResource meta.ResourceObject) error {
	if _, err := kubernetes.ResourceC(kubeClient).CreateIfNotExists(crResource); err != nil {
		return fmt.Errorf("Error creating service: %v", err)
	}
	return nil
}

func cliInstall(namespace string, replicas int, image v1alpha1.Image, cliName string, persistence bool) error {
	cmd := []string{"install", cliName}

	if persistence {
		cmd = append(cmd, "--enable-persistence")
	}

	cmd = append(cmd, "--image", framework.ConvertImageToImageTag(image))
	cmd = append(cmd, "--replicas", strconv.Itoa(replicas))

	_, err := ExecuteCliCommandInNamespace(namespace, cmd...)
	return err
}
