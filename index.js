const Fastify = require('fastify')
const codeRunner = require('./code-runner')
const slack = require('./slack')

const fastify = Fastify({ logger: true })
fastify.register(require('fastify-helmet'))
fastify.register(require('fastify-formbody'))

fastify.get('/', async () => {
  return { hello: 'world' }
})

fastify.post('/run', async (req, res) => {
  if (!req.body) return
  req.log.info(req.body)

  const {
    trigger_id: triggerId,
    text,
    response_url: responseUrl
  } = req.body

  let language, code
  try {
    ({ language, code } = codeRunner.parseText(text))
  } catch (e) {
    req.log.error(e)
    res.code(200).send({ text: e.message })
    return
  }

  if (!code) { // provide modal for them to fill in code
    res.code(200).send() // private acknowledgement to Slack
    const slackRes = await slack.sendModal({ triggerId, language })
    req.log.info(slackRes, 'Slack modal response')
    return
  }

  res.code(200).send({ response_type: 'in_channel' }) // tell Slack we got it, publicly
  let output
  try {
    output = await codeRunner.run(language, code)
    output = codeRunner.escapeCodeBlock(output)
  } catch (e) {
    req.log.error(e)
    // these errors occur during execution setup; compile/run errors are stuffed into `output`
    await slack.sendChannelResponse({ responseUrl, text: 'Sorry! Unable to setup execution environment :(' })
    return
  }
  req.log.info({ output })

  await slack.sendChannelResponse({ responseUrl, text: output ? '```' + output + '```' : '```[No output]```' })
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
