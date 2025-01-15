import { useQueryWithStore } from 'react-admin'

export const useStat = (
  type: string,
  from: number,
  to: number,
  count: number,
) => {
  const { data, loading } = useQueryWithStore({
    type: 'getMany',
    resource: 'stats',
    payload: {
      type,
      from,
      to,
      start: 0,
      end: count,
    },
  })

  return [data ?? [], loading]
}
