const fs = require('fs')

const defaultServiceName = 'graphql-server'
const defaultServicePort = '5000'
const defaultGraphiQlEnabled = 'false'
const defaultJaegerAgentHost = 'localhost'
const defaultJaegerAgentPort = '6832'
const defaultJaegerLogSpans = 'false'
const defaultNatsServers = 'nats://localhost:4222'
const defaultNatsUser = 'client'
const defaultNatsPassword = 'pass'
const defaultSiteServiceAddr = 'localhost:4010'
const defaultSensorServiceAddr = 'localhost:4020'
const defaultSwitchServiceAddr = 'localhost:4030'

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
      graphiQlEnabled: this._getValue('GRAPHIQL_ENABLED', defaultGraphiQlEnabled) === 'true',
      jaegerAgentHost: this._getValue('JAEGER_AGENT_HOST', defaultJaegerAgentHost),
      jaegerAgentPort: parseInt(this._getValue('JAEGER_AGENT_PORT', defaultJaegerAgentPort)),
      jaegerLogSpans: this._getValue('JAEGER_LOG_SPANS', defaultJaegerLogSpans) === 'true',
      natsServers: this._getValue('NATS_SERVERS', defaultNatsServers).split(','),
      natsUser: this._getValue('NATS_USER', defaultNatsUser),
      natsPassword: this._getValue('NATS_PASSWORD', defaultNatsPassword),
      siteServiceAddr: this._getValue('SITE_SERVICE_ADDR', defaultSiteServiceAddr),
      sensorServiceAddr: this._getValue('SENSOR_SERVICE_ADDR', defaultSensorServiceAddr),
      switchServiceAddr: this._getValue('SWITCH_SERVICE_ADDR', defaultSwitchServiceAddr)
    }

    return Promise.resolve(config)
  }
}

module.exports = ConfigProvider
