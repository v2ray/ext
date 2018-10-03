def _go_command(ctx):
  output = ctx.attr.output
  if ctx.attr.os == "windows":
    output = output + ".exe"

  output_file = ctx.actions.declare_file(output)
  pkg = ctx.attr.pkg
  command = "go build -o '%s' %s" % (output_file.path, pkg)

  env_dict = {
    "CGO_ENABLED": "0"
  }
  
  if ctx.attr.os:
    env_dict["GOOS"] = ctx.attr.os
  if ctx.attr.arch:
    env_dict["GOARCH"] = ctx.attr.arch
  if ctx.attr.mips:
    env_dict["GOMIPS"] = ctx.attr.mips
  if ctx.attr.mips64:
    env_dict["GOMIPS64"] = ctx.attr.mips64
  if ctx.attr.arm:
    env_dict["GOARM"] = ctx.attr.arm
  if ctx.attr.arm64:
    env_dict["GOARM64"] = ctx.attr.arm64

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
    'os': attr.string(),
    'arch': attr.string(),
    'mips': attr.string(),
    'mips64': attr.string(),
    'arm': attr.string(),
    'arm64': attr.string(),
  },
  executable = True,
)
