import { MouseEvent, useState } from 'react'
import PropTypes from 'prop-types'
import { useDispatch } from 'react-redux'
import { useTranslate } from 'react-admin'
import { IconButton, Menu, MenuItem } from '@material-ui/core'
import { makeStyles } from '@material-ui/core/styles'
import MoreVertIcon from '@material-ui/icons/MoreVert'
import clsx from 'clsx'
import {
  playNext,
  addTracks,
  setTrack,
  openAddToPlaylist,
  openExtendedInfoDialog,
  openDownloadMenu,
  DOWNLOAD_MENU_SONG,
  openShareMenu,
  // @ts-expect-error importing js in ts
} from '../actions'
import config from '../config'
// @ts-expect-error importing js in ts
import { formatBytes } from '../utils'
import { listenToTrack, SongRecord } from './listenToTrack'

const useStyles = makeStyles({
  noWrap: {
    whiteSpace: 'nowrap',
  },
})

export const ListenContextMenu = ({
  record,
  onAddToPlaylist,
  className,
}: {
  record?: SongRecord
  onAddToPlaylist: (id: string) => void
  className: string
}) => {
  const classes = useStyles()
  const dispatch = useDispatch()
  const translate = useTranslate()
  const [anchorEl, setAnchorEl] = useState<Element | null>(null)
  const options = {
    playNow: {
      enabled: true,
      label: translate('resources.song.actions.playNow'),
      action: (record: SongRecord) => dispatch(setTrack(record)),
    },
    playNext: {
      enabled: true,
      label: translate('resources.song.actions.playNext'),
      action: (record: SongRecord) =>
        dispatch(playNext({ [record.id]: record })),
    },
    addToQueue: {
      enabled: true,
      label: translate('resources.song.actions.addToQueue'),
      action: (record: SongRecord) =>
        dispatch(addTracks({ [record.id]: record })),
    },
    addToPlaylist: {
      enabled: true,
      label: translate('resources.song.actions.addToPlaylist'),
      action: (record: SongRecord) =>
        dispatch(
          openAddToPlaylist({
            selectedIds: [record.id],
            onSuccess: (id: string) => onAddToPlaylist(id),
          }),
        ),
    },
    share: {
      enabled: config.enableSharing,
      label: translate('ra.action.share'),
      action: (record: SongRecord) =>
        dispatch(openShareMenu([record.id], 'song', record.title)),
    },
    download: {
      enabled: config.enableDownloads,
      label: `${translate('ra.action.download')} (${formatBytes(record!.size)})`,
      action: (record: SongRecord) =>
        dispatch(openDownloadMenu(record, DOWNLOAD_MENU_SONG)),
    },
    info: {
      enabled: true,
      label: translate('resources.song.actions.info'),
      action: (record: SongRecord) => dispatch(openExtendedInfoDialog(record)),
    },
  }

  const handleClick = (e: MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(e.currentTarget)
    e.stopPropagation()
  }

  const handleClose = (e: MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(null)
    e.stopPropagation()
  }

  const handleItemClick = (e: MouseEvent<HTMLLIElement>) => {
    e.preventDefault()
    setAnchorEl(null)
    const key = e.currentTarget.getAttribute('value') as keyof typeof options
    options[key].action(listenToTrack(record!))
    e.stopPropagation()
  }

  const open = Boolean(anchorEl)

  return (
    <span className={clsx(classes.noWrap, className)}>
      <IconButton onClick={handleClick} size={'small'}>
        <MoreVertIcon fontSize={'small'} />
      </IconButton>
      <Menu
        id={'menu' + record!.id}
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
      >
        {Object.keys(options).map(
          (key) =>
            options[key as keyof typeof options].enabled && (
              <MenuItem value={key} key={key} onClick={handleItemClick}>
                {options[key as keyof typeof options].label}
              </MenuItem>
            ),
        )}
      </Menu>
    </span>
  )
}

ListenContextMenu.propTypes = {
  resource: PropTypes.string.isRequired,
  record: PropTypes.object.isRequired,
  onAddToPlaylist: PropTypes.func,
}

ListenContextMenu.defaultProps = {
  onAddToPlaylist: () => {},
  record: {},
  resource: 'song',
  addLabel: true,
}
