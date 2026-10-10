Output demo.gif
Require go
Set Shell "zsh"
Set Theme "Dracula"
Set FontSize 18
Set Width 960
Set Height 480
Set Padding 18
Set WindowBar "Colorful"
Set BorderRadius 12
Set TypingSpeed 65ms

Hide
Type "clear; export HOME=$(mktemp -d /tmp/done-demo.XXXXXX); export GOCACHE=/tmp/done-gocache; export GOMODCACHE=/Users/sanaa/go/pkg/mod; go build -o /tmp/done-demo-bin ./cmd/done; clear; env -u NO_COLOR TERM=xterm-256color COLORTERM=truecolor /tmp/done-demo-bin"
Enter
Sleep 2s
Show
Sleep 500ms
Type "a"
Type "get housewarming gift"
Enter
Sleep 500ms
Type "a"
Type "call mom"
Enter
Sleep 500ms
Type "a"
Type "write a letter"
Enter
Sleep 500ms
Type "a"
Type "place ikea order"
Enter
Sleep 500ms
Up
Sleep 350ms
Up
Sleep 350ms
Up
Sleep 500ms
Down
Sleep 350ms
Down
Sleep 350ms
Down
Sleep 500ms
Space
Sleep 2s
Type "q"
Sleep 3s
