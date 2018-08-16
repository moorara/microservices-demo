const _ = require('lodash')

const Logger = require('../utils/logger')

class SwitchService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SwitchService')

    this.store = {
      switches: [
        { id: '3333-3333', siteId: 'aaaa-aaaa', name: 'Light', state: 'OFF', states: ['OFF', 'ON'] },
        { id: '4444-4444', siteId: 'bbbb-bbbb', name: 'Light', state: 'OFF', states: ['OFF', 'ON'] }
      ]
    }
  }

  create (context, input) {
    let swtch = Object.assign({}, input)
    swtch.id = _.uniqueId()
    this.store.switches.push(swtch)
    return Promise.resolve(swtch)
  }

  all (context, siteId) {
    const switches = this.store.switches.filter(s => s.siteId === siteId)
    return Promise.resolve(switches)
  }

  get (context, id) {
    const swtch = Object.assign({}, this.store.switches.find(s => s.id === id))
    return Promise.resolve(swtch)
  }

  update (context, id, { state }) {
    const swtch = this.store.switches.find(s => s.id === id)
    swtch.state = state
    const updated = Object.assign({}, swtch)
    return Promise.resolve(updated)
  }

  delete (context, id) {
    _.remove(this.store.switches, s => s.id === id)
    return Promise.resolve()
  }
}

module.exports = SwitchService
