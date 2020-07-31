const nats = require('nats')
const chalk = require('chalk')
const { v4: uuidv4 } = require('uuid')

const Config = require('./config')

const subject = 'asset_service'
const queue = 'workers'

class MockAssetService {
  constructor (config) {
    this.store = {
      alarms: [],
      cameras: []
    }
    this.conn = nats.connect({
      servers: config.natsServers,
      user: config.natsUser,
      pass: config.natsPassword
    })
  }

  createAlarm (request, reply) {
    const id = uuidv4()
    const alarm = Object.assign({ id }, request.input)
    const response = { kind: request.kind, alarm }
    const msg = JSON.stringify(response)
    this.conn.publish(reply, msg)
    this.store.alarms.push(alarm)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  allAlarm (request, reply) {
    const alarms = this.store.alarms.filter(a => a.siteId === request.siteId)
    const response = { kind: request.kind, alarms }
    const msg = JSON.stringify(response)
    this.conn.publish(reply, msg)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  getAlarm (request, reply) {
    const alarm = this.store.alarms.find(a => a.id === request.id)
    const response = { kind: request.kind, alarm }
    const msg = JSON.stringify(response)
    this.conn.publish(reply, msg)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  updateAlarm (request, reply) {
    const alarm = this.store.alarms.find(a => a.id === request.id)
    Object.assign(alarm, request.input)
    const response = { kind: request.kind, updated: true }
    const msg = JSON.stringify(response)
    this.conn.publish(reply, msg)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  deleteAlarm (request, reply) {
    this.store.alarms = this.store.alarms.filter(a => a.id !== request.id)
    const response = { kind: request.kind, deleted: true }
    const msg = JSON.stringify(response)
    this.conn.publish(reply, msg)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  createCamera (request, reply) {
    const id = uuidv4()
    const camera = Object.assign({ id }, request.input)
    const response = { kind: request.kind, camera }
    const msg = JSON.stringify(response)

    this.store.cameras.push(camera)
    this.conn.publish(reply, msg)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  allCamera (request, reply) {
    const cameras = this.store.cameras.filter(c => c.siteId === request.siteId)
    const response = { kind: request.kind, cameras }
    const msg = JSON.stringify(response)
    this.conn.publish(reply, msg)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  getCamera (request, reply) {
    const camera = this.store.cameras.find(c => c.id === request.id)
    const response = { kind: request.kind, camera }
    const msg = JSON.stringify(response)
    this.conn.publish(reply, msg)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  updateCamera (request, reply) {
    const camera = this.store.cameras.find(c => c.id === request.id)
    Object.assign(camera, request.input)
    const response = { kind: request.kind, updated: true }
    const msg = JSON.stringify(response)
    this.conn.publish(reply, msg)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  deleteCamera (request, reply) {
    this.store.cameras = this.store.cameras.filter(c => c.id !== request.id)
    const response = { kind: request.kind, deleted: true }
    const msg = JSON.stringify(response)
    this.conn.publish(reply, msg)

    console.log(chalk.green(`Sent message ${msg}`))
  }

  start () {
    console.log(chalk.green('Mock NATS API Service started ...'))

    this.sid = this.conn.subscribe(subject, { queue }, (msg, reply) => {
      console.log(chalk.blue(`Received message ${msg}`))

      let request

      try {
        request = JSON.parse(msg)
      } catch (err) {
        console.log(chalk.red(`Invalid request ${msg}`))
      }

      switch (request.kind) {
        case 'createAlarm':
          this.createAlarm(request, reply)
          break
        case 'allAlarm':
          this.allAlarm(request, reply)
          break
        case 'getAlarm':
          this.getAlarm(request, reply)
          break
        case 'updateAlarm':
          this.updateAlarm(request, reply)
          break
        case 'deleteAlarm':
          this.deleteAlarm(request, reply)
          break
        case 'createCamera':
          this.createCamera(request, reply)
          break
        case 'allCamera':
          this.allCamera(request, reply)
          break
        case 'getCamera':
          this.getCamera(request, reply)
          break
        case 'updateCamera':
          this.updateCamera(request, reply)
          break
        case 'deleteCamera':
          this.deleteCamera(request, reply)
          break
        default:
          console.log(chalk.yellow(`Unknown request ${request.kind}`))
      }
    })
  }

  stop () {
    this.conn.unsubscribe(this.sid)
    this.conn.close()
  }
}

const config = new Config()
const assetService = new MockAssetService(config)
assetService.start()
