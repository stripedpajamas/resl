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

  if (!req.body) return
  req.log.info(req.body)

  const { text, response_url: responseUrl } = req.body
  const { language, code } = parseText(text)

  req.log.info({ language, code })

  let output
  try {
    output = await run(language, code)
    output = escapeCodeBlock(output)
  } catch (e) {
    req.log.error(e)
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
  req.log.info({ output })

  await got.post(responseUrl, {
    json: {
      response_type: 'in_channel',
      text: output ? '```' + output + '```' : '```[No output]```'
    }
  })
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
  let code = text.slice(firstBreakIdx + 1).trim()

  // remove possible backticks
  while (code[0] === '`') code = code.slice(1, code.length - 1)
  return { language, code }
}

const escapeCodeBlock = (text) => {
  return text.replace('`', '\\`')
}

const start = async () => {
  try {
    const [port = 3000, address = '127.0.0.1'] = process.argv.slice(2)
    await fastify.listen(port, address)
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}

start()
