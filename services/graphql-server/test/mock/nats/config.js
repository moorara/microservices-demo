const fs = require('fs')

const defaultNatsServers = 'nats://localhost:4222'
const defaultNatsUser = 'client'
const defaultNatsPassword = 'pass'
const defaultAssetSubject = 'asset_service'
const defaultAssetQueue = 'workers'

class Config {
  constructor () {
    this.natsServers = this._getValue('NATS_SERVERS', defaultNatsServers).split(',')
    this.natsUser = this._getValue('NATS_USER', defaultNatsUser)
    this.natsPassword = this._getValue('NATS_PASSWORD', defaultNatsPassword)
    this.assetSubject = this._getValue('ASSET_SUBJECT', defaultAssetSubject)
    this.assetQueue = this._getValue('ASSET_QUEUE', defaultAssetQueue)
  }

  _getValue (name, defaultValue) {
    let value = process.env[name]
    if (value) {
      return value
    }

    let filepath = process.env[name + '_FILE']
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
