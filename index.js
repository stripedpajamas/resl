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
    const slackRes = await slack.sendChannelResponse({ responseUrl, text: 'Sorry! Unable to setup execution environment :(' })
    req.log.info(slackRes, 'Slack channel response')
    return
  }
  req.log.info({ output })

  const slackRes = await slack.sendChannelResponse({ responseUrl, text: output ? '```' + output + '```' : '```[No output]```' })
  req.log.info(slackRes, 'Slack channel response')
})

fastify.post('/modal', async (req, res) => {
  const payload = JSON.parse(req.body.payload || '{}')

  const { type, view, response_urls: responseUrls, user } = payload
  if (type !== 'view_submission' || !view || !responseUrls || !responseUrls.length || !user) {
    res.status(400).send()
    return
  }
  const { id: userId } = user
  const { response_url: responseUrl } = responseUrls.pop()

  req.log.info(view)
  const { state = {}, private_metadata: langName } = view
  const { values = {} } = state
  const { main_block: mainBlock = {} } = values
  const { code_input: codeInput = {} } = mainBlock
  const { value: rawCode } = codeInput

  if (!rawCode || !langName) {
    res.status(200).send({
      response_action: 'errors',
      errors: { main_block: 'Code is required to run' }
    })
    return
  }

  let code, language
  try {
    code = codeRunner.parseCode(rawCode)
    language = codeRunner.getLanguageByName(langName)
  } catch (e) {
    req.log.error(e)
    res.status(500).send()
    return
  }

  res.status(200).send({ response_action: 'clear' }) // this closes the modal

  let output
  try {
    output = await codeRunner.run(language, code)
    output = codeRunner.escapeCodeBlock(output)
  } catch (e) {
    req.log.error(e)
    // these errors occur during execution setup; compile/run errors are stuffed into `output`
    const slackRes = await slack.sendChannelResponse({ responseUrl, text: 'Sorry! Unable to setup execution environment :(' })
    req.log.info(slackRes, 'Slack channel response')
    return
  }
  req.log.info({ output })

  const tripleBackticks = '```' // including their code in the response text
  const text = `<@${userId}>\n${tripleBackticks}${code}${tripleBackticks}\n${tripleBackticks}${output || '[No output]'}${tripleBackticks}`

  const slackRes = await slack.sendChannelResponse({ responseUrl, text })
  req.log.info(slackRes, 'Slack channel response')
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
