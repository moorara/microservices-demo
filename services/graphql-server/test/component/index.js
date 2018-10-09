/* eslint-env mocha */
const http = require('http')
const axios = require('axios')
const should = require('should')

const serviceUrl = process.env.SERVICE_URL || 'http://localhost:5000'

describe('GraphQL Server', function () {
  this.slow(1000)

  describe('General', () => {
    const client = axios.create({
      httpAgent: new http.Agent({ keepAlive: true }),
      baseURL: serviceUrl
    })

    it('GET /liveness', done => {
      client.get('/liveness').then(res => {
        res.status.should.equal(200)
        done()
      }).catch(done)
    })
    it('GET /readiness', done => {
      client.get('/readiness').then(res => {
        res.status.should.equal(200)
        done()
      }).catch(done)
    })
    it('GET /metrics', done => {
      client.get('/metrics').then(res => {
        res.status.should.equal(200)
        res.data.should.match(/# TYPE nodejs_version_info gauge/)
        res.data.should.match(/# TYPE process_cpu_seconds_total counter/)
        res.data.should.match(/# TYPE process_cpu_system_seconds_total counter/)
        res.data.should.match(/# TYPE process_cpu_user_seconds_total counter/)
        res.data.should.match(/# TYPE process_resident_memory_bytes gauge/)
        res.data.should.match(/# TYPE process_virtual_memory_bytes gauge/)
        res.data.should.match(/# TYPE process_heap_bytes gauge/)
        res.data.should.match(/# TYPE process_open_fds gauge/)
        res.data.should.match(/# TYPE process_max_fds gauge/)
        res.data.should.match(/# TYPE nodejs_eventloop_lag_seconds gauge/)
        res.data.should.match(/# TYPE nodejs_active_handles_total gauge/)
        res.data.should.match(/# TYPE nodejs_active_requests_total gauge/)
        res.data.should.match(/# TYPE nodejs_heap_space_size_total_bytes gauge/)
        res.data.should.match(/# TYPE nodejs_heap_space_size_used_bytes gauge/)
        res.data.should.match(/# TYPE nodejs_heap_space_size_available_bytes gauge/)
        res.data.should.match(/# TYPE graphql_server_jaeger_traces counter/)
        res.data.should.match(/# TYPE graphql_server_jaeger_started_spans counter/)
        res.data.should.match(/# TYPE graphql_server_jaeger_finished_spans counter/)
        res.data.should.match(/# TYPE http_requests_duration_seconds histogram/)
        res.data.should.match(/# TYPE http_requests_duration_quantiles_seconds summary/)
        res.data.should.match(/# TYPE graphql_operations_latency_seconds histogram/)
        res.data.should.match(/# TYPE graphql_operations_latency_quantiles_seconds summary/)
        done()
      }).catch(done)
    })
  })

  describe('GraphQL', () => {
    const store = {
      sites: [],
      sensors: [],
      switches: [],
      alarms: [],
      cameras: []
    }

    const client = axios.create({
      httpAgent: new http.Agent({ keepAlive: true }),
      method: 'post',
      url: `${serviceUrl}/graphql`,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Accept-Encoding': 'gzip'
      }
    })

    const query = query => {
      return client.request({
        data: { query }
      })
    }

    // CREATE

    it('createSite', async () => {
      const inputs = [
        { name: 'Plant A', location: 'Ottawa, ON', priority: 0, tags: ['energy'] },
        { name: 'Plant B', location: 'Toronto, ON', priority: 0, tags: ['energy'] },
        { name: 'Plant C', location: 'Montreal, QC', priority: 0, tags: ['energy'] },
        { name: 'Plant D', location: 'Vancouver, BC', priority: 0, tags: ['energy'] }
      ]

      for (let { name, location, priority, tags } of inputs) {
        const res = await query(`
          mutation {
            createSite (input: { name: "${name}", location: "${location}", priority: ${priority}, tags: ${JSON.stringify(tags)} }) {
              id
              name
              location
              priority
              tags
            }
          }
        `)

        res.status.should.equal(200)
        const site = res.data.data.createSite

        should.exist(site.id)
        site.name.should.equal(name)
        site.location.should.equal(location)
        site.priority.should.equal(priority)
        site.tags.should.eql(tags)

        store.sites.push(site)
      }
    })

    it('createSensor', async () => {
      const inputs = [
        { name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 },
        { name: 'temperature', unit: 'fahrenheit', minSafe: -22, maxSafe: 86 },
        { name: 'pressure', unit: 'atmosphere', minSafe: 0.5, maxSafe: 1 },
        { name: 'pressure', unit: 'pascal', minSafe: 50000, maxSafe: 100000 }
      ]

      for (let i in inputs) {
        const siteId = store.sites[i].id
        const { name, unit, minSafe, maxSafe } = inputs[i]

        const res = await query(`
          mutation {
            createSensor (input: { siteId: "${siteId}", name: "${name}", unit: "${unit}", minSafe: ${minSafe}, maxSafe: ${maxSafe} }) {
              id
              siteId
              name
              unit
              minSafe
              maxSafe
            }
          }
        `)

        res.status.should.equal(200)
        const sensor = res.data.data.createSensor

        should.exist(sensor.id)
        sensor.siteId.should.equal(siteId)
        sensor.name.should.equal(name)
        sensor.unit.should.equal(unit)
        sensor.minSafe.should.equal(minSafe)
        sensor.maxSafe.should.equal(maxSafe)

        store.sensors.push(sensor)
      }
    })

    it('installSwitch', async () => {
      const inputs = [
        { name: 'light', state: 'OFF', states: ['OFF', 'ON'] },
        { name: 'light', state: 'OFF', states: ['OFF', 'ON'] },
        { name: 'light', state: 'OFF', states: ['OFF', 'ON'] },
        { name: 'light', state: 'OFF', states: ['OFF', 'ON'] }
      ]

      for (let i in inputs) {
        const siteId = store.sites[i].id
        const { name, state, states } = inputs[i]

        const res = await query(`
          mutation {
            installSwitch (input: { siteId: "${siteId}", name: "${name}", state: "${state}", states: ${JSON.stringify(states)} }) {
              id
              siteId
              name
              state
              states
            }
          }
        `)

        res.status.should.equal(200)
        const swtch = res.data.data.installSwitch

        should.exist(swtch.id)
        swtch.siteId.should.equal(siteId)
        swtch.name.should.equal(name)
        swtch.state.should.equal(state)
        swtch.states.should.eql(states)

        store.switches.push(swtch)
      }
    })

    it('createAlarm', async () => {
      const inputs = [
        { serialNo: '1001', material: 'co' },
        { serialNo: '1002', material: 'co' },
        { serialNo: '1003', material: 'co' },
        { serialNo: '1004', material: 'co' }
      ]

      for (let i in inputs) {
        const siteId = store.sites[i].id
        const { serialNo, material } = inputs[i]

        const res = await query(`
          mutation {
            createAlarm (input: { siteId: "${siteId}", serialNo: "${serialNo}", material: "${material}" }) {
              id
              siteId
              serialNo
              material
            }
          }
        `)

        res.status.should.equal(200)
        const alarm = res.data.data.createAlarm

        should.exist(alarm.id)
        alarm.siteId.should.equal(siteId)
        alarm.serialNo.should.equal(serialNo)
        alarm.material.should.equal(material)

        store.alarms.push(alarm)
      }
    })

    it('createCamera', async () => {
      const inputs = [
        { serialNo: '2001', resolution: 921600 },
        { serialNo: '2002', resolution: 921600 },
        { serialNo: '2003', resolution: 921600 },
        { serialNo: '2004', resolution: 921600 }
      ]

      for (let i in inputs) {
        const siteId = store.sites[i].id
        const { serialNo, resolution } = inputs[i]

        const res = await query(`
          mutation {
            createCamera (input: { siteId: "${siteId}", serialNo: "${serialNo}", resolution: ${resolution} }) {
              id
              siteId
              serialNo
              resolution
            }
          }
        `)

        res.status.should.equal(200)
        const camera = res.data.data.createCamera

        should.exist(camera.id)
        camera.siteId.should.equal(siteId)
        camera.serialNo.should.equal(serialNo)
        camera.resolution.should.equal(resolution)

        store.cameras.push(camera)
      }
    })

    // GET ALL

    it('sites', async () => {
      const res = await query(`
        query {
          sites {
            id
            name
            location
            priority
            tags
            sensors {
              id
              siteId
              name
              unit
              minSafe
              maxSafe
            }
            switches {
              id
              siteId
              name
              state
              states
            }
            assets {
              id
              siteId
              serialNo
              ... on Alarm {
                material
              }
              ... on Camera {
                resolution
              }
            }
          }
        }
      `)

      res.status.should.equal(200)
      const results = res.data.data.sites

      for (let site of store.sites) {
        const sensors = store.sensors.filter(s => s.siteId === site.id)
        const switches = store.switches.filter(s => s.siteId === site.id)

        const s = results.find(s => s.id === site.id)

        s.id.should.equal(site.id)
        s.name.should.equal(site.name)
        s.location.should.equal(site.location)
        s.priority.should.equal(site.priority)
        s.tags.should.eql(site.tags)

        s.sensors.should.eql(sensors)
        s.switches.should.eql(switches)

        const alarms = store.alarms.filter(a => a.siteId === site.id)
        const cameras = store.cameras.filter(c => c.siteId === site.id)
        const assets = [].concat(alarms, cameras)
        s.assets.should.eql(assets)
      }
    })

    it('sensors', async () => {
      for (let { id: siteId } of store.sites) {
        const res = await query(`
          query {
            sensors (siteId: "${siteId}") {
              id
              siteId
              name
              unit
              minSafe
              maxSafe
            }
          }
        `)

        res.status.should.equal(200)
        const results = res.data.data.sensors

        const sensors = store.sensors.filter(s => s.siteId === siteId)
        results.should.eql(sensors)
      }
    })

    it('switches', async () => {
      for (let { id: siteId } of store.sites) {
        const res = await query(`
          query {
            switches (siteId: "${siteId}") {
              id
              siteId
              name
              state
              states
            }
          }
        `)

        res.status.should.equal(200)
        const results = res.data.data.switches

        const switches = store.switches.filter(s => s.siteId === siteId)
        results.should.eql(switches)
      }
    })

    it('assets', async () => {
      for (let { id: siteId } of store.sites) {
        const res = await query(`
          query {
            assets (siteId: "${siteId}") {
              id
              siteId
              serialNo
              ... on Alarm {
                material
              }
              ... on Camera {
                resolution
              }
            }
          }
        `)

        res.status.should.equal(200)
        const results = res.data.data.assets

        const alarms = store.alarms.filter(a => a.siteId === siteId)
        const cameras = store.cameras.filter(c => c.siteId === siteId)
        const assets = [].concat(alarms, cameras)
        results.should.eql(assets)
      }
    })

    // UPDATE

    it('updateSite', async () => {
      const inputs = [
        { name: 'Plant A', location: 'Ottawa, ON, CANADA', priority: 2, tags: ['energy', 'power'] },
        { name: 'Plant B', location: 'Toronto, ON, CANADA', priority: 3, tags: ['energy', 'hydro'] },
        { name: 'Plant C', location: 'Montreal, QC, CANADA', priority: 4, tags: ['energy', 'gas'] },
        { name: 'Plant D', location: 'Vancouver, BC, CANADA', priority: 5, tags: ['energy', 'oil'] }
      ]

      for (let i in inputs) {
        const id = store.sites[i].id
        const { name, location, priority, tags } = inputs[i]

        const res = await query(`
          mutation {
            updateSite (id: "${id}", input: { name: "${name}", location: "${location}", priority: ${priority}, tags: ${JSON.stringify(tags)} }) {
              id
              name
              location
              priority
              tags
            }
          }
        `)

        res.status.should.equal(200)
        const site = res.data.data.updateSite

        site.id.should.equal(id)
        site.name.should.equal(name)
        site.location.should.equal(location)
        site.priority.should.equal(priority)
        site.tags.should.eql(tags)

        store.sites[i] = site
      }
    })

    it('updateSensor', async () => {
      const inputs = [
        { name: 'temperature', unit: 'celsius', minSafe: -20, maxSafe: 20 },
        { name: 'temperature', unit: 'fahrenheit', minSafe: -4, maxSafe: 68 },
        { name: 'pressure', unit: 'atmosphere', minSafe: 0.6, maxSafe: 0.8 },
        { name: 'pressure', unit: 'pascal', minSafe: 60000, maxSafe: 80000 }
      ]

      for (let i in inputs) {
        const { id, siteId } = store.sensors[i]
        const { name, unit, minSafe, maxSafe } = inputs[i]

        const res = await query(`
          mutation {
            updateSensor (id: "${id}", input: { siteId: "${siteId}", name: "${name}", unit: "${unit}", minSafe: ${minSafe}, maxSafe: ${maxSafe} }) {
              id
              siteId
              name
              unit
              minSafe
              maxSafe
            }
          }
        `)

        res.status.should.equal(200)
        const sensor = res.data.data.updateSensor

        sensor.id.should.equal(id)
        sensor.siteId.should.equal(siteId)
        sensor.name.should.equal(name)
        sensor.unit.should.equal(unit)
        sensor.minSafe.should.equal(minSafe)
        sensor.maxSafe.should.equal(maxSafe)

        store.sensors[i] = sensor
      }
    })

    it('setSwitch', async () => {
      const inputs = [
        { state: 'ON' },
        { state: 'ON' },
        { state: 'ON' },
        { state: 'ON' }
      ]

      for (let i in inputs) {
        const { id, siteId, name, states } = store.switches[i]
        const { state } = inputs[i]

        const res = await query(`
          mutation {
            setSwitch (id: "${id}", state: "${state}") {
              id
              siteId
              name
              state
              states
            }
          }
        `)

        res.status.should.equal(200)
        const swtch = res.data.data.setSwitch

        swtch.id.should.equal(id)
        swtch.siteId.should.equal(siteId)
        swtch.name.should.equal(name)
        swtch.state.should.equal(state)
        swtch.states.should.eql(states)

        store.switches[i] = swtch
      }
    })

    it('updateAlarm', async () => {
      const inputs = [
        { serialNo: '1001', material: 'smoke' },
        { serialNo: '1002', material: 'smoke' },
        { serialNo: '1003', material: 'smoke' },
        { serialNo: '1004', material: 'smoke' }
      ]

      for (let i in inputs) {
        const { id, siteId } = store.alarms[i]
        const { serialNo, material } = inputs[i]

        const res = await query(`
          mutation {
            updateAlarm (id: "${id}", input: { siteId: "${siteId}", serialNo: "${serialNo}", material: "${material}" }) {
              id
              siteId
              serialNo
              material
            }
          }
        `)

        res.status.should.equal(200)
        const alarm = res.data.data.updateAlarm

        alarm.id.should.equal(id)
        alarm.siteId.should.equal(siteId)
        alarm.serialNo.should.equal(serialNo)
        alarm.material.should.equal(material)

        store.alarms[i] = alarm
      }
    })

    it('updateCamera', async () => {
      const inputs = [
        { serialNo: '2001', resolution: 2073600 },
        { serialNo: '2002', resolution: 2073600 },
        { serialNo: '2003', resolution: 2073600 },
        { serialNo: '2004', resolution: 2073600 }
      ]

      for (let i in inputs) {
        const { id, siteId } = store.cameras[i]
        const { serialNo, resolution } = inputs[i]

        const res = await query(`
          mutation {
            updateCamera (id: "${id}", input: { siteId: "${siteId}", serialNo: "${serialNo}", resolution: ${resolution} }) {
              id
              siteId
              serialNo
              resolution
            }
          }
        `)

        res.status.should.equal(200)
        const camera = res.data.data.updateCamera

        camera.id.should.equal(id)
        camera.siteId.should.equal(siteId)
        camera.serialNo.should.equal(serialNo)
        camera.resolution.should.equal(resolution)

        store.cameras[i] = camera
      }
    })

    // GET

    it('site', async () => {
      for (let { id, name, location, priority, tags } of store.sites) {
        const res = await query(`
          query {
            site (id: "${id}") {
              id
              name
              location
              priority
              tags
            }
          }
        `)

        res.status.should.equal(200)
        const site = res.data.data.site

        site.id.should.equal(id)
        site.name.should.equal(name)
        site.location.should.equal(location)
        site.priority.should.equal(priority)
        site.tags.should.eql(tags)
      }
    })

    it('sensor', async () => {
      for (let { id, siteId, name, unit, minSafe, maxSafe } of store.sensors) {
        const res = await query(`
          query {
            sensor (id: "${id}") {
              id
              siteId
              name
              unit
              minSafe
              maxSafe
            }
          }
        `)

        res.status.should.equal(200)
        const sensor = res.data.data.sensor

        sensor.id.should.equal(id)
        sensor.siteId.should.equal(siteId)
        sensor.name.should.equal(name)
        sensor.unit.should.equal(unit)
        sensor.minSafe.should.equal(minSafe)
        sensor.maxSafe.should.equal(maxSafe)
      }
    })

    it('switch', async () => {
      for (let { id, siteId, name, state, states } of store.switches) {
        const res = await query(`
          query {
            switch (id: "${id}") {
              id
              siteId
              name
              state
              states
            }
          }
        `)

        res.status.should.equal(200)
        const swtch = res.data.data.switch

        swtch.id.should.equal(id)
        swtch.siteId.should.equal(siteId)
        swtch.name.should.equal(name)
        swtch.state.should.equal(state)
        swtch.states.should.eql(states)
      }
    })

    it('asset', async () => {
      for (let { id, siteId, serialNo, material } of store.alarms) {
        const res = await query(`
          query {
            asset (id: "${id}") {
              id
              siteId
              serialNo
              ... on Alarm {
                material
              }
            }
          }
        `)

        res.status.should.equal(200)
        const alarm = res.data.data.asset

        alarm.id.should.equal(id)
        alarm.siteId.should.equal(siteId)
        alarm.serialNo.should.equal(serialNo)
        alarm.material.should.equal(material)
      }

      for (let { id, siteId, serialNo, resolution } of store.cameras) {
        const res = await query(`
          query {
            asset (id: "${id}") {
              id
              siteId
              serialNo
              ... on Camera {
                resolution
              }
            }
          }
        `)

        res.status.should.equal(200)
        const camera = res.data.data.asset

        camera.id.should.equal(id)
        camera.siteId.should.equal(siteId)
        camera.serialNo.should.equal(serialNo)
        camera.resolution.should.equal(resolution)
      }
    })

    // DELETE

    it('deleteAsset', async () => {
      for (let { id } of store.alarms) {
        const res = await query(`
          mutation {
            deleteAsset (id: "${id}")
          }
        `)

        res.status.should.equal(200)
        const result = res.data.data.deleteAsset

        result.should.be.true()
      }

      for (let { id } of store.cameras) {
        const res = await query(`
          mutation {
            deleteAsset (id: "${id}")
          }
        `)

        res.status.should.equal(200)
        const result = res.data.data.deleteAsset

        result.should.be.true()
      }
    })

    it('removeSwitch', async () => {
      for (let { id } of store.switches) {
        const res = await query(`
          mutation {
            removeSwitch (id: "${id}")
          }
        `)

        res.status.should.equal(200)
        const result = res.data.data.removeSwitch

        result.should.be.true()
      }
    })

    it('deleteSensor', async () => {
      for (let { id } of store.sensors) {
        const res = await query(`
          mutation {
            deleteSensor (id: "${id}")
          }
        `)

        res.status.should.equal(200)
        const result = res.data.data.deleteSensor

        result.should.be.true()
      }
    })

    it('deleteSite', async () => {
      for (let { id } of store.sites) {
        const res = await query(`
          mutation {
            deleteSite (id: "${id}")
          }
        `)

        res.status.should.equal(200)
        const result = res.data.data.deleteSite

        result.should.be.true()
      }
    })
  })
})
