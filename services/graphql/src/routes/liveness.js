const express = require('express')

const router = express.Router()

router.use('/liveness', (req, res) => {
  res.sendStatus(200)
})

module.exports = router
