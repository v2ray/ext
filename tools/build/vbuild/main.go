package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	xbuild "v2ray.com/ext/build"
	"v2ray.com/ext/tools/build"
	"v2ray.com/ext/zip"
)

var (
	flagTargetDir    = flag.String("dir", "", "Directory to put generated files.")
	flagTargetOS     = flag.String("os", runtime.GOOS, "Target OS of this build.")
	flagTargetArch   = flag.String("arch", runtime.GOARCH, "Target CPU arch of this build.")
	flagArchive      = flag.Bool("zip", false, "Whether to make an archive of files or not.")
	flagMetadataFile = flag.String("metadata", "metadata.txt", "File to store metadata info of released packages.")
	flagSignBinary   = flag.Bool("sign", false, "Whether or not to sign the binaries.")

	binPath string
)

func createTargetDirectory(version string, goOS xbuild.OS, goArch xbuild.Arch) (string, error) {
	var targetDir string
	if len(*flagTargetDir) > 0 {
		targetDir = *flagTargetDir
	} else {
		suffix := xbuild.GetSuffix(goOS, goArch)

		targetDir = filepath.Join(binPath, "v2ray-"+version+suffix)
		if version != "custom" {
			os.RemoveAll(targetDir)
		}
	}

	err := os.MkdirAll(targetDir, os.ModeDir|0777)
	return targetDir, err
}

func getBinPath() string {
	GOPATH := os.Getenv("GOPATH")
	return filepath.Join(GOPATH, "bin")
}

func isOfficialBuild() bool {
	version := os.Getenv("TRAVIS_TAG")
	return len(version) > 0
}

func main() {
	flag.Parse()
	binPath = getBinPath()

	v2rayOS := xbuild.ParseOS(*flagTargetOS)
	v2rayArch := xbuild.ParseArch(*flagTargetArch)

	version := os.Getenv("TRAVIS_TAG")

	if len(version) == 0 {
		version = "custom"
	}

	fmt.Printf("Building V2Ray (%s) for %s %s\n", version, v2rayOS, v2rayArch)

	targetDir, err := createTargetDirectory(version, v2rayOS, v2rayArch)
	if err != nil {
		fmt.Println("Unable to create directory " + targetDir + ": " + err.Error())
	}

	targets := build.GetReleaseTargets(v2rayOS, v2rayArch)
	for _, target := range targets {
		if err := target.BuildTo(targetDir); err != nil {
			fmt.Println("Failed to build V2Ray on", v2rayArch, "for", v2rayOS, "with error", err.Error())
			return
		}

		if *flagSignBinary {
			gpgPass := os.Getenv("GPG_SIGN_PASS")
			targetFile := filepath.Join(targetDir, target.Target)
			if err := build.GPGSignFile(targetFile, gpgPass); err != nil {
				fmt.Println("Unable to sign file", targetFile, "with error", err.Error())
				return
			}
		}
	}

	if err := build.CopyAllConfigFiles(targetDir, v2rayOS); err != nil {
		fmt.Println("Unable to copy config files: " + err.Error())
	}

	if *flagArchive {
		if err := os.Chdir(binPath); err != nil {
			fmt.Printf("Unable to switch to directory (%s): %v\n", binPath, err)
		}
		suffix := xbuild.GetSuffix(v2rayOS, v2rayArch)
		zipFile := "v2ray" + suffix + ".zip"
		root := filepath.Base(targetDir)

		var zipOptions []zip.Option
		if isOfficialBuild() {
			zipOptions = append(zipOptions, zip.With7Zip())
		}
		if err := zip.CompressFolder(root, zipFile, zipOptions...); err != nil {
			fmt.Printf("Unable to create archive (%s): %v\n", zipFile, err)
			return
		}

		metaWriter, err := build.NewFileMetadataWriter(filepath.Join(binPath, *flagMetadataFile))
		if err != nil {
			fmt.Println("Failed to create metadata writer: ", err)
			return
		}

		meta, err := build.GenerateFileMetadata(zipFile)
		if err != nil {
			fmt.Println("Failed to generate metadata for file: ", zipFile, err)
			return
		}

		metaWriter.Append(meta)
		metaWriter.Close()
	}
}
