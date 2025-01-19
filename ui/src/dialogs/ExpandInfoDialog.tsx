import PropTypes from 'prop-types'
import { useDispatch, useSelector } from 'react-redux'
import { RecordContextProvider, useTranslate } from 'react-admin'
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
} from '@material-ui/core'
import { closeExtendedInfoDialog } from '../actions/dialogs'
import { MouseEvent } from 'react'

const ExpandInfoDialog = ({
  title,
  content,
}: {
  title?: string
  content: JSX.Element
}) => {
  // @ts-expect-error I'm not dealing with typing of state
  const { open, record } = useSelector((state) => state.expandInfoDialog)
  const dispatch = useDispatch()
  const translate = useTranslate()

  const handleClose = (e: MouseEvent<HTMLButtonElement>) => {
    dispatch(closeExtendedInfoDialog())
    e.stopPropagation()
  }

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      aria-labelledby="info-dialog-album"
      fullWidth={true}
      maxWidth={'sm'}
    >
      <DialogTitle id="info-dialog-album">
        {translate(title || 'resources.song.actions.info')}
      </DialogTitle>
      <DialogContent>
        {record && (
          <RecordContextProvider value={record}>
            {content}
          </RecordContextProvider>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
          {translate('ra.action.close')}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

ExpandInfoDialog.propTypes = {
  title: PropTypes.string,
  content: PropTypes.object.isRequired,
}

export default ExpandInfoDialog
