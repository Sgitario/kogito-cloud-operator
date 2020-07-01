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

package flag

import (
	"fmt"
	"github.com/kiegroup/kogito-cloud-operator/pkg/apis/app/v1alpha1"
	"github.com/kiegroup/kogito-cloud-operator/pkg/util"
	"github.com/spf13/cobra"
)

const defaultDeployRuntime = string(v1alpha1.QuarkusRuntimeType)

var (
	runtimeTypeValidEntries = []string{string(v1alpha1.QuarkusRuntimeType), string(v1alpha1.SpringbootRuntimeType)}
)

// RuntimeFlags is common properties used to configure runtime properties
type RuntimeFlags struct {
	Runtime string
}

// AddRuntimeFlags adds the Runtime flags to the given command
func AddRuntimeFlags(command *cobra.Command, flags *RuntimeFlags) {
	command.Flags().StringVarP(&flags.Runtime, "runtime", "r", defaultDeployRuntime, "The runtime which should be used to build the Service. Valid values are 'quarkus' or 'springboot'. Default to '"+defaultDeployRuntime+"'.")
}

// CheckRuntimeArgs validates the RuntimeFlags flags
func CheckRuntimeArgs(flags *RuntimeFlags) error {
	if !util.Contains(flags.Runtime, runtimeTypeValidEntries) {
		return fmt.Errorf("runtime not valid. Valid runtimes are %s. Received %s", runtimeTypeValidEntries, flags.Runtime)
	}
	return nil
}
