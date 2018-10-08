/* eslint-env mocha */
require('should')
const sinon = require('sinon')

const resolvers = require('./asset')

describe('assetResolvers', () => {
  let span, logger
  let assetService, _assetService
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

    assetService = {
      getAsset () {},
      getAssets () {},
      createAlarm () {},
      updateAlarm () {},
      createCamera () {},
      updateCamera () {},
      deleteAsset () {}
    }
    _assetService = sinon.mock(assetService)

    context = { span, logger, assetService }
    info = {}
  })

  afterEach(() => {
    _assetService.restore()
  })

  describe('Query', () => {
    let id

    describe('asset', () => {
      it('should reject with an error when service request fails', done => {
        id = 'aaaa-aaaa'
        const err = new Error('get error')
        _assetService.expects('getAsset').withArgs({ span: context.span }, id).rejects(err)
        resolvers.Query.asset(null, { id }, context, info).catch(e => {
          e.should.eql(err)
          _assetService.verify()
          done()
        })
      })
      it('should resolve with an alarm when asset is an alarm', done => {
        id = 'aaaa-aaaa'
        const asset = { id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
        _assetService.expects('getAsset').withArgs({ span: context.span }, id).resolves(asset)
        resolvers.Query.asset(null, { id }, context, info).then(result => {
          result.should.eql(asset)
          _assetService.verify()
          done()
        }).catch(done)
      })
      it('should resolve with a camera when asset is a camera', done => {
        id = 'bbbb-bbbb'
        const asset = { id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
        _assetService.expects('getAsset').withArgs({ span: context.span }, id).resolves(asset)
        resolvers.Query.asset(null, { id }, context, info).then(result => {
          result.should.eql(asset)
          _assetService.verify()
          done()
        }).catch(done)
      })
    })

    describe('assets', () => {
      let siteId

      it('should reject with an error when service request fails', done => {
        siteId = '1111-1111'
        const err = new Error('all error')
        _assetService.expects('getAssets').withArgs({ span: context.span }, siteId).rejects(err)
        resolvers.Query.assets(null, { siteId }, context, info).catch(e => {
          e.should.eql(err)
          _assetService.verify()
          done()
        })
      })
      it('should resolve with an array of alarms and cameras', done => {
        siteId = '1111-1111'
        const assets = [
          { id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'smoke' },
          { id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
        ]
        _assetService.expects('getAssets').withArgs({ span: context.span }, siteId).resolves(assets)
        resolvers.Query.assets(null, { siteId }, context, info).then(result => {
          result.should.eql(assets)
          _assetService.verify()
          done()
        }).catch(done)
      })
    })
  })

  describe('Mutation', () => {
    describe('createAlarm', () => {
      let input

      it('should reject with an error when service request fails', done => {
        input = { siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
        const err = new Error('create alarm error')
        _assetService.expects('createAlarm').withArgs({ span: context.span }, input).rejects(err)
        resolvers.Mutation.createAlarm(null, { input }, context, info).catch(e => {
          e.should.eql(err)
          _assetService.verify()
          done()
        })
      })
      it('should resolve with an alarm when service request succeeds', done => {
        input = { siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
        const alarm = { id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
        _assetService.expects('createAlarm').withArgs({ span: context.span }, input).resolves(alarm)
        resolvers.Mutation.createAlarm(null, { input }, context, info).then(result => {
          result.should.eql(alarm)
          _assetService.verify()
          done()
        }).catch(done)
      })
    })

    describe('updateAlarm', () => {
      let id, input

      it('should reject with an error when service request fails', done => {
        id = 'aaaa-aaaa'
        input = { siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
        const err = new Error('update alarm error')
        _assetService.expects('updateAlarm').withArgs({ span: context.span }, id, input).rejects(err)
        resolvers.Mutation.updateAlarm(null, { id, input }, context, info).catch(e => {
          e.should.eql(err)
          _assetService.verify()
          done()
        })
      })
      it('should resolve with an alarm when service request succeeds', done => {
        id = 'aaaa-aaaa'
        input = { siteId: '1111-1111', serialNo: '1001', material: 'co' }
        const alarm = { id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'co' }
        _assetService.expects('updateAlarm').withArgs({ span: context.span }, id, input).resolves(alarm)
        resolvers.Mutation.updateAlarm(null, { id, input }, context, info).then(result => {
          result.should.eql(alarm)
          _assetService.verify()
          done()
        }).catch(done)
      })
    })

    describe('createCamera', () => {
      let input

      it('should reject with an error when service request fails', done => {
        input = { siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
        const err = new Error('create camera error')
        _assetService.expects('createCamera').withArgs({ span: context.span }, input).rejects(err)
        resolvers.Mutation.createCamera(null, { input }, context, info).catch(e => {
          e.should.eql(err)
          _assetService.verify()
          done()
        })
      })
      it('should resolve with a camera when service request succeeds', done => {
        input = { siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
        const camera = { id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
        _assetService.expects('createCamera').withArgs({ span: context.span }, input).resolves(camera)
        resolvers.Mutation.createCamera(null, { input }, context, info).then(result => {
          result.should.eql(camera)
          _assetService.verify()
          done()
        }).catch(done)
      })
    })

    describe('updateCamera', () => {
      let id, input

      it('should reject with an error when service request fails', done => {
        id = 'bbbb-bbbb'
        input = { siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
        const err = new Error('update camera error')
        _assetService.expects('updateCamera').withArgs({ span: context.span }, id, input).rejects(err)
        resolvers.Mutation.updateCamera(null, { id, input }, context, info).catch(e => {
          e.should.eql(err)
          _assetService.verify()
          done()
        })
      })
      it('should resolve with a camera when service request succeeds', done => {
        id = 'bbbb-bbbb'
        input = { siteId: '1111-1111', serialNo: '2001', resolution: 1920000 }
        const camera = { id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 1920000 }
        _assetService.expects('updateCamera').withArgs({ span: context.span }, id, input).resolves(camera)
        resolvers.Mutation.updateCamera(null, { id, input }, context, info).then(result => {
          result.should.eql(camera)
          _assetService.verify()
          done()
        }).catch(done)
      })
    })

    describe('deleteAsset', () => {
      let id

      it('should reject with an error when service request fails', done => {
        id = 'aaaa-aaaa'
        const err = new Error('delete asset error')
        _assetService.expects('deleteAsset').withArgs({ span: context.span }, id).rejects(err)
        resolvers.Mutation.deleteAsset(null, { id }, context, info).catch(e => {
          e.should.eql(err)
          _assetService.verify()
          done()
        })
      })
      it('should resolve with true when service request succeeds', done => {
        id = 'bbbb-bbbb'
        _assetService.expects('deleteAsset').withArgs({ span: context.span }, id).resolves(true)
        resolvers.Mutation.deleteAsset(null, { id }, context, info).then(result => {
          result.should.be.true()
          _assetService.verify()
          done()
        }).catch(done)
      })
    })
  })

  describe('Asset', () => {
    describe('resolveType', () => {
      let obj, context, info

      beforeEach(() => {
        context = null
        info = {}
      })

      it('should return Alarm type when asset is an alarm object', () => {
        obj = { id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
        resolvers.Asset.__resolveType(obj, context, info).should.equal('Alarm')
      })
      it('should return Camera type when asset is a camera object', () => {
        obj = { id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
        resolvers.Asset.__resolveType(obj, context, info).should.equal('Camera')
      })
    })
  })

  describe('Alarm', () => {
    let alarm

    beforeEach(() => {
      alarm = { id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
    })

    it('should return id', () => {
      resolvers.Alarm.id(alarm).should.equal(alarm.id)
    })
    it('should return siteId', () => {
      resolvers.Alarm.siteId(alarm).should.equal(alarm.siteId)
    })
    it('should return serialNo', () => {
      resolvers.Alarm.serialNo(alarm).should.equal(alarm.serialNo)
    })
    it('should return material', () => {
      resolvers.Alarm.material(alarm).should.equal(alarm.material)
    })
  })

  describe('Camera', () => {
    let camera

    beforeEach(() => {
      camera = { id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
    })

    it('should return id', () => {
      resolvers.Camera.id(camera).should.equal(camera.id)
    })
    it('should return siteId', () => {
      resolvers.Camera.siteId(camera).should.equal(camera.siteId)
    })
    it('should return serialNo', () => {
      resolvers.Camera.serialNo(camera).should.equal(camera.serialNo)
    })
    it('should return resolution', () => {
      resolvers.Camera.resolution(camera).should.equal(camera.resolution)
    })
  })
})
