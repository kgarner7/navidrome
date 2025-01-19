import ShowChartIcon from '@material-ui/icons/ShowChart'
import ShowChartOutlinedIcon from '@material-ui/icons/ShowChartOutlined'
import ListenList from './ListenList'
import DynamicMenuIcon from '../layout/DynamicMenuIcon'

export default {
  list: ListenList,
  icon: (
    <DynamicMenuIcon
      path="listen"
      icon={ShowChartOutlinedIcon}
      activeIcon={ShowChartIcon}
    />
  ),
}
