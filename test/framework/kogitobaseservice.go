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
	apps "k8s.io/api/apps/v1"
)

// ProtoKogitoService defines a service definition
type ProtoKogitoService struct {
	label       string
	namespace   string
	replicas    int
	serviceName string
}

// WaitForService waits that the service has a certain number of replicas
func (service *ProtoKogitoService) WaitForService(timeoutInMin int) error {
	return WaitForOnOpenshift(service.namespace, service.label+" running", timeoutInMin,
		func() (bool, error) {
			deployment, err := service.GetDeployment()
			if err != nil {
				return false, err
			}
			if deployment == nil {
				return false, nil
			}
			return deployment.Status.Replicas == int32(service.replicas) && deployment.Status.AvailableReplicas == int32(service.replicas), nil
		})
}

// GetDeployment retrieves the running service deployment
func (service *ProtoKogitoService) GetDeployment() (*apps.Deployment, error) {
	return GetDeployment(service.namespace, service.serviceName)
}
