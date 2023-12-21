import { BrowserWindow, Menu, Tray as ElectronTray, app, Notification } from 'electron'
import { join } from 'path'
import { DaemonHelper } from '../../src/helper/daemon'

export class Tray {
  private win: BrowserWindow

  constructor(win: BrowserWindow) {
    this.win = win
  }

  loadTray() {
    new ElectronTray(join(process.env.VITE_PUBLIC, 'logo.ico')).setContextMenu(this.getTrayTemplate())
  }

  static showTrayNotification() {
    new Notification({
      title: 'RadicalVPN',
      body: 'RadicalVPN is running in the background',
    }).show()
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
          await new DaemonHelper().disconnectFromServer()

          this.win.destroy()
          app.quit()
        },
      },
    ])
  }
}
