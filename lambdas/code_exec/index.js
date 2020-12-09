const cp = require('child_process')
const fs = require('fs/promises')
const path = require('path')

exports.handler = async (event) => {
  const { language, code } = event
  
  const fp = await writeCodeFile('index.js', code)
  const out = await runCode(fp)

  return out
}

async function writeCodeFile(fileName, code) {
	console.log("Writing Code File")

	const filepath = path.join("/tmp", fileName)

	console.log("Writing file: %s \n", filepath)

  await fs.writeFile(filepath, code)

	console.log("Created file to execute code: %s \n", filepath)

	return filepath
}

async function runCode(filepath) {
  console.log("Running Code")

  return runAndCollectOutput('node', [filepath], {
    timeout: 10000
  })
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

