import { Card, Grid, makeStyles, TextField, useTheme } from '@material-ui/core'
import Annotation from 'chartjs-plugin-annotation'
import {
  Chart,
  CategoryScale,
  Colors,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js'
import { addDays, endOfDay, format, parse } from 'date-fns'
import {
  createRef,
  useCallback,
  useLayoutEffect,
  useMemo,
  useState,
} from 'react'
import { linkToRecord } from 'react-admin'
import { useSelector } from 'react-redux'
import { useHistory, useLocation } from 'react-router-dom'

// @ts-expect-error importing js model in ts
import httpClient from '../dataProvider/httpClient'

import BarChartWithImage from './BarChartWithImage'
import GenreChart from './GenreChart'
import BufferedNumberInput from './BufferedNumberInput'

Chart.register(
  Annotation,
  CategoryScale,
  Colors,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
)

Chart.defaults.plugins.title.font = {
  size: 20,
  ...Chart.defaults.plugins.title.font,
}

const MUI_DATE_FORMAT = 'yyyy-MM-dd'

const useStyles = makeStyles({
  card: {
    padding: 15,
  },
})

const Stat = () => {
  const cardRef = createRef<HTMLDivElement>()
  const [width, setWidth] = useState<number | undefined>()
  const theme = useTheme()

  // @ts-expect-error admin does in fact exist. THis is ra-admin
  const open = useSelector((state) => state.admin.ui.sidebarOpen)
  const classes = useStyles()

  const history = useHistory()
  const { search } = useLocation()
  // @ts-expect-error activity is a react-admin prop
  const refreshData = useSelector((state) => state?.activity?.refresh)

  const [start, end, count] = useMemo(() => {
    const params = new URLSearchParams(search)
    const now = new Date()

    const start = params.has('start')
      ? parse(params.get('start')!, MUI_DATE_FORMAT, now)
      : addDays(now, -7)
    const end = params.has('end')
      ? parse(params.get('end')!, MUI_DATE_FORMAT, now)
      : now
    const count = params.has('count') ? Number(params.get('count')) : 5

    return [start, end, count]
  }, [search])

  const [startTs, startFormat, endTs, endFormat] = useMemo(() => {
    return [
      start.getTime(),
      format(start, MUI_DATE_FORMAT),
      endOfDay(end).getTime(),
      format(end, MUI_DATE_FORMAT),
    ]
  }, [end, start])

  const setParam = useCallback(
    (k: 'start' | 'end' | 'count', val: string) => {
      const search = new URLSearchParams({
        start: startFormat,
        end: endFormat,
        count: count.toString(),
        [k]: val,
      })
      history.replace({ pathname: '/stats', search: search.toString() })
    },
    [count, endFormat, history, startFormat],
  )

  useLayoutEffect(() => {
    Chart.defaults.color = theme.palette.text.primary
    Chart.defaults.font.family = theme.typography.fontFamily
  }, [theme])

  useLayoutEffect(() => {
    const updateSize = () => {
      const sidebar = document.querySelector('.MuiDrawer-root')
      if (sidebar) {
        if (open) {
          setWidth(window.screen.width - 240 - 60)
        } else {
          setWidth(window.screen.width - 55 - 60)
        }
      } else {
        setWidth(window.screen.width - 60)
      }
    }

    window.addEventListener('resize', updateSize)
    updateSize()

    return () => {
      window.removeEventListener('resize', updateSize)
    }
  }, [cardRef, open])

  return (
    <Card className={classes.card} style={{ width }}>
      <Grid container spacing={2}>
        <Grid item xs>
          <TextField
            fullWidth
            variant="filled"
            label="Start date"
            type="date"
            value={startFormat}
            inputProps={{ max: endFormat }}
            onChange={(elem) => setParam('start', elem.currentTarget.value)}
          />
        </Grid>
        <Grid item xs>
          <TextField
            fullWidth
            variant="filled"
            label="End date"
            type="date"
            value={endFormat}
            inputProps={{ min: startFormat }}
            onChange={(elem) => setParam('end', elem.currentTarget.value)}
          />
        </Grid>
        <Grid item xs>
          <BufferedNumberInput
            value={count}
            setValue={(value) => setParam('count', value.toString())}
          />
        </Grid>
      </Grid>

      <BarChartWithImage
        count={count}
        from={startTs}
        to={endTs}
        type="album"
        labelKey="name"
        title="Top albums"
        route={(elem) => linkToRecord('album', elem.id, 'show')}
      />
      <BarChartWithImage
        count={count}
        from={startTs}
        to={endTs}
        type="artist"
        labelKey="name"
        title="Top artists"
        route={(elem) => linkToRecord('artist', elem.id, 'show')}
      />
      <BarChartWithImage
        count={count}
        from={startTs}
        to={endTs}
        type="song"
        labelKey="title"
        title="Top songs"
      />
      <GenreChart count={count} from={startTs} to={endTs} />
    </Card>
  )
}

export default Stat
