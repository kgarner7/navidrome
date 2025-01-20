import PropTypes from 'prop-types'
import { useRecordContext } from 'react-admin'
import { formatDuration } from '../utils/formatters'

export const DurationField = ({ source, ...rest }: { source: string }) => {
  const record = useRecordContext(rest)
  if (!record) return null

  try {
    return <span>{formatDuration(record[source])}</span>
  } catch (e) {
    // eslint-disable-next-line no-console
    console.log('Error in DurationField! Record:', record)
    return <span>00:00</span>
  }
}

DurationField.propTypes = {
  label: PropTypes.string,
  record: PropTypes.object,
  source: PropTypes.string.isRequired,
}

DurationField.defaultProps = {
  addLabel: true,
}
