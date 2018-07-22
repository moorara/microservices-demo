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
    delete process.env.MONGO_URL
    delete process.env.MONGO_URL_FILE
    delete process.env.MONGO_USER
    delete process.env.MONGO_USER_FILE
    delete process.env.MONGO_PASS
    delete process.env.MONGO_PASS_FILE
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
        c.mongoUrl.should.equal('mongodb://localhost:27017/sites')
        should.equal(c.mongoUser, null)
        should.equal(c.mongoPass, null)
        should.equal(c.jaegerAgentHost, 'localhost')
        should.equal(c.jaegerAgentPort, '6832')
        should.equal(c.jaegerLogSpans, false)
        done()
      }).catch(done)
    })

    it('should return values from environment variables', done => {
      process.env.SERVICE_NAME = 'my-service'
      process.env.SERVICE_PORT = '10000'
      process.env.MONGO_URL = 'mongodb://mongo.mlab.com:27017'
      process.env.MONGO_USER = 'user'
      process.env.MONGO_PASS = 'pass'
      process.env.JAEGER_AGENT_HOST = 'jaeger-agent'
      process.env.JAEGER_AGENT_PORT = '6832'
      process.env.JAEGER_LOG_SPANS = true

      config.getConfig().then(c => {
        c.serviceName.should.equal('my-service')
        c.servicePort.should.equal(10000)
        c.mongoUrl.should.equal('mongodb://mongo.mlab.com:27017/sites')
        c.mongoUser.should.equal('user')
        c.mongoPass.should.equal('pass')
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
      process.env.MONGO_URL_FILE = urlFile.name
      process.env.MONGO_USER_FILE = userFile.name
      process.env.MONGO_PASS_FILE = passFile.name
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
        c.mongoUrl.should.equal('mongodb://user:pass@mongo1:27017,mongo2:27017,mongo3:27017/sites')
        c.mongoUser.should.equal('root')
        c.mongoPass.should.equal('toor')
        c.jaegerAgentHost.should.equal('jaeger-agent')
        c.jaegerAgentPort.should.equal(6832)
        c.jaegerLogSpans.should.equal(true)
        done()
      }).catch(done)
    })
  })
})
