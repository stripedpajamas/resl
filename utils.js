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
  const { executeCmd, extension } = config
  if (!executeCmd || !extension) return ''
  return `${executeCmd} ${buildFileName(name, extension)}`
}

const buildFileName = (name, extension) => {
  return `${name}.${extension}`
}

const cleanUpFiles = (config, name) => {
  const deletes = config.outputExtensions.map(ext => (
    new Promise((resolve, reject) => {
      fs.unlink(`${name}.${ext}`)
        .then(() => resolve())
        .catch(err => reject(err))
    })
  ))

  return Promise.resolve(deletes)
}

module.exports = {
  buildArgs,
  buildFileName,
  cleanUpFiles
}
