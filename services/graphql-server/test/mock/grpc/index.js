const _ = require('lodash')
const grpc = require('grpc')
const chalk = require('chalk')

const store = require('./data')
const { proto } = require('../../../src/grpc')

class SwitchServiceServer {
  constructor (options) {
    options = options || {}
    this.port = options.port || process.env.PORT || 4500
    this.server = options.server || new grpc.Server()
    this.store = options.store || store
  }

  logCall (method, metadata, req, res) {
    console.log(`${chalk.green(new Date().toISOString())}`)
    console.log(`  ${chalk.red.bold(method)}`)

    const map = metadata.getMap()
    console.log(`    ${chalk.cyan.bold('metadata')}`)
    for (let key in map) {
      console.log(`      ${key}: ${JSON.stringify(map[key])}`)
    }

    console.log(`    ${chalk.yellow.bold('request')}`)
    for (let field in req) {
      console.log(`      ${field}: ${JSON.stringify(req[field])}`)
    }

    console.log(`    ${chalk.blue.bold('response')}`)
    for (let field in res) {
      console.log(`      ${field}: ${JSON.stringify(res[field])}`)
    }
  }

  start () {
    this.server.addService(proto.SwitchService.service, this)
    const serviceAddr = `0.0.0.0:${this.port}`
    const credentials = grpc.ServerCredentials.createInsecure()
    this.server.bind(serviceAddr, credentials)

    console.log(chalk.green(`Mock gRPC API Server Listening on ${this.port} ...`))
    this.server.start()
  }

  installSwitch (call, callback) {
    const metadata = call.metadata
    const { siteId, name, state, states } = call.request

    const id = _.uniqueId()
    const swtch = { id, siteId, name, state, states }
    this.store.switches.push(swtch)

    this.logCall('installSwitch', metadata, call.request, swtch)
    callback(null, swtch)
  }

  removeSwitch (call, callback) {
    const metadata = call.metadata
    const { id } = call.request

    _.remove(this.store.switches, s => s.id === id)

    this.logCall('removeSwitch', metadata, call.request, {})
    callback(null, {})
  }

  getSwitch (call, callback) {
    const metadata = call.metadata
    const { id } = call.request

    const swtch = this.store.switches.find(s => s.id === id)

    this.logCall('getSwitch', metadata, call.request, swtch)
    callback(null, swtch)
  }

  getSwitches (call) {
    const metadata = call.metadata
    const { siteId } = call.request

    const switches = this.store.switches.filter(s => s.siteId === siteId)
    switches.forEach(swtch => {
      call.write(swtch)
    })

    this.logCall('getSwitches', metadata, call.request, switches)
    call.end()
  }

  setSwitch (call, callback) {
    const metadata = call.metadata
    const { id, state } = call.request

    const swtch = this.store.switches.find(s => s.id === id)
    swtch.state = state

    this.logCall('setSwitch', metadata, call.request, swtch)
    callback(null, swtch)
  }
}

const server = new SwitchServiceServer()
server.start()
