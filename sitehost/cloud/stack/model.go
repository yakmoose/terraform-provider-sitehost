package stack

type (
	DockerFileService struct {
		Build       string            `yaml:"build,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Command     []string          `yaml:"command,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		EnvFile     []string          `yaml:"env_file,omitempty"`
		Restart     string            `yaml:"restart,omitempty"`
		// looks like a map, but it's an array of things
		// yaml parser won't treat it as a map
		Labels  []string `yaml:"labels,omitempty"`
		Volumes []string `yaml:"volumes,omitempty"`
	}

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
)
