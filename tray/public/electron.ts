const { app, BrowserWindow, shell, ipcMain, Tray } = require('electron')
const path  = require('path')
const isDev = require('electron-is-dev')
const url = require('url')

let mainWindow
let tray

const gotTheLock = app.requestSingleInstanceLock()
if (!gotTheLock) {
  app.quit()
} else {
  app.on('second-instance', (event, commandLine, workingDirectory) => {
    // Someone tried to run a second instance, we should focus our window.

    // Protocol handler for windows
    // commandLine: An array of the second instance’s (command line / deep linked) arguments
    if (process.platform === 'win32') {
      // Keep only command line / deep linked arguments
      console.log(commandLine.slice(1))
    }

    if (mainWindow) {
      if (mainWindow.isMinimized()) {
        mainWindow.restore()
      }
      mainWindow.focus()
    }
  })
}

const getWindowPosition = () => {
  const windowBounds = mainWindow.getBounds()
  const trayBounds = tray.getBounds()

  // Center window horizontally below the tray icon
  const x = Math.round(trayBounds.x + (trayBounds.width / 2) - (windowBounds.width / 2))

  // Position window 4 pixels vertically below the tray icon
  const y = Math.round(trayBounds.y + trayBounds.height + 4)

  return { x: x, y: y }
}

const createTray = () => {
  const pathname = path.join(__dirname, 'assets', 'tray.png')
  tray = new Tray(pathname)
  tray.on('right-click', toggleWindow)
  tray.on('double-click', toggleWindow)
  tray.on('click', toggleWindow)
}

const showWindow = () => {
  const position = getWindowPosition()
  mainWindow.setPosition(position.x, position.y, false)
  mainWindow.show()
  mainWindow.focus()
}

const toggleWindow = () => {
  if (mainWindow.isVisible()) {
    mainWindow.hide()
  } else {
    showWindow()
  }
}

const createWindow = () => {
  mainWindow = new BrowserWindow({
    width: 300,
    height: 450,
    show: false,
    frame: false,
    fullscreenable: false,
    transparent: true,
    movable: false,
    resizable: false,
    minimizable: false,
    maximizable: false,
    skipTaskbar: true,
  })

  mainWindow.loadURL(isDev ? 'http://localhost:3000' : url.format({
    pathname: path.join(__dirname, '../build', 'index.html'),
    protocol: 'file:',
    slashes: true
  }))

  if (isDev) {
    const {
      default: installExtension,
      REACT_DEVELOPER_TOOLS
    } = require('electron-devtools-installer')

    installExtension(REACT_DEVELOPER_TOOLS)
      .then(name => {
        console.log(`Added Extension: ${name}`)
      })
      .catch(err => {
        console.log('An error occurred: ', err)
      })
  }

  // Emitted when the window is closed.
  mainWindow.on('closed', function () {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    mainWindow = null
    tray = null
  })

  mainWindow.once('ready-to-show', () => {
    ipcMain.on('open-external-window', (event, arg) => {
      shell.openExternal(arg)
    })
  })
}

app.on('ready', () => {
  createTray()
  createWindow()
})

app.on('window-all-closed', () => {
  app.quit()
})

app.on('activate', () => {
  if (mainWindow === null) {
    createWindow()
  }
})

ipcMain.on('load-page', (event, arg) => {
  mainWindow.loadURL(arg)
})
