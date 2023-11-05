import axios from 'axios'
import { readFileSync } from 'node:fs'

export class DaemonHelper {
  public getCredentials() {
    let rawCredentials: string

    switch (process.platform) {
      case 'darwin':
        rawCredentials = readFileSync('/Library/Application Support/RadicalVPN/service.txt', 'utf-8')
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

  public async connectToServer(nodeId: string) {
    const credentials = this.getCredentials()

    try {
      const resp = await axios.post(
        `http://localhost:${credentials.port}/local/connect`,
        {
          node: nodeId,
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
}
