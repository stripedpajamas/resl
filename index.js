const Fastify = require('fastify')
const run = require('./runner')

const fastify = Fastify({ logger: true })

fastify.get('/', async () => {
  return { hello: 'world' }
})

fastify.post('/run', async (req) => {
  if (!req.body) return {}

  const { code } = req.body
  return { output: await run(code) }
})

const start = async () => {
  try {
    await fastify.listen(3000)
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}

start()
