const { app, BrowserWindow } = require('electron')
const { spawn } = require('child_process')
const fs = require('fs')
const path = require('path')


const url = 'http://localhost:51234'

function createWindow() {
  const win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false,
    },
  })

  const goBinary = getGoBinary()
  const server = spawn(goBinary)
  win.loadURL(url)

  win.on('close', (event) => {
    server.kill()
  })
}

app.whenReady().then(() => {
  createWindow()
  
  app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
      app.quit()
    }
  })
  
  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

function getGoBinary() {
  let goBinaryPath
  let tempBinaryPath

  // Determine the platform and set the server binary path accordingly
  switch (process.platform) {
  case 'win32':
    goBinaryPath = path.join(__dirname, 'binaries', 'server-windows.exe')
    tempBinaryPath = path.join(app.getPath('temp'), 'server-windows.exe')
    break;
  case 'darwin':
    goBinaryPath = path.join(__dirname, 'binaries', 'server-macos')
    tempBinaryPath = path.join(app.getPath('temp'), 'server-macos')
    break;
  case 'linux':
    goBinaryPath = path.join(__dirname, 'binaries', 'server-linux')
    tempBinaryPath = path.join(app.getPath('temp'), 'server-linux')
    break;
  default:
    console.error('Unsupported platform');
    app.quit();
    return;
  }

  fs.copyFileSync(goBinaryPath, tempBinaryPath);
  fs.chmodSync(tempBinaryPath, 0o755);

  return tempBinaryPath
}
