import React from 'react'
import { Route } from 'react-router-dom'
import Personal from './personal/Personal'
import Stat from './stat'

const routes = [
  <Route exact path="/personal" render={() => <Personal />} key={'personal'} />,
  <Route exact path="/stats" render={() => <Stat />} key={'stat'} />,
]

export default routes
