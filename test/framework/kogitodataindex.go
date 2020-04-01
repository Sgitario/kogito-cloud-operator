// Copyright 2019 Red Hat, Inc. and/or its affiliates
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

// KogitoDataIndexInstall returns the installation service for the Kogito Data Index
func KogitoDataIndexInstall(namespace string, installerType InstallerType, replicas int) error {
	image := getDataIndexImage()
	cliName := "data-index"
	crResource := buildCrDataIndex(namespace, int32(replicas), image)
	persistence := false
	return InstallService(namespace, getDataIndexServiceName(), replicas, installerType, image, cliName, crResource, persistence)
}

// KogitoDataIndexWaitForService wait for Kogito Data Index to be deployed
func KogitoDataIndexWaitForService(namespace string, replicas int, timeoutInMin int) error {
	return WaitForService(namespace, getDataIndexServiceName(), replicas, timeoutInMin)
}

func getDataIndexServiceName() string {
	return infrastructure.DefaultDataIndexName
}

func getDataIndexImage() v1alpha1.Image {
	var imageTag = infrastructure.DefaultDataIndexImageFullTag
	if len(config.GetDataIndexImageTag()) > 0 {
		imageTag = config.GetDataIndexImageTag()
	}

	return BuildImage(imageTag)
}

func buildCrDataIndex(namespace string, replicas int32, image v1alpha1.Image) meta.ResourceObject {
	resource := &v1alpha1.KogitoDataIndex{
		ObjectMeta: metav1.ObjectMeta{
			Name:      infrastructure.DefaultDataIndexName,
			Namespace: namespace,
		},
		Spec: v1alpha1.KogitoDataIndexSpec{
			KogitoServiceSpec: v1alpha1.KogitoServiceSpec{
				Replicas: &replicas,
				Image:    image,
			},
		},
		Status: v1alpha1.KogitoDataIndexStatus{
			KogitoServiceStatus: v1alpha1.KogitoServiceStatus{
				ConditionsMeta: v1alpha1.ConditionsMeta{
					Conditions: []v1alpha1.Condition{},
				},
			},
		},
	}

	return resource
}
