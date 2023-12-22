import { BrowserWindow, Menu, Tray as ElectronTray, app, Notification, nativeImage } from 'electron'
import { join } from 'path'

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
          this.win.show()
        },
      },
      {
        label: 'Quit',
        click: async () => {
          //await new DaemonHelper().disconnectFromServer()

          this.win.destroy()
          app.quit()
        },
      },
    ])
  }
}
