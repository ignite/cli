package xgenny

// SourceModification describes modified and created files in the source code after a run.
type SourceModification struct {
	modified map[string]struct{}
	created  map[string]struct{}
}

func NewSourceModification() SourceModification {
	return SourceModification{
		make(map[string]struct{}),
		make(map[string]struct{}),
	}
}

// ModifiedFiles returns the modified files of the source modification.
func (sm SourceModification) ModifiedFiles() (modifiedFiles []string) {
	for modified := range sm.modified {
		modifiedFiles = append(modifiedFiles, modified)
	}
	return
}

// CreatedFiles returns the created files of the source modification.
func (sm SourceModification) CreatedFiles() (createdFiles []string) {
	for created := range sm.created {
		createdFiles = append(createdFiles, created)
	}
	return
}

// AppendModifiedFiles appends modified files in the source modification that are not already documented.
func (sm *SourceModification) AppendModifiedFiles(modifiedFiles ...string) {
	for _, modifiedFile := range modifiedFiles {
		_, alreadyModified := sm.modified[modifiedFile]
		_, alreadyCreated := sm.created[modifiedFile]
		if !alreadyModified && !alreadyCreated {
			sm.modified[modifiedFile] = struct{}{}
		}
	}
}

// AppendCreatedFiles appends a created files in the source modification that are not already documented.
func (sm *SourceModification) AppendCreatedFiles(createdFiles ...string) {
	for _, createdFile := range createdFiles {
		_, alreadyModified := sm.modified[createdFile]
		_, alreadyCreated := sm.created[createdFile]
		if !alreadyModified && !alreadyCreated {
			sm.created[createdFile] = struct{}{}
		}
	}
}

// Merge merges new source modification to an existing one.
func (sm *SourceModification) Merge(newSm SourceModification) {
	sm.AppendModifiedFiles(newSm.ModifiedFiles()...)
	sm.AppendCreatedFiles(newSm.CreatedFiles()...)
}
