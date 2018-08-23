const path = require('path')
const grpc = require('grpc')
const protoLoader = require('@grpc/proto-loader')

const PROTO_PATH = path.join(__dirname, '../grpc/switch.proto')

// https://grpc.io/docs/tutorials/basic/node.html
// https://www.npmjs.com/package/@grpc/proto-loader
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: false, // true preserves field names and false changes them to camelCase
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
})
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition)
const proto = protoDescriptor.proto

module.exports = { proto }
