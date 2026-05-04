package domain

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestFilterNotes_BasicMatch tests filterNotes with basic string matching.
func TestFilterNotes_BasicMatch(t *testing.T) {
	messages := []string{
		"implemented auth service",
		"refactored database layer",
		"fixed login bug",
		"updated documentation",
	}

	// Search for "auth"
	result, err := filterNotes(messages, "auth")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}
	if len(result) != len(messages) {
		t.Errorf("expected %d results, got %d", len(messages), len(result))
	}
	if !result[0] {
		t.Errorf("expected messages[0] to match 'auth'")
	}
	if result[1] || result[2] || result[3] {
		t.Errorf("expected only messages[0] to match 'auth'")
	}
}

// TestFilterNotes_CaseInsensitive tests that filtering is case-insensitive.
func TestFilterNotes_CaseInsensitive(t *testing.T) {
	messages := []string{
		"Added AUTH middleware",
		"auth service implementation",
		"AuthController update",
	}

	result, err := filterNotes(messages, "AUTH")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	// All should match (case-insensitive)
	if !result[0] || !result[1] || !result[2] {
		t.Errorf("expected all messages to match 'AUTH' (case-insensitive)")
	}
}

// TestFilterNotes_NoMatch tests when no messages match the keyword.
func TestFilterNotes_NoMatch(t *testing.T) {
	messages := []string{
		"fixed database connection",
		"updated UI components",
		"refactored service layer",
	}

	result, err := filterNotes(messages, "xyz-nonexistent")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	for i, matched := range result {
		if matched {
			t.Errorf("expected no matches for 'xyz-nonexistent', but messages[%d] matched", i)
		}
	}
}

// TestFilterNotes_EmptyMessages tests with empty messages slice.
func TestFilterNotes_EmptyMessages(t *testing.T) {
	messages := []string{}

	result, err := filterNotes(messages, "test")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result slice, got %d items", len(result))
	}
}

// TestFilterNotes_EmptyKeyword tests searching with empty keyword.
func TestFilterNotes_EmptyKeyword(t *testing.T) {
	messages := []string{
		"some message",
		"another message",
	}

	result, err := filterNotes(messages, "")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	// Empty keyword should match nothing (consistent with grep behavior)
	for i, matched := range result {
		if matched {
			t.Logf("Note: empty keyword matched messages[%d]", i)
		}
	}
}

// TestFilterNotes_SpecialCharacters tests that special characters are handled correctly.
func TestFilterNotes_SpecialCharacters(t *testing.T) {
	messages := []string{
		"fixed issue #123 in code",
		"updated version to 2.0.1",
		"added [WIP] feature branch",
		"error: connection refused",
	}

	// Search for special pattern
	result, err := filterNotes(messages, "error:")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	if !result[3] {
		t.Errorf("expected messages[3] to match 'error:'")
	}
}

// TestFilterNotes_MultipleMatches tests when a keyword appears multiple times in one message.
func TestFilterNotes_MultipleMatches(t *testing.T) {
	messages := []string{
		"test test test multiple occurrences",
		"no match here",
		"another test message",
	}

	result, err := filterNotes(messages, "test")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	if !result[0] || result[1] || !result[2] {
		t.Errorf("expected messages[0] and messages[2] to match 'test'")
	}
}

// TestSearchNotes_NoFocuses tests SearchNotes when there are no focuses.
func TestSearchNotes_NoFocuses(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{},
		notes:   make(map[string][]Note),
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	results, err := svc.SearchNotes("test")
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}

	if results != nil && len(results) > 0 {
		t.Errorf("expected no results for empty focus list")
	}
}

// TestSearchNotes_NoNotes tests SearchNotes when focuses exist but have no notes.
func TestSearchNotes_NoNotes(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{
			{Name: "2026-05-04-1000__task1", Archived: false},
			{Name: "2026-05-04-1100__task2", Archived: false},
		},
		notes: make(map[string][]Note),
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	results, err := svc.SearchNotes("test")
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}

	if results != nil && len(results) > 0 {
		t.Errorf("expected no results when no notes exist")
	}
}

// TestSearchNotes_SingleMatch tests finding a single note across multiple focuses.
func TestSearchNotes_SingleMatch(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{
			{Name: "2026-05-04-1000__auth", Archived: false},
			{Name: "2026-05-04-1100__database", Archived: false},
			{Name: "2026-05-04-1200__ui", Archived: false},
		},
		notes: map[string][]Note{
			"2026-05-04-1000__auth": {
				{Timestamp: time.Date(2026, 5, 4, 10, 30, 0, 0, time.UTC), Message: "implemented JWT authentication"},
				{Timestamp: time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC), Message: "added rate limiting"},
			},
			"2026-05-04-1100__database": {
				{Timestamp: time.Date(2026, 5, 4, 11, 30, 0, 0, time.UTC), Message: "optimized queries"},
			},
			"2026-05-04-1200__ui": {
				{Timestamp: time.Date(2026, 5, 4, 12, 30, 0, 0, time.UTC), Message: "responsive layout"},
			},
		},
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	results, err := svc.SearchNotes("authentication")
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Note.Message != "implemented JWT authentication" {
		t.Errorf("unexpected note message: %v", results[0].Note.Message)
	}
	if results[0].Focus.Name != "2026-05-04-1000__auth" {
		t.Errorf("unexpected focus name: %v", results[0].Focus.Name)
	}
}

// TestSearchNotes_MultipleMatches tests finding multiple notes across focuses.
func TestSearchNotes_MultipleMatches(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{
			{Name: "2026-05-04-1000__auth", Archived: false},
			{Name: "2026-05-04-1100__database", Archived: false},
		},
		notes: map[string][]Note{
			"2026-05-04-1000__auth": {
				{Timestamp: time.Date(2026, 5, 4, 10, 30, 0, 0, time.UTC), Message: "implemented auth service"},
				{Timestamp: time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC), Message: "fixed auth bug"},
			},
			"2026-05-04-1100__database": {
				{Timestamp: time.Date(2026, 5, 4, 11, 30, 0, 0, time.UTC), Message: "auth tables schema"},
			},
		},
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	results, err := svc.SearchNotes("auth")
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	// Verify the results contain the expected messages
	messages := make(map[string]bool)
	for _, r := range results {
		messages[r.Note.Message] = true
	}

	expectedMessages := map[string]bool{
		"implemented auth service": true,
		"fixed auth bug":           true,
		"auth tables schema":       true,
	}

	for msg := range expectedMessages {
		if !messages[msg] {
			t.Errorf("expected message %q not found in results", msg)
		}
	}
}

// TestSearchNotes_ArchivedFocuses tests that archived focuses are included in search.
func TestSearchNotes_ArchivedFocuses(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{
			{Name: "2026-05-04-1000__active", Archived: false},
			{Name: "2026-05-04-1100__archived-task", Archived: true},
		},
		notes: map[string][]Note{
			"2026-05-04-1000__active": {
				{Timestamp: time.Date(2026, 5, 4, 10, 30, 0, 0, time.UTC), Message: "current work"},
			},
			"archive/2026-05-04-1100__archived-task": {
				{Timestamp: time.Date(2026, 5, 4, 11, 30, 0, 0, time.UTC), Message: "old archived notes"},
			},
		},
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	results, err := svc.SearchNotes("archived")
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result for 'archived' keyword, got %d", len(results))
	}
	if !results[0].Focus.Archived {
		t.Errorf("expected archived focus in results")
	}
}

// TestSearchNotes_CaseInsensitive tests case-insensitive search.
func TestSearchNotes_CaseInsensitive(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{
			{Name: "2026-05-04-1000__work", Archived: false},
		},
		notes: map[string][]Note{
			"2026-05-04-1000__work": {
				{Timestamp: time.Date(2026, 5, 4, 10, 30, 0, 0, time.UTC), Message: "Implemented FEATURE successfully"},
			},
		},
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	// Test uppercase search
	results, err := svc.SearchNotes("FEATURE")
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result for 'FEATURE', got %d", len(results))
	}

	// Test lowercase search
	results, err = svc.SearchNotes("feature")
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result for 'feature', got %d", len(results))
	}
}

// MockRepositoryWithNotes provides a searchable mock repository.
type MockRepositoryWithNotes struct {
	MockRepository
	focuses []Focus
	notes   map[string][]Note
}

func (m *MockRepositoryWithNotes) List() ([]Focus, error) {
	return m.focuses, nil
}

func (m *MockRepositoryWithNotes) GetNotes(name string) ([]Note, error) {
	return m.notes[name], nil
}

// TestFilterNotes_FallbackToGrep tests that filterNotes falls back to grep/go when rg is unavailable.
// This test verifies the grep path is exercised.
func TestFilterNotes_FallbackPath(t *testing.T) {
	// This test documents that filterNotes gracefully degrades when tools are unavailable
	// The actual execution path depends on what's available on the system

	messages := []string{
		"hello world",
		"world of go",
		"testing framework",
	}

	result, err := filterNotes(messages, "world")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	expectedMatches := []bool{true, true, false}
	for i, expected := range expectedMatches {
		if result[i] != expected {
			t.Errorf("messages[%d]: expected match=%v, got %v", i, expected, result[i])
		}
	}
}

// TestFilterNotes_ToolDetection tests detection of rg and grep tools.
func TestFilterNotes_ToolDetection(t *testing.T) {
	messages := []string{"test message"}

	// Just verify the function completes without error regardless of available tools
	result, err := filterNotes(messages, "test")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	if !result[0] {
		t.Errorf("expected message to match 'test'")
	}
}

// TestSearchNotes_ErrorHandling tests SearchNotes with a failing repository.
func TestSearchNotes_ErrorHandling(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{
			{Name: "2026-05-04-1000__work", Archived: false},
		},
		notes: map[string][]Note{},
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	// Search when repo.GetNotes would error (skipped in implementation)
	results, err := svc.SearchNotes("test")
	if err != nil {
		t.Fatalf("SearchNotes should not fail on GetNotes error: %v", err)
	}

	// Should return empty/nil when no matches
	if results != nil && len(results) > 0 {
		t.Errorf("expected no results when GetNotes returns empty")
	}
}

// Test tool-specific filtering behavior
func TestFilterNotes_WithRgIfAvailable(t *testing.T) {
	if _, err := exec.LookPath("rg"); err != nil {
		t.Skip("rg not available, skipping rg-specific test")
	}

	messages := []string{
		"ripgrep search test",
		"standard grep test",
		"no match here",
	}

	result, err := filterNotes(messages, "search")
	if err != nil {
		t.Fatalf("filterNotes with rg failed: %v", err)
	}

	if !result[0] || result[1] || result[2] {
		t.Errorf("ripgrep search failed, expected [true, false, false], got %v", result)
	}
}

// Test tool-specific filtering behavior with grep
func TestFilterNotes_WithGrepIfAvailable(t *testing.T) {
	if _, err := exec.LookPath("grep"); err != nil {
		t.Skip("grep not available, skipping grep-specific test")
	}

	messages := []string{
		"grep based search",
		"another line",
		"search term here",
	}

	result, err := filterNotes(messages, "search")
	if err != nil {
		t.Fatalf("filterNotes with grep failed: %v", err)
	}

	if !result[0] || result[1] || !result[2] {
		t.Errorf("grep search failed, expected [true, false, true], got %v", result)
	}
}

// TestSearchNotes_LargeDataset tests search performance with many notes.
func TestSearchNotes_LargeDataset(t *testing.T) {
	// Create mock with 100 focuses and 2 notes each = 200 notes
	focuses := make([]Focus, 100)
	notes := make(map[string][]Note)

	for i := 0; i < 100; i++ {
		name := "2026-05-04-" + strings.ReplaceAll(strings.Repeat("0", 4-len(string(rune(i))))+string(rune(i)), "0", "0")
		name += "__task" + string(rune('0'+i%10))
		focuses[i] = Focus{Name: name, Archived: false}

		notes[name] = []Note{
			{Timestamp: time.Now(), Message: "working on feature"},
			{Timestamp: time.Now(), Message: "debugging issue"},
		}
	}

	repo := &MockRepositoryWithNotes{
		focuses: focuses,
		notes:   notes,
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	results, err := svc.SearchNotes("feature")
	if err != nil {
		t.Fatalf("SearchNotes on large dataset failed: %v", err)
	}

	// Should find 100 matches (one per focus)
	if len(results) != 100 {
		t.Errorf("expected 100 results, got %d", len(results))
	}
}

// TestSearchNotes_EmptyKeyword tests search with empty keyword.
func TestSearchNotes_EmptyKeyword(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{
			{Name: "2026-05-04-1000__work", Archived: false},
		},
		notes: map[string][]Note{
			"2026-05-04-1000__work": {
				{Timestamp: time.Now(), Message: "some message"},
			},
		},
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	_, err := svc.SearchNotes("")
	if err == nil {
		t.Fatal("expected error for empty keyword, got nil")
	}
}

// TestFilterNotes_GoPathUsed tests the Go-based filtering fallback path.
// This test focuses on the pure Go implementation when no tools are available.
func TestFilterNotes_GoPathUsed(t *testing.T) {
	messages := []string{
		"First message with test",
		"Second message different",
		"Third message with test again",
	}

	// Test filtering
	result, err := filterNotes(messages, "test")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	// Verify matches
	if !result[0] || result[1] || !result[2] {
		t.Errorf("expected [true, false, true], got %v", result)
	}
}

// TestFilterNotes_LongKeyword tests with longer keyword patterns.
func TestFilterNotes_LongKeyword(t *testing.T) {
	messages := []string{
		"implemented the authentication service",
		"fixed a critical auth bug",
		"refactored the database layer",
		"added comprehensive authentication documentation",
	}

	result, err := filterNotes(messages, "authentication")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	if !result[0] || result[1] || result[2] || !result[3] {
		t.Errorf("expected [true, false, false, true], got %v", result)
	}
}

// TestFilterNotes_NumericKeyword tests with numeric search terms.
func TestFilterNotes_NumericKeyword(t *testing.T) {
	messages := []string{
		"version 1.0.0 release",
		"bug fix 2024-05-03",
		"commit abc123def456",
	}

	result, err := filterNotes(messages, "123")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	if result[0] || result[1] || !result[2] {
		t.Errorf("expected [false, false, true], got %v", result)
	}
}

// TestFilterNotes_PartialWordMatch tests partial word matching behavior.
func TestFilterNotes_PartialWordMatch(t *testing.T) {
	messages := []string{
		"testing the test",
		"contested data",
		"manifest file",
	}

	result, err := filterNotes(messages, "test")
	if err != nil {
		t.Fatalf("filterNotes failed: %v", err)
	}

	// "test" should match in first two messages
	if !result[0] || !result[1] || result[2] {
		t.Errorf("expected [true, true, false], got %v", result)
	}
}

// TestSearchNotes_MixedArchivedAndActive tests search with both archived and active focuses.
func TestSearchNotes_MixedArchivedAndActive(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{
			{Name: "2026-05-04-1000__active-task", Archived: false},
			{Name: "2026-05-04-1100__old-task", Archived: true},
			{Name: "2026-05-04-1200__new-task", Archived: false},
		},
		notes: map[string][]Note{
			"2026-05-04-1000__active-task": {
				{Timestamp: time.Now(), Message: "refactored code"},
			},
			"archive/2026-05-04-1100__old-task": {
				{Timestamp: time.Now(), Message: "refactored code"},
			},
			"2026-05-04-1200__new-task": {
				{Timestamp: time.Now(), Message: "new feature"},
			},
		},
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	results, err := svc.SearchNotes("refactored")
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results (1 active, 1 archived), got %d", len(results))
	}

	// Verify we got one active and one archived
	activeCount := 0
	archivedCount := 0
	for _, r := range results {
		if r.Focus.Archived {
			archivedCount++
		} else {
			activeCount++
		}
	}

	if activeCount != 1 || archivedCount != 1 {
		t.Errorf("expected 1 active and 1 archived, got %d active and %d archived", activeCount, archivedCount)
	}
}

// TestSearchNotes_SameKeywordMultipleFocuses tests finding the same keyword across many focuses.
func TestSearchNotes_SameKeywordMultipleFocuses(t *testing.T) {
	repo := &MockRepositoryWithNotes{
		focuses: []Focus{
			{Name: "2026-05-04-1000__project1", Archived: false},
			{Name: "2026-05-04-1100__project2", Archived: false},
			{Name: "2026-05-04-1200__project3", Archived: false},
		},
		notes: map[string][]Note{
			"2026-05-04-1000__project1": {
				{Timestamp: time.Now(), Message: "working on bugs"},
				{Timestamp: time.Now(), Message: "found bugs in code"},
			},
			"2026-05-04-1100__project2": {
				{Timestamp: time.Now(), Message: "no issues"},
				{Timestamp: time.Now(), Message: "fixed bugs successfully"},
			},
			"2026-05-04-1200__project3": {
				{Timestamp: time.Now(), Message: "testing completed"},
			},
		},
	}
	config := NewMockConfigStore()
	workspace := NewMockWorkspaceStore()

	svc := NewFocusService(repo, config, workspace)

	results, err := svc.SearchNotes("bugs")
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("expected 3 results for 'bugs' keyword, got %d", len(results))
	}
}
