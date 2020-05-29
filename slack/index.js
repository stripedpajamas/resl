const got = require('got')
const getModal = require('./modal')

const { TOKEN } = process.env
const VIEWS_OPEN_URL = 'https://slack.com/api/views.open'

async function sendChannelResponse ({ responseUrl, text }) {
  const body = await got.post(responseUrl, {
    json: {
      response_type: 'in_channel',
      text
    }
  })

  return body
}

async function sendModal ({ triggerId, language }) {
  const { langName, placeholder } = language
  const body = await got.post(VIEWS_OPEN_URL, {
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

module.exports = { sendChannelResponse, sendModal }
