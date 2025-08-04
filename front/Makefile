setup:
	rm -rf ./node_modules
	bun install

dev:
	bun farm dev

build: setup license
	bun farm build

license:
	bun x license-checker-rseidelsohn --json \
		| jq 'to_entries | map({key: .key, value: (.value | del(.path, .licenseFile))}) | from_entries' \
		> LICENSE.json