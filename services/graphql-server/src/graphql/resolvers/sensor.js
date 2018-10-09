const resolvers = {
  Query: {
    async sensor (_, { id }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.sensorService.get(ctx, id)
      } catch (err) {
        context.logger.error('Error on sensor:', err)
        throw err
      }
    },

    async sensors (_, { siteId }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.sensorService.all(ctx, siteId)
      } catch (err) {
        context.logger.error('Error on sensors:', err)
        throw err
      }
    }
  },

  Mutation: {
    async createSensor (_, { input }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.sensorService.create(ctx, input)
      } catch (err) {
        context.logger.error('Error on createSensor:', err)
        throw err
      }
    },

    async updateSensor (_, { id, input }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.sensorService.update(ctx, id, input)
      } catch (err) {
        context.logger.error('Error on updateSensor:', err)
        throw err
      }
    },

    async deleteSensor (_, { id }, context, info) {
      try {
        const ctx = { span: context.span }
        await context.sensorService.delete(ctx, id)
        return true
      } catch (err) {
        context.logger.error('Error on deleteSensor:', err)
        throw err
      }
    }
  },

  Sensor: {
    id: sensor => sensor.id,
    siteId: sensor => sensor.siteId,
    name: sensor => sensor.name,
    unit: sensor => sensor.unit,
    minSafe: sensor => sensor.minSafe,
    maxSafe: sensor => sensor.maxSafe
  }
}

module.exports = resolvers
