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

  const { text, response_url: responseUrl } = req.body

  const { language, code } = parseText(text)
  const output = await run(language, code)

  got.post(responseUrl, {
    json: {
      text: '```' + output + '```'
    }
  })
})

// useful for testing
fastify.post('/done', async (req, res) => {
  req.log.info(req.body)
  return {}
})

function parseText (text) {
  const [language, ...codeWords] = text.split(' ')
  let code = codeWords.join(' ')
  while (code[0] === '`') code = code.slice(1, code.length - 1)
  return { language, code }
}

const start = async () => {
  try {
    await fastify.listen(3000)
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}

start()
