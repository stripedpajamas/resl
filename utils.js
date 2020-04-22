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
  if (!installCmd || !packages || !packages.length) return null
  return [installCmd, ...packages, '> /dev/null 2>&1'].join(' ')
}

const buildCompileCommand = (config, name) => {
  const { compileCmd, compileExtension } = config
  if (!compileCmd || !compileExtension) return null
  return [compileCmd, buildFileName(name, compileExtension), '> /dev/null 2>&1'].join(' ')
}

const buildExecuteCommand = (config, name) => {
  const { executeCmd, executeExtension } = config
  if (!executeCmd || !executeExtension) return null
  return [executeCmd, buildFileName(name, executeExtension)].join(' ')
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

const getTemplate = (templateName) => require(`./templates/${templateName}`)

module.exports = {
  buildArgs,
  buildFileName,
  cleanUpFiles,
  getTemplate
}
