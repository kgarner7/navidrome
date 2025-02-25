import { useQueryWithStore } from 'react-admin'

export interface Stat {
  count: number
  id: string
  name: string
}

export const useStat = (
  stat: string,
  from: number,
  to: number,
  page: number,
  perPage: number,
): [Stat[], boolean, number] => {
  const toString = to.toString()

  const { data, loading, total } = useQueryWithStore({
    type: 'getList',
    resource: 'stats',
    payload: {
      pagination: {
        page,
        perPage,
      },
      filter: {
        stat: stat,
        from,
        // WHY ARE YOU LIKE THIS TRUNCATING THE LAST PART
        to: toString + toString[toString.length - 1],
      },
      sort: { field: 'count' },
    },
  })

  return [data ?? [], loading, loading ? -1 : total!]
}
