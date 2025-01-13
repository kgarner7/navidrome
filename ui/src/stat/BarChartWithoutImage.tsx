import { useMemo } from 'react'
import { Bar } from 'react-chartjs-2'
import { defaultOptions, makeOptions } from './options'

interface Stat {
  [k: string]: string | number
  count: number
}

interface BarChartProps {
  data: Stat[]
  labelKey: string
  title: string
}

const BarChartWithoutImage = ({ data, labelKey, title }: BarChartProps) => {
  const [values, labels] = useMemo(() => {
    const labels: string[] = new Array(data.length)
    const values: number[] = new Array(data.length)

    for (const [idx, stat] of data.entries()) {
      labels[idx] = stat[labelKey] as string
      values[idx] = stat.count
    }

    return [values, labels]
  }, [data, labelKey])

  return (
    <Bar
      options={makeOptions(false, title)}
      updateMode="show"
      data={{
        datasets: [{ data: values }],
        labels,
      }}
    />
  )
}

export default BarChartWithoutImage
