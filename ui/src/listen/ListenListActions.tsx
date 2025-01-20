import { sanitizeListRestProps, TopToolbar } from 'react-admin'
// @ts-expect-error importing js in tsx
import { ToggleFieldsMenu } from '../common'

export const ListenListActions = ({
  className,
  ...rest
}: {
  className?: string
}) => {
  return (
    <TopToolbar className={className} {...sanitizeListRestProps(rest)}>
      <ToggleFieldsMenu resource="listen" />
    </TopToolbar>
  )
}
