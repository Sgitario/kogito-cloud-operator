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
	"github.com/kiegroup/kogito-cloud-operator/pkg/apis/app/v1alpha1"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/meta"
	"github.com/kiegroup/kogito-cloud-operator/pkg/infrastructure"
	"github.com/kiegroup/kogito-cloud-operator/test/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KogitoManagementConsoleInstall returns the installation service for the Kogito Management Console
func KogitoManagementConsoleInstall(namespace string, installerType InstallerType, replicas int) error {
	image := getManagementConsoleImage()
	cliName := "mgmt-console"
	crResource := buildCrManagementConsole(namespace, int32(replicas), image)
	persistence := false
	return InstallService(namespace, getManagementConsoleServiceName(), replicas, installerType, image, cliName, crResource, persistence)
}

// KogitoManagementConsoleWaitForService wait for Kogito Management Console to be deployed
func KogitoManagementConsoleWaitForService(namespace string, replicas int, timeoutInMin int) error {
	return WaitForService(namespace, getManagementConsoleServiceName(), replicas, timeoutInMin)
}

func getManagementConsoleServiceName() string {
	return infrastructure.DefaultMgmtConsoleName
}

func getManagementConsoleImage() v1alpha1.Image {
	var imageTag = infrastructure.DefaultMgmtConsoleImageFullTag
	if len(config.GetManagementConsoleImageTag()) > 0 {
		imageTag = config.GetManagementConsoleImageTag()
	}

	return BuildImage(imageTag)
}

func buildCrManagementConsole(namespace string, replicas int32, image v1alpha1.Image) meta.ResourceObject {
	resource := &v1alpha1.KogitoMgmtConsole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      infrastructure.DefaultMgmtConsoleName,
			Namespace: namespace,
		},
		Spec: v1alpha1.KogitoMgmtConsoleSpec{
			KogitoServiceSpec: v1alpha1.KogitoServiceSpec{
				Replicas: &replicas,
				Image:    image,
			},
		},
		Status: v1alpha1.KogitoMgmtConsoleStatus{
			KogitoServiceStatus: v1alpha1.KogitoServiceStatus{
				ConditionsMeta: v1alpha1.ConditionsMeta{
					Conditions: []v1alpha1.Condition{},
				},
			},
		},
	}

	return resource
}
