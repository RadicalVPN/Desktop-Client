import { BrowserWindow, Menu, Tray as ElectronTray, app, Notification, nativeImage, net } from 'electron'
import { join } from 'path'
import { DaemonCredentials } from '../../src/helper/credentials'

export class Tray {
  private win: BrowserWindow

  constructor(win: BrowserWindow) {
    this.win = win
  }

  loadTray() {
    new ElectronTray(this.getTrayIcon()).setContextMenu(this.getTrayTemplate())
  }

  static showTrayNotification() {
    new Notification({
      title: 'RadicalVPN',
      body: 'RadicalVPN is running in the background',
    }).show()
  }

  private getTrayIcon() {
    let path: string

    switch (process.platform) {
      case 'win32':
        path = 'tray/windows/logo.ico'
        break
      case 'darwin':
        path = 'tray/darwin/logo2Template.png'
        break
    }

    return nativeImage.createFromPath(join(process.env.VITE_PUBLIC, path))
  }

  private getTrayTemplate() {
    return Menu.buildFromTemplate([
      {
        label: 'Show RadicalVPN',
        click: () => {
          if (process.platform === 'darwin') {
            app.dock.show()
          }

          this.win.show()
        },
      },
      {
        label: 'Quit',
        click: async () => {
          try {
            await this.disconnectFromServer()
          } catch (e) {
            console.error('failed to disconnect from vpn', e)
          }

          this.win.destroy()
          app.quit()
        },
      },
    ])
  }

  /**
   * Special implementation of DaemonHelper().disconnectFromServer()
   * Uses the electron main process to send a request to the daemon
   */
  private async disconnectFromServer() {
    return new Promise((resolve, reject) => {
      const credentials = DaemonCredentials.getCredentials()

      try {
        const req = net.request({
          method: 'POST',
          url: `http://localhost:${credentials.port}/local/disconnect`,
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          headers: {
            'x-radical-daemon-secret': credentials.secret,
          },
        })

        req.on('response', (resp) => {
          resolve(resp)
        })

        req.end()
      } catch (e) {
        reject(e)
      }
    })
  }
}
