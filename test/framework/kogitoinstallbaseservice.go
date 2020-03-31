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

// ProtoKogitoInstallService defines a service that can be installed.
type ProtoKogitoInstallService struct {
	ProtoKogitoService
	installerType   InstallerType
	persistence     bool
	imageTag        string
	cliName         string
	BuildCrResource func(service ProtoKogitoInstallService) meta.ResourceObject
}

// Install the Kogito Service component
func (service *ProtoKogitoInstallService) Install() error {
	GetLogger(service.namespace).Infof("%s install %s with %d replicas and persistence %v", service.label, service.installerType, service.replicas, service.persistence)
	switch service.installerType {
	case CLIInstallerType:
		return service.cliInstall()
	case CRInstallerType:
		return service.crInstall()
	default:
		panic(fmt.Errorf("Unknown installer type %s", service.installerType))
	}
}

func (service *ProtoKogitoInstallService) crInstall() error {
	resource := service.BuildCrResource(*service)

	if _, err := kubernetes.ResourceC(kubeClient).CreateIfNotExists(resource); err != nil {
		return fmt.Errorf("Error creating Kogito jobs service: %v", err)
	}
	return nil
}

func (service *ProtoKogitoInstallService) cliInstall() error {
	cmd := []string{"install", service.cliName}

	if service.persistence {
		cmd = append(cmd, "--enable-persistence")
	}

	cmd = append(cmd, "--image", framework.ConvertImageToImageTag(service.getImageTag()))
	cmd = append(cmd, "--replicas", strconv.Itoa(service.replicas))

	_, err := ExecuteCliCommandInNamespace(service.namespace, cmd...)
	return err
}

func (service *ProtoKogitoInstallService) getImageTag() v1alpha1.Image {
	image := framework.ConvertImageTagToImage(service.imageTag)
	if len(image.Tag) == 0 || image.Tag == "latest" {
		image.Tag = config.GetServicesImageVersion()
	}

	return image
}
