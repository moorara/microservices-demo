import React from 'react'

const AboutPage = (props) => (
  <section className="hero is-light is-large is-bold">
    <div className="hero-body">
      <div className="container">
        <div className="content">
          <h1 className="title">About</h1>
          <h2 className="subtitle">Control Center Application</h2>
          <p>
            This is a demo application for monitoring hypothetical remote sites in real-time.
            Each site has a number of sensors sending real-time values continuously.
          </p>
          <p>This demo application is built using some <strong>best-practices</strong> and cool technologies.</p>
          <dl>
            <dt>Architecture</dt>
            <dd><strong>Microservices</strong></dd>
            <dd><strong>Docker, Swarm, Kubernetes</strong></dd>
            <dt>Front-End</dt>
            <dd><strong>React, Redux</strong></dd>
            <dd><strong>Bulma</strong></dd>
            <dd><strong>Webpack</strong></dd>
            <dt>Back-End</dt>
            <dd><strong>Go</strong></dd>
            <dd><strong>Node.js, Express.js</strong></dd>
            <dt>Database</dt>
            <dd><strong>MongoDB, ArangoDB</strong></dd>
            <dd><strong>Postgres</strong></dd>
            <dd><strong>RabbitMQ, Redis</strong></dd>
            <dt>API</dt>
            <dd><strong>REST</strong></dd>
            <dt>Gateways</dt>
            <dd><strong>Traefik, Caddy</strong></dd>
            <dt>Monitoring</dt>
            <dd><strong>Fluentd, ElasticSearch, Kibana</strong></dd>
            <dd><strong>Prometheus, Grafana</strong></dd>
          </dl>
        </div>
      </div>
    </div>
  </section>
)

export default AboutPage
