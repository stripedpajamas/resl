const fs = require('fs').promises
const path = require('path')
const ulid = require('ulid')
const mkdir = require('make-dir')
const del = require('del')
const docker = require('./docker')
const languages = require('./languages.json')

async function setup (config, code) {
  const { fileName, template } = config

  const codeToRun = template
    ? require(`./templates/${template}`)(code)
    : code

  // create execution directory
  const executionDirName = ulid.ulid()
  const executionDir = await mkdir(executionDirName)

  // write source code file into dir
  await fs.writeFile(path.join(executionDir, fileName), codeToRun)

  return { name: executionDirName, executionDir }
}

async function compile (config, name, executionDir) {
  const { image, compileCmd } = config

  if (!compileCmd) return
  const { output, code } = await docker.runWithArgs(image, name, executionDir, compileCmd)

  if (code) {
    throw new Error(output)
  }
}

async function run (config, name, executionDir) {
  const { image, runCmd } = config
  const { output, code } = await docker.runWithArgs(image, name, executionDir, runCmd)

  if (code) {
    throw new Error(output)
  }

  return output
}

async function cleanup (executionDir) {
  return del(executionDir)
}

module.exports = async function (language, code) {
  const config = languages[language]
  if (typeof code !== 'string' || !config) return ''

  const { name, executionDir } = await setup(config, code)

  try {
    await compile(config, name, executionDir)
    const output = await run(config, name, executionDir)
    await cleanup(executionDir)

    return output
  } catch (e) {
    try {
      await cleanup(executionDir)
    } catch (_) {} // don't care if cleanup doesn't work
    return e.message
  }
}
