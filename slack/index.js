const got = require('got')
const { createHmac } = require('crypto')
const getModal = require('./modal')

const { TOKEN } = process.env
const VIEWS_OPEN_URL = 'https://slack.com/api/views.open'

function verifyRequest (req, res, done) {
  const { SIGNING_SECRET } = process.env
  const { body, headers } = req
  const { raw: rawBody } = body
  const { 'x-slack-request-timestamp': timestamp, 'x-slack-signature': sig } = headers

  if (Math.abs((Date.now() / 1000) - timestamp) > 60 * 5) {
    // if request is not dated within 5 mins (past or future) reject it
    req.log.warn(`Invalid timestamp: ${timestamp}`)
    return res.code(401).send()
  }

  const validationString = `v0:${timestamp}:${rawBody}`
  const hash = createHmac('sha256', SIGNING_SECRET)
    .update(validationString)
    .digest('hex')

  if (`v0=${hash}` !== sig) {
    req.log.warn('Invalid signature. Wanted %s but got %s', hash, sig)
    return res.code(401).send()
  }

  done()
}

async function sendChannelResponse ({ responseUrl, text }) {
  const { body } = await got.post(responseUrl, {
    json: {
      response_type: 'in_channel',
      text
    }
  })

  return body
}

async function sendModal ({ triggerId, language }) {
  const { langName, placeholder } = language
  const { body } = await got.post(VIEWS_OPEN_URL, {
    json: {
      trigger_id: triggerId,
      view: getModal(langName, placeholder)
    },
    headers: {
      authorization: `Bearer ${TOKEN}`
    }
  })

  return body
}

module.exports = { sendChannelResponse, sendModal, verifyRequest }
