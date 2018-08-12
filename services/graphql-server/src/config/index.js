const fs = require('fs')

const defaultServiceName = 'graphql-server'
const defaultServicePort = '5000'
const defaultJaegerAgentHost = 'localhost'
const defaultJaegerAgentPort = '6832'
const defaultJaegerLogSpans = 'false'

class ConfigProvider {
  constructor (options) {
    options = options || {}
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
      jaegerAgentHost: this._getValue('JAEGER_AGENT_HOST', defaultJaegerAgentHost),
      jaegerAgentPort: parseInt(this._getValue('JAEGER_AGENT_PORT', defaultJaegerAgentPort)),
      jaegerLogSpans: this._getValue('JAEGER_LOG_SPANS', defaultJaegerLogSpans) === 'true'
    }

    return Promise.resolve(config)
  }
}

module.exports = ConfigProvider