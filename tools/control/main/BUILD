load("//bazel:build.bzl", "foreign_go_binary")
load("//bazel:gpg.bzl", "gpg_sign")

package(default_visibility=["//visibility:public"])

foreign_go_binary(
  name = "v2ctl_windows_amd64",
  pkg = "v2ray.com/ext/tools/control/main",
  output = "v2ctl",
  os = "windows",
  arch = "amd64",
)

gpg_sign(
  name = "v2ctl_windows_amd64_sig",
  base = ":v2ctl_windows_amd64",
)

foreign_go_binary(
  name = "v2ctl_windows_x86",
  pkg = "v2ray.com/ext/tools/control/main",
  output = "v2ctl",
  os = "windows",
  arch = "x86",
)

foreign_go_binary(
  name = "v2ctl_linux_amd64",
  pkg = "v2ray.com/ext/tools/control/main",
  output = "v2ctl",
  os = "linux",
  arch = "amd64",
)
