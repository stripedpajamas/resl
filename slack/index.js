const got = require('got')
const getModal = require('./modal')

const { TOKEN } = process.env
const VIEWS_OPEN_URL = 'https://slack.com/api/views.open'

const sendChannelResponse = ({ responseUrl, text }) => got.post(responseUrl, {
  json: {
    response_type: 'in_channel',
    text
  }
})

const sendModal = ({ triggerId, language }) => got.post(VIEWS_OPEN_URL, {
  json: {
    triggerId,
    view: getModal(language.name, language.placeholder)
  },
  headers: {
    authorization: `Bearer ${TOKEN}`
  }
})

module.exports = { sendChannelResponse, sendModal }
