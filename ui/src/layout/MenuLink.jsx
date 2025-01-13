import React, { forwardRef } from 'react'
import { MenuItemLink, useTranslate } from 'react-admin'
import { makeStyles } from '@material-ui/core'
import TuneIcon from '@material-ui/icons/Tune'

const useStyles = makeStyles((theme) => ({
  menuItem: {
    color: theme.palette.text.secondary,
  },
}))

const MenuLink = forwardRef(
  ({ onClick, sidebarIsOpen, dense, icon, link, text }, ref) => {
    const translate = useTranslate()
    const classes = useStyles()
    return (
      <MenuItemLink
        ref={ref}
        to={link}
        primaryText={translate(text)}
        leftIcon={icon}
        onClick={onClick}
        className={classes.menuItem}
        sidebarIsOpen={sidebarIsOpen}
        dense={dense}
      />
    )
  },
)

MenuLink.displayName = 'MenuLink'

export default MenuLink
