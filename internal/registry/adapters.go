// Upstream adapter-repo path constants and the projects/ excluded-child set.
package registry

const (
	TVLAdapterPath = "DefiLlama-Adapters/projects/"
	DMSAdapterPath = "dimension-adapters/"
)

// excludedAdapterChildren: subfolders under DefiLlama-Adapters/projects/ that are not adapters.
var excludedAdapterChildren = map[string]struct{}{
	"helper":   {},
	"treasury": {},
	"entities": {},
	"config":   {},
	"stacks":   {},
	"test":     {},
}

// IsExcludedAdapterChild reports whether name is a DefiLlama-Adapters/projects/ child to skip.
func IsExcludedAdapterChild(name string) bool {
	_, skip := excludedAdapterChildren[name]
	return skip
}
