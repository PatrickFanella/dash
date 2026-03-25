package importer

type DashyConfig struct {
	Sections []DashySection `yaml:"sections"`
}

type DashySection struct {
	Name        string           `yaml:"name"`
	Icon        string           `yaml:"icon"`
	DisplayData DashyDisplayData `yaml:"displayData"`
	Items       []DashyItem      `yaml:"items"`
	Widgets     []DashyWidget    `yaml:"widgets"`
}

type DashyDisplayData struct {
	Collapsed bool  `yaml:"collapsed"`
	Cols      int32 `yaml:"cols"`
}

type DashyItem struct {
	Title          string `yaml:"title"`
	URL            string `yaml:"url"`
	Icon           string `yaml:"icon"`
	Description    string `yaml:"description"`
	StatusCheck    *bool  `yaml:"statusCheck"`
	StatusCheckURL string `yaml:"statusCheckUrl"`
}

type DashyWidget struct {
	Type string `yaml:"type"`
}
