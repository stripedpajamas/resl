const fs = require('fs/promises')
const path = require('path')
const execa = require('execa')
const ULID = require('ulid')

exports.handler = async ({ code, props }) => {
	if (props.warmup) return { temp: '9000' }

  console.log('Running Code Handler')

  const { extension, runCmd } = props

	console.log({ extension, runCmd })

  const filePath = createFilePath('/tmp', extension)
  await writeCodeFile(filePath, code)

	const output = await runCode(runCmd, [filePath])
	await deleteCodeFile(filePath)

	return output
}

const createFilePath = (folder, extension) => {
  const file = `${ULID.ulid()}.${extension}`
  return path.join(folder, file)
}

const writeCodeFile = async (filePath, code) => {
	console.log(`Writing file: ${filePath}`)

  await fs.writeFile(filePath, code)

	console.log(`Created file to execute code: ${filePath}`)
}

const deleteCodeFile = async (filePath) => {
  console.log(`Deleting file: ${filePath}`)
  await fs.unlink(filePath)
}

const runCode = async (cmd, args) => {
  console.log('Running Code')

	const subprocess = execa(cmd, args, { all: true })

	const timeout = setTimeout(() => {
		subprocess.kill('SIGTERM', {
			forceKillAfterTimeout: 10000
		});
	}, 8000);

  try {
		const { all: output } = await subprocess
		clearTimeout(timeout)

		return { output }
  } catch (err) {
		if (err.killed || err.isCanceled) {
			return { output: 'Execution timed out' }
		}

		return { output: err.all }
	}
}

