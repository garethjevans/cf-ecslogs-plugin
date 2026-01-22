package ecslogs

// ECSLog represents a log entry in ECS (Elastic Common Schema) format
type ECSLog struct {
	Timestamp string   `json:"@timestamp"`
	Message   string   `json:"message"`
	Log       *Log     `json:"log,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Labels    *Labels  `json:"labels,omitempty"`
	Process   *Process `json:"process,omitempty"`
	Error     *Error   `json:"error,omitempty"`
	Service   *Service `json:"service,omitempty"`
	Cloud     *Cloud   `json:"cloud,omitempty"`
	Host      *Host    `json:"host,omitempty"`
	ECS       *ECS     `json:"ecs,omitempty"`
}

// Log contains log-specific metadata
type Log struct {
	Level  string `json:"level,omitempty"`
	Logger string `json:"logger,omitempty"`
}

// Labels contains additional label information
type Labels struct {
	AppGUID     string `json:"app_guid,omitempty"`
	AppName     string `json:"app_name,omitempty"`
	OrgGUID     string `json:"org_guid,omitempty"`
	OrgName     string `json:"org_name,omitempty"`
	SpaceGUID   string `json:"space_guid,omitempty"`
	SpaceName   string `json:"space_name,omitempty"`
	InstanceID  string `json:"instance_id,omitempty"`
	ProcessID   string `json:"process_id,omitempty"`
	ProcessType string `json:"process_type,omitempty"`
	SourceType  string `json:"source_type,omitempty"`
}

// Process contains process information
type Process struct {
	Type   string  `json:"type,omitempty"`
	Index  int     `json:"index,omitempty"`
	PID    int     `json:"pid,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}

// Thread contains thread information
type Thread struct {
	Name string `json:"name,omitempty"`
}

// Error contains error information
type Error struct {
	Message    string `json:"message,omitempty"`
	Type       string `json:"type,omitempty"`
	StackTrace string `json:"stack_trace,omitempty"`
}

// Service contains service information
type Service struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
}

// Cloud contains cloud provider information
type Cloud struct {
	Provider string `json:"provider,omitempty"`
	Region   string `json:"region,omitempty"`
	Zone     string `json:"availability_zone,omitempty"`
}

// Host contains host information
type Host struct {
	Name string `json:"name,omitempty"`
	IP   string `json:"ip,omitempty"`
}

// ECS contains ECS version information
type ECS struct {
	Version string `json:"version,omitempty"`
}
