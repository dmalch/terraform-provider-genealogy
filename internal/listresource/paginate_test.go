package listresource

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	. "github.com/onsi/gomega"
)

func TestPaginate(t *testing.T) {
	t.Run("Yields every element across pages in order", func(t *testing.T) {
		RegisterTestingT(t)
		pages := map[int][]string{
			1: {"a", "b"},
			2: {"c"},
		}
		fetchPage := func(_ context.Context, page int) ([]string, int, error) {
			return pages[page], 3, nil
		}
		project := func(s string) (list.ListResult, bool) {
			return list.ListResult{DisplayName: s}, true
		}
		onError := func(err error) list.ListResult {
			return list.ListResult{Diagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic("fetch failed", err.Error()),
			}}
		}

		seq := Paginate(t.Context(), fetchPage, onError, project)

		var got []string
		for r := range seq {
			got = append(got, r.DisplayName)
		}
		Expect(got).To(Equal([]string{"a", "b", "c"}))
	})

	t.Run("Stops calling fetchPage when consumer cancels", func(t *testing.T) {
		RegisterTestingT(t)
		var fetched []int
		fetchPage := func(_ context.Context, page int) ([]string, int, error) {
			fetched = append(fetched, page)
			return []string{fmt.Sprintf("item-%d-1", page), fmt.Sprintf("item-%d-2", page)}, 1000, nil
		}
		project := func(s string) (list.ListResult, bool) {
			return list.ListResult{DisplayName: s}, true
		}
		onError := func(err error) list.ListResult {
			return list.ListResult{Diagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic("fetch failed", err.Error()),
			}}
		}

		seq := Paginate(t.Context(), fetchPage, onError, project)

		count := 0
		for range seq {
			count++
			if count == 1 {
				break
			}
		}

		Expect(count).To(Equal(1))
		Expect(fetched).To(Equal([]int{1}))
	})

	t.Run("Emits one error diagnostic and stops when fetchPage returns an error", func(t *testing.T) {
		RegisterTestingT(t)
		boom := errors.New("api exploded")
		fetchPage := func(_ context.Context, _ int) ([]string, int, error) {
			return nil, 0, boom
		}
		project := func(s string) (list.ListResult, bool) {
			return list.ListResult{DisplayName: s}, true
		}
		onError := func(err error) list.ListResult {
			return list.ListResult{Diagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic("fetch failed", err.Error()),
			}}
		}

		seq := Paginate(t.Context(), fetchPage, onError, project)

		var results []list.ListResult
		for r := range seq {
			results = append(results, r)
		}

		Expect(results).To(HaveLen(1))
		Expect(results[0].Diagnostics.HasError()).To(BeTrue())
		Expect(results[0].Diagnostics[0].Detail()).To(ContainSubstring("api exploded"))
	})

	t.Run("Stops on an empty page even when total is non-zero", func(t *testing.T) {
		RegisterTestingT(t)
		fetchPage := func(_ context.Context, _ int) ([]string, int, error) {
			return []string{}, 5, nil
		}
		project := func(s string) (list.ListResult, bool) {
			return list.ListResult{DisplayName: s}, true
		}
		onError := func(err error) list.ListResult {
			return list.ListResult{Diagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic("fetch failed", err.Error()),
			}}
		}

		seq := Paginate(t.Context(), fetchPage, onError, project)

		count := 0
		for range seq {
			count++
		}
		Expect(count).To(Equal(0))
	})
}
