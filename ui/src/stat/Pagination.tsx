import { Input, InputAdornment, TextField } from '@material-ui/core'

interface PaginationProps {
  setPage: (value: number) => void

  count: number
  page: number
  total: number
}

export const Pagination = ({
  setPage,
  count,
  page,
  total,
}: PaginationProps) => {
  const max = Math.ceil(total / count)

  return (
    <div>
      <Input
        type="number"
        inputProps={{ min: 1, max }}
        value={page}
        onChange={(e) => setPage(Number(e.currentTarget.value))}
        startAdornment={<InputAdornment position="start">Page</InputAdornment>}
        endAdornment={
          <InputAdornment position="end">
            of {max} ({total} total)
          </InputAdornment>
        }
        fullWidth
      />{' '}
    </div>
  )
}
