export interface SongRecord {
  id: string
  size: number
  title: string
}

export const listenToTrack = (record: object): SongRecord => {
  // @ts-expect-error yes, file id will exist
  return { ...record, id: record.fileId }
}
