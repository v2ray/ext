load("//bazel:build.bzl", "foreign_go_binary")
load("//bazel:gpg.bzl", "gpg_sign")

def gen_targets(matrix):
  for (os, arch) in matrix:
    bin_name = "v2ctl_" + os + "_" + arch
    foreign_go_binary(
      name = bin_name,
      pkg = "v2ray.com/ext/tools/control/main",
      output = "v2ctl",
      os = os,
      arch = arch,
    )

    gpg_sign(
      name = bin_name + "_sig",
      base = ":" + bin_name,
    )
