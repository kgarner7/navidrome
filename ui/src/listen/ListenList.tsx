import { makeStyles, useMediaQuery } from '@material-ui/core'
import { useMemo } from 'react'
import {
  Datagrid,
  FunctionField,
  Identifier,
  List,
  ListProps,
  TextField,
} from 'react-admin'
import { useDispatch } from 'react-redux'
// @ts-expect-error importing untyped js file in ts
import { setTrack } from '../actions'
import {
  ArtistLinkField,
  DurationField,
  SongInfo,
  useSelectedFields,
  // @ts-expect-error importing untyped js file in ts
} from '../common'
import ExpandInfoDialog from '../dialogs/ExpandInfoDialog'
// @ts-expect-error JS. Not porting to ts
import { AlbumLinkField } from '../song/AlbumLinkField'
import { ListenListActions } from './ListenListActions'
import { fromUnixTime } from 'date-fns'
import { listenToTrack } from './listenToTrack'
import { ListenContextMenu } from './ListenContextMenu'

const useStyles = makeStyles({
  contextHeader: {
    marginLeft: '3px',
    marginTop: '-2px',
    verticalAlign: 'text-top',
  },
  row: {
    '&:hover': {
      '& $contextMenu': {
        visibility: 'visible',
      },
      '& $ratingField': {
        visibility: 'visible',
      },
    },
  },
  contextMenu: {
    visibility: 'hidden',
  },
})

const ListenList = (props: ListProps) => {
  const classes = useStyles()
  const dispatch = useDispatch()
  // @ts-expect-error i'm not typing theme
  const isDesktop = useMediaQuery((theme) => theme.breakpoints.up('md'))

  const handleRowClick = (
    _id: Identifier,
    _basePath: string,
    record: object,
  ) => {
    dispatch(setTrack(listenToTrack(record)))
    return ''
  }

  const toggleableFields = useMemo(() => {
    return {
      album: isDesktop && <AlbumLinkField source="album" sortByOrder={'ASC'} />,
      albumArtist: <ArtistLinkField source="albumArtist" />,
      duration: <DurationField source="duration" />,
    }
  }, [isDesktop])

  const columns = useSelectedFields({
    resource: 'listen',
    columns: toggleableFields,
    defaultOff: ['albumArtist', 'artist'],
  })

  return (
    <>
      <List
        {...props}
        sort={{ field: 'listened_at', order: 'DESC' }}
        exporter={false}
        actions={<ListenListActions />}
        bulkActionButtons={false}
      >
        <Datagrid rowClick={handleRowClick} classes={{ row: classes.row }}>
          <FunctionField
            source="submission_time"
            render={(r) => fromUnixTime(r!.submissionTime).toLocaleString()}
          />
          <TextField source="title" />

          {columns}
          <ListenContextMenu className={classes.contextMenu} />
        </Datagrid>
      </List>
      <ExpandInfoDialog content={<SongInfo />} />
    </>
  )
}

export default ListenList
