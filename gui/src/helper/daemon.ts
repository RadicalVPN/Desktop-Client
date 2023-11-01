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
}
