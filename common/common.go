package common

var (
	// OutputFields for custom format output
	OutputFields string
	// OutputFormat for custom format output
	OutputFormat string
	// RegionSet picks the region to connect to, if you use this option it will use it over the default region
	RegionSet string
	// DefaultYes : automatic yes to prompts; assume \"yes\" as answer to all prompts and run non-interactively
	DefaultYes bool
	// PrettySet : Prints the json output in pretty format
	PrettySet bool
)
