import { sanitizeListRestProps, TopToolbar } from 'react-admin'
import { useMediaQuery } from '@material-ui/core'
// @ts-expect-error importing js in tsx
import { ToggleFieldsMenu } from '../common'

export const ListenListActions = ({
  className,
  ...rest
}: {
  className?: string
}) => {
  // @ts-expect-error i'm not typing theme
  const isNotSmall = useMediaQuery((theme) => theme.breakpoints.up('sm'))
  return (
    <TopToolbar className={className} {...sanitizeListRestProps(rest)}>
      {isNotSmall && <ToggleFieldsMenu resource="listen" />}
    </TopToolbar>
  )
}
