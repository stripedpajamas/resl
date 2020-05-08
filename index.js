const Fastify = require('fastify')
const { runHandler } = require('./code-runner')

const fastify = Fastify({ logger: true })
fastify.register(require('fastify-helmet'))
fastify.register(require('fastify-formbody'))

fastify.get('/', async () => {
  return { hello: 'world' }
})

fastify.post('/run', async (req, res) => {
  res.code(200).send({ response_type: 'in_channel' }) // tell Slack we got it

  const { body, log } = req
  runHandler({ body, log })
})

// useful for testing
fastify.post('/done', async (req, res) => {
  req.log.info(req.body)
  return {}
})

const start = async () => {
  try {
    const [port = 3000, address = '127.0.0.1'] = process.argv.slice(2)
    await fastify.listen(port, address)
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}

module.exports = { start }
