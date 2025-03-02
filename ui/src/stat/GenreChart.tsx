import { useMemo, useState } from 'react'
import { Bar } from 'react-chartjs-2'

import { makeOptions } from './options'
import { useStat } from './useStat'
import { Pagination } from './Pagination'

interface BarChartProps {
  count: number
  from: number
  to: number
}

const GenreChart = ({ count, from, to }: BarChartProps) => {
  const [page, setPage] = useState(1)
  const [genres, loading, total] = useStat('genre', from, to, page, count)

  const [values, labels] = useMemo(() => {
    const labels: string[] = new Array(genres.length)
    const values: number[] = new Array(genres.length)

    for (const [idx, stat] of genres.entries()) {
      labels[idx] = stat.name
      values[idx] = stat.count
    }

    return [values, labels]
  }, [genres])

  return (
    <>
      <Bar
        options={makeOptions(false, `Top genres: (${count} / ${total})`)}
        updateMode="none"
        data={{
          datasets: [{ data: values }],
          labels,
        }}
      />
      {(!loading || total !== -1) && (
        <Pagination setPage={setPage} count={count} page={page} total={total} />
      )}
    </>
  )
}

export default GenreChart
