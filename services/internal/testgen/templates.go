package testgen

const JestTemplate = `import { %s } from '%s';

describe('%s', () => {
  %s
});
`

const TestCaseTemplate = `it('should %s', () => {
    %s
  });`

const AsyncTestCaseTemplate = `it('should %s', async () => {
    %s
  });`

func GenerateTestSkeleton(funcName, filePath string, isAsync bool) string {
	testCase := "// Add test implementation"
	template := TestCaseTemplate
	
	if isAsync {
		template = AsyncTestCaseTemplate
	}

	return template
}
