const nats = require('nats')
const chalk = require('chalk')

const Config = require('./config')

class MockAssetService {
  constructor (config) {
    this.nats = nats.connect({
      servers: config.natsServers,
      user: config.natsUser,
      pass: config.natsPassword
    })
    this.subject = config.assetSubject
    this.queue = config.assetQueue
  }

  start () {
    console.log(chalk.green(`Mock NATS API Service started ...`))

    this.subscription = this.nats.subscribe(this.subject, { queue: this.queue }, msg => {
      console.log(chalk.blue(`Received message ${msg}`))
    })
  }

  stop () {
    this.nats.unsubscribe(this.subscription)
    this.nats.close()
  }
}

const config = new Config()
const assetService = new MockAssetService(config)
assetService.start()
