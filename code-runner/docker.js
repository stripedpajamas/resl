const cp = require('child_process')

const buildDockerRunArgs = (image, name, executionDir, cmd) => {
  return [
    'run',
    '--rm', // delete after running
    '--name', // name the container for easier manual deletion, if necessary
    name,
    '--network', // use the following network setting for the container
    'none', // no networking
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

function runAndCollectOutput (cmd, args, opts = {}) {
  return new Promise((resolve, reject) => {
    let output = Buffer.alloc(0)
    const child = cp.spawn(cmd, args)

    let timeout
    if (opts.timeout) {
      timeout = setTimeout(() => {
        child.kill()
        const error = new Error('Execution timed out')
        error.code = 'TIMEOUT'
        reject(error)
      }, opts.timeout)
    }

    child.stdout.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })
    child.stderr.on('data', (chunk) => { output = Buffer.concat([output, chunk]) })

    child.on('error', (e) => { // error event is emitted when child could not be spawned
      clearTimeout(timeout)
      reject(e)
    })
    child.on('exit', (code) => {
      clearTimeout(timeout)
      resolve({ output: output.toString(), code })
    })
  })
}

async function runManualCleanup (name) {
  try {
    await runAndCollectOutput('docker', ['rm', '-f', name])
  } catch (_) {}
}

async function runWithArgs (image, name, executionDir, cmd) {
  try {
    const output = await runAndCollectOutput(
      'docker',
      buildDockerRunArgs(image, name, executionDir, cmd),
      { timeout: 30000 }
    )
    return output
  } catch (e) {
    if (e.code === 'TIMEOUT') {
      await runManualCleanup(name)
      throw e
    }
  }
}

function pull (image) {
  return runAndCollectOutput('docker', ['pull', image])
}

module.exports = { runWithArgs, pull }
