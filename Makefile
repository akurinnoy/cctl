PREFIX ?= $(HOME)/.local/bin

.PHONY: build test check install clean

build:
	cd cctl-data && go build -o cctl-data .

test:
	cd cctl-data && go test ./...

check:
	@command -v jq >/dev/null  || { echo "Missing: jq  — brew install jq"; exit 1; }
	@command -v fzf >/dev/null || { echo "Missing: fzf — brew install fzf"; exit 1; }
	@[ -d /Applications/iTerm.app ] || { echo "Missing: iTerm2 — https://iterm2.com"; exit 1; }

install: build check
	mkdir -p $(PREFIX)
	install -m 755 cctl $(PREFIX)/cctl
	install -m 755 cctl-data/cctl-data $(PREFIX)/cctl-data
	chmod +x $(PREFIX)/cctl $(PREFIX)/cctl-data
	@echo "Installed cctl and cctl-data to $(PREFIX)"
	@echo "Make sure $(PREFIX) is in your PATH"

clean:
	rm -f cctl-data/cctl-data
