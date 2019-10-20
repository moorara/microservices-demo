/* eslint-env mocha */
const should = require('should')

const Site = require('../../../models/site')

describe('Site', () => {
  let SiteModel

  beforeEach(() => {
    SiteModel = new Site().Model
  })

  describe('SiteModel', () => {
    it('should create a new site with name and location', () => {
      const site = new SiteModel({
        name: 'New Site',
        location: 'Ottawa, ON'
      })

      should.exist(site.id)
      site.name.should.equal('New Site')
      site.location.should.equal('Ottawa, ON')
      should.exist(site.createdAt)
      should.exist(site.updatedAt)
    })

    it('should create a new site with name, location, and tags', () => {
      const site = new SiteModel({
        name: 'New Site',
        location: 'Ottawa, ON',
        tags: ['hydro', 'power']
      })

      should.exist(site.id)
      site.name.should.equal('New Site')
      site.location.should.equal('Ottawa, ON')
      site.tags[0].should.eql('hydro')
      site.tags[1].should.eql('power')
      should.exist(site.createdAt)
      should.exist(site.updatedAt)
    })

    it('should create a new site with name, location, tags, and priority', () => {
      const site = new SiteModel({
        name: 'New Site',
        location: 'Ottawa, ON',
        tags: ['hydro', 'power'],
        priority: 3
      })

      should.exist(site.id)
      site.name.should.equal('New Site')
      site.location.should.equal('Ottawa, ON')
      site.tags[0].should.eql('hydro')
      site.tags[1].should.eql('power')
      site.priority.should.equal(3)
      should.exist(site.createdAt)
      should.exist(site.updatedAt)
    })

    it('should create a new site and format name properly', () => {
      const site = new SiteModel({
        name: 'new site',
        location: 'Ottawa, ON'
      })

      should.exist(site.id)
      site.name.should.equal('new site')
      site.location.should.equal('Ottawa, ON')
      should.exist(site.createdAt)
      should.exist(site.updatedAt)
      should.equal(site.formatName(), 'New Site')
    })
  })

  describe('toJSON', () => {
    it('should return only public properties of site', () => {
      const site = new SiteModel({
        name: 'New Site',
        location: 'Ottawa, ON',
        tags: ['hydro', 'power'],
        priority: 3
      }).toJSON()

      should.exist(site.id)
      site.name.should.equal('New Site')
      site.location.should.equal('Ottawa, ON')
      site.tags[0].should.eql('hydro')
      site.tags[1].should.eql('power')
      site.priority.should.equal(3)
    })
  })
})
