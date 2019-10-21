/* eslint-env mocha */
const should = require('should')

const { proto } = require('.')

describe('grpc', () => {
  describe('proto', () => {
    it('should include the SwitchService', () => {
      should.exist(proto.SwitchService)
      should.exist(proto.SwitchService.service.InstallSwitch)
      should.exist(proto.SwitchService.service.RemoveSwitch)
      should.exist(proto.SwitchService.service.GetSwitch)
      should.exist(proto.SwitchService.service.GetSwitches)
      should.exist(proto.SwitchService.service.SetSwitch)
    })
  })
})
