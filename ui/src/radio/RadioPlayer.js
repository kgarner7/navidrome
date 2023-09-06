import clsx from 'clsx'
import IcecastMetadataPlayer from 'icecast-metadata-player'
import { useCallback, useEffect, useRef, useState } from 'react'

import Slider from 'rc-slider/lib/Slider'

import {
  AnimatePauseIcon,
  AnimatePlayIcon,
  CloseIcon,
  DeleteIcon,
  VolumeMuteIcon,
  VolumeUnmuteIcon,
} from 'navidrome-music-player/es/components/Icon'
import RadioTitle from './RadioTitle'
import { useDispatch, useSelector } from 'react-redux'

import { clearQueue } from '../actions'
import { useMediaQuery } from '@material-ui/core'
import RadioPlayerMobile from './RadioPlayerMobile'
import subsonic from '../subsonic'
import config from '../config'
import { sendNotification } from '../utils'
import { useTranslate } from 'react-admin'
import RadioDialog from './RadioDialog'

const DEFAULT_ICON = {
  pause: <AnimatePauseIcon />,
  play: <AnimatePlayIcon />,
  destroy: <CloseIcon />,
  close: <CloseIcon />,
  delete: <DeleteIcon size={24} />,
  volume: <VolumeUnmuteIcon size={26} />,
  mute: <VolumeMuteIcon size={26} />,
}

const MIN_TIME_BETWEEN_SCROBBLE_MS = 30 * 1000
const SCROBBLE_DELAY_MS = 4 * 60 * 1000

const RadioPlayer = ({
  className,
  cover,
  icon = {},
  locale,
  homePageUrl,
  id,
  name,
  streamUrl,
  theme,
}) => {
  const dispatch = useDispatch()
  const audioRef = useRef()
  const translate = useTranslate()

  const [cast, setCast] = useState(null)
  const [currentStream, setCurrentStream] = useState(null)
  const [loading, setLoading] = useState(false)
  const [metadata, setMetadata] = useState({})
  const [playing, setPlaying] = useState(false)
  const [savedVolume, setSavedVolume] = useState(1)
  const [volume, setVolume] = useState(1)
  const [open, setOpen] = useState(false)

  const isMobile = useMediaQuery(
    '(max-width: 768px) and (orientation : portrait)'
  )
  const showNotifications = useSelector(
    (state) => state.settings.notifications || false
  )

  const Spin = () => <span className="loading group">{icon.loading}</span>

  const iconMap = { ...DEFAULT_ICON, ...icon, loading: <Spin /> }

  const mapListenToBar = (vol) => Math.sqrt(vol)
  const mapBarToListen = (vol) => vol ** 2

  useEffect(() => {
    const streamChanged = currentStream !== streamUrl

    if (cast && !streamChanged && cast.state !== 'stopped') {
      return
    }

    if (!config.enableProxy) {
      const node = audioRef.current

      if (node && streamChanged) {
        node.crossOrigin = 'anonymous'
        node.src = streamUrl

        node.play()
        setPlaying(true)

        return () => {
          node.src = ''
        }
      } else {
        return
      }
    }

    if (cast) {
      cast.stop()
      cast.detachAudioElement()
    }

    if (streamUrl) {
      const player = new IcecastMetadataPlayer(streamUrl, {
        onMetadata: (data) => {
          console.log(data)
          if (data.StreamTitle) {
            const split = data.StreamTitle.split(' - ')

            let artist, title

            if (split.length === 1) {
              title = split[0]
            } else {
              artist = split[0]
              title = split.slice(1).join(' - ')
            }

            setMetadata({ artist, title })
          }
        },
        onPlay: () => {
          setLoading(false)
          setPlaying(true)
        },
        onError: (message, error) => {
          console.error(message, error)
        },
        icyDetectionTimeout: 20000,
        enableLogging: true, // set this to true for dev
        audioElement: audioRef.current,
        playbackMethod: 'mediasource',
        // metadataTypes: ['icy'],
      })

      player.id = Math.random()
      player.play()

      setCast(player)
      setLoading(true)
    } else {
      setCast(null)
    }

    setCurrentStream(streamUrl)
    setMetadata({})
  }, [audioRef, cast, currentStream, streamUrl])

  useEffect(() => {
    if (metadata.title) {
      let scrobbled = false

      const currentUpdate = new Date()

      const { artist, fix, title } = metadata

      if ('mediaSession' in navigator) {
        navigator.mediaSession.metadata = new MediaMetadata({
          album: name,
          artist,
          title,
        })
      }

      if (showNotifications && !fix) {
        const body = artist
          ? `${artist} - ${name}`
          : translate('resources.radio.message.noArtistNotif')

        sendNotification(title, body, cover)
      }

      if (artist) {
        subsonic.scrobbleRadio(artist, title, false)

        const timeout = setTimeout(() => {
          scrobbled = true
          subsonic.scrobbleRadio(artist, title, true)
        }, SCROBBLE_DELAY_MS)

        return () => {
          const now = new Date()

          if (!scrobbled) {
            clearTimeout(timeout)
            if (now - currentUpdate > MIN_TIME_BETWEEN_SCROBBLE_MS || fix) {
              subsonic.scrobbleRadio(artist, title, true)
            }
          }
        }
      }
    }
  }, [cover, metadata, name, showNotifications, translate])

  useEffect(() => {
    const audio = audioRef.current

    if (audio) {
      function volumeChange() {
        const { volume } = audio
        setVolume(mapListenToBar(volume))
      }

      audio.addEventListener('volumechange', volumeChange)

      return () => {
        audio.removeEventListener('volumechange', volumeChange)
      }
    }
  }, [audioRef])

  const setAudioVolume = useCallback((volumeBarVal) => {
    if (audioRef.current) {
      audioRef.current.volume = mapBarToListen(volumeBarVal)

      setSavedVolume(volumeBarVal)
      setVolume(volumeBarVal)
    }
  }, [])

  const mute = useCallback(() => {
    const audio = audioRef.current

    if (audio) {
      setVolume(0)
      setSavedVolume(audio.volume)

      audio.volume = 0
    }
  }, [])

  const resetVolume = useCallback(() => {
    setAudioVolume(mapListenToBar(savedVolume || 0.1))
  }, [savedVolume, setAudioVolume])

  const openModal = () => {
    setOpen(true)
  }

  const closeModal = () => {
    setOpen(false)
  }

  const togglePlay = useCallback(() => {
    const audio = audioRef.current

    if (audio) {
      if (audio.paused) {
        audio.play()
        setPlaying(true)
      } else {
        audio.pause()
        setPlaying(false)
      }
    }
  }, [])

  const coverClick = useCallback(() => {
    window.location.href = `#/radio/${id}/show`
  }, [id])

  const stopPlaying = useCallback(() => {
    audioRef.current.src = ''
    dispatch(clearQueue())
  }, [dispatch])

  return (
    <div
      className={clsx(
        'react-jinke-music-player-main',
        {
          'light-theme': theme === 'light',
          'dark-theme': theme === 'dark',
        },
        className
      )}
    >
      {isMobile && (
        <RadioPlayerMobile
          cover={cover}
          icon={iconMap}
          id={id}
          loading={loading}
          locale={locale}
          metadata={metadata}
          name={name}
          onClose={stopPlaying}
          onCoverClick={coverClick}
          onFix={openModal}
          onPlay={togglePlay}
          playing={playing}
        />
      )}
      {!isMobile && (
        <div className={clsx('music-player-panel', 'translate')}>
          <section className="panel-content">
            {cover && (
              <div
                className={clsx('img-content', 'img-rotate', {
                  'img-rotate-pause': !playing || !cover,
                })}
                style={{ backgroundImage: `url(${cover})` }}
                onClick={() => coverClick()}
              />
            )}
            <div className="progress-bar-content">
              {metadata.title && (
                <span className="audio-title" title={metadata.title}>
                  <RadioTitle
                    id={id}
                    isMobile={false}
                    metadata={metadata}
                    name={name}
                    onFix={openModal}
                  />
                </span>
              )}
            </div>
            <div className="player-content">
              <span className="group">
                {loading ? (
                  <span
                    className="group loading-icon"
                    title={locale.loadingText}
                  >
                    {iconMap.loading}
                  </span>
                ) : (
                  <span
                    className="group play-btn"
                    onClick={togglePlay}
                    title={
                      playing ? locale.clickToPauseText : locale.clickToPlayText
                    }
                  >
                    {playing ? iconMap.pause : iconMap.play}
                  </span>
                )}
              </span>

              <span className="group play-sounds" title={locale.volumeText}>
                {volume === 0 ? (
                  <span className="sounds-icon" onClick={resetVolume}>
                    {iconMap.mute}
                  </span>
                ) : (
                  <span className="sounds-icon" onClick={mute}>
                    {iconMap.volume}
                  </span>
                )}
                <Slider
                  value={volume}
                  onChange={setAudioVolume}
                  className="sound-operation"
                  min={0}
                  max={1}
                  step={0.01}
                />
              </span>
              <span
                title={locale.destroyText}
                className="group destroy-btn"
                onClick={stopPlaying}
              >
                {iconMap.destroy}
              </span>
            </div>
          </section>
        </div>
      )}

      {streamUrl && <audio ref={audioRef} />}
      <RadioDialog
        open={open}
        onClose={closeModal}
        setMetadata={setMetadata}
        title={metadata.title}
      />
    </div>
  )
}

export default RadioPlayer
