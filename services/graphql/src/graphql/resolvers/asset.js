const resolvers = {
  Query: {
    async asset (_, { id }, context, info) {
      try {
        const ctx = { span: context.span }
        const getAlarm = context.assetService.getAlarm(ctx, id)
        const getCamera = context.assetService.getCamera(ctx, id)
        const [alarm, camera] = await Promise.all([getAlarm, getCamera])
        return alarm || camera
      } catch (err) {
        context.logger.error('Error on asset:', err)
        throw err
      }
    },

    async assets (_, { siteId }, context, info) {
      try {
        const ctx = { span: context.span }
        const allAlarm = context.assetService.allAlarm(ctx, siteId)
        const allCamera = context.assetService.allCamera(ctx, siteId)
        const [alarms, cameras] = await Promise.all([allAlarm, allCamera])
        return [].concat(alarms, cameras)
      } catch (err) {
        context.logger.error('Error on assets:', err)
        throw err
      }
    }
  },

  Mutation: {
    async createAlarm (_, { input }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.assetService.createAlarm(ctx, input)
      } catch (err) {
        context.logger.error('Error on createAlarm:', err)
        throw err
      }
    },

    async updateAlarm (_, { id, input }, context, info) {
      try {
        const ctx = { span: context.span }
        const updated = await context.assetService.updateAlarm(ctx, id, input)
        return updated ? await context.assetService.getAlarm(ctx, id) : null
      } catch (err) {
        context.logger.error('Error on updateAlarm:', err)
        throw err
      }
    },

    async createCamera (_, { input }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.assetService.createCamera(ctx, input)
      } catch (err) {
        context.logger.error('Error on createCamera:', err)
        throw err
      }
    },

    async updateCamera (_, { id, input }, context, info) {
      try {
        const ctx = { span: context.span }
        const updated = await context.assetService.updateCamera(ctx, id, input)
        return updated ? await context.assetService.getCamera(ctx, id) : null
      } catch (err) {
        context.logger.error('Error on updateCamera:', err)
        throw err
      }
    },

    async deleteAsset (_, { id }, context, info) {
      try {
        const ctx = { span: context.span }
        const deleteAlarm = context.assetService.deleteAlarm(ctx, id)
        const deleteCamera = context.assetService.deleteCamera(ctx, id)
        const [alarmDeleted, cameraDeleted] = await Promise.all([deleteAlarm, deleteCamera])
        return alarmDeleted || cameraDeleted
      } catch (err) {
        context.logger.error('Error on deleteAsset:', err)
        throw err
      }
    }
  },

  Asset: {
    __resolveType (obj, context, info) {
      if (obj.material) {
        return 'Alarm'
      }
      if (obj.resolution) {
        return 'Camera'
      }
      return null
    }
  },

  Alarm: {
    id: alarm => alarm.id,
    siteId: alarm => alarm.siteId,
    serialNo: alarm => alarm.serialNo,
    material: alarm => alarm.material
  },

  Camera: {
    id: camera => camera.id,
    siteId: camera => camera.siteId,
    serialNo: camera => camera.serialNo,
    resolution: camera => camera.resolution
  }
}

module.exports = resolvers
