package assets

import (
	"os"
	"path/filepath"
)

/*platformDefinedPath defines the asset path
designted by execution platform*/
var platformDefinedPath string

/*GetEffectiveAssetsPath returns the path
where asset files can be located.

Currently, only platform can set such an lookup path,
and a path defined by platform overrides default one.

If no asset path have been defined, this function returns
filepath.Dir(os.Executable(){0})
*/
func GetEffectiveAssetsPath() (string, error) {
	if platformDefinedPath != "" {
		return platformDefinedPath, nil
	}
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

/*SetPlatformDefinedPath set the path to lookup assets
Only Developers in the Platform role should be able to
set assets load path with this method as the path defined
here is never checked.
*/
func SetPlatformDefinedPath(pdp string) {
	platformDefinedPath = pdp
}
