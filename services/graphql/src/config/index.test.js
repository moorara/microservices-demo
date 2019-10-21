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
    delete process.env.GRAPHIQL_ENABLED
    delete process.env.GRAPHIQL_ENABLED_FILE
    delete process.env.JAEGER_AGENT_HOST
    delete process.env.JAEGER_AGENT_HOST_FILE
    delete process.env.JAEGER_AGENT_PORT
    delete process.env.JAEGER_AGENT_PORT_FILE
    delete process.env.JAEGER_LOG_SPANS
    delete process.env.JAEGER_LOG_SPANS_FILE
    delete process.env.NATS_SERVERS
    delete process.env.NATS_SERVERS_FILE
    delete process.env.NATS_USER
    delete process.env.NATS_USER_FILE
    delete process.env.NATS_PASSWORD
    delete process.env.NATS_PASSWORD_FILE
    delete process.env.SITE_SERVICE_ADDR
    delete process.env.SITE_SERVICE_ADDR_FILE
    delete process.env.SENSOR_SERVICE_ADDR
    delete process.env.SENSOR_SERVICE_ADDR_FILE
    delete process.env.SWITCH_SERVICE_ADDR
    delete process.env.SWITCH_SERVICE_ADDR_FILE
  })

  describe('getConfig', () => {
    it('should return default values', done => {
      config.getConfig().then(c => {
        c.serviceName.should.equal('graphql-server')
        c.servicePort.should.equal(5000)
        c.graphiQlEnabled.should.equal(false)
        c.jaegerAgentHost.should.equal('localhost')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(false)
        c.natsServers.should.eql(['nats://localhost:4222'])
        c.natsUser.should.equal('client')
        c.natsPassword.should.equal('pass')
        c.siteServiceAddr.should.equal('localhost:4010')
        c.sensorServiceAddr.should.equal('localhost:4020')
        c.switchServiceAddr.should.equal('localhost:4030')
        done()
      }).catch(done)
    })

    it('should return values from environment variables', done => {
      process.env.SERVICE_NAME = 'my-service'
      process.env.SERVICE_PORT = '10000'
      process.env.GRAPHIQL_ENABLED = 'true'
      process.env.JAEGER_AGENT_HOST = 'jaeger-agent'
      process.env.JAEGER_AGENT_PORT = '6832'
      process.env.JAEGER_LOG_SPANS = true
      process.env.NATS_SERVERS = 'nats://nats-0:4222,nats://nats-1:4222,nats://nats-2:4222'
      process.env.NATS_USER = 'nats_client'
      process.env.NATS_PASSWORD = 'password!'
      process.env.SITE_SERVICE_ADDR = 'site-service:4010'
      process.env.SENSOR_SERVICE_ADDR = 'sensor-service:4020'
      process.env.SWITCH_SERVICE_ADDR = 'switch-service:4030'

      config.getConfig().then(c => {
        c.serviceName.should.equal('my-service')
        c.servicePort.should.equal(10000)
        c.graphiQlEnabled.should.equal(true)
        c.jaegerAgentHost.should.equal('jaeger-agent')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(true)
        c.natsServers.should.eql(['nats://nats-0:4222', 'nats://nats-1:4222', 'nats://nats-2:4222'])
        c.natsUser.should.equal('nats_client')
        c.natsPassword.should.equal('password!')
        c.siteServiceAddr.should.equal('site-service:4010')
        c.sensorServiceAddr.should.equal('sensor-service:4020')
        c.switchServiceAddr.should.equal('switch-service:4030')
        done()
      }).catch(done)
    })

    it('should return values from files', done => {
      // tmp will clean up after itself!
      // See https://raszi.github.io/node-tmp
      const nameFile = tmp.fileSync()
      const portFile = tmp.fileSync()
      const grahpiqlFile = tmp.fileSync()
      const agentHostFile = tmp.fileSync()
      const agentPortFile = tmp.fileSync()
      const logSpansFile = tmp.fileSync()
      const natsServersFile = tmp.fileSync()
      const natsUserFile = tmp.fileSync()
      const natsPassFile = tmp.fileSync()
      const siteFile = tmp.fileSync()
      const sensorFile = tmp.fileSync()
      const switchFile = tmp.fileSync()

      process.env.SERVICE_NAME_FILE = nameFile.name
      process.env.SERVICE_PORT_FILE = portFile.name
      process.env.GRAPHIQL_ENABLED_FILE = grahpiqlFile.name
      process.env.JAEGER_AGENT_HOST_FILE = agentHostFile.name
      process.env.JAEGER_AGENT_PORT_FILE = agentPortFile.name
      process.env.JAEGER_LOG_SPANS_FILE = logSpansFile.name
      process.env.NATS_SERVERS_FILE = natsServersFile.name
      process.env.NATS_USER_FILE = natsUserFile.name
      process.env.NATS_PASSWORD_FILE = natsPassFile.name
      process.env.SITE_SERVICE_ADDR_FILE = siteFile.name
      process.env.SENSOR_SERVICE_ADDR_FILE = sensorFile.name
      process.env.SWITCH_SERVICE_ADDR_FILE = switchFile.name

      fs.writeFileSync(nameFile.fd, 'new-service')
      fs.writeFileSync(portFile.fd, '20000')
      fs.writeFileSync(grahpiqlFile.fd, 'true')
      fs.writeFileSync(agentHostFile.fd, 'jaeger-agent')
      fs.writeFileSync(agentPortFile.fd, '6832')
      fs.writeFileSync(logSpansFile.fd, 'true')
      fs.writeFileSync(natsServersFile.fd, 'nats://nats-0:4222,nats://nats-1:4222,nats://nats-2:4222')
      fs.writeFileSync(natsUserFile.fd, 'nats_client')
      fs.writeFileSync(natsPassFile.fd, 'password!')
      fs.writeFileSync(siteFile.fd, 'site-service:4010')
      fs.writeFileSync(sensorFile.fd, 'sensor-service:4020')
      fs.writeFileSync(switchFile.fd, 'switch-service:4030')

      config.getConfig().then(c => {
        c.serviceName.should.equal('new-service')
        c.servicePort.should.equal(20000)
        c.graphiQlEnabled.should.equal(true)
        c.jaegerAgentHost.should.equal('jaeger-agent')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(true)
        c.natsServers.should.eql(['nats://nats-0:4222', 'nats://nats-1:4222', 'nats://nats-2:4222'])
        c.natsUser.should.equal('nats_client')
        c.natsPassword.should.equal('password!')
        c.siteServiceAddr.should.equal('site-service:4010')
        c.sensorServiceAddr.should.equal('sensor-service:4020')
        c.switchServiceAddr.should.equal('switch-service:4030')
        done()
      }).catch(done)
    })
  })
})
