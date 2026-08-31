package registry

import "testing"

func TestRegistryLoadsAndFiltersExpoIntegrations(t *testing.T) {
	t.Parallel()

	reg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(reg.Integrations) < 7 {
		t.Fatalf("integration count = %d, want >= 7", len(reg.Integrations))
	}

	items, err := Filter("expo-mobile", "claude")
	if err != nil {
		t.Fatalf("Filter() error = %v", err)
	}
	seen := map[string]bool{}
	for _, item := range items {
		seen[item.ID] = true
	}
	for _, required := range []string{"expo-skills", "callstack-agent-skills", "gstack", "superpowers", "ui-ux-pro-max"} {
		if !seen[required] {
			t.Fatalf("filtered integrations missing %q", required)
		}
	}
	if seen["spec-kit"] {
		t.Fatalf("spec-kit should not be included in expo-mobile profile")
	}
}

func TestRecommendedIntegrationsSortFirst(t *testing.T) {
	t.Parallel()

	items, err := Filter("fullstack-saas", "claude")
	if err != nil {
		t.Fatalf("Filter() error = %v", err)
	}
	seenOptional := false
	for _, item := range items {
		if !item.Recommended {
			seenOptional = true
			continue
		}
		if seenOptional {
			t.Fatalf("recommended integration %q sorted after optional integration", item.ID)
		}
	}
}
