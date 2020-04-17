const fs = require('fs')
const cp = require('child_process')
const util = require('util')
const ulid = require('ulid')

const writeFile = util.promisify(fs.writeFile)
const deleteFile = util.promisify(fs.unlink)

const buildArgs = (name) => [
  'run',
  '--rm',
  '--name',
  name,
  '-v',
  `${process.cwd()}:/usr/src/app`,
  '-w',
  '/usr/src/app',
  'node:12-alpine',
  'node',
  `${name}.js`
]

module.exports = function run (code) {
  if (typeof code !== 'string') return ''
  return new Promise((resolve, reject) => {
    const name = ulid.ulid()
    writeFile(`${name}.js`, code).then(() => {
      let output = Buffer.alloc(0)

      const docker = cp.spawn('docker', buildArgs(name))
      docker.stdout.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })
      docker.stderr.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })

      docker.on('error', (err) => { reject(err) })
      docker.on('exit', () => {
        deleteFile(`${name}.js`).then(() => {
          resolve(output.toString())
        })
      })
    })
  })
}
