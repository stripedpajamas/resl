const fs = require('fs/promises')
const path = require('path')
const execa = require('execa')
const ULID = require('ulid')

exports.handler = async ({ code, props }) => {
  console.log('Running Code Handler')

  const { extension, runCmd } = props

	console.log({ extension, runCmd })

  const filePath = createFilePath('/tmp', extension)
  await writeCodeFile(filePath, code)

  let output

  try {
    output = await runCode(runCmd, [filePath])
  } catch (err) {
    output = err.all
  }
  
  await deleteCodeFile(filePath)
  return { output }
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

	const { all: output } = await execa(cmd, args, { all: true })

  return output
}
