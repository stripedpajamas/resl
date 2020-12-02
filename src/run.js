const cp = require('child_process')

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
      resolve({ output: output.toString(), exitCode: code })
    })
  })
}
exports.run = async (event) => {
  console.log(event)
  const { output, exitCode } = await runAndCollectOutput('node', ['-v'], {})

  return { output, exitCode }

  // writes code into file
  // runs code
  // collects output
  // returns output
}
