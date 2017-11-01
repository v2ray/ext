package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"v2ray.com/ext/tools/build"
)

var (
	flagTargetDir    = flag.String("dir", "", "Directory to put generated files.")
	flagTargetOS     = flag.String("os", runtime.GOOS, "Target OS of this build.")
	flagTargetArch   = flag.String("arch", runtime.GOARCH, "Target CPU arch of this build.")
	flagArchive      = flag.Bool("zip", false, "Whether to make an archive of files or not.")
	flagMetadataFile = flag.String("metadata", "metadata.txt", "File to store metadata info of released packages.")
	flagSignBinary   = flag.Bool("sign", false, "Whether or not to sign the binaries.")
	flagEncryptedZip = flag.Bool("encrypt", false, "Also generate encrypted zip files.")

	binPath string
)

func createTargetDirectory(version string, goOS build.GoOS, goArch build.GoArch) (string, error) {
	var targetDir string
	if len(*flagTargetDir) > 0 {
		targetDir = *flagTargetDir
	} else {
		suffix := build.GetSuffix(goOS, goArch)

		targetDir = filepath.Join(binPath, "v2ray-"+version+suffix)
		if version != "custom" {
			os.RemoveAll(targetDir)
		}
	}

	err := os.MkdirAll(targetDir, os.ModeDir|0777)
	return targetDir, err
}

func getTargetFile(name string, goOS build.GoOS) string {
	suffix := ""
	if goOS == build.Windows {
		suffix += ".exe"
	}
	return name + suffix
}

func getBinPath() string {
	GOPATH := os.Getenv("GOPATH")
	return filepath.Join(GOPATH, "bin")
}

func main() {
	flag.Parse()
	binPath = getBinPath()

	v2rayOS := build.ParseOS(*flagTargetOS)
	v2rayArch := build.ParseArch(*flagTargetArch)

	version := os.Getenv("TRAVIS_TAG")

	if len(version) == 0 {
		version = "custom"
	}

	fmt.Printf("Building V2Ray (%s) for %s %s\n", version, v2rayOS, v2rayArch)

	targetDir, err := createTargetDirectory(version, v2rayOS, v2rayArch)
	if err != nil {
		fmt.Println("Unable to create directory " + targetDir + ": " + err.Error())
	}

	targetFile := getTargetFile("v2ray", v2rayOS)
	targetFileFull := filepath.Join(targetDir, targetFile)
	if err := build.BuildV2RayCore(targetFileFull, v2rayOS, v2rayArch, false); err != nil {
		fmt.Println("Unable to build V2Ray: " + err.Error())
		return
	}
	if v2rayOS == build.Windows {
		if err := build.BuildV2RayCore(filepath.Join(targetDir, "w"+targetFile), v2rayOS, v2rayArch, true); err != nil {
			fmt.Println("Unable to build V2Ray no console: " + err.Error())
			return
		}
	}

	confUtil := getTargetFile("v2ctl", v2rayOS)
	confUtilFull := filepath.Join(targetDir, confUtil)
	if err := build.GoBuild("v2ray.com/ext/tools/control/main", confUtilFull, v2rayOS, v2rayArch, ""); err != nil {
		fmt.Println("Unable to build V2Ray control: " + err.Error())
		return
	}

	if *flagSignBinary {
		gpgPass := os.Getenv("GPG_SIGN_PASS")
		//if err != nil {
		//	fmt.Println("Unable get GPG pass: " + err.Error())
		//	return
		//}
		if err := build.GPGSignFile(targetFileFull, gpgPass); err != nil {
			fmt.Println("Unable to sign V2Ray binary: " + err.Error())
			return
		}

		if err := build.GPGSignFile(confUtilFull, gpgPass); err != nil {
			fmt.Println("Unable to sign control util: " + err.Error())
			return
		}

		if v2rayOS == build.Windows {
			if err := build.GPGSignFile(filepath.Join(targetDir, "w"+targetFile), gpgPass); err != nil {
				fmt.Println("Unable to sign V2Ray no console: " + err.Error())
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
		suffix := build.GetSuffix(v2rayOS, v2rayArch)
		zipFile := "v2ray" + suffix + ".zip"
		root := filepath.Base(targetDir)
		err = build.SevenZipFolder(root, zipFile)
		if err != nil {
			fmt.Printf("Unable to create archive (%s): %v\n", zipFile, err)
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

		if *flagEncryptedZip {
			if err := build.SevenZipBuild(root, "vencrypted"+suffix+".7z", meta.Checksum()); err != nil {
				fmt.Println("Failed to generate encrypted zip file.")
			}
		}
	}
}
