/* eslint-env mocha */
const should = require('should')

const Link = require('../../../models/link')

describe('Link', () => {
  let LinkModel

  beforeEach(() => {
    LinkModel = new Link().Model
  })

  describe('LinkModel', () => {
    it('should create a new link with url and title', () => {
      let link = new LinkModel({
        url: 'https://nodejs.org',
        title: 'Node.js'
      })

      should.exist(link.id)
      link.url.should.equal('https://nodejs.org')
      link.title.should.equal('Node.js')
      should.exist(link.createdAt)
      should.exist(link.updatedAt)
    })

    it('should create a new link with url, title, and tags', () => {
      let link = new LinkModel({
        url: 'https://golang.org',
        title: 'Golang',
        tags: ['go']
      })

      should.exist(link.id)
      link.url.should.equal('https://golang.org')
      link.title.should.equal('Golang')
      link.tags[0].should.eql('go')
      should.exist(link.createdAt)
      should.exist(link.updatedAt)
    })

    it('should create a new link with url, title, tags, and rank', () => {
      let link = new LinkModel({
        url: 'https://www.docker.com',
        title: 'Docker',
        tags: ['container'],
        rank: 3
      })

      should.exist(link.id)
      link.url.should.equal('https://www.docker.com')
      link.title.should.equal('Docker')
      link.tags[0].should.eql('container')
      link.rank.should.equal(3)
      should.exist(link.createdAt)
      should.exist(link.updatedAt)
    })

    it('should create a new link and format title properly', () => {
      let link = new LinkModel({
        url: 'https://kubernetes.io',
        title: 'kubernetes website'
      })

      should.exist(link.id)
      link.url.should.equal('https://kubernetes.io')
      link.title.should.equal('kubernetes website')
      should.equal(link.formatTitle(), 'Kubernetes Website')
    })
  })

  describe('toJSON', () => {
    it('should return only public properties of link', () => {
      let link = new LinkModel({
        url: 'https://www.hashicorp.com',
        title: 'HashiCorp',
        rank: 5
      }).toJSON()

      should.exist(link.id)
      link.url.should.equal('https://www.hashicorp.com')
      link.title.should.equal('HashiCorp')
      link.tags.length.should.equal(0)
      link.rank.should.equal(5)
      should.not.exist(link.createdAt)
      should.not.exist(link.updatedAt)
    })
  })
})
