package model

// Report describes the requested checks only, never the security of an application.
type Report struct {
	SchemaVersion string       `json:"schema_version"`
	Stage         string       `json:"stage"`
	Complete      bool         `json:"complete"`
	Checks        []string     `json:"checks"`
	Coverage      Coverage     `json:"coverage"`
	Findings      []Finding    `json:"findings"`
	Diagnostics   []Diagnostic `json:"diagnostics"`
	Exclusions    []Exclusion  `json:"exclusions"`
	Limitations   []string     `json:"limitations"`
}

type Coverage struct {
	FilesScanned  int   `json:"files_scanned"`
	BytesScanned  int64 `json:"bytes_scanned"`
	SourcesParsed int   `json:"sources_parsed"`
}

type Finding struct {
	ID             string `json:"id"`
	RuleID         string `json:"rule_id"`
	UpstreamRuleID string `json:"upstream_rule_id"`
	Engine         string `json:"engine"`
	EngineVersion  string `json:"engine_version"`
	Severity       string `json:"severity"`
	Confidence     string `json:"confidence"`
	Path           string `json:"path"`
	Line           int    `json:"line"`
	Column         int    `json:"column"`
	Explanation    string `json:"explanation"`
	Evidence       string `json:"evidence"`
	Remediation    string `json:"remediation"`
}

type Diagnostic struct {
	Code    string `json:"code"`
	Path    string `json:"path,omitempty"`
	Line    int    `json:"line,omitempty"`
	Message string `json:"message"`
}

type Exclusion struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

func (r Report) ExitCode() int {
	if !r.Complete {
		return 2
	}
	for _, f := range r.Findings {
		if (f.Severity == "high" || f.Severity == "critical") && f.Confidence == "high" {
			return 1
		}
	}
	return 0
}
