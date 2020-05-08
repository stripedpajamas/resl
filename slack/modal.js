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
  private_metadata: language,
  blocks: [
    {
      block_id: 'main_block',
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
    },
    {
      block_id: 'response_block',
      type: 'input',
      optional: true,
      label: {
        type: 'plain_text',
        text: 'Select a channel to post the result in'
      },
      element: {
        action_id: 'conversation_select_action',
        type: 'conversations_select',
        default_to_current_conversation: true,
        response_url_enabled: true
      }
    }
  ]
})
