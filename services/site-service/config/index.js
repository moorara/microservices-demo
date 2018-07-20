const fs = require('fs')
const Promise = require('bluebird')

const Logger = require('../util/logger')

const dbName = 'sites'
const defaultServiceName = 'site-service'
const defaultServicePort = '4010'
const defaultMongoURL = 'mongodb://localhost:27017'
const defaultJaegerAgentHost = 'localhost'
const defaultJaegerAgentPort = '6831'
const defaultJaegerReporterLogSpans = 'false'

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
      serviceName: this._getValue('SERVICE_NAME', defaultServiceName),
      servicePort: parseInt(this._getValue('SERVICE_PORT', defaultServicePort)),
      mongoUrl: this._getValue('MONGO_URL', defaultMongoURL) + `/${dbName}`,
      mongoUser: this._getValue('MONGO_USER'),
      mongoPass: this._getValue('MONGO_PASS'),
      jaegerAgentHost: this._getValue('JAEGER_AGENT_HOST', defaultJaegerAgentHost),
      jaegerAgentPort: parseInt(this._getValue('JAEGER_AGENT_PORT', defaultJaegerAgentPort)),
      jaegerReporterLogSpans: this._getValue('JAEGER_REPORTER_LOG_SPANS', defaultJaegerReporterLogSpans) === 'true'
    }

    return Promise.resolve(config)
  }
}

module.exports = ConfigProvider
