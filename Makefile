.PHONY: 

snapshot:
	goreleaser --snapshot --skip-publish --rm-dist
