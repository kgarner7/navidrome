import { makeStyles } from '@material-ui/core'
import { fromUnixTime } from 'date-fns'
import { useMemo } from 'react'
import {
  Datagrid,
  Filter,
  FilterProps,
  FunctionField,
  Identifier,
  List,
  ListProps,
  Pagination,
  SearchInput,
  TextField,
} from 'react-admin'
import { useDispatch } from 'react-redux'
// @ts-expect-error importing untyped js file in ts
import { setTrack } from '../actions'
import {
  ArtistLinkField,
  SongInfo,
  useSelectedFields,
  // @ts-expect-error importing untyped js file in ts
} from '../common'
import { DurationField } from '../common/DurationField'
import ExpandInfoDialog from '../dialogs/ExpandInfoDialog'
import { AlbumLinkField } from '../song/AlbumLinkField'
import { ListenContextMenu } from './ListenContextMenu'
import { ListenListActions } from './ListenListActions'
import { listenToTrack } from './listenToTrack'

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
    },
  },
  contextMenu: {
    visibility: 'hidden',
  },
})

const ListenFilter = (props: Omit<FilterProps, 'children'>) => {
  return (
    <Filter {...props} variant={'outlined'}>
      <SearchInput source="title" alwaysOn />
    </Filter>
  )
}

const ListenList = (props: ListProps) => {
  const classes = useStyles()
  const dispatch = useDispatch()

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
      album: <AlbumLinkField source="album" sortByOrder={'ASC'} />,
      albumArtist: <ArtistLinkField source="albumArtist" />,
      duration: <DurationField source="duration" />,
    }
  }, [])

  const columns = useSelectedFields({
    resource: 'listen',
    columns: toggleableFields,
    defaultOff: ['albumArtist', 'artist'],
  })

  return (
    <>
      <List
        {...props}
        sort={{ field: 'submission_time', order: 'DESC' }}
        exporter={false}
        actions={<ListenListActions />}
        filters={<ListenFilter />}
        bulkActionButtons={false}
        pagination={<Pagination rowsPerPageOptions={[25, 50, 100]} />}
        perPage={25}
      >
        <Datagrid rowClick={handleRowClick} classes={{ row: classes.row }}>
          <FunctionField
            source="submission_time"
            render={(r) => r && r.submissionTime}
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
