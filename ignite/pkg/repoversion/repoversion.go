package repoversion

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Version struct {
	Tag  string
	Hash string
}

func Determine(path string) (v Version, err error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return Version{}, err
	}

	tags, err := repo.Tags()
	if err != nil {
		return Version{}, err
	}

	head, err := repo.Head()
	if err != nil {
		return Version{}, err
	}

	cIter, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		return Version{}, err
	}

	var (
		tag             string
		tagHash         string
		commitIndex     int
		taggerTimestamp int64
	)

	idMap := make(map[string]int)

	err = cIter.ForEach(func(c *object.Commit) error {
		idMap[c.Hash.String()] = commitIndex
		commitIndex++

		return nil
	})

	if err != nil {
		return Version{}, err
	}

	err = tags.ForEach(func(t *plumbing.Reference) error {
		obj, err := repo.TagObject(t.Hash())

		if err == nil {

			_, exists := idMap[obj.Target.String()]

			if !exists {
				return nil
			}

			if taggerTimestamp < obj.Tagger.When.Unix() {
				taggerTimestamp = obj.Tagger.When.Unix()

				tag = strings.TrimPrefix(obj.Name, "v")
				tagHash = obj.Target.String()
			}
		} else {
			_, exists := idMap[t.Hash().String()]

			if !exists {
				return nil
			}

			commit, err := repo.CommitObject(t.Hash())
			if err != nil {
				return err
			}

			if taggerTimestamp < commit.Committer.When.Unix() {
				taggerTimestamp = commit.Committer.When.Unix()

				tag = strings.TrimPrefix(t.Name().Short(), "v")
				tagHash = t.Hash().String()
			}
		}

		return nil
	})

	if err != nil {
		return Version{}, err
	}

	const subHashLen int = 8

	tagHashIndex := idMap[tagHash]
	headHashText := head.Hash().String()
	subHeadHash := headHashText

	if len(headHashText) > subHashLen {
		subHeadHash = subHeadHash[:subHashLen]
	}

	v.Tag = tag
	v.Hash = headHashText

	if tagHashIndex > 0 {
		v.Tag = fmt.Sprintf("%s-%s", tag, subHeadHash)
	}

	return v, nil
}
