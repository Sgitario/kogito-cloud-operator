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

	"github.com/kiegroup/kogito-cloud-operator/pkg/apis/app/v1alpha1"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/kubernetes"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/meta"
	"github.com/kiegroup/kogito-cloud-operator/pkg/infrastructure"
	"github.com/kiegroup/kogito-cloud-operator/test/config"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// KogitoJobsServiceInstallationHandler returns the installation service for the Kogito Jobs Service
func KogitoJobsServiceInstallationHandler(namespace string, installerType InstallerType, replicas int, persistence bool) *ProtoKogitoInstallService {
	return &ProtoKogitoInstallService{
		installerType:      installerType,
		persistence:        persistence,
		imageTag:           getJobsServiceImageTag(),
		cliName:            "jobs-service",
		BuildCrResource:    buildCrJobsService,
		ProtoKogitoService: *KogitoJobsService(namespace, replicas),
	}
}

// KogitoJobsService returns the service for the Kogito Jobs Service
func KogitoJobsService(namespace string, replicas int) *ProtoKogitoService {
	return &ProtoKogitoService{
		label:       "Kogito Jobs Service",
		namespace:   namespace,
		replicas:    replicas,
		serviceName: infrastructure.DefaultJobsServiceName,
	}
}

// SetKogitoJobsServiceReplicas sets the number of replicas for the Kogito Jobs Service
func SetKogitoJobsServiceReplicas(namespace string, nbPods int32) error {
	GetLogger(namespace).Infof("Set Kogito jobs service replica number to %d", nbPods)
	kogitoJobsService, err := GetKogitoJobsService(namespace)
	if err != nil {
		return err
	} else if kogitoJobsService == nil {
		return fmt.Errorf("No Kogito jobs service found in namespace %s", namespace)
	}
	kogitoJobsService.Spec.Replicas = &nbPods
	return kubernetes.ResourceC(kubeClient).Update(kogitoJobsService)
}

// GetKogitoJobsService retrieves the running jobs service
func GetKogitoJobsService(namespace string) (*v1alpha1.KogitoJobsService, error) {
	service := &v1alpha1.KogitoJobsService{}
	if exists, err := kubernetes.ResourceC(kubeClient).FetchWithKey(types.NamespacedName{Name: infrastructure.DefaultJobsServiceName, Namespace: namespace}, service); err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("Error while trying to look for Kogito jobs service: %v ", err)
	} else if !exists {
		return nil, nil
	}
	return service, nil
}

func getJobsServiceImageTag() string {
	if len(config.GetJobsServiceImageTag()) > 0 {
		return config.GetJobsServiceImageTag()
	}

	return infrastructure.DefaultJobsServiceImageFullTag
}

func buildCrJobsService(service ProtoKogitoInstallService) meta.ResourceObject {
	replicas := int32(service.replicas)
	resource := &v1alpha1.KogitoJobsService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      infrastructure.DefaultJobsServiceName,
			Namespace: service.namespace,
		},
		Spec: v1alpha1.KogitoJobsServiceSpec{
			KogitoServiceSpec: v1alpha1.KogitoServiceSpec{
				Replicas: &replicas,
				Image:    service.getImageTag(),
			},
		},
		Status: v1alpha1.KogitoJobsServiceStatus{
			KogitoServiceStatus: v1alpha1.KogitoServiceStatus{
				ConditionsMeta: v1alpha1.ConditionsMeta{
					Conditions: []v1alpha1.Condition{},
				},
			},
		},
	}

	if service.persistence {
		resource.Spec.InfinispanProperties = v1alpha1.InfinispanConnectionProperties{
			UseKogitoInfra: true,
		}
	}

	return resource
}
