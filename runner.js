const fs = require('fs').promises
const cp = require('child_process')
const ulid = require('ulid')
const languages = require('./languages')
const utils = require('./utils')

module.exports = function run (language, code) {
  const config = languages[language]
  if (typeof code !== 'string' || !config) return ''

  const { executeExtension, compileExtension } = config

  return new Promise((resolve, reject) => {
    const name = ulid.ulid()
    const fileExtension = compileExtension || executeExtension
    const fileName = utils.buildFileName(name, fileExtension)

    fs.writeFile(fileName, code).then(() => {
      let output = Buffer.alloc(0)
      const docker = cp.spawn('docker', utils.buildArgs(config, name, fileName))
      docker.stdout.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })
      docker.stderr.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })

      docker.on('error', (err) => {
        utils.cleanUpFiles(config, name)
          .then(() => reject(err))
          .catch(() => reject(err))
      })
      docker.on('exit', () => {
        console.log(output.toString())
        utils.cleanUpFiles(config, name)
          .then(() => resolve(output.toString()))
          .catch(err => reject(err))
      })
    })
  })
}
