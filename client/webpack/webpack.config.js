const path = require('path')

module.exports = {
  target: 'web',
  entry: [path.join(__dirname, '../index.jsx')],
  output: {
    path: path.join(__dirname, '../../public/js'),
    publicPath: '/js/',
    filename: 'bundle.js'
  },
  module: {
    rules: [
      {
          test: /\.js(x?)$/,
          exclude: /node_modules/,
          use: [
            {
              loader: 'babel-loader',
              query: {
                presets:['env', 'react'],
                cacheDirectory: true
              }
            }
          ]
      }
    ]
  },
  resolve: {
    modules: [
      // path.join(__dirname, "src"),
      'node_modules'
    ],
    extensions: ['.js', '.jsx']
  },
  plugins: [
  ],
  cache: true
}
