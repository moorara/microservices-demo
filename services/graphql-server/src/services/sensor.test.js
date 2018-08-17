/* eslint-env mocha */
const should = require('should')
const opentracing = require('opentracing')

const SensorService = require('./sensor')

describe('SensorService', () => {
  let config, options

  beforeEach(() => {
    config = {}
    options = {
      logger: {
        debug () {},
        verbose () {},
        info () {},
        warn () {},
        error () {},
        fatal () {}
      },
      histogram: { observe () {} },
      summary: { observe () {} },
      tracer: new opentracing.MockTracer()
    }
  })

  describe('constructor', () => {
    it('should create a new service with defaults', () => {
      const service = new SensorService(config, { tracer: options.tracer })
      should.exist(service.logger)
      should.exist(service.histogram)
      should.exist(service.summary)
      should.exist(service.tracer)
    })
    it('should create a new service with provided options', () => {
      const service = new SensorService(config, options)
      service.logger.should.equal(options.logger)
      service.histogram.should.equal(options.histogram)
      service.summary.should.equal(options.summary)
      service.tracer.should.equal(options.tracer)
    })
  })

  describe('create', () => {
    let service, context

    beforeEach(() => {
      service = new SensorService(config, options)
      context = {}
    })

    it('should create and persist a new sensor', done => {
      const input = { siteId: 'aaaa-aaaa', name: 'pressure', unit: 'atmosphere', minSafe: 0.5, maxSafe: 1.0 }
      service.create(context, input).then(sensor => {
        should.exist(sensor.id)
        sensor.siteId.should.equal(input.siteId)
        sensor.name.should.equal(input.name)
        sensor.unit.should.equal(input.unit)
        sensor.minSafe.should.equal(input.minSafe)
        sensor.maxSafe.should.equal(input.maxSafe)
        done()
      }).catch(done)
    })
    it('should create and persist a new sensor', done => {
      const input = { siteId: 'bbbb-bbbb', name: 'pressure', unit: 'pascal', minSafe: 50000.0, maxSafe: 100000.0 }
      service.create(context, input).then(sensor => {
        should.exist(sensor.id)
        sensor.siteId.should.equal(input.siteId)
        sensor.name.should.equal(input.name)
        sensor.unit.should.equal(input.unit)
        sensor.minSafe.should.equal(input.minSafe)
        sensor.maxSafe.should.equal(input.maxSafe)
        done()
      }).catch(done)
    })
  })

  describe('all', () => {
    let service, context

    beforeEach(() => {
      service = new SensorService(config, options)
      context = {}
    })

    it('should return all sensors of a site', done => {
      const siteId = 'aaaa-aaaa'
      service.all(context, siteId).then(sensors => {
        sensors.should.have.length(1)
        done()
      }).catch(done)
    })
    it('should return all sensors of a site', done => {
      const siteId = 'bbbb-bbbb'
      service.all(context, siteId).then(sensors => {
        sensors.should.have.length(1)
        done()
      }).catch(done)
    })
  })

  describe('get', () => {
    let service, context

    beforeEach(() => {
      service = new SensorService(config, options)
      context = {}
    })

    it('should return a sensor by id', done => {
      const id = '1111-1111'
      service.get(context, id).then(sensor => {
        sensor.id.should.equal(id)
        sensor.siteId.should.equal('aaaa-aaaa')
        sensor.name.should.equal('temperature')
        sensor.unit.should.equal('celsius')
        sensor.minSafe.should.equal(-30.0)
        sensor.maxSafe.should.equal(30.0)
        done()
      }).catch(done)
    })
    it('should return a sensor by id', done => {
      const id = '2222-2222'
      service.get(context, id).then(sensor => {
        sensor.id.should.equal(id)
        sensor.siteId.should.equal('bbbb-bbbb')
        sensor.name.should.equal('temperature')
        sensor.unit.should.equal('fahrenheit')
        sensor.minSafe.should.equal(-22.0)
        sensor.maxSafe.should.equal(86.0)
        done()
      }).catch(done)
    })
  })

  describe('update', () => {
    let service, context

    beforeEach(() => {
      service = new SensorService(config, options)
      context = {}
    })

    it('should update a sensor by id', done => {
      const id = '1111-1111'
      const input = { siteId: 'aaaa-aaaa', name: 'pressure', unit: 'atmosphere', minSafe: 0.6, maxSafe: 0.9 }
      service.update(context, id, input).then(sensor => {
        sensor.id.should.equal(id)
        sensor.siteId.should.equal(input.siteId)
        sensor.name.should.equal(input.name)
        sensor.unit.should.equal(input.unit)
        sensor.minSafe.should.equal(input.minSafe)
        sensor.maxSafe.should.equal(input.maxSafe)
        done()
      }).catch(done)
    })
    it('should update a sensor by id', done => {
      const id = '2222-2222'
      const input = { siteId: 'bbbb-bbbb', name: 'pressure', unit: 'pascal', minSafe: 60000.0, maxSafe: 90000.0 }
      service.update(context, id, input).then(sensor => {
        sensor.id.should.equal(id)
        sensor.siteId.should.equal(input.siteId)
        sensor.name.should.equal(input.name)
        sensor.unit.should.equal(input.unit)
        sensor.minSafe.should.equal(input.minSafe)
        sensor.maxSafe.should.equal(input.maxSafe)
        done()
      }).catch(done)
    })
  })

  describe('delete', () => {
    let service, context

    beforeEach(() => {
      service = new SensorService(config, options)
      context = {}
    })

    it('should delete a sensor by id', done => {
      const id = '1111-1111'
      service.delete(context, id).then(() => {
        done()
      }).catch(done)
    })
    it('should delete a sensor by id', done => {
      const id = '2222-2222'
      service.delete(context, id).then(() => {
        done()
      }).catch(done)
    })
  })
})
