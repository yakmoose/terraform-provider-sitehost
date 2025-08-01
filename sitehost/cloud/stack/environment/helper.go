package environment

import (
	"fmt"

	"github.com/sitehostnz/gosh/pkg/models"
)

func createEnvironmentVariableChangeSet(oldValues, newValues interface{}) []models.EnvironmentVariable {
	var environmentVariables []models.EnvironmentVariable

	newV, ok := newValues.(map[string]interface{})
	if !ok {
		return environmentVariables
	}

	oldV, ok := oldValues.(map[string]interface{})
	if !ok {
		return environmentVariables
	}

	// things that exist in new need to be added or updated.
	for k, v := range newV {
		// if the content is different or does not exist, then we need to update it.
		if oldV[k] != newV[k] {
			environmentVariables = append(
				environmentVariables,
				models.EnvironmentVariable{Name: k, Content: fmt.Sprint(v)},
			)
		}
	}

	// removals - if something does not exist in the new map, then we need to remove it.
	for k := range oldV {
		if _, exists := newV[k]; !exists {
			environmentVariables = append(
				environmentVariables,
				models.EnvironmentVariable{Name: k, Content: ""},
			)
		}
	}

	return environmentVariables
}
