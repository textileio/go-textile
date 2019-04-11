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
        electron: 'electron'
      }
    }
  }
}
