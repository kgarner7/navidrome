// These defaults are only used in development mode. When bundled in the app,
// the __APP_CONFIG__ object is dynamically filled by the ServeIndex function,
// in the /server/app/serve_index.go
const defaultConfig = {
  version: 'dev',
  firstTime: false,
  baseURL: '',
  variousArtistsId: '03b645ef2100dfc42fa9785ea3102295', // See consts.VariousArtistsID in consts.go
  // Login backgrounds from https://unsplash.com/collections/1065384/music-wallpapers
  loginBackgroundURL: 'https://source.unsplash.com/collection/1065384/1600x900',
  maxSidebarPlaylists: 100,
  enableTranscodingConfig: true,
  enableDownloads: true,
  enableFavourites: true,
  losslessFormats: 'FLAC,WAV,ALAC,DSF',
  welcomeMessage: '',
  gaTrackingId: '',
  devActivityPanel: true,
  enableStarRating: true,
  defaultTheme: 'Dark',
  defaultLanguage: '',
  defaultUIVolume: 100,
  enableUserEditing: true,
  enableSharing: true,
  shareURL: '',
  defaultDownloadableShare: true,
  devSidebarPlaylists: true,
  lastFMEnabled: true,
  listenBrainzEnabled: true,
  enableExternalServices: true,
  enableCoverAnimation: true,
  devShowArtistPage: true,
  enableReplayGain: true,
  defaultDownsamplingFormat: 'opus',
  publicBaseUrl: '/share',
}

let config: typeof defaultConfig

declare global {
  interface Window {
    __APP_CONFIG__?: string
    __SHARE_INFO__?: string
  }
}

try {
  const appConfig = JSON.parse(window.__APP_CONFIG__ ?? '{}')
  config = {
    ...defaultConfig,
    ...appConfig,
  }
} catch (e) {
  config = defaultConfig
}

export let shareInfo: unknown

try {
  shareInfo = JSON.parse(window.__SHARE_INFO__ ?? 'null')
} catch (e) {
  shareInfo = null
}

export default config