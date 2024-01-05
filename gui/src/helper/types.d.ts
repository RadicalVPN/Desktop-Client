interface IDaemonVersionInfo {
  version: string
  nightly: {
    isNightly: string
    isOutdated: string
  }
  release: {
    isRelease: string
    isOutdated: string
  }
}
