const got = require('got')
const run = require('./runner')

const parseText = (text) => {
  // [lang]\s[code]
  const firstBreakIdx = text.split('').findIndex(x => /\s/.test(x))
  const language = text.slice(0, firstBreakIdx)
  let code = text.slice(firstBreakIdx + 1).trim()

  // remove possible backticks
  while (code[0] === '`') code = code.slice(1, code.length - 1)
  return { language, code }
}

const escapeCodeBlock = (text) => {
  return text.replace('`', '\\`')
}

const runHandler = async ({ body, log }) => {
  if (!body) return
  log.info(body)

  const { text, response_url: responseUrl } = body
  const { language, code } = parseText(text)

  log.info({ language, code })

  let output
  try {
    output = await run(language, code)
    output = escapeCodeBlock(output)
  } catch (e) {
    log.error(e)
    // these errors occur during execution setup;
    // compile/run errors are stuffed into `output`
    await got.post(responseUrl, {
      json: {
        response_type: 'in_channel',
        text: 'Sorry! Unable to setup execution environment :('
      }
    })
    return
  }
  log.info({ output })

  await got.post(responseUrl, {
    json: {
      response_type: 'in_channel',
      text: output ? '```' + output + '```' : '```[No output]```'
    }
  })
}

module.exports = { runHandler }
