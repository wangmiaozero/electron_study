module.exports = {
  pluginOptions: {
    electronBuilder: {
      builderOptions: {
        appId: 'com.ccsdesgin.note',
        productName: '创云笔记',
        copyright: 'Copyright © 2022',
        directories: {
          output: './dist_electron'
        },
        extraResources: {
          from: './server',
          to: 'server'
        },
        win: {
          icon: './public/favicon.ico',
          target: [
            {
              target: 'nsis',
              arch: [ 'x64' ]
            }
          ]
        },
        nsis: {
          oneClick: false,
          allowElevation: true,
          allowToChangeInstallationDirectory: true
        }
      },
      preload: './src/preload.js'
    }
  }
}