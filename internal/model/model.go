package model

import "time"

// Profile describes a product delivery profile. Profiles intentionally define
// outcomes and gates, not a mandatory application framework.
type Profile struct {
	ID                      string            `json:"id"`
	DisplayName             string            `json:"displayName"`
	Description             string            `json:"description"`
	DefaultStack            map[string]string `json:"defaultStack"`
	RequiredArtifacts       []string          `json:"requiredArtifacts"`
	QualityGates            []string          `json:"qualityGates"`
	RecommendedIntegrations []string          `json:"recommendedIntegrations"`
}

// Manifest is the durable control record written to every managed project.
type Manifest struct {
	SchemaVersion string     `json:"schemaVersion"`
	Project       Project    `json:"project"`
	Governance    Governance `json:"governance"`
	State         State      `json:"state"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type Project struct {
	Name        string `json:"name"`
	Idea        string `json:"idea"`
	Profile     string `json:"profile"`
	Mode        string `json:"mode"`
	TargetAgent string `json:"targetAgent"`
}

type Governance struct {
	ApprovalMode       string   `json:"approvalMode"`
	SourceOfTruth      []string `json:"sourceOfTruth"`
	MaxResearchAgents  int      `json:"maxResearchAgents"`
	MaxReviewCycles    int      `json:"maxReviewCycles"`
	AllowAutoMerge     bool     `json:"allowAutoMerge"`
	AllowProductionRun bool     `json:"allowProductionRun"`
}

type State struct {
	Phase  string `json:"phase"`
	Status string `json:"status"`
}

// Integration is an external skill or plugin that may be installed separately.
// Auto Company never silently installs third-party code.
type Integration struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Repository  string              `json:"repository"`
	License     string              `json:"license"`
	Maturity    string              `json:"maturity"`
	Role        string              `json:"role"`
	Profiles    []string            `json:"profiles"`
	Agents      []string            `json:"agents"`
	Recommended bool                `json:"recommended"`
	Install     map[string][]string `json:"install"`
	Notes       string              `json:"notes"`
}

type IntegrationRegistry struct {
	SchemaVersion string        `json:"schemaVersion"`
	Integrations  []Integration `json:"integrations"`
}
