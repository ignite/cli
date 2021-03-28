module.exports = {
  configureWebpack: {
    module: {
      rules: [
        {
          test: /\.js?$/,
          exclude: /node_modules\/(?!(@tendermint)\/).*/,
          loader: 'babel-loader'
        },
      ]
    }
  }
}