package persistence

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	. "github.com/Masterminds/squirrel"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/pocketbase/dbx"
)

type tagRepository struct {
	*baseTagRepository
}

func NewTagRepository(ctx context.Context, db dbx.Builder) model.TagRepository {
	return &tagRepository{
		baseTagRepository: newBaseTagRepository(ctx, db, nil), // nil = no filter, works with all tags
	}
}

func (r *tagRepository) Add(libraryID int, tags ...model.Tag) error {
	conf := model.TagMappings()

	for chunk := range slices.Chunk(tags, 200) {
		sq := Insert(r.tableName).Columns("id", "tag_name", "tag_value", "song").
			Suffix("on conflict (id) do nothing")

		// Create library_tag entries for library filtering
		libSq := Insert("library_tag").Columns("tag_id", "library_id", "album_count", "media_file_count").
			Suffix("on conflict (tag_id, library_id) do nothing")

		hasValues := false
		for _, t := range chunk {
			c := conf[t.TagName]

			if !c.Album && !c.Song {
				continue
			}

			hasValues = true
			song := conf[t.TagName].Song
			sq = sq.Values(t.ID, t.TagName, t.TagValue, song)
			libSq = libSq.Values(t.ID, libraryID, 0, 0)
		}

		if hasValues {
			_, err := r.executeSQL(sq)
			if err != nil {
				return err
			}

			_, err = r.executeSQL(libSq)
			if err != nil {
				return fmt.Errorf("adding library_tag entries: %w", err)
			}
		}
	}
	return nil
}

// UpdateCounts updates the library_tag table with per-library statistics.
// Only genres are being updated for now.
func (r *tagRepository) UpdateCounts() error {
	template := `
INSERT INTO library_tag (tag_id, library_id, %[1]s_count)
SELECT jt.value as tag_id, %[1]s.library_id, count(distinct %[1]s.id) as %[1]s_count
FROM %[1]s
JOIN json_tree(%[1]s.tags, '$.genre') as jt ON jt.atom IS NOT NULL AND jt.key = 'id'
JOIN tag ON tag.id = jt.value
GROUP BY jt.value, %[1]s.library_id
ON CONFLICT (tag_id, library_id) 
DO UPDATE SET %[1]s_count = excluded.%[1]s_count;
`

	for _, table := range []string{"album", "media_file"} {
		start := time.Now()
		query := Expr(fmt.Sprintf(template, table))
		c, err := r.executeSQL(query)
		log.Debug(r.ctx, "Updated library tag counts", "table", table, "elapsed", time.Since(start), "updated", c)
		if err != nil {
			return fmt.Errorf("updating %s library tag counts: %w", table, err)
		}
	}
	return nil
}

func (r *tagRepository) purgeUnused() error {
	del := Delete(r.tableName).Where(`	
	id not in (select jt.value
	from media_file left join json_tree(media_file.tags, '$') as jt
	where atom is not null
	  and key = 'id'
	UNION 
	select jt.value
	from media_file left join json_tree(media_file.tags, '$') as jt
	where atom is not null
	  and key = 'id')
`)
	c, err := r.executeSQL(del)

	conf := model.TagMappings()
	allowedTags := []model.TagName{}

	for tag, tagConf := range conf {
		if tagConf.Album || tagConf.Song {
			allowedTags = append(allowedTags, tag)
		}
	}

	if err != nil {
		err = fmt.Errorf("error purging unused tags: %w", err)
	}

	confDel := Delete(r.tableName).Where(NotEq{"tag_name": allowedTags})
	c2, err2 := r.executeSQL(confDel)

	if err2 != nil {
		err2 = fmt.Errorf("error removing non-configured tags: %w", err)
	}

	combined := errors.Join(err, err2)
	if combined != nil {
		return combined
	}

	totalCount := c + c2

	if totalCount > 0 {
		log.Debug(r.ctx, "Purged unused tags", "totalDeleted", totalCount)
	}
	return err
}

var _ model.ResourceRepository = &tagRepository{}
