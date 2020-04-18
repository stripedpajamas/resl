const Fastify = require('fastify')
const got = require('got')
const run = require('./runner')

const fastify = Fastify({ logger: true })
fastify.register(require('fastify-helmet'))
fastify.register(require('fastify-formbody'))

fastify.get('/', async () => {
  return { hello: 'world' }
})

fastify.post('/run', async (req, res) => {
  res.code(200).send({ response_type: 'in_channel' }) // tell Slack we got it

  try {
    if (!req.body) return
    req.log.info(req.body)

    const { text, response_url: responseUrl } = req.body
    const { language, code } = parseText(text)
    const output = await run(language, code)
    const escapedOutput = escapeCodeBlock(output)

    await got.post(responseUrl, {
      json: {
        response_type: 'in_channel',
        text: escapedOutput ? '```' + escapedOutput + '```' : '```[No output]```'
      }
    })
  } catch (e) {
    req.log.error(e)
  }
})

// useful for testing
fastify.post('/done', async (req, res) => {
  req.log.info(req.body)
  return {}
})

const parseText = (text) => {
  // [lang]\s[code]
  const firstBreakIdx = text.split('').findIndex(x => /\s/.test(x))
  const language = text.slice(0, firstBreakIdx)
  let code = text.slice(firstBreakIdx + 1)

  // remove possible backticks
  while (code[0] === '`') code = code.slice(1, code.length - 1)
  return { language, code }
}

const escapeCodeBlock = (text) => {
  return text.replace('`', '\\`')
}

const start = async () => {
  try {
    const [port, address] = process.argv.slice(2)
    await fastify.listen(port, address)
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}

start()
