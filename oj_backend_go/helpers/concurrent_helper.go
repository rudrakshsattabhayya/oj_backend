package helpers

import (
	"mime/multipart"
	"sync"
)

type FuncWithArgs struct {
	ID       string
	Function interface{}
	Args     []interface{}
}

type FuncResult struct {
	Res interface{}
	Err     error
}

func RunConcurrently(funcs []FuncWithArgs) map[string]FuncResult {
	var wg sync.WaitGroup
	results := make(map[string]FuncResult)
	var mu sync.Mutex

	wg.Add(len(funcs))

	for _, f := range funcs {
		go func(f FuncWithArgs) {
			defer wg.Done()
			var res1 interface{}
			var err error

			switch fn := f.Function.(type) {
			case func(...interface{}) (interface{}, interface{}):
				// res1, err = fn(f.Args...)
			case func(file *multipart.FileHeader, bucketName string) (string, error):
				if len(f.Args) == 2 {
					res1, err = fn(f.Args[0].(*multipart.FileHeader), f.Args[1].(string))
				}
			}

			mu.Lock()
			results[f.ID] = FuncResult{Res: res1, Err: err}
			mu.Unlock()
		}(f)
	}

	wg.Wait()
	return results
}
