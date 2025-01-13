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
  useState,
} from 'react'
import { linkToRecord, Loading } from 'react-admin'
import { useSelector } from 'react-redux'

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
  const [start, setStart] = useState(
    format(addDays(new Date(), -7), MUI_DATE_FORMAT),
  )
  const [end, setEnd] = useState(format(new Date(), MUI_DATE_FORMAT))
  // @ts-expect-error admin does in fact exist. THis is ra-admin
  const open = useSelector((state) => state.admin.ui.sidebarOpen)
  const [count, setCount] = useState(5)
  const classes = useStyles()

  const fetchData = useCallback(async () => {
    const startMs = parse(start, MUI_DATE_FORMAT, new Date()).getTime()
    const endMs = endOfDay(parse(end, MUI_DATE_FORMAT, new Date())).getTime()

    const toFetch = BASE_TO_FETCH.map((item) =>
      httpClient(
        `/api/stats/${item}?from=${startMs}&to=${endMs}&_start=0&_end=${count}`,
      ).then((resp: { json: unknown }) => resp.json),
    )
    const data = await Promise.all(toFetch)
    setStats(data)
  }, [count, end, start])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  useLayoutEffect(() => {
    Chart.defaults.color = theme.palette.text.primary
    Chart.defaults.font.family = theme.typography.fontFamily
  }, [theme])

  useLayoutEffect(() => {
    const updateSize = () => {
      const sidebar = document.querySelector('.MuiDrawer-root')
      if (sidebar) {
        if (open) {
          setWidth(window.screen.width - 240 - 30)
        } else {
          setWidth(window.screen.width - 55 - 30)
        }
      } else {
        setWidth(window.screen.width - 30)
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
            value={start}
            inputProps={{ max: end }}
            onChange={(elem) => setStart(elem.currentTarget.value)}
          />
        </Grid>
        <Grid item xs>
          <TextField
            fullWidth
            variant="filled"
            label="End date"
            type="date"
            value={end}
            inputProps={{ min: start }}
            onChange={(elem) => setEnd(elem.currentTarget.value)}
          />
        </Grid>
        <Grid item xs>
          <BufferedNumberInput value={count} setValue={setCount} />
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
