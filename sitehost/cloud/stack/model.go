// Package stack represents interactions with a stack on sitehose, this is the model.
package stack

type (
	// DockerFileService represents a DockerFile.
	DockerFileService struct {
		Build       string   `yaml:"build,omitempty"`
		Image       string   `yaml:"image,omitempty"`
		Command     []string `yaml:"command,omitempty"`
		Ports       []string `yaml:"ports,omitempty"`
		Environment []string `yaml:"environment,omitempty"`
		EnvFile     []string `yaml:"env_file,omitempty"`
		Restart     string   `yaml:"restart,omitempty"`
		// looks like a map, but it's an array of things
		// yaml parser won't treat it as a map
		Labels  []string `yaml:"labels,omitempty"`
		Volumes []string `yaml:"volumes,omitempty"`
	}

	// Compose represents a docker Compose.
	Compose struct {
		Version  string                       `yaml:"version"`
		Services map[string]DockerFileService `yaml:"services"`

		Networks map[string]struct {
			Driver string `yaml:"driver,omitempty"`
		} `yaml:"networks,omitempty"`

		Volumes map[string]struct {
			Driver string `yaml:"driver,omitempty"`
		} `yaml:"volumes,omitempty"`
	}

	// ParsedStackName represents the components of a parsed stack name, it is used primarily in importing so we can handle multiple import formats.
	ParsedStackName struct {
		ServerName string
		Project    string
		Service    string
	}
)
