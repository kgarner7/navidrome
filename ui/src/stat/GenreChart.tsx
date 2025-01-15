import { useMemo } from 'react'
import { Loading } from 'react-admin'
import { Bar } from 'react-chartjs-2'

import { makeOptions } from './options'
import { useStat } from './useStat'

interface BarChartProps {
  count: number
  from: number
  to: number
}

const GenreChart = ({ count, from, to }: BarChartProps) => {
  const [genres, loading] = useStat('genre', from, to, count)

  const [values, labels] = useMemo(() => {
    const labels: string[] = new Array(genres.length)
    const values: number[] = new Array(genres.length)

    for (const [idx, stat] of genres.entries()) {
      labels[idx] = stat.name
      values[idx] = stat.count
    }

    return [values, labels]
  }, [genres])

  if (loading) {
    return <Loading />
  }

  return (
    <Bar
      options={makeOptions(false, 'Top genres')}
      updateMode="show"
      data={{
        datasets: [{ data: values }],
        labels,
      }}
    />
  )
}

export default GenreChart
