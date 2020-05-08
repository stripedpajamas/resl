const cp = require('child_process')

const buildDockerRunArgs = (image, executionDir, cmd) => {
  return [
    'run',
    '--rm', // delete after running
    '--stop-timeout', // stop the container after
    '10', // 10 seconds
    '-v', // mount the following path
    `${executionDir}:/usr/src/app`, // host_path:container_path
    '-w', // set the container's working directory to the following path
    '/usr/src/app',
    image, // use this image from Docker Hub
    'sh', // run in `sh`
    '-c', // run the following command
    cmd
  ]
}

function runAndCollectOutput (cmd, args) {
  return new Promise((resolve, reject) => {
    let output = Buffer.alloc(0)
    const child = cp.spawn(cmd, args)

    child.stdout.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })
    child.stderr.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })

    child.on('error', (e) => reject(e)) // error event is emitted when child could not be spawned
    child.on('exit', (code) => resolve({ output: output.toString(), code }))
  })
}

function runWithArgs (image, executionDir, cmd) {
  return runAndCollectOutput('docker', buildDockerRunArgs(image, executionDir, cmd))
}

function pull (image) {
  return runAndCollectOutput('docker', ['pull', image])
}

module.exports = { runWithArgs, pull }
