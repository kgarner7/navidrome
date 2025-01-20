import PropTypes from 'prop-types'
import { Link } from 'react-admin'

export const AlbumLinkField = ({
  record,
}: {
  record?: { albumId: string; album: string }
  source: string
}) =>
  record ? (
    <Link
      to={`/album/${record.albumId}/show`}
      onClick={(e) => e.stopPropagation()}
    >
      {record.album}
    </Link>
  ) : null

AlbumLinkField.propTypes = {
  sortBy: PropTypes.string,
  sortByOrder: PropTypes.oneOf(['ASC', 'DESC']),
}

AlbumLinkField.defaultProps = {
  addLabel: true,
}
