package registry

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/aleka7sk/auto-company/internal/assets"
	"github.com/aleka7sk/auto-company/internal/model"
)

const registryPath = "registry/integrations.json"

func Load() (model.IntegrationRegistry, error) {
	var result model.IntegrationRegistry
	data, err := assets.FS.ReadFile(registryPath)
	if err != nil {
		return result, fmt.Errorf("read embedded integration registry: %w", err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("parse integration registry: %w", err)
	}
	return result, nil
}

func Filter(profile, agent string) ([]model.Integration, error) {
	reg, err := Load()
	if err != nil {
		return nil, err
	}
	items := make([]model.Integration, 0)
	for _, item := range reg.Integrations {
		if profile != "" && !contains(item.Profiles, profile) && !contains(item.Profiles, "all") {
			continue
		}
		if agent != "" && !contains(item.Agents, agent) && !contains(item.Agents, "all") {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Recommended == items[j].Recommended {
			return items[i].ID < items[j].ID
		}
		return items[i].Recommended
	})
	return items, nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
