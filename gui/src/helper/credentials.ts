import { readFileSync } from 'fs'

export class DaemonCredentials {
  static getCredentials() {
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
}
