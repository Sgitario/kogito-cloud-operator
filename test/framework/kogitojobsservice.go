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

// KogitoJobsServiceInstall returns the installation service for the Kogito Jobs Service
func KogitoJobsServiceInstall(namespace string, installerType InstallerType, replicas int, persistence bool) error {
	image := getJobsServiceImage()
	cliName := "jobs-service"
	crResource := buildCrJobsService(namespace, int32(replicas), image, persistence)
	return InstallService(namespace, getJobsServiceName(), replicas, installerType, image, cliName, crResource, persistence)
}

// KogitoJobsServiceWaitForService wait for Kogito Jobs Service to be deployed
func KogitoJobsServiceWaitForService(namespace string, replicas int, timeoutInMin int) error {
	return WaitForService(namespace, getJobsServiceName(), replicas, timeoutInMin)
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
	if exists, err := kubernetes.ResourceC(kubeClient).FetchWithKey(types.NamespacedName{Name: getJobsServiceName(), Namespace: namespace}, service); err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("Error while trying to look for Kogito jobs service: %v ", err)
	} else if !exists {
		return nil, nil
	}
	return service, nil
}

func getJobsServiceName() string {
	return infrastructure.DefaultJobsServiceName
}

func getJobsServiceImage() v1alpha1.Image {
	imageTag := infrastructure.DefaultJobsServiceImageFullTag
	if len(config.GetJobsServiceImageTag()) > 0 {
		imageTag = config.GetJobsServiceImageTag()
	}

	return BuildImage(imageTag)
}

func buildCrJobsService(namespace string, replicas int32, image v1alpha1.Image, persistence bool) meta.ResourceObject {
	resource := &v1alpha1.KogitoJobsService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      infrastructure.DefaultJobsServiceName,
			Namespace: namespace,
		},
		Spec: v1alpha1.KogitoJobsServiceSpec{
			KogitoServiceSpec: v1alpha1.KogitoServiceSpec{
				Replicas: &replicas,
				Image:    image,
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

	if persistence {
		resource.Spec.InfinispanProperties = v1alpha1.InfinispanConnectionProperties{
			UseKogitoInfra: true,
		}
	}

	return resource
}
