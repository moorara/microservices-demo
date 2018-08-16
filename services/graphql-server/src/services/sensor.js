const _ = require('lodash')

const Logger = require('../utils/logger')

class SensorService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SensorService')

    this.store = {
      sensors: [
        { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30.0, maxSafe: 30.0 },
        { id: '2222-2222', siteId: 'bbbb-bbbb', name: 'temperature', unit: 'fahrenheit', minSafe: -22.0, maxSafe: 86.0 }
      ]
    }
  }

  create (context, input) {
    let sensor = Object.assign({}, input)
    sensor.id = _.uniqueId()
    this.store.sensors.push(sensor)
    return Promise.resolve(sensor)
  }

  all (context, siteId) {
    const sensors = this.store.sensors.filter(s => s.siteId === siteId)
    return Promise.resolve(sensors)
  }

  get (context, id) {
    const sensor = Object.assign({}, this.store.sensors.find(s => s.id === id))
    return Promise.resolve(sensor)
  }

  update (context, id, input) {
    const sensor = Object.assign({}, { id }, input)
    for (let i in this.store.sensors) {
      if (this.store.sensors[i].id === id) {
        this.store.sensors[i] = sensor
        const updated = Object.assign({}, sensor)
        return Promise.resolve(updated)
      }
    }
  }

  delete (context, id) {
    _.remove(this.store.sensors, s => s.id === id)
    return Promise.resolve()
  }
}

module.exports = SensorService
