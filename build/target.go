package build

type Target interface {
	BuildTo(directory string) error
}
