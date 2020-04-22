const fs = require('fs').promises

const buildArgs = (config, name) => {
  const { image } = config
  return [
    'run',
    '--rm',
    '--name',
    name,
    '-v',
    `${process.cwd()}:/usr/src/app`,
    '-w',
    '/usr/src/app',
    image,
    'sh',
    '-c',
    buildRunCommand(config, name)
  ]
}

const buildRunCommand = (config, name) => {
  const output = [
    buildInstallCommand(config),
    buildCompileCommand(config, name),
    buildExecuteCommand(config, name)
  ]

  return output.filter(cmd => cmd).join(' && ')
}

const buildInstallCommand = (config) => {
  const { installCmd, packages } = config
  if (!installCmd || !packages || !packages.length) return ''
  return [installCmd, ...packages].join(' ')
}

const buildCompileCommand = (config, name) => {
  const { compileCmd, compileExtension } = config
  if (!compileCmd || !compileExtension) return ''
  return `${compileCmd} ${buildFileName(name, compileExtension)}`
}

const buildExecuteCommand = (config, name) => {
  const { executeCmd, executeExtension } = config
  if (!executeCmd || !executeExtension) return ''
  return `${executeCmd} ${buildFileName(name, executeExtension)}`
}

const buildFileName = (name, extension) => {
  return `${name}.${extension}`
}

const cleanUpFiles = (config, name) => {
  const deletes = config.outputExtensions.map(ext => {
    return fs.unlink(buildFileName(name, ext))
  })

  return Promise.all(deletes)
}

module.exports = {
  buildArgs,
  buildFileName,
  cleanUpFiles
}
