/* eslint-env mocha */
const fs = require('fs')
const tmp = require('tmp')
const should = require('should')

const ConfigProvider = require('../../../config')

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
    delete process.env.MONGO_URI
    delete process.env.MONGO_URI_FILE
    delete process.env.MONGO_USERNAME
    delete process.env.MONGO_USERNAME_FILE
    delete process.env.MONGO_PASSWORD
    delete process.env.MONGO_PASSWORD_FILE
    delete process.env.JAEGER_AGENT_HOST
    delete process.env.JAEGER_AGENT_HOST_FILE
    delete process.env.JAEGER_AGENT_PORT
    delete process.env.JAEGER_AGENT_PORT_FILE
    delete process.env.JAEGER_LOG_SPANS
    delete process.env.JAEGER_LOG_SPANS_FILE
  })

  describe('getConfig', () => {
    it('should return default values', done => {
      config.getConfig().then(c => {
        c.serviceName.should.equal('site-service')
        c.servicePort.should.equal(4010)
        c.mongoUri.should.equal('mongodb://localhost:27017/sites')
        should.equal(c.mongoUsername, null)
        should.equal(c.mongoPassword, null)
        c.jaegerAgentHost.should.equal('localhost')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(false)
        done()
      }).catch(done)
    })

    it('should return values from environment variables', done => {
      process.env.SERVICE_NAME = 'my-service'
      process.env.SERVICE_PORT = '10000'
      process.env.MONGO_URI = 'mongodb://mongo.mlab.com:27017'
      process.env.MONGO_USERNAME = 'user'
      process.env.MONGO_PASSWORD = 'pass'
      process.env.JAEGER_AGENT_HOST = 'jaeger-agent'
      process.env.JAEGER_AGENT_PORT = '6832'
      process.env.JAEGER_LOG_SPANS = true

      config.getConfig().then(c => {
        c.serviceName.should.equal('my-service')
        c.servicePort.should.equal(10000)
        c.mongoUri.should.equal('mongodb://mongo.mlab.com:27017/sites')
        c.mongoUsername.should.equal('user')
        c.mongoPassword.should.equal('pass')
        c.jaegerAgentHost.should.equal('jaeger-agent')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(true)
        done()
      }).catch(done)
    })

    it('should return values from files', done => {
      // tmp will clean up after itself!
      // See https://raszi.github.io/node-tmp
      const nameFile = tmp.fileSync()
      const portFile = tmp.fileSync()
      const urlFile = tmp.fileSync()
      const userFile = tmp.fileSync()
      const passFile = tmp.fileSync()
      const agentHostFile = tmp.fileSync()
      const agentPortFile = tmp.fileSync()
      const logSpansFile = tmp.fileSync()

      process.env.SERVICE_NAME_FILE = nameFile.name
      process.env.SERVICE_PORT_FILE = portFile.name
      process.env.MONGO_URI_FILE = urlFile.name
      process.env.MONGO_USERNAME_FILE = userFile.name
      process.env.MONGO_PASSWORD_FILE = passFile.name
      process.env.JAEGER_AGENT_HOST_FILE = agentHostFile.name
      process.env.JAEGER_AGENT_PORT_FILE = agentPortFile.name
      process.env.JAEGER_LOG_SPANS_FILE = logSpansFile.name

      fs.writeFileSync(nameFile.fd, 'new-service')
      fs.writeFileSync(portFile.fd, '20000')
      fs.writeFileSync(urlFile.fd, 'mongodb://user:pass@mongo1:27017,mongo2:27017,mongo3:27017')
      fs.writeFileSync(userFile.fd, 'root')
      fs.writeFileSync(passFile.fd, 'toor')
      fs.writeFileSync(agentHostFile.fd, 'jaeger-agent')
      fs.writeFileSync(agentPortFile.fd, '6832')
      fs.writeFileSync(logSpansFile.fd, 'true')

      config.getConfig().then(c => {
        c.serviceName.should.equal('new-service')
        c.servicePort.should.equal(20000)
        c.mongoUri.should.equal('mongodb://user:pass@mongo1:27017,mongo2:27017,mongo3:27017/sites')
        c.mongoUsername.should.equal('root')
        c.mongoPassword.should.equal('toor')
        c.jaegerAgentHost.should.equal('jaeger-agent')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(true)
        done()
      }).catch(done)
    })
  })
})
