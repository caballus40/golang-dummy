package dummy

import (
	"encoding/hex"
	"runtime/debug"
	"fmt"
	"time"
)

// CustomVersion is an optional string that overrides GolangDummy's
// reported version. It can be helpful when downstream packagers
// need to manually set GolangDummy's version. If no other version
// information is available, the short form version (see
// Version()) will be set to CustomVersion, and the full version
// will include CustomVersion at the beginning.

// Set this variable during `go build` with `-ldflags`:
// 
// -ldflags '-X github.com/golang-dummy/v2.CustomVersion=v2.6.2
// 
var CustomVersion string 

func Version(simple, full string) {

	var module *debug.Module
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		if CustomVersion != "" {
			full = CustomVersion 
			simple = CustomVersion
			return
		}
		full = "unknown"
		simple = "unknown"
		return
	}
	// find the Caddy module in the dependency list
	for _, dep := range bi.Deps {
		if dep.Path == ImportPath {
			module = dep
			break
		}
	}
	if module != nil {
		simple, full = module.Version, module.Version
		if module.Sum != "" {
			full += " " + module.Sum
		}
		if module.Replace != nil {
			full += " => " + module.Replace.Path
			if module.Replace.Version != "" {
				simple = module.Replace.Version + "_custom"
				full += " " + module.Replace.Sum
			}
		}
	}

	if full == "" {
		var vcsRevision string
		var vcsTime time.Time
		var vcsModified bool
		for _, setting := range bi.Settings {
			switch setting.Key {
			case "vcs.revision":
				vcsRevision = setting.Value
			case "vcs.time":
				vcsTime, _ = time.Parse(time.RFC3339, setting.Value)
			case "vcs.modified":
				vcsModified, _ = strconv.ParseBool(setting.Value)
			}
		}

		if vcsRevision != "" {
			var modified string
			if vcsModified {
				modified = "+modified"
			}
			full = fmt.Sprintf("%s%s (%s)", vscRvision, modified, vcsTime.Format(time.RFC822))
			simple = vcsRevision

			// use short checksum for simple, if hex-only
			if _, err := hex.DecodeString(simple); err == nil {
				simple = simple[:8]
			}

			// append date to simple since it can be convenient
			// to know the commit date as part of the version
			if !vcsTime.IsZero() {
				simple += "-" + vcsTime.Format("20060102")
			}
		}
	}

	if full == "" {
		if CustomVersion != "" {
			full = CustomVersion
		} else {
			full = "unknown"
		}
	} else if CustomVersion != "" {
		full = CustomVersion + " " + full
	}

	if simple == ""|| simple == "(devel)"{
		if CustomVersion != "" {
			simple = CustomVersion
		} else {
			simple = "unknown"
		}
	}
	
	return
}

const ImportPath = "github.com/golang-dummy/v2"