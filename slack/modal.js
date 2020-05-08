module.exports = (language, placeholder) => ({
  type: 'modal',
  title: {
    type: 'plain_text',
    text: 'RESL',
    emoji: true
  },
  submit: {
    type: 'plain_text',
    text: 'Run Code',
    emoji: true
  },
  close: {
    type: 'plain_text',
    text: 'Cancel',
    emoji: true
  },
  blocks: [
    {
      type: 'input',
      element: {
        type: 'plain_text_input',
        action_id: 'code_input',
        multiline: true,
        placeholder: {
          type: 'plain_text',
          text: placeholder
        }
      },
      label: {
        type: 'plain_text',
        text: `Enter ${language} here`
      },
      hint: {
        type: 'plain_text',
        text: 'Wrapping your code in backticks is optional'
      }
    }
  ]
})
