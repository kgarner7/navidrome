import PropTypes from 'prop-types'
import { Link } from 'react-admin'
import { useDispatch } from 'react-redux'
import { closeExtendedInfoDialog } from '../actions/dialogs'

export const AlbumLinkField = ({
  record,
}: {
  record?: { albumId: string; album: string }
  source: string
}) => {
  const dispatch = useDispatch()

  if (!record) return null

  record ? (
    <Link
      to={`/album/${record.albumId}/show`}
      onClick={(e) => {
        e.stopPropagation()
        dispatch(closeExtendedInfoDialog())
      }}
    >
      {record.album}
    </Link>
  ) : null
}

AlbumLinkField.propTypes = {
  sortBy: PropTypes.string,
  sortByOrder: PropTypes.oneOf(['ASC', 'DESC']),
}

AlbumLinkField.defaultProps = {
  addLabel: true,
}
