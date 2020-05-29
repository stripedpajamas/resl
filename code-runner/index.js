const run = require('./runner')
const languages = require('./languages.json')

const getLanguage = (lang) => {
  if (languages[lang]) {
    return languages[lang]
  }
  throw new Error('Unsupported language')
}

const getLanguageByName = (langName) => {
  for (const lang in languages) {
    if (languages[lang].langName === langName) return languages[lang]
  }
  throw new Error('Unsupported language')
}

const parseText = (text) => {
  // [lang]\s[code]
  const firstBreakIdx = text.split('').findIndex(x => /\s/.test(x))
  if (firstBreakIdx < 0) {
    return { language: getLanguage(text.trim()), code: '' }
  }

  const language = getLanguage(text.slice(0, firstBreakIdx))
  const code = parseCode(text.slice(firstBreakIdx + 1).trim())

  return { language, code }
}

const parseCode = (text) => {
  // remove possible backticks and unescape html sequences
  let code = text
    .replace(/&amp;/g, '&')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
  while (code[0] === '`') {
    code = code.slice(1, code.length - 1)
  }
  return code
}

const escapeCodeBlock = (text) => {
  return text.replace('`', '\\`')
}

module.exports = { run, escapeCodeBlock, parseText, parseCode, getLanguage, getLanguageByName }
