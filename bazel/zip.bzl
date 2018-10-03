def _pkg_zip_impl(ctx):
  dir = ctx.actions.declare_directory(ctx.label.name, sibling=ctx.outputs.out)

  commands = []
  commands.append("mkdir -p " + dir.path)
  commands.extend(["cp " + dep.path + " " + dir.path + "/" + dep.basename for dep in ctx.files.srcs])
  commands.append("zip -rj " + ctx.outputs.out.path + " " + dir.path)
  ctx.actions.run_shell(
    #command = "zip -j " + ctx.outputs.out.path + " " + " ".join([dep.path for dep in ctx.files.srcs]),
    command = " && ".join(commands),
    inputs = ctx.files.srcs,
    outputs = [dir, ctx.outputs.out],
    progress_message = "Creating .zip archive",
    mnemonic = "Zip",
  )

pkg_zip = rule(
  implementation = _pkg_zip_impl,
  attrs = {
    "out": attr.string(),
    "srcs": attr.label_list(allow_files=True),
  },
  outputs = {
    "out": "%{out}",
  },
)
