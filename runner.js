const fs = require('fs').promises
const cp = require('child_process')
const ulid = require('ulid')

const buildArgs = (config, name) => {
  const { image, cmd, extension } = config
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
    cmd,
    `${name}${extension}`
  ]
}

const parseLanguage = (language) => {
  return new Promise((resolve, reject) => {
    fs.readFile('./languages.json').then(contents => {
      const languages = JSON.parse(contents)
      const config = languages[language]
      resolve(config)
    })
  })
}

const executeCode = (config, code) => {
  return new Promise((resolve, reject) => {
    const name = ulid.ulid()
    const fileName = `${name}${config.extension}`
    fs.writeFile(fileName, code).then(() => {
      let output = Buffer.alloc(0)
      const docker = cp.spawn('docker', buildArgs(config, name))
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

module.exports = function run (language, code) {
  if (typeof code !== 'string') return ''
  return parseLanguage(language)
    .then(config => {
      if (!config) return ''
      return executeCode(config, code)
    })
}
