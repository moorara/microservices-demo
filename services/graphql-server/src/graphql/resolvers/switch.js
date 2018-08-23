const resolvers = {
  Query: {
    async switch (_, { id }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.switchService.getSwitch(ctx, id)
      } catch (err) {
        context.logger.error(err)
        throw err
      }
    },

    async switches (_, { siteId }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.switchService.getSwitches(ctx, siteId)
      } catch (err) {
        context.logger.error(err)
        throw err
      }
    }
  },

  Mutation: {
    async installSwitch (_, { input }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.switchService.installSwitch(ctx, input)
      } catch (err) {
        context.logger.error(err)
        throw err
      }
    },

    async setSwitch (_, { id, state }, context, info) {
      try {
        const ctx = { span: context.span }
        return await context.switchService.setSwitch(ctx, id, { state })
      } catch (err) {
        context.logger.error(err)
        throw err
      }
    },

    async removeSwitch (_, { id }, context, info) {
      try {
        const ctx = { span: context.span }
        await context.switchService.removeSwitch(ctx, id)
        return true
      } catch (err) {
        context.logger.error(err)
        throw err
      }
    }
  },

  Switch: {
    id: swtch => swtch.id,
    siteId: swtch => swtch.siteId,
    name: swtch => swtch.name,
    state: swtch => swtch.state,
    states: swtch => swtch.states
  }
}

module.exports = resolvers
