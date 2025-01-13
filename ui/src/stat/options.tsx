import { PluginOptionsByType } from 'chart.js'
import { AnnotationOptions } from 'chartjs-plugin-annotation'
import { ComponentProps } from 'react'
import { Bar } from 'react-chartjs-2'

export const makeOptions = (
  showHover: boolean,
  title: string,
  annotations?: AnnotationOptions<'label'>[],
  plugins?: PluginOptionsByType<'bar'>,
): ComponentProps<typeof Bar>['options'] => {
  return {
    indexAxis: 'y',
    onHover:
      showHover === false
        ? undefined
        : (event, chartElement) => {
            event.native.target.style.cursor = chartElement[0]
              ? 'pointer'
              : 'default'
          },
    plugins: {
      annotation: { annotations },
      legend: { display: false },
      title: { display: true, text: title },
      ...plugins,
    },
    responsive: true,
    scales: {
      y: {
        ticks: {
          callback(_, index) {
            const text = this.getLabelForValue(index)

            if (text.length > 30) {
              return text.substring(0, 30) + '...'
            } else {
              return text
            }
          },
        },
      },
    },
  }
}
