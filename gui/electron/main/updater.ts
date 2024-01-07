import { BrowserWindow, app, dialog, net, shell } from 'electron'
import { DaemonCredentials } from '../../src/helper/credentials'

export class Updater {
  private win: BrowserWindow

  constructor(win: BrowserWindow) {
    this.win = win
  }

  async checkForUpdates() {
    const info = await this.getVersionInfo()
    let message = 'A new version of RadicalVPN is available.\n'

    if (info.nightly) {
      message += `\nYou are running a nightly version.\n`
    }

    message += `\nWould you like to download it?`

    if (info.outdated) {
      const result = dialog.showMessageBoxSync({
        title: 'RadicalVPN Update',
        type: 'question',
        buttons: ['Yes', 'No'],
        message,
      })

      // 0 = Yes
      // 1 = No
      if (result === 0) {
        shell.openExternal('https://github.com/RadicalVPN/Desktop-Client/releases')

        this.win.destroy()
        app.quit()
      }
    }
  }

  private async getVersionInfo(): Promise<{
    currentVersion: string
    nightly: boolean
    release: boolean
    outdated: boolean
  }> {
    return new Promise((resolve, reject) => {
      const credentials = DaemonCredentials.getCredentials()

      try {
        const req = net.request({
          method: 'GET',
          url: `http://localhost:${credentials.port}/version`,
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          headers: {
            'x-radical-daemon-secret': credentials.secret,
          },
        })

        req.on('response', (res) => {
          let data = ''

          res.on('data', (chunk) => {
            data += chunk
          })

          res.on('end', () => {
            resolve(JSON.parse(data))
          })
        })

        req.end()
      } catch (e) {
        resolve({
          currentVersion: 'N/A',
          nightly: false,
          release: false,
          outdated: false,
        })
      }
    })
  }
}
