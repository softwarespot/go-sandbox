# ssh service

## URLs

- https://pkg.go.dev/github.com/gliderlabs/ssh
- https://github.com/shazow/ssh-chat
- https://medium.com/@alexfoleydevops/building-an-ssh-chatroom-with-go-6df65facd6cb
- https://shazow.net/posts/ssh-how-does-it-even/

## Usage

```
ssh -o "StrictHostKeyChecking=no" -o "UserKnownHostsFile=/dev/null" -o "ForwardAgent=no" -p 2222 localhost
```
