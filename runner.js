const fs = require('fs').promises
const cp = require('child_process')
const ulid = require('ulid')
const languages = require('./languages')

const buildArgs = (config, name, file) => {
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
    buildRunCommand(config, file),
  ]
}

const buildRunCommand = (config, file) => {
  const cmd = [];
  const { runCmd, compileCmd, installCmd } = config

  installCmd && cmd.push(installCmd)

  compileCmd && cmd.push(`${compileCmd} ${file}`)

  cmd.push(`${runCmd} ${file}`)
  
  const output = cmd.join(' && ')
  
  return output
}

module.exports = function run (language, code) {
  const config = languages[language]
  if (typeof code !== 'string' || !config) return ''
  return new Promise((resolve, reject) => {
    const name = ulid.ulid()
    const fileName = `${name}.${config.extension}`
    fs.writeFile(fileName, code).then(() => {
      let output = Buffer.alloc(0)
      const docker = cp.spawn('docker', buildArgs(config, name, fileName))
      docker.stdout.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })
      docker.stderr.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })

      docker.on('error', (err) => { reject(err) })
      docker.on('exit', () => {
        fs.unlink(fileName).then(() => {
          resolve(output.toString())
        })
      })
    })
  })
}
