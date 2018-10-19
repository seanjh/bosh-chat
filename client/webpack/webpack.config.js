const path = require('path')

module.exports = {
  mode: 'development',
  entry: './client/main.js',
  output: {
    path: path.resolve(__dirname, '../../static'),
    filename: 'app.js',
    publicPath: '/static/'
  }
}
