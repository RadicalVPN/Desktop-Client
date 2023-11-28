import axios from 'axios'
import { readFileSync } from 'node:fs'
import { spawn } from 'node:child_process'

export interface ParsedLog {
  level: string
  color: string
  message: string
}

export class DaemonHelper {
  public getCredentials() {
    let rawCredentials: string

    switch (process.platform) {
      case 'darwin':
        rawCredentials = readFileSync('/Library/Application Support/RadicalVPN/service.txt', 'utf-8')
        break
      case 'win32':
        rawCredentials = readFileSync('C:\\Program Files\\RadicalVPN\\service.txt', 'utf-8')
        break
      default:
        rawCredentials = '1234|dummy'
    }

    const credentials = rawCredentials.split('|')
    return {
      port: parseInt(credentials[0]),
      secret: credentials[1],
    }
  }

  public async isAuthed() {
    const credentials = this.getCredentials()

    try {
      const resp = await axios.get(`http://localhost:${credentials.port}/`, {
        headers: {
          'x-radical-daemon-secret': credentials.secret,
        },
        validateStatus: () => true,
      })

      return resp.status === 200
    } catch {
      return false
    }
  }

  public async getServerList() {
    const credentials = this.getCredentials()

    try {
      const resp = await axios.get(`http://localhost:${credentials.port}/server`, {
        headers: {
          'x-radical-daemon-secret': credentials.secret,
        },
        validateStatus: () => true,
      })

      return resp.data
    } catch {
      return []
    }
  }

  public async connectToServer(nodeLocation: string, privacyFirewall: string) {
    const credentials = this.getCredentials()

    try {
      const resp = await axios.post(
        `http://localhost:${credentials.port}/local/connect`,
        {
          nodeLocation: nodeLocation,
          privacyFirewall,
        },
        {
          headers: {
            'x-radical-daemon-secret': credentials.secret,
          },
          validateStatus: () => true,
        },
      )

      return {
        status: resp.status === 200,
        data: resp.data,
      }
    } catch {
      return {
        status: false,
        data: {},
      }
    }
  }

  public async disconnectFromServer() {
    const credentials = this.getCredentials()

    try {
      const resp = await axios.post(
        `http://localhost:${credentials.port}/local/disconnect`,
        {},
        {
          headers: {
            'x-radical-daemon-secret': credentials.secret,
          },
          validateStatus: () => true,
        },
      )

      return resp.status === 200
    } catch {
      return false
    }
  }

  /**
   * Get the connection state of the daemon
   * true = vpn is connected
   * false = vpn is disconnected
   * @returns {Promise<boolean>}
   */
  public async getConnectionState(): Promise<boolean> {
    const credentials = this.getCredentials()

    try {
      const resp = await axios.get(`http://localhost:${credentials.port}/local/connected`, {
        headers: {
          'x-radical-daemon-secret': credentials.secret,
        },
        validateStatus: () => true,
      })

      return resp.data.connected || false
    } catch {
      return false
    }
  }

  public async getDaemonVersion(): Promise<string> {
    const credentials = this.getCredentials()

    try {
      const resp = await axios.get(`http://localhost:${credentials.port}/version`, {
        headers: {
          'x-radical-daemon-secret': credentials.secret,
        },
        validateStatus: () => true,
      })

      return resp.data.version || 'Unknown'
    } catch (e) {
      console.error(e)
      return 'Unknown'
    }
  }

  public async isDaemonInstallRequired(): Promise<boolean> {
    return new Promise((resolve, reject) => {
      if (process.platform != 'darwin') {
        resolve(false)
        return
      }

      const cmd = spawn(
        '/Applications/RadicalVPN.app/Contents/MacOS/RadicalVPN-Installer.app/Contents/MacOS/RadicalVPN-Installer',
        ['--install-required'],
      )

      cmd.stdout.on('data', (data) => {
        console.log('data', data.toString())
      })
      cmd.stderr.on('data', (err) => {
        reject(err)
      })

      cmd.on('error', (err) => {
        reject(err)
      })

      cmd.on('exit', (code) => {
        resolve(code === 0)
      })
    })
  }

  public async installDaemon(): Promise<boolean> {
    return new Promise((resolve, reject) => {
      const cmd = spawn(
        '/Applications/RadicalVPN.app/Contents/MacOS/RadicalVPN-Installer.app/Contents/MacOS/RadicalVPN-Installer',
        ['--install'],
      )

      cmd.stdout.on('data', (data) => {
        console.log('data', data.toString())
      })

      cmd.stderr.on('data', (err) => {
        reject(err)
      })

      cmd.on('error', (err) => {
        reject(err)
      })

      cmd.on('exit', (code) => {
        resolve(code === 0)
      })
    })
  }

  public async daemonIsStarted(): Promise<boolean> {
    const credentials = this.getCredentials()

    try {
      const resp = await axios.get(`http://localhost:${credentials.port}/ping`, {
        headers: {
          'x-radical-daemon-secret': credentials.secret,
        },
        validateStatus: () => true,
      })

      return resp.status === 200
    } catch (e) {
      return false
    }
  }

  public async getLogs(): Promise<ParsedLog[]> {
    let logs: string

    switch (process.platform) {
      case 'darwin':
        logs = readFileSync('/Library/Application Support/RadicalVPN/radicalvpn.log', 'utf-8')
        break
      case 'win32':
        logs = readFileSync('C:\\Program Files\\RadicalVPN\\radicalvpn.log', 'utf-8')
        break
      default:
        logs = ''
    }

    const logsArr = logs.split('\n')

    return logsArr.map((log) => {
      const logLvl = log.split(' ')[4]?.toLowerCase() || 'info'

      return {
        level: logLvl,
        color:
          logLvl === 'erro'
            ? 'danger'
            : logLvl === 'trac'
              ? 'danger'
              : logLvl === 'warn'
                ? 'warning'
                : logLvl === 'debu'
                  ? 'secondary'
                  : 'primary',
        message: log,
      }
    })
  }
}
