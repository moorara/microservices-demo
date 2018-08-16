/* eslint-env mocha */
require('should')
const fs = require('fs')
const tmp = require('tmp')

const ConfigProvider = require('.')

describe('ConfigProvider', () => {
  let config

  beforeEach(() => {
    config = new ConfigProvider()
  })

  afterEach(() => {
    delete process.env.SERVICE_NAME
    delete process.env.SERVICE_NAME_FILE
    delete process.env.SERVICE_PORT
    delete process.env.SERVICE_PORT_FILE
    delete process.env.JAEGER_AGENT_HOST
    delete process.env.JAEGER_AGENT_HOST_FILE
    delete process.env.JAEGER_AGENT_PORT
    delete process.env.JAEGER_AGENT_PORT_FILE
    delete process.env.JAEGER_LOG_SPANS
    delete process.env.JAEGER_LOG_SPANS_FILE
    delete process.env.SITE_SERVICE_ADDR
    delete process.env.SITE_SERVICE_ADDR_FILE
    delete process.env.SENSOR_SERVICE_ADDR
    delete process.env.SENSOR_SERVICE_ADDR_FILE
    delete process.env.SWITCH_SERVICE_ADDR
    delete process.env.SWITCH_SERVICE_ADDR_FILE
    delete process.env.GRAPHIQL_ENABLED
    delete process.env.GRAPHIQL_ENABLED_FILE
  })

  describe('getConfig', () => {
    it('should return default values', done => {
      config.getConfig().then(c => {
        c.serviceName.should.equal('graphql-server')
        c.servicePort.should.equal(5000)
        c.jaegerAgentHost.should.equal('localhost')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(false)
        c.siteServiceAddr.should.equal('localhost:4010')
        c.sensorServiceAddr.should.equal('localhost:4020')
        c.switchServiceAddr.should.equal('localhost:4030')
        c.graphiQlEnabled.should.equal(false)
        done()
      }).catch(done)
    })

    it('should return values from environment variables', done => {
      process.env.SERVICE_NAME = 'my-service'
      process.env.SERVICE_PORT = '10000'
      process.env.JAEGER_AGENT_HOST = 'jaeger-agent'
      process.env.JAEGER_AGENT_PORT = '6832'
      process.env.JAEGER_LOG_SPANS = true
      process.env.SITE_SERVICE_ADDR = 'site-service:4010'
      process.env.SENSOR_SERVICE_ADDR = 'sensor-service:4020'
      process.env.SWITCH_SERVICE_ADDR = 'switch-service:4030'
      process.env.GRAPHIQL_ENABLED = 'true'

      config.getConfig().then(c => {
        c.serviceName.should.equal('my-service')
        c.servicePort.should.equal(10000)
        c.jaegerAgentHost.should.equal('jaeger-agent')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(true)
        c.siteServiceAddr.should.equal('site-service:4010')
        c.sensorServiceAddr.should.equal('sensor-service:4020')
        c.switchServiceAddr.should.equal('switch-service:4030')
        c.graphiQlEnabled.should.equal(true)
        done()
      }).catch(done)
    })

    it('should return values from files', done => {
      // tmp will clean up after itself!
      // See https://raszi.github.io/node-tmp
      const nameFile = tmp.fileSync()
      const portFile = tmp.fileSync()
      const agentHostFile = tmp.fileSync()
      const agentPortFile = tmp.fileSync()
      const logSpansFile = tmp.fileSync()
      const siteFile = tmp.fileSync()
      const sensorFile = tmp.fileSync()
      const switchFile = tmp.fileSync()
      const grahpiqlFile = tmp.fileSync()

      process.env.SERVICE_NAME_FILE = nameFile.name
      process.env.SERVICE_PORT_FILE = portFile.name
      process.env.JAEGER_AGENT_HOST_FILE = agentHostFile.name
      process.env.JAEGER_AGENT_PORT_FILE = agentPortFile.name
      process.env.JAEGER_LOG_SPANS_FILE = logSpansFile.name
      process.env.SITE_SERVICE_ADDR_FILE = siteFile.name
      process.env.SENSOR_SERVICE_ADDR_FILE = sensorFile.name
      process.env.SWITCH_SERVICE_ADDR_FILE = switchFile.name
      process.env.GRAPHIQL_ENABLED_FILE = grahpiqlFile.name

      fs.writeFileSync(nameFile.fd, 'new-service')
      fs.writeFileSync(portFile.fd, '20000')
      fs.writeFileSync(agentHostFile.fd, 'jaeger-agent')
      fs.writeFileSync(agentPortFile.fd, '6832')
      fs.writeFileSync(logSpansFile.fd, 'true')
      fs.writeFileSync(siteFile.fd, 'site-service:4010')
      fs.writeFileSync(sensorFile.fd, 'sensor-service:4020')
      fs.writeFileSync(switchFile.fd, 'switch-service:4030')
      fs.writeFileSync(grahpiqlFile.fd, 'true')

      config.getConfig().then(c => {
        c.serviceName.should.equal('new-service')
        c.servicePort.should.equal(20000)
        c.jaegerAgentHost.should.equal('jaeger-agent')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(true)
        c.siteServiceAddr.should.equal('site-service:4010')
        c.sensorServiceAddr.should.equal('sensor-service:4020')
        c.switchServiceAddr.should.equal('switch-service:4030')
        c.graphiQlEnabled.should.equal(true)
        done()
      }).catch(done)
    })
  })
})
