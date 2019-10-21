const fs = require('fs')

const defaultNatsServers = 'nats://localhost:4222'
const defaultNatsUser = 'client'
const defaultNatsPassword = 'pass'

class Config {
  constructor () {
    this.natsServers = this._getValue('NATS_SERVERS', defaultNatsServers).split(',')
    this.natsUser = this._getValue('NATS_USER', defaultNatsUser)
    this.natsPassword = this._getValue('NATS_PASSWORD', defaultNatsPassword)
  }

  _getValue (name, defaultValue) {
    let value = process.env[name]
    if (value) {
      return value
    }

    const filepath = process.env[name + '_FILE']
    if (filepath) {
      value = fs.readFileSync(filepath)
      if (value) {
        return value.toString()
      }
    }

    return defaultValue
  }
}

module.exports = Config
