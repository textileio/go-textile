module.exports = {
  babel: {
    presets: ['mobx'],
    plugins: [
      ['@babel/plugin-proposal-decorators', { 'legacy': true }],
      ['@babel/plugin-proposal-class-properties', { 'loose': true }]
    ]
  },
  webpack: {
    configure: {
      externals: {
        ed25519: 'ed25519'
      }
      // output: {
      //   // eslint-disable-next-line no-path-concat
      //   path: __dirname + '/build',
      //   publicPath: 'app'
      // }
    }
  }
}
