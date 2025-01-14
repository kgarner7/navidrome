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
  useEffect,
  useLayoutEffect,
  useMemo,
  useState,
} from 'react'
import { linkToRecord, Loading } from 'react-admin'
import { useSelector } from 'react-redux'
import { useHistory, useLocation } from 'react-router-dom'

// @ts-expect-error importing js model in ts
import httpClient from '../dataProvider/httpClient'

import BarChartWithImage from './BarChartWithImage'
import BarChartWithoutImage from './BarChartWithoutImage'
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

const BASE_TO_FETCH = ['album', 'artist', 'genre', 'song']
const MUI_DATE_FORMAT = 'yyyy-MM-dd'

const useStyles = makeStyles({
  card: {
    padding: 15,
  },
})

const Stat = () => {
  const cardRef = createRef<HTMLDivElement>()
  const [stats, setStats] = useState([])
  const [width, setWidth] = useState<number | undefined>()
  const theme = useTheme()

  // @ts-expect-error admin does in fact exist. THis is ra-admin
  const open = useSelector((state) => state.admin.ui.sidebarOpen)
  const classes = useStyles()

  const history = useHistory()
  const { search } = useLocation()
  // @ts-expect-error activity is a react-admin prop
  const refreshData = useSelector((state) => state?.activity?.refresh)

  const state = useMemo(() => {
    const params = new URLSearchParams(search)
    const start =
      params.get('start') || format(addDays(new Date(), -7), MUI_DATE_FORMAT)
    const end = params.get('end') || format(new Date(), MUI_DATE_FORMAT)
    const count = params.has('count') ? params.get('count')! : '5'

    return { start, end, count }
  }, [search])

  const setParam = useCallback(
    (k: keyof typeof state, val: string) => {
      const search = new URLSearchParams({
        ...state,
        [k]: val,
      })
      history.replace({ pathname: '/stats', search: search.toString() })
    },
    [history, state],
  )

  const fetchData = useCallback(async () => {
    const startMs = parse(state.start, MUI_DATE_FORMAT, new Date()).getTime()
    const endMs = endOfDay(
      parse(state.end, MUI_DATE_FORMAT, new Date()),
    ).getTime()

    const toFetch = BASE_TO_FETCH.map((item) =>
      httpClient(
        `/api/stats/${item}?from=${startMs}&to=${endMs}&_start=0&_end=${state.count}`,
      ).then((resp: { json: unknown }) => resp.json),
    )
    const data = await Promise.all(toFetch)
    setStats(data)
  }, [state.count, state.end, state.start])

  useEffect(() => {
    fetchData()
  }, [fetchData, refreshData])

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

  if (stats.length === 0) {
    return <Loading />
  }

  return (
    <Card className={classes.card} style={{ width }}>
      <Grid container spacing={2}>
        <Grid item xs>
          <TextField
            fullWidth
            variant="filled"
            label="Start date"
            type="date"
            value={state.start}
            inputProps={{ max: state.end }}
            onChange={(elem) => setParam('start', elem.currentTarget.value)}
          />
        </Grid>
        <Grid item xs>
          <TextField
            fullWidth
            variant="filled"
            label="End date"
            type="date"
            value={state.end}
            inputProps={{ min: state.start }}
            onChange={(elem) => setParam('end', elem.currentTarget.value)}
          />
        </Grid>
        <Grid item xs>
          <BufferedNumberInput
            value={Number(state.count)}
            setValue={(value) => setParam('count', value.toString())}
          />
        </Grid>
      </Grid>

      <BarChartWithImage
        data={stats[0]}
        labelKey="name"
        title="Top albums"
        route={(elem) => linkToRecord('album', elem.id, 'show')}
      />
      <BarChartWithImage
        data={stats[1]}
        labelKey="name"
        title="Top artists"
        route={(elem) => linkToRecord('artist', elem.id, 'show')}
      />
      <BarChartWithImage data={stats[3]} labelKey="title" title="Top songs" />
      <BarChartWithoutImage
        data={stats[2]}
        labelKey="name"
        title="Top genres"
      />
    </Card>
  )
}

export default Stat
