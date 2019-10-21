const resolvers = {
  Query: {
    async site (_, { id }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.siteService.get(ctx, id)
      } catch (err) {
        context.logger.error('Error on site:', err)
        throw err
      }
    },

    async sites (_, args, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.siteService.all(ctx, args)
      } catch (err) {
        context.logger.error('Error on sites:', err)
        throw err
      }
    }
  },

  Mutation: {
    async createSite (_, { input }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.siteService.create(ctx, input)
      } catch (err) {
        context.logger.error('Error on createSite:', err)
        throw err
      }
    },

    async updateSite (_, { id, input }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.siteService.update(ctx, id, input)
      } catch (err) {
        context.logger.error('Error on updateSite:', err)
        throw err
      }
    },

    async deleteSite (_, { id }, context, info) {
      try {
        const ctx = { span: context.span }
        await context.siteService.delete(ctx, id)
        return true
      } catch (err) {
        context.logger.error('Error on deleteSite:', err)
        throw err
      }
    }
  },

  Site: {
    id: site => site.id,
    name: site => site.name,
    location: site => site.location,
    priority: site => site.priority,
    tags: site => site.tags,

    sensors: (site, args, context, info) => {
      const ctx = { span: context.span }
      return context.sensorService.all(ctx, site.id)
    },

    switches: (site, args, context, info) => {
      const ctx = { span: context.span }
      return context.switchService.getSwitches(ctx, site.id)
    },

    assets: async (site, args, context, info) => {
      const ctx = { span: context.span }
      const allAlarm = context.assetService.allAlarm(ctx, site.id)
      const allCamera = context.assetService.allCamera(ctx, site.id)
      const [alarms, cameras] = await Promise.all([allAlarm, allCamera])
      return [].concat(alarms, cameras)
    }
  }
}

module.exports = resolvers
