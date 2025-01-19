import PropTypes from 'prop-types'
import { useLocation } from 'react-router-dom'
import { createElement } from 'react'

const DynamicMenuIcon = ({
  icon,
  activeIcon,
  path,
}: {
  icon: unknown
  activeIcon?: unknown
  path: string
}) => {
  const location = useLocation()

  if (!activeIcon) {
    return createElement(icon as '', { 'data-testid': 'icon' })
  }

  return location.pathname.startsWith('/' + path)
    ? createElement(activeIcon as 'input', { 'data-testid': 'activeIcon' })
    : createElement(icon as 'input', { 'data-testid': 'icon' })
}

DynamicMenuIcon.propTypes = {
  path: PropTypes.string.isRequired,
  icon: PropTypes.object.isRequired,
  activeIcon: PropTypes.object,
}

export default DynamicMenuIcon
