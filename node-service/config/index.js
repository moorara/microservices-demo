const fs = require('fs')
const Promise = require('bluebird')

const Logger = require('../util/logger')

const dbName = 'toys'
const defaultServicePort = '4020'
const defaultMongoURL = 'mongodb://localhost:27017'

class ConfigProvider {
  constructor (options) {
    options = options || {}
    this.logger = options.logger || new Logger('ConfigProvider')
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

  getConfig () {
    const config = {
      servicePort: this._getValue('SERVICE_PORT', defaultServicePort),
      mongoUrl: this._getValue('MONGO_URL', defaultMongoURL) + `/${dbName}`,
      mongoUser: this._getValue('MONGO_USER'),
      mongoPass: this._getValue('MONGO_PASS')
    }

    return Promise.resolve(config)
  }
}

module.exports = ConfigProvider
