package gdiff

import "testing"

func Test_verifyNoDiffsForIdenticalStrings(t *testing.T) {
	verify := func(st SequenceType) {
		diff := MyersDiffer().Diff("a man, a plan, a canal: panama", "a man, a plan, a canal: panama", st)
		if len(diff.Edits()) > 0 {
			t.Error("expected no edits for identical strings")
		}
	}
	verify(CHAR_SPLIT)
	verify(WORD_SPLIT)
	verify(LINE_SPLIT)
}

func Test_verifySingleWordDiff(t *testing.T) {
	verify := func(st SequenceType) {
		diff := MyersDiffer().Diff("a man, a plan, a canal: panama", "a man, a plan, my canal: panama", st)
		edits := diff.Edits()
		if len(edits) != 2 {
			t.Error("expected one diff (two edits)")
		}
		if edits[0].Type != DELETE || edits[0].Start != 4 {
			t.Error("expected first edit to be a deletion of word 5")
		}
		if edits[1].Type != INSERT || edits[1].Start != 4 {
			t.Error("expected second edit to be an insertion of word 5")
		}
	}
	verify(WORD_SPLIT)
}
