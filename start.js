// prime the docker images and then start the web server

const server = require('.')
const docker = require('./docker')
const languageDefs = require('./languages.json')

async function main () {
  if (process.argv[4] !== 'nopull') {
    console.log('Pulling Docker images...')

    for (const lang in languageDefs) {
      const { output, code } = await docker.pull(languageDefs[lang].image)

      console.log(output)

      if (code) {
        throw new Error(`Failed to pull Docker image for language ${lang}`)
      }
    }

    console.log('\n\nStarting web server...')
  }

  server.start()
}

main()
