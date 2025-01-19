import config from '../config'
import { useAlbumsPerPage } from './useAlbumsPerPage'

export const useGetHandleArtistClick = (width: string) => {
  const [perPage] = useAlbumsPerPage(width)
  return (id: string) => {
    return config.devShowArtistPage && id !== config.variousArtistsId
      ? `/artist/${id}/show`
      : `/album?filter={"artist_id":"${id}"}&order=ASC&sort=max_year&displayedFilters={"compilation":true}&perPage=${perPage}`
  }
}
