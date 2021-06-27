#! /bin/zsh
go build -o ./air/refactoring ./cmd/refactoring 2> >(egrep -v '(NSUser)' | grep -v '^# fyne.io/fyne/v2/app$' >&2)
