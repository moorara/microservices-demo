# react-client

## Commands

| Command                  | Description                                                   |
|--------------------------|---------------------------------------------------------------|
| `yarn run lint`          | Runs linter with configured rules                             |
| `yarn run lint:watch`    | Runs linter and watches for changes                           |
| `yarn run test`          | Runs unit tests                                               |
| `yarn run test:watch`    | Runs unit tests and watches for changes                       |
| `yarn run test:coverage` | Runs unit tests and generate coverage reports                 |
| `yarn run test:update`   | Runs unit tests and updates `jest snapshots`                  |
| `yarn run dev:api`       | Starts development API server using `json-server`             |
| `yarn run dev:webpack`   | Starts development server using built-in `webpack dev server` |
| `yarn run dev:node`      | Starts development server using `node.js`                     |
| `yarn run build:webpack` | Builds application for production using `webpack cli`         |
| `yarn run build:node`    | Builds application for production using `node.js`             |
| `yarn run prod:node`     | Runs a production server for serving the application          |
| `make docker`            | Builds Docker image                                           |
| `make save-images`       | Saves built Docker images                                     |
| `make up`                | Runs the service locally in containers                        |
| `make down`              | Stops and removes local containers                            |

## Guides

### Snapshot Testing

When it comes to snapshot testing React components with `jest`, there are two options:

  1. Snapshot testing with `react-test-renderer`
  2. Snapshot testing with `enzyme` and `enzyme-to-json`
