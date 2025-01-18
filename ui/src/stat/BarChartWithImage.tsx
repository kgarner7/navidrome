import type { AnnotationOptions } from 'chartjs-plugin-annotation'
import type { Chart, Plugin } from 'chart.js'
import { useCallback, useMemo, useRef } from 'react'
import { Loading, useRedirect } from 'react-admin'
import { Bar } from 'react-chartjs-2'

// @ts-expect-error Importing a JS module with no typing. I do not want to fix these types
import subsonic from '../subsonic'
import { makeOptions } from './options'
import { useStat } from './useStat'

interface Stat {
  [k: string]: string | number
  count: number
}

interface BarChartProps {
  count: number
  from: number
  title: string
  to: number
  type: 'album' | 'artist' | 'song'
  route?: (element: Stat) => string
}

const augmentStat = (type: BarChartProps['type'], stat: object) => {
  switch (type) {
    case 'song':
      return { album: 1, ...stat }
    case 'album':
      return { artist: 1, ...stat }
    case 'artist':
      return stat
  }
}

const BarChartWithImage = ({
  count,
  from,
  title,
  to,
  type,
  route,
}: BarChartProps) => {
  const heightRef = useRef(0)
  const barRef = useRef<Chart<'bar', number[], string>>()
  const redirect = useRedirect()

  const [data, loading] = useStat(type, from, to, count)

  const plugin = useCallback(() => {
    const data: Plugin = {
      id: 'update-height',
      beforeDraw: (chart) => {
        const newHeight = chart
          .getDatasetMeta(0)
          .data[0].getProps(['height'], true).height

        if (heightRef.current !== 0 && newHeight !== heightRef.current) {
          chart.update()
          heightRef.current = newHeight
        }
      },
    }

    return data
  }, [])

  const [annotations, values, labels] = useMemo(() => {
    const annotations: AnnotationOptions<'label'>[] = new Array(data.length)
    const labels: string[] = new Array(data.length)
    const values: number[] = new Array(data.length)

    for (const [idx, stat] of data.entries()) {
      annotations[idx] = {
        type: 'label',
        content: (ctx) => {
          const size = Math.round(
            ctx.chart.getDatasetMeta(0).data[0].getProps(['height'], true)
              .height as number,
          )
          const img = new Image(size, size)
          img.src = subsonic.getCoverArtUrl(augmentStat(type, stat), 300)
          return img
        },
        position: { x: 'start' },
        xValue: 0,
        yValue: stat.name,
      }
      labels[idx] = stat.name as string
      values[idx] = stat.count
    }

    return [annotations, values, labels]
  }, [data, type])

  const options = useMemo(() => {
    const ops = makeOptions(route !== undefined, title, annotations, {
      // @ts-expect-error custom plugins are exported
      custom: plugin,
    })!
    if (route) {
      ops.onClick = (_, element) => {
        if (element.length > 0) {
          redirect(route(data[element[0].index]))
        }
      }
    }

    return ops
  }, [annotations, data, plugin, redirect, route, title])

  if (loading) {
    return <Loading />
  }

  return (
    <Bar
      ref={barRef}
      options={options}
      updateMode="show"
      data={{
        datasets: [{ data: values }],
        labels,
      }}
    />
  )
}

export default BarChartWithImage
