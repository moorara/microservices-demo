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
    delete process.env.SERVICE_PORT
    delete process.env.SERVICE_PORT_FILE
    delete process.env.MONGO_URL
    delete process.env.MONGO_URL_FILE
    delete process.env.MONGO_USER
    delete process.env.MONGO_USER_FILE
    delete process.env.MONGO_PASS
    delete process.env.MONGO_PASS_FILE
  })

  describe('getConfig', () => {
    it('should return default values', done => {
      config.getConfig().then(c => {
        c.servicePort.should.equal('4020')
        c.mongoUrl.should.equal('mongodb://localhost:27017/toys')
        should.equal(c.mongoUser, null)
        should.equal(c.mongoPass, null)
        done()
      }).catch(done)
    })

    it('should return values from environment variables', done => {
      process.env.SERVICE_PORT = '7000'
      process.env.MONGO_URL = 'mongodb://mongo.mlab.com:27017'
      process.env.MONGO_USER = 'user'
      process.env.MONGO_PASS = 'pass'

      config.getConfig().then(c => {
        c.servicePort.should.equal('7000')
        c.mongoUrl.should.equal('mongodb://mongo.mlab.com:27017/toys')
        c.mongoUser.should.equal('user')
        c.mongoPass.should.equal('pass')
        done()
      }).catch(done)
    })

    it('should return values from files', done => {
      // tmp will clean up after itself!
      // See https://raszi.github.io/node-tmp/
      let portFile = tmp.fileSync()
      let urlFile = tmp.fileSync()
      let userFile = tmp.fileSync()
      let passFile = tmp.fileSync()

      process.env.SERVICE_PORT_FILE = portFile.name
      process.env.MONGO_URL_FILE = urlFile.name
      process.env.MONGO_USER_FILE = userFile.name
      process.env.MONGO_PASS_FILE = passFile.name

      fs.writeFileSync(portFile.fd, '9000')
      fs.writeFileSync(urlFile.fd, 'mongodb://user:pass@mongo1:27017,mongo2:27017,mongo3:27017')
      fs.writeFileSync(userFile.fd, 'root')
      fs.writeFileSync(passFile.fd, 'toor')

      config.getConfig().then(c => {
        c.servicePort.should.equal('9000')
        c.mongoUrl.should.equal('mongodb://user:pass@mongo1:27017,mongo2:27017,mongo3:27017/toys')
        c.mongoUser.should.equal('root')
        c.mongoPass.should.equal('toor')
        done()
      }).catch(done)
    })
  })
})
