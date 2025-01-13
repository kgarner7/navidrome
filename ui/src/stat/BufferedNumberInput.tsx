import { makeStyles, TextField } from '@material-ui/core'
import { useState } from 'react'

interface CountInputProps {
  setValue: (value: number) => void
  value: number
}

const useStyles = makeStyles({
  input: {
    '& input[type=number]': {
      '-moz-appearance': 'textfield',
    },
    '& input[type=number]::-webkit-outer-spin-button': {
      '-webkit-appearance': 'none',
      margin: 0,
    },
    '& input[type=number]::-webkit-inner-spin-button': {
      '-webkit-appearance': 'none',
      margin: 0,
    },
  },
})

const BufferedNumberInput = ({ setValue, value }: CountInputProps) => {
  const classes = useStyles()
  const [count, setCount] = useState(value)

  return (
    <TextField
      className={classes.input}
      fullWidth
      variant="filled"
      label="Number to fetch"
      type="number"
      value={count}
      inputProps={{ min: 5, max: 15 }}
      onChange={(elem) => setCount(Number(elem.currentTarget.value))}
      onBlur={(elem) => setValue(Number(elem.currentTarget.value))}
    />
  )
}

export default BufferedNumberInput
