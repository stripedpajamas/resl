const test = require('ava')
const sinon = require('sinon')
const qs = require('qs')
const slack = require('../slack')
const codeRunner = require('../code-runner')
const { buildServer } = require('../')

// slack escapes html sequences
function slackifyText (text) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}

function buildModalPayload (rawCode, langName) {
  return JSON.stringify({
    type: 'view_submission',
    view: {
      state: {
        values: {
          main_block: {
            code_input: {
              value: rawCode
            }
          }
        }
      },
      private_metadata: langName
    },
    response_urls: [{ response_url: 'responseUrl1' }],
    user: {
      id: 'userId1'
    }
  })
}

test.beforeEach((t) => {
  const mockedCodeRunner = { ...codeRunner, run: sinon.stub() }
  const mockedSlack = sinon.stub({ ...slack })
  mockedSlack.verifyRequest.callsArg(2) // calls done()

  t.context.server = buildServer({ codeRunner: mockedCodeRunner, slack: mockedSlack, logger: false })

  t.context.slack = mockedSlack
  t.context.runner = mockedCodeRunner
})

test('GET `/` 200', async (t) => {
  const { server } = t.context
  const res = await server.inject({
    method: 'GET',
    url: '/'
  })

  t.is(res.statusCode, 200)
  t.deepEqual(JSON.parse(res.payload), { hello: 'world' })
})

test('POST `/run` 400 empty body', async (t) => {
  const { server } = t.context
  const res = await server.inject({
    method: 'POST',
    url: '/run',
    headers: { 'content-type': 'application/x-www-form-urlencoded' }
  })
  t.is(res.statusCode, 400)
})

test('POST `/run` 200 unsupported/empty language', async (t) => {
  const { server } = t.context
  const res = await server.inject({
    method: 'POST',
    url: '/run',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    payload: qs.stringify({
      trigger_id: 'triggerId',
      text: '',
      response_url: 'responseUrl'
    })
  })

  t.is(res.statusCode, 200)
  t.deepEqual(JSON.parse(res.payload), { text: 'Unsupported language' })
})

test('POST `/run` 200 no code means modal', async (t) => {
  const { server, slack } = t.context
  const res = await server.inject({
    method: 'POST',
    url: '/run',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    payload: qs.stringify({
      trigger_id: 'triggerId',
      text: 'js',
      response_url: 'responseUrl'
    })
  })

  t.is(res.statusCode, 200)
  t.is(res.payload, '') // empty response

  t.deepEqual(slack.sendModal.lastCall.args, [
    { triggerId: 'triggerId', language: codeRunner.getLanguage('js') }
  ])
})

test('POST `/run` 200 run failure', async (t) => {
  const { server, runner, slack } = t.context

  runner.run.rejects(new Error('Run failure'))

  const res = await server.inject({
    method: 'POST',
    url: '/run',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    payload: qs.stringify({
      trigger_id: 'triggerId',
      text: 'js `console.log("hi")`',
      response_url: 'responseUrl'
    })
  })

  t.is(res.statusCode, 200)
  t.deepEqual(JSON.parse(res.payload), { response_type: 'in_channel' })

  t.deepEqual(slack.sendChannelResponse.lastCall.args, [
    { responseUrl: 'responseUrl', text: 'Sorry! Unable to setup execution environment :(' }
  ])
})

test('POST `/run` 200 code will run', async (t) => {
  const { server, slack, runner } = t.context

  const tests = [
    {
      text: 'js `console.log("hi")`',
      expectedLanguage: 'js',
      expectedCode: 'console.log("hi")',
      codeOutput: 'hi',
      expectedReturnedOutput: '```hi```'
    },
    {
      text: 'js ```console.log("hi")```',
      expectedLanguage: 'js',
      expectedCode: 'console.log("hi")',
      codeOutput: 'hi',
      expectedReturnedOutput: '```hi```'
    },
    {
      text: 'js ```const a = async () => "hi"; a()```',
      expectedLanguage: 'js',
      expectedCode: 'const a = async () => "hi"; a()',
      codeOutput: '',
      expectedReturnedOutput: '```[No output]```'
    },
    {
      text: 'py ```print("hello")```',
      expectedLanguage: 'py',
      expectedCode: 'print("hello")',
      codeOutput: 'hello',
      expectedReturnedOutput: '```hello```'
    }
  ]

  for (const testCase of tests) {
    const { text, expectedLanguage, expectedCode, codeOutput, expectedReturnedOutput } = testCase

    runner.run.resolves(codeOutput)

    const res = await server.inject({
      method: 'POST',
      url: '/run',
      headers: { 'content-type': 'application/x-www-form-urlencoded' },
      payload: qs.stringify({
        trigger_id: 'triggerId',
        text: slackifyText(text),
        response_url: 'responseUrl'
      })
    })

    t.is(res.statusCode, 200)
    t.deepEqual(JSON.parse(res.payload), { response_type: 'in_channel' })

    t.deepEqual(runner.run.lastCall.args, [
      codeRunner.getLanguage(expectedLanguage),
      expectedCode
    ])
    t.deepEqual(slack.sendChannelResponse.lastCall.args, [
      { responseUrl: 'responseUrl', text: expectedReturnedOutput }
    ])
  }
})

test('POST `/modal` 400 bad input', async (t) => {
  const { server } = t.context

  const badPayloads = [
    { type: 'not view submission' },
    { type: 'view_submission' },
    { type: 'view_submission', view: {} },
    { type: 'view_submission', view: {}, response_urls: [] },
    { type: 'view_submission', view: {}, response_urls: ['r'] }
  ]

  for (const payload of badPayloads) {
    const res = await server.inject({
      method: 'POST',
      url: '/modal',
      headers: { 'content-type': 'application/x-www-form-urlencoded' },
      payload
    })
    t.is(res.statusCode, 400)
  }
})

test('POST `/modal` 200 no code', async (t) => {
  const { server } = t.context

  const res = await server.inject({
    method: 'POST',
    url: '/modal',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    payload: qs.stringify({ payload: buildModalPayload('', 'JavaScript') })
  })

  t.is(res.statusCode, 200)
  t.deepEqual(JSON.parse(res.payload), {
    response_action: 'errors',
    errors: { main_block: 'Code and language are required to run' }
  })
})

test('POST `/modal` 200 no lang', async (t) => {
  const { server } = t.context

  const res = await server.inject({
    method: 'POST',
    url: '/modal',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    payload: qs.stringify({ payload: buildModalPayload('code is here', '') })
  })

  t.is(res.statusCode, 200)
  t.deepEqual(JSON.parse(res.payload), {
    response_action: 'errors',
    errors: { main_block: 'Code and language are required to run' }
  })
})

test('POST `/modal` 200 cannot parse code or lang', async (t) => {
  const { server } = t.context

  const res = await server.inject({
    method: 'POST',
    url: '/modal',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    payload: qs.stringify({ payload: buildModalPayload('good code', 'bad lang') })
  })

  t.is(res.statusCode, 500)
})

test('POST `/modal` 200 run failure', async (t) => {
  const { server, runner, slack } = t.context

  runner.run.rejects(new Error('Run failure'))

  const res = await server.inject({
    method: 'POST',
    url: '/modal',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    payload: qs.stringify({ payload: buildModalPayload('hello', 'JavaScript') })
  })

  t.is(res.statusCode, 200)
  t.deepEqual(JSON.parse(res.payload), { response_action: 'clear' })

  t.deepEqual(slack.sendChannelResponse.lastCall.args, [
    { responseUrl: 'responseUrl1', text: 'Sorry! Unable to setup execution environment :(' }
  ])
})

test('POST `/modal` 200 will run code', async (t) => {
  const { server, runner, slack } = t.context

  const tests = [
    {
      code: '`let a = 123`',
      langName: 'JavaScript',
      codeOutput: '',
      expectedCode: 'let a = 123',
      expectedLanguage: 'js',
      expectedReturnedOutput: '<@userId1>\n```let a = 123```\n```[No output]```'
    },
    {
      code: '```console.log("hello world")```',
      langName: 'JavaScript',
      codeOutput: 'hello world',
      expectedCode: 'console.log("hello world")',
      expectedLanguage: 'js',
      expectedReturnedOutput: '<@userId1>\n```console.log("hello world")```\n```hello world```'
    }
  ]

  for (const testCase of tests) {
    const { code, langName, codeOutput, expectedCode, expectedLanguage, expectedReturnedOutput } = testCase

    runner.run.resolves(codeOutput)

    const res = await server.inject({
      method: 'POST',
      url: '/modal',
      headers: { 'content-type': 'application/x-www-form-urlencoded' },
      payload: qs.stringify({ payload: buildModalPayload(code, langName) })
    })

    t.is(res.statusCode, 200)
    t.deepEqual(JSON.parse(res.payload), { response_action: 'clear' })

    t.deepEqual(runner.run.lastCall.args, [
      codeRunner.getLanguage(expectedLanguage),
      expectedCode
    ])
    t.deepEqual(slack.sendChannelResponse.lastCall.args, [
      { responseUrl: 'responseUrl1', text: expectedReturnedOutput }
    ])
  }
})
