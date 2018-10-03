def _go_command(ctx):
  output = ctx.attr.output
  if ctx.attr.os == "windows":
    output = output + ".exe"

  output_file = ctx.actions.declare_file(ctx.attr.os + "/" + ctx.attr.arch + "/" + output)
  pkg = ctx.attr.pkg

  ld_flags = "-s -w"
  if ctx.attr.ld:
    ld_flags = ld_flags + " " + ctx.attr.ld

  options = [
    "go",
    "build",
    "-a", # force rebuild all. see https://github.com/golang/go/issues/27236
    "-o", output_file.path,
    "-compiler", "gc",
    "-gcflags", "-trimpath=${GOPATH}/src",
    "-asmflags", "-trimpath=${GOPATH}/src",
    "-ldflags", "'%s'" % ld_flags,
    pkg,
  ]

  command = " ".join(options)

  env_dict = {
    "CGO_ENABLED": "0"
  }
  
  env_dict["GOOS"] = ctx.attr.os
  env_dict["GOARCH"] = ctx.attr.arch
  if ctx.attr.mips: # https://github.com/golang/go/issues/27260
    env_dict["GOMIPS"] = ctx.attr.mips
    env_dict["GOMIPS64"] = ctx.attr.mips
    env_dict["GOMIPSLE"] = ctx.attr.mips
    env_dict["GOMIPS64LE"] = ctx.attr.mips
  if ctx.attr.arm:
    env_dict["GOARM"] = ctx.attr.arm

  for key, value in env_dict.items():
    print(key + " = " + value)

  ctx.actions.run_shell(
    outputs = [output_file],
    command = command,
    use_default_shell_env = True,
    env = env_dict,
  )
  runfiles = ctx.runfiles(files = [output_file])
  return [DefaultInfo(executable = output_file, runfiles = runfiles)]


foreign_go_binary = rule(
  _go_command,
  attrs = {
    'pkg': attr.string(),
    'output': attr.string(),
    'os': attr.string(mandatory=True),
    'arch': attr.string(mandatory=True),
    'mips': attr.string(),
    'arm': attr.string(),
    'ld': attr.string(),
  },
  executable = True,
)
