/* eslint-env mocha */
const sinon = require('sinon')
const should = require('should')

const resolvers = require('./sensor')

describe('sensorResolvers', () => {
  let span, logger
  let sensorService, _sensorService
  let context, info

  beforeEach(() => {
    span = {}
    logger = {
      debug () {},
      verbose () {},
      info () {},
      warn () {},
      error () {},
      fatal () {}
    }

    sensorService = {
      get () {},
      all () {},
      create () {},
      update () {},
      delete () {},
    }
    _sensorService = sinon.mock(sensorService)

    context = { span, logger, sensorService }
    info = {}
  })

  afterEach(() => {
    _sensorService.restore()
  })

  describe('Query', () => {
    describe('sensor', () => {
      let id

      it('should throw an error when service.get fails', done => {
        id = '1111-1111'
        const err = new Error("get error")
        _sensorService.expects('get').withArgs({ span }, id).rejects(err)
        resolvers.Query.sensor(null, { id }, context, info).catch(e => {
          e.should.eql(err)
          _sensorService.verify()
          done()
        })
      })
      it('should return a promise that resolves to a sensor', done => {
        id = '2222-2222'
        const sensor = { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30.0, maxSafe: 30.0 }
        _sensorService.expects('get').withArgs({ span }, id).resolves(sensor)
        resolvers.Query.sensor(null, { id }, context, info).then(s => {
          s.should.eql(sensor)
          _sensorService.verify()
          done()
        }).catch(done)
      })
    })

    describe('sensors', () => {
      let siteId

      it('should throw an error when service.all fails', done => {
        siteId = 'aaaa-aaaa'
        const err = new Error("all error")
        _sensorService.expects('all').withArgs({ span }, siteId).rejects(err)
        resolvers.Query.sensors(null, { siteId }, context, info).catch(e => {
          e.should.eql(err)
          _sensorService.verify()
          done()
        })
      })
      it('should return a promise that resolves to an array sensors', done => {
        siteId = 'bbbb-bbbb'
        const sensors = [{ id: '1111-1111', siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30.0, maxSafe: 30.0 }]
        _sensorService.expects('all').withArgs({ span }, siteId).resolves(sensors)
        resolvers.Query.sensors(null, { siteId }, context, info).then(s => {
          s.should.eql(sensors)
          _sensorService.verify()
          done()
        }).catch(done)
      })
    })
  })

  describe('Mutation', () => {
    describe('createSensor', () => {
      let input

      it('should throw an error when service.create fails', done => {
        input = {}
        const err = new Error("create error")
        _sensorService.expects('create').withArgs({ span }, input).rejects(err)
        resolvers.Mutation.createSensor(null, { input }, context, info).catch(e => {
          e.should.eql(err)
          _sensorService.verify()
          done()
        })
      })
      it('should return a promise that resolves to the new sensor', done => {
        input = { siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30.0, maxSafe: 30.0 }
        const sensor = Object.assign({}, { id: '2222-2222' }, input)
        _sensorService.expects('create').withArgs({ span }, input).resolves(sensor)
        resolvers.Mutation.createSensor(null, { input }, context, info).then(s => {
          s.should.eql(sensor)
          _sensorService.verify()
          done()
        }).catch(done)
      })
    })

    describe('updateSensor', () => {
      let id, input

      it('should throw an error when service.update fails', done => {
        id = '1111-1111'
        input = {}
        const err = new Error("update error")
        _sensorService.expects('update').withArgs({ span }, id, input).rejects(err)
        resolvers.Mutation.updateSensor(null, { id, input }, context, info).catch(e => {
          e.should.eql(err)
          _sensorService.verify()
          done()
        })
      })
      it('should return a promise that resolves to the updated sensor', done => {
        id = '2222-2222'
        input = { siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30.0, maxSafe: 30.0 }
        const sensor = Object.assign({}, { id }, input)
        _sensorService.expects('update').withArgs({ span }, id, input).resolves(sensor)
        resolvers.Mutation.updateSensor(null, { id, input }, context, info).then(s => {
          s.should.eql(sensor)
          _sensorService.verify()
          done()
        }).catch(done)
      })
    })

    describe('deleteSensor', () => {
      let input

      it('should throw an error when service.delete fails', done => {
        id = '1111-1111'
        const err = new Error("delete error")
        _sensorService.expects('delete').withArgs({ span }, id).rejects(err)
        resolvers.Mutation.deleteSensor(null, { id }, context, info).catch(e => {
          e.should.eql(err)
          _sensorService.verify()
          done()
        })
      })
      it('should return a promise that resolves true', done => {
        id = '2222-2222'
        _sensorService.expects('delete').withArgs({ span }, id).resolves(true)
        resolvers.Mutation.deleteSensor(null, { id }, context, info).then(result => {
          result.should.be.true()
          _sensorService.verify()
          done()
        }).catch(done)
      })
    })
  })

  describe('Sensor', () => {
    let sensor

    beforeEach(() => {
      sensor = { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30.0, maxSafe: 30.0 }
    })

    describe('id', () => {
      it('should return sensor id', () => {
        const id = resolvers.Sensor.id(sensor)
        should.equal(id, sensor.id)
      })
    })
    describe('siteId', () => {
      it('should return sensor siteId', () => {
        const siteId = resolvers.Sensor.siteId(sensor)
        should.equal(siteId, sensor.siteId)
      })
    })
    describe('name', () => {
      it('should return sensor name', () => {
        const name = resolvers.Sensor.name(sensor)
        should.equal(name, sensor.name)
      })
    })
    describe('unit', () => {
      it('should return sensor unit', () => {
        const unit = resolvers.Sensor.unit(sensor)
        should.equal(unit, sensor.unit)
      })
    })
    describe('minSafe', () => {
      it('should return sensor minSafe', () => {
        const minSafe = resolvers.Sensor.minSafe(sensor)
        should.equal(minSafe, sensor.minSafe)
      })
    })
    describe('maxSafe', () => {
      it('should return sensor maxSafe', () => {
        const maxSafe = resolvers.Sensor.maxSafe(sensor)
        should.equal(maxSafe, sensor.maxSafe)
      })
    })
  })
})
