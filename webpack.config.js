const PATH = process.cwd()
const HtmlWebpackLoader = require('html-webpack-plugin')

module.exports = {
  entry: PATH + '/src/views/index',
  mode: 'development',
  devtool: 'source-map',
  output: {
    path: PATH + '/dist',
    filename: 'main.bundle.js',
  },
  resolve: {
    extensions: ['.ts', '.tsx', '.js', '.json'],
  },
  module: {
    rules: [
      {
        test: /\.(ts|js)x?$/,
        exclude: /node_modules/,
        loader: 'babel-loader'
      }
    ]
  },
  plugins: [
    new HtmlWebpackLoader({
      template: PATH + '/src/views/index.html'
    })
  ]
}
