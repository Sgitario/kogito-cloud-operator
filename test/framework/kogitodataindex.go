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

// KogitoDataIndexInstallationHandler returns the installation service for the Kogito Data Index
func KogitoDataIndexInstallationHandler(namespace string, installerType InstallerType, replicas int) *ProtoKogitoInstallService {
	return &ProtoKogitoInstallService{
		installerType:      installerType,
		persistence:        false,
		imageTag:           getDataIndexImageTag(),
		cliName:            "data-index",
		BuildCrResource:    buildCrManagementConsole,
		ProtoKogitoService: *KogitoDataIndex(namespace, replicas),
	}
}

// KogitoDataIndex returns the service for the Kogito Data Index
func KogitoDataIndex(namespace string, replicas int) *ProtoKogitoService {
	return &ProtoKogitoService{
		label:       "Kogito Data Index",
		namespace:   namespace,
		replicas:    replicas,
		serviceName: infrastructure.DefaultDataIndexName,
	}
}

func getDataIndexImageTag() string {
	if len(config.GetDataIndexImageTag()) > 0 {
		return config.GetDataIndexImageTag()
	}

	return infrastructure.DefaultDataIndexImageFullTag
}

func buildCrDataIndex(service ProtoKogitoInstallService) meta.ResourceObject {
	replicas := int32(service.replicas)
	resource := &v1alpha1.KogitoDataIndex{
		ObjectMeta: metav1.ObjectMeta{
			Name:      infrastructure.DefaultDataIndexName,
			Namespace: service.namespace,
		},
		Spec: v1alpha1.KogitoDataIndexSpec{
			KogitoServiceSpec: v1alpha1.KogitoServiceSpec{
				Replicas: &replicas,
				Image:    service.getImageTag(),
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
